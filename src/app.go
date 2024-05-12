package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const title = "go-2048"

var logger = log.Default()
var screenWidth = 640
var screenHeight = 480

func main() {
	logger.Println("Starting game")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(title)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
