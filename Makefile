VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
MODULE   := github.com/chrplr/shuffle-go
LDFLAGS  := -s -w -X=$(MODULE)/internal/version.Version=$(VERSION)

PAPER_DIR := paper
PAPER     := shuffle-paper

.PHONY: all binaries shuffle-cli shuffle-gui shuffle-gio paper-pdf clean help

all: binaries

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

## ── Go binaries ──────────────────────────────────────────────────────────────

binaries: shuffle-cli shuffle-gui shuffle-gio ## Build all three binaries

shuffle-cli: ## Build the command-line interface
	go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/shuffle-cli

shuffle-gui: ## Build the Fyne graphical interface
	go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/shuffle-gui

shuffle-gio: ## Build the Gio graphical interface
	go build -ldflags="$(LDFLAGS)" -o $@ ./cmd/shuffle-gio

## ── Paper ────────────────────────────────────────────────────────────────────

paper-pdf: ## Compile the LaTeX paper to PDF
	cd $(PAPER_DIR) && \
	pdflatex $(PAPER) && \
	biber    $(PAPER) && \
	pdflatex $(PAPER) && \
	pdflatex $(PAPER)

## ── Housekeeping ─────────────────────────────────────────────────────────────

clean: ## Remove binaries and LaTeX auxiliary files
	rm -f shuffle-cli shuffle-gui shuffle-gio
	cd $(PAPER_DIR) && \
	rm -f $(PAPER).aux $(PAPER).bbl $(PAPER).bcf $(PAPER).blg \
	      $(PAPER).log $(PAPER).out $(PAPER).run.xml $(PAPER).toc \
	      $(PAPER).lof $(PAPER).lot $(PAPER).fls $(PAPER).fdb_latexmk
