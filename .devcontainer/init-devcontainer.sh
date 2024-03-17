go install github.com/bufbuild/buf/cmd/buf@latest && \
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

cd /workspaces/web && \
    sudo chown -R vscode:vscode node_modules && \
    npm ci
