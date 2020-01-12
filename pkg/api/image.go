package api

import (
	"errors"
	"regexp"
	"strings"
)

// ImageName represent a docker container template name.
type ImageName struct {
	Author  string
	Name    string
	Version string
}

// ParseFrom sets an imagename fields from a fully qualified name.
func (im *ImageName) ParseFrom(imageName string) error {
	re := regexp.MustCompile(`(?P<Author>[a-z-_]+)\/(?P<Name>[a-z-_]+)(?P<Version>:[a-z0-9\-.]+)?`)
	match := re.FindStringSubmatch(imageName)

	params := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			params[name] = match[i]
		}
	}

	im.Author = params["Author"]
	if im.Author == "" {
		return errors.New("author is empty")
	}
	im.Name = params["Name"]
	if im.Author == "" {
		return errors.New("name is empty")
	}
	im.Version = params["Version"]
	if im.Version == "" {
		im.Version = "latest"
	} else if strings.HasPrefix(im.Version, ":") {
		im.Version = im.Version[1:]
	}

	return nil
}

// GetID returns the formated id from the author & name fields.
func (im *ImageName) GetID() string {
	return im.Author + "__" + im.Name
}
