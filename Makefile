
start-application-server:
	@echo "Running the application on port $(PORT) with ID $(ID)..."
	go run cmd/cli/main.go application-server start --port $(PORT) --id $(ID)
