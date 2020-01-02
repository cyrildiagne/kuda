package config

import (
	"strings"

	config "github.com/cyrildiagne/kuda/pkg/kuda/config"
)

// GetDockerfileArtifactName returns a consistent docker artifact name for a
// given API & user config.
func GetDockerfileArtifactName(userCfg config.UserConfig, apiName string) string {
	// Removes the "-dev" suffix if present.
	if strings.HasSuffix(apiName, "-dev") {
		apiName = apiName[:len(apiName)-4]
	}
	return userCfg.Deployer.Skaffold.DockerRegistry + "/" + apiName
}
