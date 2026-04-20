GO = GOEXPERIMENT=jsonv2,greenteagc go


.PHONY: tidy vendor install install-gotests tools update-tools test generate generate-ui generate-swagger fix run build


install-tailwind-cli:
	mkdir -p $(HOME)/.local/bin && \
	curl -L -o $(HOME)/.local/bin/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/v4.2.1/tailwindcss-linux-x64 && \
	chmod +x $(HOME)/.local/bin/tailwindcss


install-ui-tools: install-tailwind-cli


install-go-tools:
	$(GO) install tool


install: tidy install-gotests install-go-tools install-ui-tools
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# If you want to use vscode to generate tests, you will need to have gotests installed on the machine
# because you can not get vscode to use `go tool gotests ...`
install-gotests:
	go install github.com/cweill/gotests/gotests@latest


tools:
	@$(GO) tool


update-tools:
	$(GO) get tool


tidy:
	$(GO) mod tidy


vendor: tidy
	$(GO) mod vendor


test:
	CGO_ENABLED=1 $(GO) test -race -failfast -v ./...


generate-swagger:
	$(GO) tool swag init --dir ./cmd,./models --output ./docs


generate-ui:
	$(GO) tool templ generate


migrations:
	~/go/bin/migrate -path ./internal/store/migrations -database "postgres://local:dev@localhost:5432/edgedb?sslmode=disable" up


generate: generate-ui generate-swagger


fix:
	@echo "running go fix a few times"
	@$(GO) fix ./... && $(GO) fix ./... && $(GO) fix ./...


run: generate fix
	$(GO) run ./cmd


run-reload: generate fix
	$(GO) tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run ./cmd"


build: generate fix
	$(GO) build -o ./dist/edgedb-server ./cmd
