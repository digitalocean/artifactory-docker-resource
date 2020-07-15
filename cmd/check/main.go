package main

import (
	"log"

	resource "github.com/digitalocean/artifactory-docker-resource"
	rlog "github.com/digitalocean/concourse-resource-library/log"
	jlog "github.com/jfrog/jfrog-client-go/utils/log"
)

func main() {
	input := rlog.WriteStdin()
	defer rlog.Close()

	jlog.SetLogger(jlog.NewLogger(jlog.DEBUG, log.Writer()))

	var request resource.CheckRequest
	err := request.Read(input)
	if err != nil {
		log.Fatalf("failed to read request input: %s", err)
	}

	err = request.Source.Validate()
	if err != nil {
		log.Fatalf("invalid source config: %s", err)
	}

	response, err := resource.Check(request)
	if err != nil {
		log.Fatalf("failed to perform check: %s", err)
	}

	err = response.Write()
	if err != nil {
		log.Fatalf("failed to write response to stdout: %s", err)
	}

	log.Println("check complete")
}
