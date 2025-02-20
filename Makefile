run-integration-tests:
	cd tests/integration && go test -tags integration -v -race -timeout 30s ./... -run TestMain