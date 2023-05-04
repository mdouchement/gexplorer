package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mdouchement/gexplorer/gioexplorer"
)

type UI struct {
	window *app.Window
	fonts  []font.FontFace
	theme  *material.Theme
	ops    op.Ops

	explorer     *gioexplorer.Explorer
	chooseFile   widget.Clickable
	chooseFiles  widget.Clickable
	chooseImage  widget.Clickable
	chooseImages widget.Clickable
	createFile   widget.Clickable
	list         widget.List
	filenames    []string
}

func main() {
	ui := &UI{
		window: app.NewWindow(
			app.Title("gio-explorer"),
			app.Size(unit.Dp(1024), unit.Dp(860)),
		),
		fonts:     gofont.Collection(),
		list:      widget.List{List: layout.List{Axis: layout.Vertical}},
		filenames: make([]string, 0, 1),
	}
	ui.theme = material.NewTheme(ui.fonts)
	ui.explorer = gioexplorer.NewExplorer(ui.window)

	//

	go func() {
		err := ui.loop()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()
}

func (ui *UI) loop() error {
	for e := range ui.window.Events() {
		ui.explorer.ListenEvents(e)

		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ui.ops, e)

			err := ui.events()
			if err != nil {
				log.Fatal(err)
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return magins().Layout(gtx, ui.layoutButton)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return magins().Layout(gtx, ui.layoutList)
					},
				),
			)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

func (ui *UI) events() error {
	if ui.createFile.Clicked() {
		go func() {
			filename, err := ui.explorer.CreateFile("default-name.txt")
			if err != nil {
				err = fmt.Errorf("failed creating image file: %w", err)
				fmt.Println(err)
				return
			}

			ui.filenames = ui.filenames[:0]
			ui.filenames = append(ui.filenames, filename)
		}()
	}

	if ui.chooseFile.Clicked() {
		go func() {
			filename, err := ui.explorer.ChooseFile()
			if err != nil {
				err = fmt.Errorf("failed opening file: %w", err)
				fmt.Println(err)
				return
			}

			ui.filenames = ui.filenames[:0]
			ui.filenames = append(ui.filenames, filename)
		}()
	}

	if ui.chooseFiles.Clicked() {
		go func() {
			var err error

			ui.filenames, err = ui.explorer.ChooseFiles()
			if err != nil {
				err = fmt.Errorf("failed opening files: %w", err)
				fmt.Println(err)
				return
			}
		}()
	}

	if ui.chooseImage.Clicked() {
		go func() {
			filename, err := ui.explorer.ChooseFile("png", "jpeg", "jpg")
			if err != nil {
				err = fmt.Errorf("failed opening image file: %w", err)
				fmt.Println(err)
				return
			}

			ui.filenames = ui.filenames[:0]
			ui.filenames = append(ui.filenames, filename)
		}()
	}

	if ui.chooseImages.Clicked() {
		go func() {
			var err error

			ui.filenames, err = ui.explorer.ChooseFiles("png", "jpeg", "jpg")
			if err != nil {
				err = fmt.Errorf("failed opening image files: %w", err)
				fmt.Println(err)
				return
			}
		}()
	}

	return nil
}

func (ui *UI) layoutButton(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(ui.theme, &ui.createFile, "CreateFile")
						return btn.Layout(gtx)
					},
				)
			},
		),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(ui.theme, &ui.chooseFile, "ChooseFile")
						return btn.Layout(gtx)
					},
				)
			},
		),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(ui.theme, &ui.chooseFiles, "ChooseFiles")
						return btn.Layout(gtx)
					},
				)
			},
		),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(ui.theme, &ui.chooseImage, "ChooseImage")
						return btn.Layout(gtx)
					},
				)
			},
		),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(ui.theme, &ui.chooseImages, "ChooseImages")
						return btn.Layout(gtx)
					},
				)
			},
		),
	)
}

func (ui *UI) layoutList(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx, material.Body1(ui.theme, "List of files").Layout)
			},
		),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return magins().Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						border := widget.Border{
							Color:        color.NRGBA{A: 255},
							Width:        unit.Dp(1),
							CornerRadius: unit.Dp(4),
						}
						return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(4),
								Bottom: unit.Dp(4),
								Right:  unit.Dp(4),
								Left:   unit.Dp(4),
							}
							return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.List(ui.theme, &ui.list).Layout(gtx, len(ui.filenames), func(gtx layout.Context, i int) layout.Dimensions {
									label := material.Body1(ui.theme, ui.filenames[i])
									if i%2 == 1 {
										label.Color = color.NRGBA{A: 255, R: 0, G: 168, B: 107}
									}
									return label.Layout(gtx)
								})
							})
						})
					},
				)
			},
		),
	)
}

func magins() layout.Inset {
	return layout.Inset{
		Top:    unit.Dp(2),
		Bottom: unit.Dp(2),
		Right:  unit.Dp(2),
		Left:   unit.Dp(2),
	}
}
