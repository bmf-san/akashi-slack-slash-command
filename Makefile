.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run test
	$(eval TESTENV := $(shell cat .env.test))
	$(TESTENV) \
	go test -race -v -cover ./...

.PHONY: deploy
deploy: ## deploy
	rm -rf vendor
	GOPROXY=direct GOSUMDB=off go mod tidy
	go mod vendor
	gcloud functions deploy go-slack-app-on-gae-boilerplate \
		--entry-point Slash \
		--runtime go116 \
		--trigger-http \
		--env-vars-file ./.env.yaml  \
		--region asia-northeast1 \
		--timeout 540s \
		--source .
	gcloud functions describe go-slack-app-on-gae-boilerplate --region asia-northeast1