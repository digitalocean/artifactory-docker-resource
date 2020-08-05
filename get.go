package resource

import (
	"log"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

// Get performs the get operation for the resource
func Get(req GetRequest, dir string) (GetResponse, error) {
	var res GetResponse

	if req.Version.Empty() {
		log.Println("request version is empty")
		return res, nil
	}

	c, err := newClient(req.Source)
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println(dir)

	item, err := c.SearchItem(req.Version.ArtifactoryPath())
	if err != nil {
		rlog.StdErr("failed to search", err)
		log.Println(err)
		return res, err
	}

	res = GetResponse{
		Version: req.Version,
		Metadata: metadata(
			artifactory.Artifact{
				File: utils.FileInfo{ArtifactoryPath: req.Version.ArtifactoryPath()},
				Item: item,
			}),
	}

	log.Println("fetched metadata:", item)

	if req.Params.SkipDownload {
		return res, nil
	}

	err = c.PullImage(dir, req.Params.Format, req.OCIRepository(), req.Version.Tag, req.Version.Digest)
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println("pulled image:", req.Version.ImageTag())

	return res, nil
}
