PAPER = shuffle-paper

.PHONY: pdf clean

pdf:
	pdflatex $(PAPER)
	biber    $(PAPER)
	pdflatex $(PAPER)
	pdflatex $(PAPER)

clean:
	rm -f $(PAPER).aux $(PAPER).bbl $(PAPER).bcf $(PAPER).blg \
	      $(PAPER).log $(PAPER).out $(PAPER).run.xml $(PAPER).toc \
	      $(PAPER).pdf
