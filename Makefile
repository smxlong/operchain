.PHONY: check
check: generate tidy fmt vet lint test

.PHONY: ci-check
ci-check: check
	git diff --exit-code || (echo "Uncommitted changes found. Run 'make check' and commit the changes." && exit 1)

.PHONY: generate
generate:
	go generate ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -coverprofile=coverage.txt -covermode=atomic -v ./...

.PHONY: dev-setup
dev-setup:
	cp dev/hooks/pre-commit .git/hooks/pre-commit
