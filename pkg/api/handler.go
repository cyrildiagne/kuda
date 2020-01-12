package api

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebaseAuth "firebase.google.com/go/auth"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Env stores our application-wide configuration.
type Env struct {
	GCPProjectID   string
	DockerRegistry string
	DB             *firestore.Client
	Auth           *firebaseAuth.Client
}

// GetDockerImagePath returns the fully qualified URL of a docker image on GCR
func (e *Env) GetDockerImagePath(im ImageName) string {
	return "gcr.io/" + e.GCPProjectID + "/" + im.GetID()
}

// Handler takes a configured Env and a function matching our signature.
type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			fmt.Fprintf(w, "%v\n", e.Error())
			http.Error(w, e.Error(), e.Status())
			break
		default:
			log.Printf("Internal Error - %s", e)
			fmt.Fprintf(w, "%v\n", e.Error())
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}
