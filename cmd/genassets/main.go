package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	if err := os.MkdirAll("assets", 0o755); err != nil {
		panic(err)
	}

	// player.png — blue/purple character sprite 32×32
	player := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			// Body area: a centered 16×24 block in purple
			if x >= 8 && x < 24 && y >= 4 && y < 28 {
				player.Set(x, y, color.RGBA{R: 100, G: 50, B: 200, A: 255})
			} else {
				player.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}
	writeImage("assets/player.png", player)

	// grass.png — green-toned 32×32 tile
	grass := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			g := uint8(120 + (x+y)%30) // slight variation
			grass.Set(x, y, color.RGBA{R: 34, G: g, B: 34, A: 255})
		}
	}
	writeImage("assets/grass.png", grass)

	// path.png — brown/tan-toned 32×32 tile
	path := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			r := uint8(180 + (x+y)%20)
			g := uint8(140 + (x+y)%15)
			path.Set(x, y, color.RGBA{R: r, G: g, B: 80, A: 255})
		}
	}
	writeImage("assets/path.png", path)
}

func writeImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
