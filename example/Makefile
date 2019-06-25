# Require image path at Google Cloud Registry
IMAGE = 

all: build docker-build docker-push deploy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/app .

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker push $(IMAGE)

deploy:
	gcloud beta run deploy edge-exporter-example --image $(IMAGE) --region us-central1
