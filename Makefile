GO = GOEXPERIMENT=jsonv2,greenteagc go


tidy:
	$(GO) mod tidy


vendor:
	$(GO) mod tidy
	$(GO) mod vendor