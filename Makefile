all: build

pylint:
		pylint --rcfile=pylintrc --disable=I *.py test/*.py

build:
		./setup.py build

test:
		PYTHONPATH=$(PYTHONPATH):. ./setup.py test

clean:
		./setup.py clean -a

jenkins: pylint build test

.PHONY: all pylint build test clean jenkins
