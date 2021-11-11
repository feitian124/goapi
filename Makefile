PKG = github.com/feitian124/goapi
COMMIT = $$(git describe --tags --always)

# get date
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
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

default: check

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

test:
	go test ./database/drivers/mysql -v -coverprofile=coverage.out -covermode=count

test_cover: test
	go tool cover -html=coverage.out

build:
	go build -ldflags="$(BUILD_LDFLAGS)"

install:
	go install github.com/xo/usql@v0.9.4

check:
	@gofmt -s -w ${GO_FILES}
	@go vet $(VET_PACKAGES)

info:
	@echo "COMMIT：${COMMIT}"
	@echo "OSNAME: ${OSNAME}"
	@echo "DATE: ${DATE}"
	@echo -e "\nPACKAGES:"
	@echo ${PACKAGES}
	@echo -e "\nVET_PACKAGES:"
	@echo ${VET_PACKAGES}
	@echo -e "\nGO_FILES:"
	@echo ${GO_FILES}

.PHONY: default check test mysql postgres sqlite info