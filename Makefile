# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)