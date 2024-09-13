DOCKER_IMAGES := $(shell find misc/docker/* -maxdepth 1 -type d -not -name goreleaser -printf "%f\n")
REPOSITORY_OWNER ?= 

build:
	CGO_ENABLED=0 go build -o ./bin/libra ./cmd/libra

watch:
	go run github.com/cortesi/modd/cmd/modd@latest

docker-images: $(foreach tag,$(DOCKER_IMAGES), docker-image-$(tag))

docker-image-%:
	docker build -t libra:$*-latest -f misc/docker/$*/Dockerfile .

release:
	curl -sfL https://goreleaser.com/static/run | REPOSITORY_OWNER=$(REPOSITORY_OWNER) bash -s -- release --auto-snapshot --clean