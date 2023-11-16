package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
)

var (
	// Color
	headEven = color.RGBA{0x80, 0x80, 0x80, 0xFF}
	bodyEven = color.RGBA{0xC0, 0xC0, 0xC0, 0xFF}
	headOdd  = color.RGBA{0xA9, 0xA9, 0xA9, 0xFF}
	bodyOdd  = color.RGBA{0xDC, 0xDC, 0xDC, 0xFF}
	// Window
	scr        screen.Screen
	window     screen.Window
	windowSize = image.Point{}
	lineOffset = 0
	lineHeight = 30
	headWidth  = 50
	// Font
	fontHeight = 13
	fontWidth  = 7
	// Games
	games []Game
)

type Game struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Args string `json:"args"`
}

func main() {
	// Load Games
	f, err := os.ReadFile("./games.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(f, &games)

	// Create Window
	driver.Main(func(s screen.Screen) {
		scr = s

		var err error
		window, err = s.NewWindow(&screen.NewWindowOptions{
			Width:  640,
			Height: 480,
			Title:  "Select Game",
		})
		if err != nil {
			panic(err)
		}
		defer window.Release()

		for {
			event := window.NextEvent()

			switch e := event.(type) {
			case lifecycle.Event:
				log.Printf("LifeCycle   : %s => %s\n", e.From, e.To)
				if e.To == lifecycle.StageDead {
					return
				}
			case size.Event:
				log.Printf("SizeChange  : %dx%d(Rotation:%d)\n", e.HeightPx, e.WidthPx, e.Orientation)
				windowSize = image.Point{e.WidthPx, e.HeightPx}
				UpdateGrid()
			case mouse.Event:
				switch e.Button {
				case 1: // Click
					// Check Click End
					if e.Direction != 2 {
						continue
					}
					// Game Check
					index := int(e.Y)/lineHeight + lineOffset
					if len(games)-1 < index {
						continue
					}
					// Move WorkDir
					game := games[index]
					os.Chdir(filepath.Dir(game.Path))
					err := exec.Command(game.Path, strings.Split(game.Args, " ")...).Start()
					if err != nil {
						panic(err)
					}
					return
				case -1: // Page Up
					if lineOffset == 0 {
						continue
					}
					lineOffset--
					UpdateGrid()
				case -2: // Page Down
					lineOffset++
					UpdateGrid()
				default:
					continue
				}
				fmt.Printf("%#v\n", e)
			}
		}
	})
}

func UpdateGrid() {
	// Buffer
	buf, err := scr.NewBuffer(windowSize)
	if err != nil {
		log.Fatal(err)
	}
	// Draw
	drawer := buf.RGBA()
	bounds := drawer.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			// Head
			if x < headWidth {
				// Background Color
				if (y/lineHeight+lineOffset)%2 == 0 {
					drawer.SetRGBA(x, y, headEven)
				} else {
					drawer.SetRGBA(x, y, headOdd)
				}
				// Numbering
				if y%lineHeight == 0 {
					d := font.Drawer{
						Dst:  drawer,
						Src:  image.Black,
						Face: basicfont.Face7x13,
						Dot:  fixed.P(0, y-(fontHeight-lineHeight)),
					}
					d.DrawString(fmt.Sprintf("% 4d", y/lineHeight+lineOffset+1))
				}
				continue
			}
			// Body
			// Background Color
			if (y/lineHeight+lineOffset)%2 == 0 {
				drawer.SetRGBA(x, y, bodyEven)
			} else {
				drawer.SetRGBA(x, y, bodyOdd)
			}
			// Game Title
			if y%lineHeight == 0 {
				d := font.Drawer{
					Dst:  drawer,
					Src:  image.Black,
					Face: basicfont.Face7x13,
					Dot:  fixed.P(headWidth+fontWidth, y-(fontHeight-lineHeight)),
				}
				if len(games)-1 < y/lineHeight+lineOffset {
					continue
				}
				game := games[y/lineHeight+lineOffset]
				d.DrawString(fmt.Sprintf("% -20s(%s)", game.Name, game.Path))
			}

		}
	}
	// Write
	window.Upload(image.Point{}, buf, bounds)
	window.Publish()
}
