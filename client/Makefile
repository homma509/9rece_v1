DOCKER_YAML=-f ./deployments/docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=9rece-client docker-compose $(DOCKER_YAML)

docker-build:
	$(DOCKER) build ${ARGS}
docker-up:
	$(DOCKER) up

serve:
	$(DOCKER) run --service-ports node npm run serve
lint:
	$(DOCKER) run node npm run lint
build:
	$(DOCKER) run node npm run build
clean:
	./scripts/clean.sh

deploy-client:
	$(DOCKER) run sls ./scripts/deploy.sh client ${STAGE}
deploy-dist:
	$(DOCKER) run sls ./scripts/deploy.sh dist ${STAGE}

remove-client:
	$(DOCKER) run sls ./scripts/remove.sh client ${STAGE}
remove-dist:
	$(DOCKER) run sls ./scripts/remove.sh dist ${STAGE}
