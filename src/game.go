package main

import (
	"bytes"
	"image/color"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	board [4][4]piece // [Y][X] coordinates
	font  *text.GoTextFaceSource
	ready bool
}

var width = 320
var height = 240
var blockDim = 40

var fontSize = 12.0

func (g *Game) init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return
	}
	g.font = s
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBoard(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func (g *Game) Update() error {
	if !g.ready {
		g.ready = true
		g.init()
		g.board[2][2] = piece{value: 2, hasMerged: false}
		g.board[3][3] = piece{value: 2, hasMerged: false}
	}
	hasMoved := false
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.shiftBoardLeftRight(0, 4, 1) // TODO: Convert these values into some enum.
		hasMoved = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.shiftBoardLeftRight(3, -1, -1)
		hasMoved = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.shiftBoardUpDown(0, 4, 1)
		hasMoved = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.shiftBoardUpDown(3, -1, -1)
		hasMoved = true
	}

	if hasMoved {
		g.placeRandomPiece()
		g.resetMerges()
	}
	return nil
}

func (g *Game) placeRandomPiece() {
	pieceReferences := make([]*piece, 0)
	for y := 0; y < 4; y += 1 {
		for x := 0; x < 4; x += 1 {
			piece := &g.board[y][x]
			if piece.value == 0 {
				pieceReferences = append(pieceReferences, piece)
			}
		}
	}

	var p *piece

	// Pull out either the first piece or a random one.
	if len(pieceReferences) == 0 {
		// This is the future losing state - but need to verify no merges can be done.
		return
	} else if len(pieceReferences) == 1 {
		p = pieceReferences[0]
	} else {
		p = pieceReferences[rand.Intn(len(pieceReferences)-1)]
	}

	// Algorithm for 2048, at least what I found, was 90% 2 and 10% 4 for new tiles.
	r := rand.Intn(101)
	if r < 90 {
		p.value = 2
	} else {
		p.value = 4
	}
}

func (g *Game) resetMerges() {
	for y := 0; y < 4; y += 1 {
		for x := 0; x < 4; x += 1 {
			g.board[y][x].hasMerged = false
		}
	}
}

func (g *Game) shiftBoardUpDown(yStart int, yEnd int, yInc int) {
	for x := 0; x < 4; x += 1 {
		actionsPerformed := false
		finished := false
		for !finished {
			for y := yStart; y != yEnd; y += yInc {
				currentPiece := g.board[y][x]
				canMoveDown := currentPiece.value == 0 && y+yInc != yEnd && g.board[y+yInc][x].value != 0
				canCombine := y+yInc != yEnd && g.board[y+yInc][x].value == currentPiece.value &&
					!currentPiece.hasMerged && !g.board[y+yInc][x].hasMerged && currentPiece.value != 0

				// Movement states
				if canMoveDown {
					actionsPerformed = actionsPerformed || true
					g.board[y][x] = g.board[y+yInc][x]
					g.board[y+yInc][x] = piece{value: 0, hasMerged: false}
				}
				if canCombine {
					actionsPerformed = actionsPerformed || true
					g.board[y][x].value = g.board[y][x].value * 2
					g.board[y][x].hasMerged = true
					g.board[y+yInc][x].value = 0
				}

				// End states.
				if y == yEnd+-yInc && actionsPerformed {
					actionsPerformed = false
				} else if y == yEnd+-yInc && !actionsPerformed {
					finished = true
				}
			}
		}
	}
}

func (g *Game) shiftBoardLeftRight(xStart int, xEnd int, xInc int) {
	for y := 0; y < 4; y += 1 {
		actionsPerformed := false
		finished := false
		for !finished {
			for x := xStart; x != xEnd; x += xInc {
				currentPiece := g.board[y][x]
				canMoveDown := currentPiece.value == 0 && x+xInc != xEnd && g.board[y][x+xInc].value != 0
				canCombine := x+xInc != xEnd && g.board[y][x+xInc].value == currentPiece.value &&
					!currentPiece.hasMerged && !g.board[y][x+xInc].hasMerged && currentPiece.value != 0

				// Movement states
				if canMoveDown {
					actionsPerformed = actionsPerformed || true
					g.board[y][x] = g.board[y][x+xInc]
					g.board[y][x+xInc] = piece{value: 0, hasMerged: false}
				}
				if canCombine {
					actionsPerformed = actionsPerformed || true
					g.board[y][x].value = g.board[y][x].value * 2
					g.board[y][x].hasMerged = true
					g.board[y][x+xInc].value = 0
				}

				// End states.
				if x == xEnd+-xInc && actionsPerformed {
					actionsPerformed = false
				} else if x == xEnd+-xInc && !actionsPerformed {
					finished = true
				}
			}
		}
	}
}

func (g *Game) drawText(screen *ebiten.Image, xPosition float64,
	yPosition float64, xIndex int, yIndex int) {
	currentPiece := g.board[yIndex][xIndex]
	if currentPiece.value == 0 {
		return
	}
	options := text.DrawOptions{}
	options.GeoM.Translate(xPosition, yPosition)
	options.ColorScale.ScaleWithColor(color.Black)
	options.LayoutOptions.PrimaryAlign = text.AlignCenter
	options.LayoutOptions.SecondaryAlign = text.AlignCenter
	text.Draw(screen, strconv.Itoa(currentPiece.value), &text.GoTextFace{
		Source: g.font,
		Size:   fontSize,
	}, &options)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	var blocksInDim = 4
	var spacing = 5
	var startingX = (width - (blocksInDim*blockDim + spacing)) / 2
	var startingY = (height - (blocksInDim*blockDim + spacing)) / 2
	for y := 0; y < 4; y += 1 {
		for x := 0; x < 4; x += 1 {
			xLoc := startingX + (x * 45)
			yLoc := startingY + (y * 45)
			color := g.board[y][x].determineColor()
			vector.DrawFilledRect(screen, float32(xLoc),
				float32(yLoc), float32(blockDim), float32(blockDim), color, false)
			g.drawText(screen, float64(xLoc+blockDim/2), float64(yLoc+blockDim/2), x, y)
		}
	}
}
