package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/utils"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	firebaseAuth "firebase.google.com/go/auth"
	"github.com/gorilla/mux"
)

var gcpProjectID string
var dockerRegistry string
var fsDb *firestore.Client
var fbAuth *firebaseAuth.Client

func checkAuthorized(namespace string, w http.ResponseWriter, r *http.Request) (int, error) {
	// Get bearer token.
	accessToken := r.Header.Get("Authorization")
	accessToken = strings.Split(accessToken, "Bearer ")[1]
	// Verify Token
	token, err := fbAuth.VerifyIDToken(context.Background(), accessToken)
	if err != nil {
		return 401, fmt.Errorf("error verifying token %v", err)
	}

	// Check if namespace has the user id as admin.
	ctx := context.Background()
	ns, err := fsDb.Collection("namespaces").Doc(namespace).Get(ctx)
	if err != nil {
		return 500, fmt.Errorf("error getting namespace info %v", err)
	}
	if !ns.Exists() {
		return 400, fmt.Errorf("namespace not found %v", namespace)
	}
	nsData := ns.Data()
	nsAdmins, hasAdmins := nsData["admins"]
	if !hasAdmins {
		return 403, fmt.Errorf("no admin found for namespace %v", namespace)
	}
	_, isAdmin := nsAdmins.(map[string]interface{})[token.UID]
	if !isAdmin {
		return 403, fmt.Errorf("user %v must be admin of %v", token.UID, namespace)
	}

	return 200, nil
}

func getNamespace(r *http.Request) (string, int, error) {
	// Retrieve namespace.
	namespace := r.FormValue("namespace")
	namespace = strings.ToValidUTF8(namespace, "")
	if namespace == "" {
		err := "error retrieving namespace"
		return "", 400, errors.New(err)
	}
	if namespace == "kuda" {
		err := "namespace cannot be kuda"
		return "", 403, errors.New(err)
	}
	return namespace, 200, nil
}

func handlePublish(w http.ResponseWriter, r *http.Request) {
	// Retrieve namespace.
	namespace, code, err := getNamespace(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Check authorizations.
	if code, err := checkAuthorized(namespace, w, r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// TODO: Check if image@version exists.
	// TODO: Mark image@version as public.
}

func handleDeploymentFromPublished(w http.ResponseWriter, r *http.Request) {
	// Retrieve namespace.
	namespace, code, err := getNamespace(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Check authorizations.
	if code, err := checkAuthorized(namespace, w, r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// TODO: Check if image@version exists.
	// TODO: Check if image@version is public.
	// TODO: Generate Knative YAML with appropriate namespace.
	// TODO: Run kubectl apply.
}

func handleDeployment(w http.ResponseWriter, r *http.Request) {
	// Set maximum upload size to 2GB.
	r.ParseMultipartForm((2 * 1000) << 20)

	// Retrieve requested namespace.
	namespace, code, err := getNamespace(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Check authorizations.
	if code, err := checkAuthorized(namespace, w, r); err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Retrieve Filename, Header and Size of the file.
	file, handler, err := r.FormFile("context")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error retrieving file", 500)
		return
	}
	defer file.Close()
	log.Printf("Building: %+v, %+v Ko, for namespace %v\n", handler.Filename, handler.Size/1024, namespace)

	// Create new temp directory.
	tempDir, err := ioutil.TempDir("", namespace)
	fmt.Println(tempDir)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating temp dir", 500)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up.

	// Extract file to temp directory.
	err = utils.Untar(tempDir, file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error extracting content", 500)
		return
	}

	// Load the manifest.
	manifestFile := filepath.FromSlash(tempDir + "/kuda.yaml")
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not load manifest", 400)
		return
	}

	// TODO: replace namespace by user ID.
	dockerArtifact := dockerRegistry + "/" + namespace + "__" + manifest.Name

	// Generate Skaffold & Knative config files.
	service := config.ServiceSummary{
		Name:           manifest.Name,
		Namespace:      namespace,
		DockerArtifact: dockerArtifact,
		BuildType: v1.BuildType{
			GoogleCloudBuild: &v1.GoogleCloudBuild{
				ProjectID: gcpProjectID,
			},
		},
	}
	folder := filepath.FromSlash(tempDir + "/.kuda")
	skaffoldFile, err := utils.GenerateSkaffoldConfigFiles(service, manifest.Deploy, folder)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not generate config files", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/event-stream")

	if err := RunSkaffold(tempDir, skaffoldFile, w); err != nil {
		http.Error(w, fmt.Sprintf("error running skaffold: %v", err), 500)
		return
	}

	// TODO: add API entry to the APIs base:
	// { meta, user, image, versions[ {tag, public, openapi},...] }

	fmt.Fprintf(w, "Deployment successful!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello!\n")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initGCP() {
	// Authenticate gcloud using application credentials.
	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file",
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error authenticating with credentials. %v\n", err)
	}

	// Get kubeconfig.
	args := []string{"container", "clusters", "get-credentials",
		"--project", gcpProjectID,
		"--region", "us-central1-a", "kuda"}
	cmd = exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("could not retrieve kubectl credentials %v\n", err)
	}
}

func initFirebase() (*firebaseAuth.Client, *firestore.Client) {
	config := &firebase.Config{ProjectID: gcpProjectID}
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting auth client: %v\n", err)
	}

	fs, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error connecting to firestore: %v\n", err)
	}

	return auth, fs
}

func main() {
	gcpProjectID = os.Getenv("KUDA_GCP_PROJECT")
	if gcpProjectID == "" {
		panic("cloud not load env var KUDA_GCP_PROJECT")
	}
	log.Println("Using project:", gcpProjectID)

	dockerRegistry = "gcr.io/" + gcpProjectID
	log.Println("Using registry:", dockerRegistry)

	initGCP()

	auth, fs := initFirebase()
	fbAuth = auth
	fsDb = fs

	port := getEnv("PORT", "8080")
	fmt.Println("Starting deployer on port", port)

	r := mux.NewRouter()
	r.HandleFunc("/", hello).Methods("GET")
	r.HandleFunc("/", handleDeployment).Methods("POST")
	http.ListenAndServe(":"+port, r)
}
