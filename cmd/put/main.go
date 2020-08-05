package main

import (
	"log"
	"os"

	resource "github.com/digitalocean/artifactory-docker-resource"
	rlog "github.com/digitalocean/concourse-resource-library/log"
)

func main() {
	input := rlog.WriteStdin()
	defer rlog.Close()

	var request resource.PutRequest
	err := request.Read(input)
	if err != nil {
		log.Fatalf("failed to read request input: %s", err)
	}

	err = request.Source.Validate()
	if err != nil {
		log.Fatalf("invalid source config: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}
	dir := os.Args[1]

	response, err := resource.Put(request, dir)
	if err != nil {
		log.Fatalf("failed to perform check: %s", err)
	}

	err = response.Write()
	if err != nil {
		log.Fatalf("failed to write response to stdout: %s", err)
	}

	log.Println("put complete")
}
