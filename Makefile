GO = GOEXPERIMENT=jsonv2,greenteagc go


.PHONY: tidy vendor


install-tailwind-cli:
	mkdir -p $(HOME)/.local/bin && \
	curl -L -o $(HOME)/.local/bin/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.18/tailwindcss-linux-x64 && \
	chmod +x $(HOME)/.local/bin/tailwindcss


install:install-tools install-tailwind-cli


install-tools:
	$(GO) install tool

tools:
	@$(GO) tool


update-tools:
	$(GO) get tool


tidy:
	$(GO) mod tidy


vendor: tidy
	$(GO) mod vendor