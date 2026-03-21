package main

import (
	"log"

	"github.com/abjeloevic88/go-game/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("ARPG")
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
