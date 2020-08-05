package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	m "github.com/digitalocean/concourse-resource-library/metadata"
)

// AQL provides the version query structure
type AQL struct {
	Raw   string `json:"raw,omitempty"`   // Raw AQL to filter versions on
	Repo  string `json:"repo,omitempty"`  // Repo within Artifactory to search
	Image string `json:"image,omitempty"` // Image to match
	Tag   string `json:"tag,omitempty"`   // Tag of image to match, defaults to `latest`
}

// UnmarshalJSON custom unmarshaller to convert PR number
func (a *AQL) UnmarshalJSON(data []byte) error {
	type Alias AQL
	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.Raw == "" && aux.Repo != "" {
		if aux.Image == "" {
			aux.Image = "*"
		}

		if aux.Tag == "" {
			aux.Tag = "latest"
		}

		aux.Raw = fmt.Sprintf(
			`{"repo": "%s", "path": {"$match": "%s/%s"}, "name": "manifest.json"}`,
			aux.Repo,
			aux.Image,
			aux.Tag,
		)
	}

	return nil
}

// SetModifiedTime appends the version modified time to the raw AQL query
func (a *AQL) SetModifiedTime(v Version) {
	if a.Raw == "" {
		return
	}

	mod := time.Now().AddDate(-2, 0, 0)

	if !v.Modified.IsZero() {
		mod = v.Modified
	}

	a.Raw = fmt.Sprintf(`%s, "modified": {"$gt": "%s"}}`, a.Raw[:len(a.Raw)-1], mod.Format(time.RFC3339Nano))
}

// Source represents the configuration for the resource
type Source struct {
	Endpoint    string `json:"endpoint"`           // Endpoint for Artifactory AQL API (leave blank for cloud)
	User        string `json:"user,omitempty"`     // User for Artifactory API with permissions to Repository
	Password    string `json:"password,omitempty"` // Password for Artifactory API with permissions to Repository
	AccessToken string `json:"access_token"`       // AccessToken for Artifactory API with permissions to Repository
	APIKey      string `json:"api_key,omitempty"`  // APIKey for Artifactory API with permissions to Repository
	Host        string `json:"host"`               // Host is used to get / put the image via the Docker v2 registry api
	AQL         AQL    `json:"aql"`                // AQL details to search for
	Proxy       bool   `json:"proxy,omitempty"`    // Proxy if you are using a proxy, defaults to false & direct access, see: https://www.jfrog.com/confluence/display/JFROG/HTTP+Settings#HTTPSettings-DockerReverseProxySettings
}

// Validate ensures that the source configuration is valid
func (s *Source) Validate() error {
	switch {
	case s.Endpoint == "":
		return errors.New("endpoint is required")
	case s.User != "" && s.Password == "" && s.APIKey == "" && s.AccessToken == "":
		return errors.New("user cannot be defined without a Password || AccessToken || APIKey")
	case s.AQL.Raw == "" && s.AQL.Repo == "" && s.AQL.Image == "*":
		return errors.New("aql cannot be defined without a Owner || Name")
	}

	return nil
}

