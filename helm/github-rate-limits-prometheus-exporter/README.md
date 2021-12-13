# Github Rate Limit prometheus exporter helm chart

This helm chart helps to install and configure [github-rate-limits-prometheus-exporter](../../README.md)

The helm chart itself is a simplified version of a generated helm chart for 'any' service. Values which can be configured can be viewed [here](values.yaml) 

To add a repository to your local helm repos
```sh
helm repo add grl-exporter https://kalgurn.github.io/github-rate-limits-prometheus-exporter-charts/
```

To install the chart
```sh
helm upgrade --install \
    release_name grl-exporter/github-rate-limits-prometheus-exporter \
    -f path_to_values/with_github_configuration.yaml
```

## Application specific configuration
GitHub PAT

```yaml
github:
  authType: pat 
  secretName: secret # Name of a secret which stores PAT
```

GitHub App
```yaml
github:
  authType: app
  appID: 1 # GitHub applicaiton ID
  installationID: 1 # GitHub App installation ID
  privateKeyPath: "/tmp" # path to which the private key will be mounted
  secretName: secret # name of a secret which stores key.pem
```

Example values file
```yaml
github:
  authType: pat
  secretName: gh-token

replicaCount: 1

image:
  repository: ghcr.io/kalgurn/grl-exporter
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: "v0.1.4"

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

## Prerequisites 

For the application to run the Kubernetes secrets should be installed in the same namespace. E.g., for GitHub App you can create a secret from the key.pem with the command below. 

```sh
kubectl create secret generic github-key --from-file=key.pem
```

It will create a secret with a name __github-key__ and a private key stored within the keys `data["key.pem"]`

For the PAT the easiest way would be

```sh
echo -n 'ghb_token' > ./token
kubectl create secret generic gh-token --from-file=token
```

The command above will create a secret __gh-token__ in the current namespace and a token stored within keys `data["token"]`
