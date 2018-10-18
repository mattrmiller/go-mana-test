## Lint: Validate golang code
lint:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
	gometalinter.v2 \
		--vendor \
		--exclude=^vendor\/ \
		--aggregate \
		--deadline=900s \
		--line-length=140 \
		--cyclo-over=15 \
		--enable-all \
		--disable=dupl \
		--disable=goimports \
		--disable=maligned \
		--disable=gocyclo \
		--disable=gochecknoinits \
		--disable=gochecknoglobals \
		./...

## Perform all tests
test:
	go test ./...
