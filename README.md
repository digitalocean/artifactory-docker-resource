# Artifactory Docker Resource

Concourse resource for triggering, getting and putting new versions of docker / container image artifacts within Artifactory repositories.

## Config

Complete source configuration details can be found in the `resource.go` file for the `Source` struct.

## Check

Checks use the `items` domain to `find` artifacts with the supplied raw [AQL](https://www.jfrog.com/confluence/display/JFROG/Artifactory+Query+Language) or repo, path & name combination. Each artifact found is
returned as its own unique version for Concourse with the `Repo`, `Path`, `Name` & `Modified` values from the Artifactory API. `Modified` is used to filter future checks to ensure that API queries stay
performant.

## Get

Get will download an artifact to the input directory defined along with metadata for the artifact. The artifact is downloaded following its internal Artifactory path, so the `resource/local-path` metadata file is useful
to provide the specific path within the input directory to the downloaded artifact.

## Put

Put supports publishing 1 or more artifacts using glob style patterns to locate artifacts to publish.

## Examples

Configure the resource type:

```yaml
resource_types:
- name: artifactory
  type: docker-image
  source:
    repository: digitalocean/artifactory-docker-resource
    tag: latest
```

Source configuration using raw AQL for `item.find`:

```yaml
resources:
- name: myapplication
  type: artifactory
  icon: application-export
  source:
    endpoint: https://example.com/artifactory/
    user: ci
    password: ((artifactory.password))
    aql:
      raw: '{"repo": "docker-local", "path": {"$match" : "myapp/myimage/*"}, "name": "manifest.json"}'
    host: artifactory.example.com
```

Source configuration using repo, path, name:

```yaml
resources:
- name: myapplication
  type: artifactory
  icon: application-export
  source:
    endpoint: https://example.com/artifactory/
    user: ci
    password: ((artifactory.password))
    aql:
      repo: docker-local
      image: myapp/myimage
      tag: '*'
```

Publishing artifacts to Artifactory:

```yaml
- put: myapplication
  params:
    repo_path: code
    pattern: built/image.tar
    target: docker-local
```
