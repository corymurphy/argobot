TIMESTAMP=tmp-$(shell date +%s )

.PHONY: update
update:
	minikube -p minikube docker-env
	minikube image build --build-env=ARGOBOT_VERSION=0.0.1 --build-env=ARGOBOT_COMMIT=f943335d0c6df3f1aaadad50b57f1acf90145975 . -t ghcr.io/corymurphy/argobot:$(TIMESTAMP)
	helm upgrade --install --kube-context minikube --create-namespace --set image.pullPolicy=Never --set image.tag=$(TIMESTAMP) --namespace argobot argobot ./charts/argobot/;

expose:
	ngrok http 8080

build:

test:
	go test ./...

build-argocli:
	docker build -t ghcr.io/corymurphy/argocli --target argocli .