v := false
html := false
go_test_flags := -tags no_test_coverage
go_test_coverage_flags := -func=coverage/coverage.out

# Add verbose flag if v=true
ifeq ($(v),true)
    go_test_flags += -v
endif

# Add 
ifeq ($(html),true)
    go_test_coverage_flags = -html=coverage/coverage.out
endif

run: compile-protos
	@go run src/main.go

build: compile-protos
	@go build -o bin/mistapi src/main.go

live-run: compile-protos
	@air --build.cmd "swag init --dir src && go build -o bin/mistapi src/main.go" \
	     --build.bin "./bin/mistapi" \
	     --build.full_bin=false \
	     --build.exclude_dir "docs,bin" \
	     --build.include_ext "go,tpl,tmpl,html"	# @air

compile-protos cp:
	@buf generate

# ----- TESTS -----
run-tests t: test-auth

test-auth:
	@echo -----------------------------------------
	@go test mistapi/src/auth -coverprofile=coverage/coverage.out  $(go_test_flags)
	@go tool cover $(go_test_coverage_flags)

# ----- FORMAT -----
lint:
	golangci-lint run --disable-all -E errcheck

lint-proto:
	@buf lint
