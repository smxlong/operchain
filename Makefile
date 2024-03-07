.PHONY: check
check: .git/hooks/pre-commit generate tidy fmt vet lint test

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

.git/hooks/pre-commit: .git/hooks dev/hooks/pre-commit
	cp dev/hooks/pre-commit .git/hooks/pre-commit

.git/hooks:
	mkdir -p .git/hooks
