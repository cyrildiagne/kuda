package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var authPage string

// AuthConfig represents the AuthConfig Document.
type AuthConfig struct {
	APIKey            string
	AuthDomain        string
	TermsOfServiceURL template.URL
	PrivacyPolicyURL  template.URL
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, authPage)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Retrieve the auth env variables.
	config := AuthConfig{
		APIKey:            os.Getenv("KUDA_AUTH_API_KEY"),
		AuthDomain:        os.Getenv("KUDA_AUTH_DOMAIN"),
		TermsOfServiceURL: template.URL(os.Getenv("KUDA_AUTH_TOS_URL")),
		PrivacyPolicyURL:  template.URL(os.Getenv("KUDA_AUTH_PP_URL")),
	}

	// Process template with values.
	t, err := template.ParseFiles("./public/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w := new(bytes.Buffer)
	t.Execute(w, config)
	authPage = w.String()

	// Setup static serving.
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Listening on port", port)
	err = http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
