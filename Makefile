# HELP =================================================================================================================
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run-examples: ### run examples 
	@cd _examples && go run .
.PHONY: run 

test: ### run test
	go test -cover -race .
.PHONY: test


compose-up: ### Run docker-compose
	docker-compose up  && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down
