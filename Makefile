include $(GOROOT)/src/Make.inc

TARG=arya
GOFILES=arya.go

include $(GOROOT)/src/Make.cmd

docs:
	@pandoc -s -w man -o arya.1 README.md
	@godoc -html > docs/arya.html