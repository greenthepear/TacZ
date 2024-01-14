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

func InitFonts() {
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

func (g *Game) InitUILayers() {
	g.CreateNewFreeLayerOnTop("pawnInfoLayer")
	g.CreateNewMatrixLayerOnTop("underAttacksLayer", generalGridSize, 5, 1, 0, float64(boardHeight)+24)
	g.CreateNewMatrixLayerOnTop("attacksLayer", generalGridSize, 5, 1, 0, float64(boardHeight)+24)
}

func (g *Game) GenAttackString() string {
	if g.selectedAttack == nil {
		return "Attacks:"
	}
	return fmt.Sprintf("%v: %v",
		g.selectedAttack.name, g.attacks[g.selectedAttack.name].desc)
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

		//Populate attack selection layer
		if g.selectedPawn.vars["canAttack"] == 1 {
			text.Draw(screen, g.GenAttackString(), fontPressStart, 0, boardHeight+24, color.White)
			selectedPawnAttacks := g.CreateAttackObjectsOf(&pawnCopy)
			for i, o := range selectedPawnAttacks {
				g.AddObjectToMatrixLayer(o, attacksLayerZ, i, 0)
			}
		}
	} else {
		g.ClearFreeLayer(pawnInfoLayerZ)
		g.ClearMatrixLayer(attacksLayerZ)
	}
}
