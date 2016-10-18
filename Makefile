SHELL := /bin/bash
REG := $(shell . .env && echo $$DOCKER_REGISTRY)

api-build:
	@pushd backend; ../dockerception.sh auzom-api api; docker tag \
	  auzom-api:latest $(REG)/api:latest; popd
api-push:
	@docker push $(REG)/api:latest
api: api-build api-push

worker-build:
	@pushd backend; ../dockerception.sh auzom-worker worker; docker tag \
	  auzom-worker:latest $(REG)/worker:latest; popd
worker-push:
	@docker push $(REG)/worker:latest
worker: worker-build worker-push

client-build:
	@pushd frontend; ../dockerception.sh auzom-client client; docker tag \
	  auzom-client:latest $(REG)/client:latest; popd
client-push:
	@docker push $(REG)/client:latest
client: client-build client-push

migrations-build:
	@pushd migrations; ../dockerception.sh auzom-migrations; popd
migrations: migrations-build
