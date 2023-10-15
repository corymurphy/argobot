# ArgoCD GitHub Bot

GitHub bot for ArgoCD. Inspired by Atlantis for Terraform.

## Deploy

```shell
helm repo add corymurphy https://corymurphy.github.io/argobot
```

## Develop

Built using [palantir/go-githubapp](https://github.com/palantir/go-githubapp)

Install minikube. [docs](https://minikube.sigs.k8s.io/docs/start/).

```shell
# enable ingress
minikube addons enable ingress
kubectl get pods -n ingress-nginx

# test
kubectl create deployment web --image=gcr.io/google-samples/hello-app:1.0
kubectl expose deployment web --type=NodePort --port=8080
minikube service web --url
```

### Deploy argocd

```shell
# deploy argocd
helm dependency update --skip-refresh charts/argocd/
helm upgrade argocd --install --values charts/argocd/values.yaml --namespace argocd --create-namespace charts/argocd/

# connect to argocd
kubectl port-forward --namespace argocd service/argocd-server 8000:80
```

### Deploy argobot

```shell
./bin/minikube_build_deploy
./bin/expose
```

## Planning

- [x] pull a git repository and store in a org/repo/pr_num directory
- [x] receive argo diff with app name and comment on a pull request
- [x] configure minikube to use local docker images
- [ ] create healthcheck endpoint
- [ ] complete helm chart
- [ ] push helm chart to public repo
- [ ] create helm chart repo docs
- [ ] support pvc and statefulset
- [ ] create dockerfile
- [ ] publish dockerfile to registry
- [ ] ensure that the argo cli works without kube access
- [ ] support authentication with the argo cli
- [ ] determine argo project from a directory
- [ ] run argo diff on a directory against a specified project

## Notes

