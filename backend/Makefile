DOCKER_YAML=-f ./deployments/docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=9rece-backend docker-compose $(DOCKER_YAML)

docker-build:
	$(DOCKER) build ${ARGS}

docker-up:
	$(DOCKER) up

lint:
	$(DOCKER) run go ./scripts/lint.sh

build:
	$(DOCKER) run go ./scripts/build.sh

clean:
	./scripts/clean.sh

deploy-api:
	$(DOCKER) run sls ./scripts/deploy.sh api ${STAGE}
deploy-func:
	$(DOCKER) run sls ./scripts/deploy.sh func ${STAGE}
deploy-infra:
	$(DOCKER) run sls ./scripts/deploy.sh infra ${STAGE}

remove-api:
	$(DOCKER) run sls ./scripts/remove.sh api ${STAGE}
remove-func:
	$(DOCKER) run sls ./scripts/remove.sh func ${STAGE}
remove-infra:
	$(DOCKER) run sls ./scripts/remove.sh infra ${STAGE}
