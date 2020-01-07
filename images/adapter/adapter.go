package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"google.golang.org/grpc"

	"istio.io/api/mixer/adapter/model/v1beta1"
	policy "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/pkg/status"
	"istio.io/istio/mixer/template/authorization"
	"istio.io/istio/mixer/template/metric"
	"istio.io/pkg/log"
)

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// Adapter supports the authorization template.
	Adapter struct {
		listener net.Listener
		server   *grpc.Server
	}
)

var _ authorization.HandleAuthorizationServiceServer = &Adapter{}
var _ metric.HandleMetricServiceServer = &Adapter{}

var fsClient *firestore.Client
var ctx context.Context

func decodeValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *policy.Value_StringValue:
		return t.StringValue
	case *policy.Value_Int64Value:
		return t.Int64Value
	case *policy.Value_DoubleValue:
		return t.DoubleValue
	case *policy.Value_IpAddressValue:
		return t.IpAddressValue
	case *policy.Value_TimestampValue:
		return t.TimestampValue
	default:
		return fmt.Sprintf("%v", in)
	}
}

func decodeValueMap(in map[string]*policy.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = decodeValue(v.GetValue())
	}
	return out
}

func decodeTimestamp(t interface{}) time.Time {
	if t == nil {
		return time.Time{}
	}
	tm := t.(*policy.TimeStamp)
	return time.Unix(tm.GetValue().Seconds, int64(tm.GetValue().Nanos))
}

// get the unix time for the given property value
func getUnixTimeStr(t time.Time) string {
	return fmt.Sprintf("%v", t.UnixNano())
}

// HandleAuthorization records metric entries
func (s *Adapter) HandleAuthorization(ctx context.Context, r *authorization.HandleAuthorizationRequest) (*v1beta1.CheckResult, error) {

	props := decodeValueMap(r.Instance.Subject.Properties)

	// Retrieve unique request ID.
	// requestID := props["request_id"].(string)
	// log.Infof("%v", requestID)

	// Get api_key from props
	apiKey := props["api_key"].(string)
	if apiKey == "" {
		return &v1beta1.CheckResult{
			Status: status.WithPermissionDenied("Unauthorized: API Key Missing"),
		}, nil
	}

	// Retrieve key in Firestore.
	keyEntry, err := fsClient.Collection("keys").Doc(apiKey).Get(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("Can't find key entry... %v", apiKey)
		return &v1beta1.CheckResult{
			Status: status.WithInvalidArgument(errMsg),
		}, nil
	}

	// Check if key has quotas.
	quotas := keyEntry.Data()["quotas"].(int64)
	if quotas < 1 {
		return &v1beta1.CheckResult{
			Status: status.WithResourceExhausted("Quotas exhausted."),
		}, nil
	}

	return &v1beta1.CheckResult{
		Status: status.OK,
	}, nil
}

// HandleMetric records metric entries
func (s *Adapter) HandleMetric(ctx context.Context, r *metric.HandleMetricRequest) (*v1beta1.ReportResult, error) {

	for _, inst := range r.Instances {
		dimensions := decodeValueMap(inst.Dimensions)

		// requestID := dimensions["request_id"].(string)

		apiKey := dimensions["api_key"].(string)

		responseCode := dimensions["response_code"].(int64)

		authorized := responseCode != 403
		if !authorized {
			return &v1beta1.ReportResult{}, nil
		}

		requestTimestamp := decodeTimestamp(dimensions["request_timestamp"])
		requestMethod := dimensions["request_method"].(string)
		requestHost := dimensions["request_host"].(string)
		requestURL := dimensions["request_url"].(string)

		req := map[string]interface{}{
			"apiKey":           apiKey,
			"requestMethod":    requestMethod,
			"requestHost":      requestHost,
			"requestURL":       requestURL,
			"requestTimestamp": requestTimestamp,
		}

		// authorized := responseCode != 403 && responseCode != 400 && responseCode != 429
		// req["authorized"] = authorized

		ipAddr := dimensions["user_ip"]
		ipValue := ipAddr.(*policy.IPAddress).Value
		req["userIP"] = net.IP(ipValue).String()
		req["userAgent"] = dimensions["user_agent"].(string)

		responseTimestamp := decodeTimestamp(dimensions["response_timestamp"])
		req["responseTimestamp"] = responseTimestamp
		req["responseSize"] = dimensions["response_size"].(int64)
		req["responseCode"] = dimensions["response_code"].(int64)

		// Add new request entry.
		_, _, err := fsClient.Collection("requests").Add(ctx, req)
		if err != nil {
			log.Errorf("add request: %s", err)
			continue
		}

		// Don't update credit if request was not authorized.
		if !authorized {
			continue
		}

		// Update quotas.
		keyEntry, err := fsClient.Collection("keys").Doc(apiKey).Get(ctx)
		if err != nil {
			log.Errorf("credit update: can't find key %v", apiKey)
			continue
		}
		quotas := keyEntry.Data()["quotas"].(int64)
		_, updateErr := keyEntry.Ref.Set(ctx, map[string]interface{}{
			"quotas": quotas - 1,
		}, firestore.MergeAll)
		if updateErr != nil {
			log.Errorf("credit update: %v", updateErr)
			continue
		}
	}

	return &v1beta1.ReportResult{}, nil
}

// Addr returns the listening address of the server
func (s *Adapter) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *Adapter) Run(shutdown chan error) {

	// Connect to Firestore
	ctx = context.Background()
	credentialsJSON := os.Getenv("FIRESTORE_CREDENTIALS")
	log.Infof("credentials file: %v", credentialsJSON)

	_, err := os.Stat(credentialsJSON)
	if os.IsNotExist(err) {
		log.Errorf("could not find credentials: %v", err)
	} else if err != nil {
		log.Errorf("error with crendentials: %v", err)
	} else {
		log.Infof("Found credentials: %v", credentialsJSON)
	}

	opt := option.WithCredentialsFile(os.Getenv("FIRESTORE_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		shutdown <- fmt.Errorf("error initializing firestore sdk: %v", err)
	}
	var fsErr error
	fsClient, fsErr = app.Firestore(ctx)
	if fsErr != nil {
		shutdown <- fmt.Errorf("error connecting to firestore %s", fsErr)
	}

	shutdown <- s.server.Serve(s.listener)
}

// Close gracefully shuts down the server; used for testing
func (s *Adapter) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	fsClient.Close()

	return nil
}

// NewAdapter creates a new IBP adapter that listens at provided port.
func NewAdapter(addr string) (Server, error) {
	// Create server
	if addr == "" {
		addr = "0"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}
	s := &Adapter{
		listener: listener,
	}
	fmt.Printf("listening on \"%v\"\n", s.Addr())
	s.server = grpc.NewServer()

	authorization.RegisterHandleAuthorizationServiceServer(s.server, s)
	metric.RegisterHandleMetricServiceServer(s.server, s)

	return s, nil
}
