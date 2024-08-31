DOCKER_IMAGES := $(shell find misc/docker/* -maxdepth 1 -type d -printf "%f\n")

build:
	CGO_ENABLED=0 go build -o ./bin/libra ./cmd/libra

watch:
	go run github.com/cortesi/modd/cmd/modd@latest

docker-images: $(foreach tag,$(DOCKER_IMAGES), docker-image-$(tag))

docker-image-%:
	docker build -t libra:$*-latest -f misc/docker/$*/Dockerfile .