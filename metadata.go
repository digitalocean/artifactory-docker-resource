package resource

import (
	"strconv"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	meta "github.com/digitalocean/concourse-resource-library/metadata"
)

func metadata(a artifactory.Artifact) meta.Metadata {
	var m meta.Metadata

	m.Add("artifactory-path", a.File.ArtifactoryPath)

	m.Add("created", a.Item.Created)
	m.Add("modified", a.Item.Modified)
	m.Add("name", a.Item.Name)
	m.Add("repo", a.Item.Repo)
	m.Add("size", strconv.FormatInt(a.Item.Size, 10))
	m.Add("type", a.Item.Type)

	m.AddJSON("properties", a.Item.Properties)

	return m
}
