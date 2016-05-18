all: build

pylint: ;

build: ;

# For now let's disable testing because we don't have the Docker image on registry
# yet.
test: ;

clean: ;

jenkins: pylint build test;

.PHONY: all pylint build test clean jenkins
