package resource

import "log"

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

	// write metadata to file for parsing

	// optionally save image to disk via `docker save`

	// optionally mount as rootfs for task steps

	return res, nil
}
