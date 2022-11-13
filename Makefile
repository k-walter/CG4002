SHELL := /bin/bash
default: grpc ecomm

ecomm:
	# build and sync
	cd eComm; \
		env GOOS=linux GOARCH=arm64 go build -o ecomm; \
		rsync -rauzvP ecomm xilinx:~/cg4002

grpc:
	# Generate go
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go get -u github.com/gin-gonic/gin
	export PATH="$(PATH):$(go env GOPATH)/bin"
	echo $PATH
	protoc --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  protos/main.proto
	
	# Generate python
	pip3 install grpcio
	python -m grpc_tools.protoc -I./protos \
	  --python_out=./iComm --grpc_python_out=./iComm \
	  protos/main.proto
	python -m grpc_tools.protoc -I./protos \
	  --python_out=./eComm/pynq --grpc_python_out=./eComm/pynq \
	  protos/main.proto
