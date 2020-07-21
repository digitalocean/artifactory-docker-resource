package resource

import (
	meta "github.com/digitalocean/concourse-resource-library/metadata"
)

// Put performs the Put operation for the resource
func Put(req PutRequest, dir string) (GetResponse, error) {
	get := GetResponse{
		Version:  Version{},
		Metadata: meta.Metadata{},
	}

	return get, nil
}
