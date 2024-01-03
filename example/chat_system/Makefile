.PHONY: run
PID_FILE := program_pid.txt


gen-protos:
	protoc --go_out=. --go-grpc_out=. protos/*.proto
buf-gen:
	buf

#TODO: setup for client-server comms
start-server:
	@cd server; nohup go run main.go > /dev/null 2>&1 & echo $$! > ../$(PID_FILE)


start-client:
	@cd client; go run .

cleanup:
	@echo "Stopping the background process"
	@kill -TERM `cat $(PID_FILE)` && rm -f $(PID_FILE) || true;
	@rm -f program_pid.txt

