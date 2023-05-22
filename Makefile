
COVERAGE_PROFILE_FILE := profile.cov
COVERAGE_HTML_FILE    := coverage.html
SHELL := bash -eo pipefail -c
# router/ginbinders is currently excluded from coverage because it's copied from gin-gonic/gin.
EXCLUDE_FROM_COVERAGE := router/ginbinders

$(COVERAGE_PROFILE_FILE): $(shell find router examples)
	@go test -v -race -failfast \
		-coverprofile=$@ \
		-coverpkg=$(shell go list ./router/... | grep -v $(EXCLUDE_FROM_COVERAGE) | paste -s -d , -) \
		./... | grep -v "coverage:" 1>&2 || rm $@

.PHONY: test-coverage
test-coverage: $(COVERAGE_PROFILE_FILE)
	@go tool cover -func $(COVERAGE_PROFILE_FILE) | grep total | awk '{print $$3}'

$(COVERAGE_HTML_FILE): $(COVERAGE_PROFILE_FILE)
	@go tool cover -html=$(COVERAGE_PROFILE_FILE) -o $(COVERAGE_HTML_FILE)

FILE ?= file0

.PHONY: show-coverage
show-coverage: $(COVERAGE_HTML_FILE)
	$(eval ANCHOR:=$(shell cat $(COVERAGE_HTML_FILE) | grep -E '<option value=".*">.*</option>' | grep $(FILE) | sed -n 's/^.*value="\(.*\)".*$$/\1/p'))
	@sed -E -i '' 's/select\("file[0-9]+"\);/select("$(ANCHOR)");/g' $(COVERAGE_HTML_FILE)
	@open $(COVERAGE_HTML_FILE)


.PHONY: test
test:
	@go test -race ./...

.PHONY: clean-coverage
clean-coverage:
	@rm -f $(COVERAGE_PROFILE_FILE) $(COVERAGE_HTML_FILE)