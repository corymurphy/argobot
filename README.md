# ArgoCD GitHub Bot

GitHub bot for ArgoCD. Inspired by Atlantis for Terraform.

## Deploy

```shell
kubectl create secret generic argobot \
  --namespace argobot \
  --content <kube-context> \
  --from-file=key.pem=path/to/github_app_privatekey \
  --from-literal=webhookSecret=<github_webhook_secret> \
  --from-literal=argocdApiKey=<argocd_api_key>

helm install \
  oci://ghcr.io/corymurphy/helm-charts/argobot argobot
```

## Develop

```shell
make test
make test-helm
make dev-install
```

### Deploy argocd

```shell
# deploy argocd
helm dependency update charts/argocd/
helm upgrade argocd --kube-context kind-kind \
  --install --namespace argocd --create-namespace \
  --values charts/argocd/values.yaml charts/argocd/
```
