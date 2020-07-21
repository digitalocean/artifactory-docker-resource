package resource

import (
	"log"
	"path/filepath"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	"github.com/digitalocean/concourse-resource-library/docker"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

// Get performs the get operation for the resource
func Get(req GetRequest, dir string) (GetResponse, error) {
	var res GetResponse

	if req.Version.Empty() {
		return res, nil
	}

	c, err := newClient(req.Source)
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println(dir)

	err = c.PullImage(req.Version.Repo, req.Version.Image())
	if err != nil {
		log.Println(err)
		return res, err
	}

	log.Println("pulled image:", req.Version.Image())

	item, err := c.SearchItem(req.Version.ArtifactoryPath())
	if err != nil {
		rlog.StdErr("failed to search", err)
		log.Println(err)
		return res, err
	}

	log.Println("fetched metadata:", item)

	res.Metadata = metadata(
		artifactory.Artifact{
			File: utils.FileInfo{ArtifactoryPath: req.Version.ArtifactoryPath()},
			Item: item,
		})

	d, err := docker.NewClient()
	img, err := d.Image(req.Version.Image())
	if err != nil {
		rlog.StdErr("failed to get image details", err)
		log.Println(err)
		return res, err
	}

	if !req.Params.SkipDownload {
		err = d.Save(filepath.Join(dir, "image.tar"), img.ID)
		if err != nil {
			rlog.StdErr("failed to save image to disk", err)
			log.Println(err)
			return res, err
		}
	}

	// TODO: mount as rootfs for task steps

	return res, nil
}
