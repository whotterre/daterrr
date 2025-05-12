build:
	@echo "Building executable binary..."
	cd cmd && go build .
run:
	@echo "Starting server...."
	cd cmd && go run .
sqlc_gen:
	@echo "Generating Go code from SQL queries..."
	sqlc generate
PHONY: build run sqlc_gen