// Version contains the version data Concourse uses to determine if a build should run
type Version struct {
	Repo     string    `json:"repo,omitempty"`
	Image    string    `json:"image,omitempty"`
	Tag      string    `json:"tag,omitempty"`
	Digest   string    `json:"digest,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

// ImageTag returns the string needed to fetch the artifact
func (v *Version) ImageTag() string {
	return fmt.Sprintf("%s:%s", v.Image, v.Tag)
}

// ArtifactoryPath builds the internal path for an manifest based on the version
func (v *Version) ArtifactoryPath() string {
	return fmt.Sprintf("%s/%s/%s/manifest.json", v.Repo, v.Image, v.Tag)
}

// Empty returns true if the version is empty
func (v *Version) Empty() bool {
	if v.Repo == "" || v.Image == "" || v.Tag == "" {
		return true
	}

	return false
}

// CheckRequest is the data struct received from Concoruse by the resource check operation
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// Read will read the json response from Concourse via stdin
func (r *CheckRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}

// CheckResponse is the data struct returned to Concourse by the resource check operation
type CheckResponse []Version

// Len returns the number of versions in the response
func (r CheckResponse) Len() int {
	return len(r)
}

// Write will write the json response to stdout for Concourse to parse
func (r CheckResponse) Write() error {
	return json.NewEncoder(os.Stdout).Encode(r)
}

// GetParameters is the configuration for a resource step
type GetParameters struct {
	Format       string `json:"format"`        // Format defaults to `rootfs` (used for task steps), to write image to tarball, use `oci`
	SkipDownload bool   `json:"skip_download"` // SkipDownload is used with `put` steps to skip `get` step that Concourse runs by default
}

// GetRequest is the data struct received from Concoruse by the resource get operation
type GetRequest struct {
	Source  Source        `json:"source"`
	Version Version       `json:"version"`
	Params  GetParameters `json:"params"`
}

// Read will read the json response from Concourse via stdin
func (r *GetRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}

// OCIRepository uses Source & Version to generate container registry compliant url
func (r *GetRequest) OCIRepository() string {
	if r.Source.Proxy {
		return fmt.Sprintf("%s/%s", r.Source.Host, r.Version.Image)
	}

	return fmt.Sprintf("%s/%s/%s", r.Source.Host, r.Version.Repo, r.Version.Image)
}

// GetResponse ...
type GetResponse struct {
	Version  Version    `json:"version"`
	Metadata m.Metadata `json:"metadata,omitempty"`
}

// Write will write the json response to stdout for Concourse to parse
func (r GetResponse) Write() error {
	return json.NewEncoder(os.Stdout).Encode(r)
}

// PutParameters for the resource
type PutParameters struct {
	Pattern        string        `json:"pattern"`               // Pattern to search inputs for image tarbull to push in GLOB form, e.g. `input/image.tar`
	Image          string        `json:"image"`                 // Image to push, e.g. owner/image
	Target         string        `json:"target,omitempty"`      // Target Artifactory repository, required if not using a Proxy
	Tags           []string      `json:"tags,omitempty"`        // Tags is the list of tags to apply to the image
	Properties     string        `json:"properties,omitempty"`  // Properties is file path containing image properties in `key=value\n` form
	Params         string        `json:"params,omitempty"`      // Params is a file path containing the Image struct in json syntax for dynamic workflows
	BuildEnv       string        `json:"build_env,omitempty"`   // BuildEnv is path to file containing build environment values in `key=value\n` form, e.g. `env > env.txt`
	EnvInclude     string        `json:"env_include,omitempty"` // EnvInclude case insensitive patterns in the form of "value1;value2;..." will be included
	EnvExclude     string        `json:"env_exclude,omitempty"` // EnvExclude case insensitive patterns in the form of "value1;value2;..." will be excluded, defaults to `*password*;*psw*;*secret*;*key*;*token*`
	RepositoryPath string        `json:"repo_path,omitempty"`   // RepositoryPath sets the path to the input containing the repository (git support only)
	Repository     string        `json:"repo,omitempty"`        // Repository set the Git repository url explicitly for compatibility with the Git resource
	Get            GetParameters `json:"get,omitempty"`         // Get parameters for explicit get step after put
}

// Parse reads a file and unmarshals into itself
func (p *PutParameters) Parse() error {
	data, err := ioutil.ReadFile(p.Params)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, p)
}

// PutRequest is the data struct received from Concoruse by the resource put operation
type PutRequest struct {
	Source Source        `json:"source"`
	Params PutParameters `json:"params"`
}

// OCIRepository uses Source & Params to generate container registry compliant url
func (r *PutRequest) OCIRepository() string {
	if r.Source.Proxy {
		return fmt.Sprintf("%s/%s", r.Source.Host, r.Params.Image)
	}

	return fmt.Sprintf("%s/%s/%s", r.Source.Host, r.Params.Target, r.Params.Image)
}

// Read will read the json response from Concourse via stdin
func (r *PutRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}
