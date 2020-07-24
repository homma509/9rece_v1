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

# clean:
# 	$(DOCKER) run node npm cache clean --force

clean:
	./scripts/clean.sh

deploy:
	# usage make deploy ARGS=dev
	$(DOCKER) run sls ./scripts/deploy.sh ${ARGS}

 remove:
	# usage make remove ARGS=dev
	$(DOCKER) run sls ./scripts/remove.sh ${ARGS}
