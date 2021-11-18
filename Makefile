PKG = github.com/feitian124/goapi
COMMIT = $$(git describe --tags --always)

# get date
OS_NAME=${shell uname -s}
ifeq ($(OS_NAME),Darwin)
	SED = gsed
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	SED = sed
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

PACKAGES=`go list ./... | grep -v /vendor/`
VET_PACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GO_FILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: fmt lint test

fmt:
	gofumpt -l -w .

lint:
	golangci-lint run --timeout "5m"

test:
	@go test ./... --cover -count=1

generate:
	@go generate ./...

build:
	@go build -ldflags="$(BUILD_LDFLAGS)"

install:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
	go install github.com/99designs/gqlgen@v0.14.0
	# go install github.com/xo/usql@v0.9.4

mysql:
	usql my://root:mypass@localhost:33306/testdb -f testdata/ddl/mysql56.sql
	usql my://root:mypass@localhost:33308/testdb -f testdata/ddl/mysql.sql
	usql my://root:mypass@localhost:33308/testdb -c "CREATE DATABASE IF NOT EXISTS relations;"
	usql my://root:mypass@localhost:33308/relations -f testdata/ddl/detect_relations.sql
	usql my://root:mypass@localhost:33308/testdb -c "CREATE DATABASE IF NOT EXISTS relations_singular;"
	usql my://root:mypass@localhost:33308/relations_singular -f testdata/ddl/detect_relations_singular.sql
	usql maria://root:mypass@localhost:33309/testdb -f testdata/ddl/maria.sql

postgres:
	usql pg://postgres:pgpass@localhost:55432/testdb?sslmode=disable -f testdata/ddl/postgres95.sql
	usql pg://postgres:pgpass@localhost:55413/testdb?sslmode=disable -f testdata/ddl/postgres.sql

sqlite:
	sqlite3 $(PWD)/testdata/testdb.sqlite3 < testdata/ddl/sqlite.sql

info:
	@echo "PKG：${PKG}"
	@echo "COMMIT：${COMMIT}"
	@echo "OS_NAME: ${OS_NAME}"
	@echo "DATE: ${DATE}"
	@echo -e "\nPACKAGES:"
	@echo ${PACKAGES}
	@echo -e "\nVET_PACKAGES:"
	@echo ${VET_PACKAGES}
	@echo -e "\nGO_FILES:"
	@echo ${GO_FILES}
	@echo -e "\ngolangci-lint:"

.PHONY: default fmt lint test build install mysql postgres sqlite info