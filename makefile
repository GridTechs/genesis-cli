GO111MODULE=on

DIRECTORIES=$(sort $(dir $(wildcard pkg/*/*/)))

MOCKS=$(foreach x, $(DIRECTORIES), mocks/$(x))

OUTPUT_DIR=./bin

.PHONY: build test test_race lint vet get mocks clean-mocks manual-mocks
.ONESHELL:

all: genesis

genesis: | prep get
	go build -o $(OUTPUT_DIR)/genesis ./cmd/genesis

prep:
	@mkdir $(OUTPUT_DIR) 2>> /dev/null | true


linux:
	@mkdir -p $(OUTPUT_DIR)/linux 2>> /dev/null | true
	GOOS=linux
	go build -o $(OUTPUT_DIR)/linux/genesis ./cmd/genesis
mac:
	@mkdir -p $(OUTPUT_DIR)/mac 2>> /dev/null | true
	GOOS=macos
	go build -o $(OUTPUT_DIR)/mac/genesis ./cmd/genesis
windows:
	@mkdir -p $(OUTPUT_DIR)/windows 2>> /dev/null | true
	GOOS=windows
	go build -o $(OUTPUT_DIR)/windows/genesis.exe ./cmd/genesis

multiplatform: linux mac windows

install:
	go install cmd/genesis/main.go

test:
	go test ./...

lint:
	golint ./...

vet:
	go vet ./...

get:
	go get ./...

clean-mocks:
	rm -rf mocks

mocks: $(MOCKS)
	
$(MOCKS): mocks/% : %
	mockery -output=$@ -dir=$^ -all
