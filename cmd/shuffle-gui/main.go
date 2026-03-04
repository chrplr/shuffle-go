// Copyright (C) 2026 Christophe Pallier
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/chrplr/shuffle-go/internal/version"
	"github.com/chrplr/shuffle-go/shuffle"
)

type shuffleApp struct {
	window      fyne.Window
	textArea    *widget.Entry
	constraints *widget.Entry
	seed        *widget.Entry
	limit       *widget.Entry
	iter        *widget.Entry
	equiprob    *widget.Check
	delim       *widget.Entry
	status      *widget.Label
}

func main() {
	versionFlag := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("shuffle-gui version %s\n", version.Version)
		os.Exit(0)
	}

	a := app.NewWithID("com.chrplr.shuffle")
	w := a.NewWindow("Shuffle-Go GUI")

	sApp := &shuffleApp{
		window: w,
		textArea: widget.NewMultiLineEntry(),
		constraints: widget.NewEntry(),
		seed:        widget.NewEntry(),
		limit:       widget.NewEntry(),
		iter:        widget.NewEntry(),
		equiprob:    widget.NewCheck("Equiprobable", nil),
		delim:       widget.NewEntry(),
		status:      widget.NewLabel("Ready"),
	}

	sApp.textArea.SetPlaceHolder("Enter data here...")
	sApp.constraints.SetPlaceHolder("e.g. 1 2 -3")
	sApp.seed.SetText("0")
	sApp.limit.SetText("0")
	sApp.iter.SetText("0")

	openBtn := widget.NewButton("Open", sApp.openFile)
	saveBtn := widget.NewButton("Save", sApp.saveFile)
	shuffleBtn := widget.NewButton("Shuffle!", sApp.runShuffle)
	helpBtn := widget.NewButton("Help", sApp.showHelp)

	topBar := container.NewHBox(openBtn, saveBtn, helpBtn)
	
	controls := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Constraints:"), sApp.constraints,
			widget.NewLabel("Seed (0=random):"), sApp.seed,
			widget.NewLabel("Limit (0=all):"), sApp.limit,
			widget.NewLabel("Max Iterations:"), sApp.iter,
			widget.NewLabel("Delimiter:"), sApp.delim,
			widget.NewLabel("Algorithm:"), sApp.equiprob,
		),
		shuffleBtn,
	)

	content := container.NewBorder(topBar, container.NewVBox(controls, sApp.status), nil, nil, sApp.textArea)
	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func (a *shuffleApp) openFile() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()
		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.textArea.SetText(string(data))
		a.status.SetText(fmt.Sprintf("Opened %s", reader.URI().Name()))
	}, a.window)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".csv"}))
	fd.Show()
}

func (a *shuffleApp) saveFile() {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		_, err = writer.Write([]byte(a.textArea.Text))
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.status.SetText(fmt.Sprintf("Saved %s", writer.URI().Name()))
	}, a.window)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".csv"}))
	fd.Show()
}

func (a *shuffleApp) runShuffle() {
	dataStr := a.textArea.Text
	if strings.TrimSpace(dataStr) == "" {
		a.status.SetText("Error: No data to shuffle")
		return
	}

	data, err := shuffle.LoadData(strings.NewReader(dataStr), a.delim.Text)
	if err != nil {
		a.status.SetText(fmt.Sprintf("Error loading data: %v", err))
		return
	}

	constraints, err := shuffle.ParseConstraints(a.constraints.Text)
	if err != nil {
		a.status.SetText(fmt.Sprintf("Error parsing constraints: %v", err))
		return
	}

	seed, _ := strconv.ParseInt(a.seed.Text, 10, 64)
	limit, _ := strconv.Atoi(a.limit.Text)
	iter, _ := strconv.Atoi(a.iter.Text)

	s := shuffle.NewShuffler(data, constraints, seed, iter, limit)

	var result [][]string
	if a.equiprob.Checked {
		result, err = s.ShuffleEquiprob()
	} else {
		result, err = s.ShuffleConstructive()
	}

	if err != nil {
		a.status.SetText(fmt.Sprintf("Shuffle failed: %v", err))
		return
	}

	delimiter := a.delim.Text
	if delimiter == "" {
		delimiter = " "
	}

	var output []string
	for _, row := range result {
		output = append(output, strings.Join(row, delimiter))
	}
	a.textArea.SetText(strings.Join(output, "\n"))
	a.status.SetText(fmt.Sprintf("Success! Processed %d lines.", len(result)))
}

func (a *shuffleApp) showHelp() {
	helpText := `Shuffle-Go GUI Help:

Constraints: Space-separated numbers.
- Positive 'n': Max consecutive repetitions of labels in that column.
- Negative '-m': Min distance (gap) between identical labels in that column.
- Zero: No constraint for that column.

Example: "1 2 -3" means:
- Col 1: Max 1 repetition (no adjacent identical labels).
- Col 2: Max 2 consecutive identical labels.
- Col 3: Min gap of 3 lines between identical labels.

Algorithms:
- Constructive (default): Fast, swaps lines to meet constraints.
- Equiprobable: Brute force, ensures every valid permutation is equally likely (slower).

Seed: Fixed number for reproducible results (0 for random).
`
	d := dialog.NewInformation("Help", helpText, a.window)
	d.Show()
}
