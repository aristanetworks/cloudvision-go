all: build

pylint:
		pylint --rcfile=pylintrc --disable=I *.py test/*.py

build:
		./setup.py build

test:
		PYTHONPATH=$(PYTHONPATH):. ./setup.py test

clean:
		./setup.py clean -a

#jenkins: pylint build test
jenkins: pylint build

push:
	# update dist
	a4 ssh dist git --git-dir=/dist/storage/ardc-config/.git --work-tree=/dist/storage/ardc-config pull --rebase

.PHONY: all pylint build clean jenkins test
