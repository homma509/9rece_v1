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

deploy:
	# usage make deploy TARGET=func STAGE=dev
	$(DOCKER) run sls ./scripts/deploy.sh ${TARGET} ${STAGE}

remove:
	# usage make remove TARGET=func STAGE=dev
	$(DOCKER) run sls ./scripts/remove.sh ${TARGET} ${STAGE}
