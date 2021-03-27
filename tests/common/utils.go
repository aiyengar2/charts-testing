package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func ValidateRancherPrefixed(image string) bool {
	return strings.HasPrefix(image, "rancher/")
}

func ValidateSystemDefaultRegistry(image string, registry string) bool {
	return strings.HasPrefix(image, fmt.Sprintf("%s/", registry))
}

func ValidateSystemDefaultRegistryAndRancher(image string, registry string) bool {
	return strings.HasPrefix(image, fmt.Sprintf("%s/rancher/", registry))
}

func ValidateRancherMirrored(image string) bool {
	return strings.HasPrefix(image, "rancher/mirrored-")
}

func ValidateImageExists(image string) bool {
	splitImage := strings.Split(image, ":")
	url := fmt.Sprintf("https://registry.hub.docker.com/v1/repositories/%s/tags/%s", splitImage[0], splitImage[1])
	status := GetURL(url)
	return status == http.StatusOK
}

// GetURL performs an HTTP.Get on a endpoint and returns an error if the resp is not 200
func GetURL(url string) int {
	client := &http.Client{}
	response, err := client.Get(url)
	if err == nil {
		return response.StatusCode
	} else {
		logrus.Warnf("There was the following error in performing a GET on %s: %s\n", url, err)
		return 0
	}
}
