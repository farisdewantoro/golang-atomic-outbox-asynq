.PHONY: swagger
swagger:
	@$(MAKE) swagger-concat
	@$(MAKE) swagger-server

swagger-concat:
	@echo "make swagger-concat"
	@rm -rf docs/api/tmp
	@mkdir -p docs/api/tmp
	./swagger mixin --output=docs/api/tmp/tmp.yml --format=yaml --keep-spec-order \
		docs/api/configs/main.yml docs/api/paths/*
	./swagger flatten docs/api/tmp/tmp.yml --output=docs/api/swagger.yaml --format=yaml
	./swagger flatten docs/api/tmp/tmp.yml --output=docs/api/swagger.json --format=json
	#@sed -i '1s@^@# Code generated by "make swagger"; DO NOT EDIT.\n@' docs/api/swagger.yaml
	@rm -f docs/api/tmp/tmp.yml

swagger-server:
	@echo "make swagger-server"
	@mkdir ./internal/generated || true
	@rm -rf ./internal/generated/api_models
	./swagger generate model \
		--allow-template-override \
		--spec=docs/api/swagger.yaml \
		--target=internal/generated \
		--model-package=api_models

swagger-install:
	if [ ! -f "./swagger" ]; then \
        wget "https://github.com/go-swagger/go-swagger/releases/download/v0.30.3/swagger_`go  env GOOS`_`go env GOARCH`" -O ./swagger && \
        chmod +x ./swagger; \
    fi


new-migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir docs/db/migrations -seq $$migration_name

db-migrate-up:
	@go run main.go db-migrate-up