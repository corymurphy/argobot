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
# start minikube
eval $(minikube -p minikube docker-env)
./bin/minikube_start

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
helm dependency update charts/argocd/
helm dependency update --skip-refresh charts/argocd/
helm upgrade argocd --kube-context kind-kind --install --values charts/argocd/values.yaml --namespace argocd --create-namespace charts/argocd/

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
- [x] create healthcheck endpoint
- [x] complete helm chart
- [x] push helm chart to public repo
- [ ] create helm chart repo docs
- [ ] support pvc and statefulset
- [x] create dockerfile
- [x] publish dockerfile to registry
- [x] ensure that the argo cli works without kube access
- [ ] support authentication with the argo cli
- [ ] determine argo project from a directory
- [ ] run argo diff on a directory against a specified project

## Notes

- mocking the github sdk https://github.com/migueleliasweb/go-github-mock
- testing https://github.com/google/go-github/blob/master/test/README.md
- unit testing http https://stackoverflow.com/questions/51120033/how-to-test-http-calls-in-go
- http handler response recorder https://ieftimov.com/posts/testing-in-go-testing-http-servers/
- testing bot - https://github.com/apps/cam-argobot
- atlantis uses `g := events.DefaultGithubRequestValidator{}` to validate github webhooks
- add step to download argo cli for testing or docker cli step