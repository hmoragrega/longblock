
.PHONY: all proto
proto:
	@cd proto && buf generate --template buf.gen.gogo.yml
	@mv github.com/hmoragrega/longblock/debug/v1/types/* debug/types && rm -rf github.com

.PHONY: test
test:
	go test -tags=debug -v -race ./...