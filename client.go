package resource

import "github.com/digitalocean/concourse-resource-library/artifactory"

func newClient(s Source) (*artifactory.Client, error) {
	return artifactory.NewClient(
		artifactory.Endpoint(s.Endpoint),
		artifactory.Authentication(s.User, s.Password, s.APIKey, s.AccessToken),
	)
}
