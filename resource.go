package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	m "github.com/digitalocean/concourse-resource-library/metadata"
)

// Image provides the version query structure
type Image struct {
	Repo  string `json:"repo,omitempty"`  // Artifactory repository to search
	Owner string `json:"owner,omitempty"` // Artifactory image owner to match
	Name  string `json:"name,omitempty"`  // Artifactory image name to match
	Tag   string `json:"tag,omitempty"`   // Artifactory image tag to match, defaults to `latest`
	Raw   string `json:"-"`
}

// UnmarshalJSON custom unmarshaller to convert PR number
func (a *Image) UnmarshalJSON(data []byte) error {
	type Alias Image
	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.Owner == "" {
		aux.Owner = "*"
	}

	if aux.Name == "" {
		aux.Name = "*"
	}

	if aux.Tag == "" {
		aux.Tag = "latest"
	}

	if aux.Raw == "" && aux.Repo != "" {
		aux.Raw = fmt.Sprintf(
			`{"repo": "%s", "path": {"$match": "%s/%s/%s"}, "name": "manifest.json"}`,
			aux.Repo,
			aux.Owner,
			aux.Name,
			aux.Tag,
		)
	}

	return nil
}

// SetModifiedTime appends the version modified time to the raw AQL query
func (a *Image) SetModifiedTime(v Version) {
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
	Image       Image  `json:"image,omitempty`     // Image details to search for
}

// Validate ensures that the source configuration is valid
func (s *Source) Validate() error {
	switch {
	case s.Endpoint == "":
		return errors.New("endpoint is required")
	case s.User != "" && s.Password == "" && s.APIKey == "" && s.AccessToken == "":
		return errors.New("user cannot be defined without a Password || AccessToken || APIKey")
	case s.Image.Raw == "" && s.Image.Repo == "" && s.Image.Owner == "*" || s.Image.Name == "*":
		return errors.New("image cannot be defined without a Owner || Name")
	}

	return nil
}

// Version contains the version data Concourse uses to determine if a build should run
type Version struct {
	Repo     string    `json:"repo,omitempty"`
	Owner    string    `json:"owner,omitempty"`
	Name     string    `json:"name,omitempty"`
	Tag      string    `json:"tag,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	// TODO: include sha???
}

// Image returns the string needed to fetch the artifact
func (v *Version) Image() string {
	return fmt.Sprintf("%s/%s:%s", v.Owner, v.Name, v.Tag)
}

// Empty returns true if the version is empty
func (v *Version) Empty() bool {
	if v.Repo == "" || v.Owner == "" || v.Name == "" || v.Tag == "" {
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
}

// PutRequest is the data struct received from Concoruse by the resource put operation
type PutRequest struct {
	Source Source        `json:"source"`
	Params PutParameters `json:"params"`
}

// Read will read the json response from Concourse via stdin
func (r *PutRequest) Read(input []byte) error {
	return json.Unmarshal(input, r)
}
