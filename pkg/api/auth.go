package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetAuthorizedNamespace returns a namespace only if user is admin.
func GetAuthorizedNamespace(env *Env, r *http.Request) (string, error) {
	// Retrieve namespace.
	namespace := r.Header.Get("x-kuda-namespace")
	namespace = strings.ToValidUTF8(namespace, "")
	if namespace == "" {
		err := errors.New("error retrieving namespace")
		return "", StatusError{400, err}
	}
	if namespace == "kuda" {
		err := errors.New("namespace cannot be kuda")
		return "", StatusError{403, err}
	}

	// Check authorizations.
	accessToken := r.Header.Get("Authorization")
	if err := CheckAuthorized(env, namespace, accessToken); err != nil {
		return "", err
	}

	return namespace, nil
}

// CheckAuthorized checks if a user is authorized to update a namespace.
func CheckAuthorized(env *Env, namespace string, accessToken string) error {
	// Get bearer token.
	accessToken = strings.Split(accessToken, "Bearer ")[1]
	// Verify Token
	UID, err := env.Auth.VerifyIDToken(accessToken)
	if err != nil {
		err = fmt.Errorf("error verifying token %v", err)
		return StatusError{401, err}
	}

	// Check if namespace has the user id as admin.
	isAdmin, err := env.DB.IsUserAdminOfNamespace(UID, namespace)
	if err != nil {
		return err
	}
	if !isAdmin {
		err := fmt.Errorf("user %v must be admin of %v", UID, namespace)
		return StatusError{403, err}
	}

	return nil
}
