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

default: test

fmt:
	gofumpt -l -w .

## lint: fmt then lint.
lint: fmt
	golangci-lint run --timeout "5m"

## test: test default database
test:
	export DATASOURCE=tidb_5_2 && go test ./... --cover -count=1

## testAll: test all supported databases one by one, currently mysql_8_0, mysql_5_7, mariadb_10_5, tidb_5_2.
testAll:
	export DATASOURCE=mysql_8_0 && go test ./... --cover -count=1
	export DATASOURCE=mysql_5_7 && go test ./... --cover -count=1
	export DATASOURCE=mariadb_10_5 && go test ./... --cover -count=1

## generate: generate graphql struct based on schema.
generate:
	@go generate ./...

build:
	@go build -ldflags="$(BUILD_LDFLAGS)"

## dev: start server with live reload. Runs `air` internally.
dev:
	@air

## install: install go modules and tools used by this project.
install:
	@echo " > install..."
	go mod tidy
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
	go install github.com/99designs/gqlgen@v0.14.0
	go install github.com/cosmtrek/air@v1.27.4
	go install -tags 'no_postgres no_oracle no_sqlserver no_sqlite3' github.com/xo/usql@v0.9.5

## tidb: init tidb. create database testdb and run ddl.
tidb:
	usql my://root@192.168.135.154:4000 -c "CREATE DATABASE IF NOT EXISTS testdb;"
	usql my://root@192.168.135.154:4000/testdb -f testdata/ddl/tidb_5_2.sql

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

## info: show make file environments.
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

# comments start with 2 # will available in help
help: Makefile
	@echo
	@echo "Choose a command to run:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: default fmt lint test build install mysql postgres sqlite info air help
