run_tests:
	go test -v ./tests

run_tests_with_coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

linter:
	golangci-lint run

push_tag:
	git tag v1.0.0
	git push --tags
