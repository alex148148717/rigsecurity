FROM golang:1.22

RUN apt-get update && apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /workspace

COPY  . .

RUN ls -alh /workspace

RUN protoc --proto_path=/workspace/internal/github_cel/interfaces --go_out=. --go-grpc_out=. /workspace/internal/github_cel/interfaces/access_information.proto


CMD ["bash"]
