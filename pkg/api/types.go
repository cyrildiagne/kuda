package api

import (
	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
)

// Env stores our application-wide configuration.
type Env struct {
	DB                DB
	ContainerRegistry ContainerRegistry
	ContainerBuilder  ContainerBuilder
	Auth              Auth
}

// DB is an interface for manipulating DB resources.
type DB interface {
	IsUserAdminOfNamespace(uid string, namespace string) (bool, error)

	UpdateAPIMetadata(imageID string, metadata *map[string]interface{}) error

	GetVersionnedAPI(ImageName) (*Version, error)
	UpdateVersionnedAPI(imageID string, version string, template *Version) error
}

// ContainerRegistry is an interface for Docker container registries.
type ContainerRegistry interface {
	GetDockerImagePath(ImageName) string
	ListImageTags(string) error
}

// Auth is an interface for authentication providers.
type Auth interface {
	VerifyIDToken(string) (string, error)
}

// ContainerBuilder is an interface for Docker container builders.
type ContainerBuilder interface {
	GetBuildType() *v1.BuildType
}

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

// Version stores an API version.
type Version struct {
	IsPublic bool             `firestore:"isPublic"`
	Version  string           `firestore:"version"`
	Manifest *latest.Manifest `firestore:"manifest"`
	// Paths    openapi.Paths    `firestore:"paths,omitempty"`
	// Paths    *openapi3.Swagger   `firestore:"openapi,omitempty"`
	// Paths openapi3.Paths `firestore:"openapi,omitempty"`
	// Paths    map[string]*openapi3.PathItem `firestore:"openapi,omitempty"`
	// Paths    *map[string]interface{} `firestore:"openapi,omitempty"`
}

// API stores an API.
type API struct {
	Author   string
	Name     string
	Image    string
	Versions []Version
}
