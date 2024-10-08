run_tests:
	go test -v ./tests

push_tag:
	git tag v0.0.5
	git push --tags
