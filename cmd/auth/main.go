package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var authPage string
var config AuthConfig

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

func refreshToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("refreshing...")
	fmt.Println(config.APIKey)
	endpoint := "https://securetoken.googleapis.com/v1/token?key=" + config.APIKey

	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		http.Error(w, "refresh_token missing", 400)
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "unknown error", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	fmt.Fprint(w, bodyString)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Retrieve the auth env variables.
	config = AuthConfig{
		APIKey:            os.Getenv("KUDA_AUTH_API_KEY"),
		AuthDomain:        os.Getenv("KUDA_AUTH_DOMAIN"),
		TermsOfServiceURL: template.URL(os.Getenv("KUDA_AUTH_TOS_URL")),
		PrivacyPolicyURL:  template.URL(os.Getenv("KUDA_AUTH_PP_URL")),
	}

	staticFolder := os.Getenv("STATIC_FOLDER")
	if staticFolder == "" {
		staticFolder = filepath.FromSlash("./web/auth")
	}

	// Process template with values.
	indexFile := filepath.FromSlash(staticFolder + "/index.html")
	t, err := template.ParseFiles(indexFile)
	if err != nil {
		log.Fatal(err)
	}
	w := new(bytes.Buffer)
	t.Execute(w, config)
	authPage = w.String()

	// Setup static serving.
	fileServer := http.FileServer(http.Dir(staticFolder))
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	// Setup refresh token endpoint
	mux.HandleFunc("/refresh", refreshToken)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Listening on port", port)
	err = http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}
