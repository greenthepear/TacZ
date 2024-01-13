package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontPressStart font.Face

func initFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	fontPressStart, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) DrawSelectedPawnInfo(screen *ebiten.Image) {
	if g.selectedPawn != nil {
		pawnCopy := *g.selectedPawn

		pawnCopy.x = 0
		pawnCopy.y = boardHeight

		infoString := fmt.Sprintf("Health: %.0f/%.0f\nMoves: %.0f/%.0f",
			pawnCopy.vars["leftHP"], pawnCopy.vars["maxHP"],
			pawnCopy.vars["leftMovement"], pawnCopy.vars["maxMovement"],
		)

		text.Draw(screen, infoString, fontPressStart, 16, boardHeight+8, color.White)

		g.AddObjectToFreeLayer(pawnInfoLayerZ, &pawnCopy)
	} else {
		g.ClearFreeLayer(pawnInfoLayerZ)
	}
}
