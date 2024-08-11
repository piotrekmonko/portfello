VERSION?=$(shell git rev-list --abbrev-commit --abbrev=4 -1 HEAD)
DB_PASSWORD?="$(shell kubectl get secret --namespace portfello db-postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)"

help:
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "  %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "%s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## Build:
build: gen ## Build local executable
	CGO_ENABLED=0 go build -ldflags='-X github.com/piotrekmonko/portfello/cmd.buildNumber=$(VERSION)' main.go

test: ## Regenerate mocks and run tests
	mockery
	go test -race -cover ./...

gen generate: ## Generate mocks and schemas
	go generate ./...
	mockery

## Debug:
mock-data: ## Populate database with example expenses
	go run main.go prov -t -n 23

dlv-serve: ## Debug inside a remote container
	go build -gcflags "all=-N -l" main.go
    dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./main serve

db-get-pass: ## Run kubectl to obtain generated postgresql password in k8
	echo $(shell kubectl get secret --namespace portfello db-postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)

db-shell: ## Run kubectl to obtain a shell into postgresql in k8
	kubectl run db-postgresql-client --rm --tty -i --restart='Never' --namespace portfello \
		--image docker.io/bitnami/postgresql:16.3.0-debian-12-r19 --env="PGPASSWORD=${DB_PASSWORD}" \
		--command -- psql --host db-postgresql -U postgres -d postgres -p 5432
