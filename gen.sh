set +xe

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go get -u github.com/gin-gonic/gin
export PATH="$PATH:$(go env GOPATH)/bin"
echo $PATH
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  protos/main.proto

python -m grpc_tools.protoc -I./protos \
  --python_out=./iComm --grpc_python_out=./iComm \
  ./protos/main.proto
