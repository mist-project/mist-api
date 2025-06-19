## Install swaggo tool to be able to run swag init

go install github.com/swaggo/swag/cmd/swag@latest


### Install protobuf compiler

```shell

brew install bufbuild/buf/buf

# On linux install, install version 3.12.4
apt install -y protobuf-compiler # Idk if you need this anymore

# Install go plugin for the protocol compiler, version 1.35.2
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install plugin for the protocol compiler, version 1.5.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# update your PATH so that the protoc compiler can find the plugin
export PATH="$PATH:$(go env GOPATH)/bin"

# install protoc validate
buf dep update
```

### Install live reloader
`go install github.com/air-verse/air@1.61.1`