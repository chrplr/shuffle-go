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
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"

	"github.com/chrplr/shuffle-go"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type shuffleApp struct {
	theme *material.Theme

	textArea    widget.Editor
	constraints widget.Editor
	seed        widget.Editor
	limit       widget.Editor
	iter        widget.Editor
	delim       widget.Editor
	equiprob    widget.Bool

	shuffleBtn widget.Clickable
	importBtn  widget.Clickable
	status     string

	explorer *explorer.Explorer
	window   *app.Window
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Shuffle-Go GIO"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	
	s := &shuffleApp{
		theme:    th,
		window:   w,
		explorer: explorer.NewExplorer(w),
	}
	s.constraints.SetText("")
	s.seed.SetText("0")
	s.limit.SetText("0")
	s.iter.SetText("0")
	s.delim.SetText("")

	var ops op.Ops
	for {
		e := w.Event()
		s.explorer.ListenEvents(e)
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			s.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (s *shuffleApp) layout(gtx C) D {
	if s.shuffleBtn.Clicked(gtx) {
		s.runShuffle()
	}
	if s.importBtn.Clicked(gtx) {
		s.importFile()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				return material.H6(s.theme, "Shuffle-Go (Gio)").Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(0.3, func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(s.inputField("Constraints:", &s.constraints)),
							layout.Rigid(s.inputField("Seed (0=random):", &s.seed)),
							layout.Rigid(s.inputField("Limit (0=all):", &s.limit)),
							layout.Rigid(s.inputField("Max Iterations:", &s.iter)),
							layout.Rigid(s.inputField("Delimiter:", &s.delim)),
							layout.Rigid(func(gtx C) D {
								return material.CheckBox(s.theme, &s.equiprob, "Equiprobable").Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Flexed(0.5, func(gtx C) D {
										return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
											return material.Button(s.theme, &s.importBtn, "Import").Layout(gtx)
										})
									}),
									layout.Flexed(0.5, func(gtx C) D {
										return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
											return material.Button(s.theme, &s.shuffleBtn, "Shuffle!").Layout(gtx)
										})
									}),
								)
							}),
						)
					})
				}),
				layout.Flexed(0.7, func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
						border := widget.Border{
							Color:        color.NRGBA{A: 0xff, R: 0xcc, G: 0xcc, B: 0xcc},
							CornerRadius: unit.Dp(4),
							Width:        unit.Dp(1),
						}
						return border.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
								editor := material.Editor(s.theme, &s.textArea, "Enter data here...")
								return editor.Layout(gtx)
							})
						})
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				return material.Body2(s.theme, s.status).Layout(gtx)
			})
		}),
	)
}

func (s *shuffleApp) inputField(label string, editor *widget.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return material.Body2(s.theme, label).Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				border := widget.Border{
					Color:        color.NRGBA{A: 0xff, R: 0xcc, G: 0xcc, B: 0xcc},
					CornerRadius: unit.Dp(4),
					Width:        unit.Dp(1),
				}
				return border.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
						e := material.Editor(s.theme, editor, "")
						return e.Layout(gtx)
					})
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		)
	}
}

func (s *shuffleApp) runShuffle() {
	dataStr := s.textArea.Text()
	if strings.TrimSpace(dataStr) == "" {
		s.status = "Error: No data to shuffle"
		return
	}

	data, err := shuffle.LoadData(strings.NewReader(dataStr), s.delim.Text())
	if err != nil {
		s.status = fmt.Sprintf("Error loading data: %v", err)
		return
	}

	constraints, err := shuffle.ParseConstraints(s.constraints.Text())
	if err != nil {
		s.status = fmt.Sprintf("Error parsing constraints: %v", err)
		return
	}

	seed, _ := strconv.ParseInt(s.seed.Text(), 10, 64)
	limit, _ := strconv.Atoi(s.limit.Text())
	iter, _ := strconv.Atoi(s.iter.Text())

	sh := shuffle.NewShuffler(data, constraints, seed, iter, limit)

	var result [][]string
	if s.equiprob.Value {
		result, err = sh.ShuffleEquiprob()
	} else {
		result, err = sh.ShuffleConstructive()
	}

	if err != nil {
		s.status = fmt.Sprintf("Shuffle failed: %v", err)
		return
	}

	delimiter := s.delim.Text()
	if delimiter == "" {
		delimiter = " "
	}

	var output []string
	for _, row := range result {
		output = append(output, strings.Join(row, delimiter))
	}
	s.textArea.SetText(strings.Join(output, "\n"))
	s.status = fmt.Sprintf("Success! Processed %d lines.", len(result))
}

func (s *shuffleApp) importFile() {
	go func() {
		file, err := s.explorer.ChooseFile()
		if err != nil {
			s.status = fmt.Sprintf("Error choosing file: %v", err)
			s.window.Invalidate()
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			s.status = fmt.Sprintf("Error reading file: %v", err)
			s.window.Invalidate()
			return
		}

		s.textArea.SetText(string(data))
		s.status = "File imported successfully"
		s.window.Invalidate()
	}()
}
