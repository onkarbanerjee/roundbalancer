
start-application-server:
	@echo "Running the application on port $(PORT) with ID $(ID)..."
	go run cmd/echo/main.go -id=$(ID) -port=$(PORT)
	go run cmd/cli/main.go application-server start --port $(PORT) --id $(ID)

mocks:
	rm -rf mocks
	go run -mod=mod go.uber.org/mock/mockgen -destination=./mocks/mock_pool.go -package=mocks -typed github.com/onkarbanerjee/roundbalancer/pkg/backends GroupOfBackends
	go run -mod=mod go.uber.org/mock/mockgen -destination=./mocks/mock_loadbalancer.go -package=mocks -typed github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer LoadBalancer

test_and_coverage:
	go test -race -coverprofile=coverage.txt ./...

.PHONY: mocks