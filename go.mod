module github.com/digitalocean/artifactory-docker-resource

go 1.14

require (
	github.com/digitalocean/concourse-resource-library v0.0.0-20200805204403-3eb5253b5085
	github.com/jfrog/jfrog-client-go v0.12.0
)

replace github.com/docker/docker v1.13.1 => github.com/docker/engine v0.0.0-20200720230453-22153d111ead
