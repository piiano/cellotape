
COVERAGE_PROFILE_FILE := profile.cov
COVERAGE_HTML_FILE    := coverage.html
SHELL := bash -euo pipefail -c


$(COVERAGE_PROFILE_FILE): $(shell find router examples)
	@go test -v -race -failfast \
		-coverprofile=$@ \
		-coverpkg=$(shell go list ./router/... | paste -s -d , -) \
		./... | grep -v "coverage:" 1>&2 || rm $@

.PHONY: test-coverage
test-coverage: $(COVERAGE_PROFILE_FILE)
	@go tool cover -func $(COVERAGE_PROFILE_FILE) | grep total | awk '{print $$3}'

$(COVERAGE_HTML_FILE): $(COVERAGE_PROFILE_FILE)
	@go tool cover -html=$(COVERAGE_PROFILE_FILE) -o $(COVERAGE_HTML_FILE)

.PHONY: show-coverage
show-coverage: $(COVERAGE_HTML_FILE)
	@open $(COVERAGE_HTML_FILE)

.PHONY: test
test:
	@go test -race -failfast ./...

.PHONY: clean-coverage
clean-coverage:
	@rm -f $(COVERAGE_PROFILE_FILE) $(COVERAGE_HTML_FILE)