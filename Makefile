run_tests:
	go test -v ./tests

linter:
	golangci-lint run

push_tag:
	git tag v1.0.0
	git push --tags
