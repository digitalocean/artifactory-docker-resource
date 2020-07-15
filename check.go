package resource

import (
	"log"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

// Check performs the check operation for the resource
func Check(req CheckRequest) (CheckResponse, error) {
	c, err := newClient(req.Source)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Source.Image.SetModifiedTime(req.Version)

	log.Println("query:", req.Source.Image.Raw)

	data, err := c.SearchItems(req.Source.Image.Raw)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := processItems(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res = selectVersions(req.Version, res)

	log.Println("version count in response:", len(res))
	log.Println("versions:", res)

	return res, nil
}

func processItems(s []utils.ResultItem) (CheckResponse, error) {
	var res CheckResponse

	for _, i := range s {
		v, err := processItem(i)
		if err != nil {
			return nil, err
		}

		res = append(res, v)
	}

	return res, nil
}

func processItem(i utils.ResultItem) (Version, error) {
	var v Version

	m, err := time.Parse(time.RFC3339, i.Modified)
	if err != nil {
		return v, err
	}

	var owner, name, tag string
	for _, prop := range i.Properties {
		switch prop.Key {
		case "docker.repoName":
			v := strings.Split(prop.Value, "/")
			owner = v[0]
			name = v[1]
		case "docker.manifest":
			tag = prop.Value
		}
	}

	v = Version{Repo: i.Repo, Owner: owner, Name: name, Tag: tag, Modified: m}

	return v, nil
}

// selectVersions handles business logic based on input version
// 	from Concourse and versions found in external resource
func selectVersions(v Version, res CheckResponse) CheckResponse {
	// If there are no new but an input version, return the input
	if len(res) == 0 && v.Repo != "" {
		log.Println("no new versions, use input version")
		res = append(res, v)

	}

	// If there are new versions and no input version, return latest new version
	if len(res) != 0 && v.Repo == "" {
		log.Println("new versions but no input version, use latest")
		res = CheckResponse{res[len(res)-1]}
	}

	return res
}
