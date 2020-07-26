DOCKER_YAML=-f ./deployments/docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=9rece-frontend docker-compose $(DOCKER_YAML)

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

deploy-front:
	$(DOCKER) run sls ./scripts/deploy.sh front ${STAGE}
deploy-dist:
	$(DOCKER) run sls ./scripts/deploy.sh dist ${STAGE}

remove-front:
	$(DOCKER) run sls ./scripts/remove.sh front ${STAGE}
remove-dist:
	$(DOCKER) run sls ./scripts/remove.sh dist ${STAGE}
