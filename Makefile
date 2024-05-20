.PHONY: release build

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

docs:
	(sleep 0.5 && open http://localhost:6060/pkg/github.com/codecrafters-io/shell-tester/internal/)
	godoc -http=:6060

release:
	git tag v$(next_version_number)
	git push origin main v$(next_version_number)

build:
	go build -o dist/main.out ./cmd/tester

test:
	go test -count=1 -p 1 -v ./internal/...

test_bash: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/bash \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}\
	]" \
	dist/main.out

test_dash: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/dash \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}\
	]" \
	dist/main.out

test_ryan: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/ryan_shell \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"} \
	]" \
	dist/main.out

test_all: test_bash test_ryan test_dash

record_fixtures:
	CODECRAFTERS_RECORD_FIXTURES=true make test

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils
