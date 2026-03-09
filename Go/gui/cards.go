package main

import (
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type CardSprites struct {
	Faces map[string]*ebiten.Image
	Back  *ebiten.Image
	W     int
	H     int
}

func loadCardSprites(path string) (CardSprites, error) {
	file, err := os.Open(path)
	if err != nil {
		return CardSprites{}, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return CardSprites{}, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// This sprite sheet is a fixed atlas:
	// 13 columns x 5 rows, each cell 18x22.
	const (
		cellW = 18
		cellH = 22
	)

	suits := []string{"Spades", "Clubs", "Diamonds", "Hearts"}
	ranks := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	faces := make(map[string]*ebiten.Image, 52)
	for row := 0; row < 4; row++ {
		for col := 0; col < 13; col++ {
			x := col * cellW
			y := row * cellH
			rect := image.Rect(x, y, x+cellW, y+cellH)
			key := ranks[col] + "_of_" + suits[row]
			faces[key] = sheet.SubImage(rect).(*ebiten.Image)
		}
	}

	// Bottom row has jokers on the left and card backs on the right.
	// Use first back at row 4, col 10.
	backRect := image.Rect(10*cellW, 4*cellH, 11*cellW, 5*cellH)
	back := sheet.SubImage(backRect).(*ebiten.Image)

	return CardSprites{
		Faces: faces,
		Back:  back,
		W:     cellW,
		H:     cellH,
	}, nil
}
