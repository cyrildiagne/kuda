package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/cyrildiagne/kuda/pkg/auth"

	"github.com/gorilla/mux"
)

var authRedirectServer *http.Server
var user *auth.User

func enableCORS(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	decoder := json.NewDecoder(r.Body)
	user = &auth.User{}
	err := decoder.Decode(user)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		// Shutdown server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := authRedirectServer.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
}

func startServer(port string) (*auth.User, error) {
	r := mux.NewRouter()
	r.HandleFunc("/", handleToken).Methods("POST")
	r.HandleFunc("/", enableCORS).Methods("OPTIONS")

	authRedirectServer = &http.Server{Addr: ":" + port, Handler: r}

	if err := authRedirectServer.ListenAndServe(); err != nil {
		if user != nil {
			return user, nil
		}
		return nil, err
	}
	return nil, errors.New("could not retrieve user auth")
}

func startLoginFlow(authURL string) (*auth.User, error) {
	port := os.Getenv("KUDA_CLI_LOGIN_PORT")
	if port == "" {
		port = "8094"
	}
	// Append redirect command.
	authURL += "?cli=" + port
	// Run command.
	args := []string{authURL}
	cmd := exec.Command("open", args...)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	user, err := startServer(port)
	if err != nil {
		return nil, err
	}
	return user, nil
}
