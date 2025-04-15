run-unit-tests:
	go mod tidy && \
	go test -v -race -timeout 30s ./...
run-integration-tests:
	cd tests/integration && \
	go mod tidy && \
	go test -tags integration -v -race -timeout 30s ./... -run TestMain