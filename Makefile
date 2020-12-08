DOCKER := docker
IMAGE := bioapi
TAG := latest
PORT := 8090

.PHONY: build
build:
	$(DOCKER) build -t $(IMAGE):$(TAG) .

.PHONY: start
start:
	$(DOCKER) run -d -p $(PORT):$(PORT) $(IMAGE):$(TAG)

.PHONY: stop
stop:
	$(DOCKER) rm -f $$(docker ps --format "{{.ID}}\t{{.Image}}" | grep $(IMAGE):$(TAG) | awk '{print $$1}')

.PHONY: restart
restart: stop start

.PHONY: rollout
rollout: build stop start
