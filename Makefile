compose-up-local:
	docker-compose -f deploy/local/docker-compose.yaml up -d

compose-down-local:
	docker-compose -f deploy/local/docker-compose.yaml down

compose-up-test:
	docker-compose -f deploy/test/docker-compose.yaml up -d

compose-down-test:
	docker-compose -f deploy/test/docker-compose.yaml down

gen-mocks:
	touch internal/mocks/a.txt && rm internal/mocks/* && mockery --dir internal --output internal/mocks --all