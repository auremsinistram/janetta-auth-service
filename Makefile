run:
	@DEBUG=1 go run ./cmd/main/main.go -configPath ./config.json

errgen:
	@errgen -i ./errors.json -c ./pkg/econst/econst.go -r ./error.response.json -e "internal server error"

update-deps:
	@go get -u ./... && go mod tidy

update-errgen:
	@go install github.com/auremsinistram/go-errgen/cmd/errgen@latest
