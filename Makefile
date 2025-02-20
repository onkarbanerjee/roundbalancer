
start-application-server:
	@echo "Running the application on port $(PORT) with ID $(ID)..."
	go run cmd/cli/main.go application-server start --port $(PORT) --id $(ID)

mocks:
	go run go.uber.org/mock/mockgen -destination=./mocks/mock_pool.go -package=mocks -typed github.com/onkarbanerjee/roundbalancer/pkg/backend Pool

.PHONY: mocks