# Github Rate Limit prometheus exporter helm chart

This helm chart helps to install and configure [github-rate-limits-prometheus-exporter](../../README.md)

The helm chart itself is a simplified version of a generated helm chart for 'any' service. Values which can be configured can be viewed [here](values.yaml) 

The application specific configuration
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
