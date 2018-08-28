.PHONY: all help test doc

all: help

help:				## Show this help
	@scripts/help.sh

test:				## Test potential bugs and race conditions
	@scripts/test.sh

doc:				## Generate docs
	@scripts/doc.sh
