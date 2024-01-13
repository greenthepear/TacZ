package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func snapXYtoGridf(squareLength, x, y float64) (float64, float64) {
	x = math.Floor(x/squareLength) * squareLength
	y = math.Floor(y/squareLength) * squareLength
	return x, y
}

func snapXYtoGrid(squareLength px, x, y int) (int, int) {
	sl := float64(squareLength)
	xf := math.Floor(float64(x)/sl) * sl
	yf := math.Floor(float64(y)/sl) * sl
	return int(xf), int(yf)
}

func (g *Game) DrawCursor(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()
	if cy > generalGridHeight*int(generalGridSize)-1 {
		return
	}
	sx, sy := snapXYtoGridf(float64(generalGridSize), float64(cx), float64(cy))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(g.imagePacks["UI"].images["cursor0"], op)
}

func (o *GameObject) CurrSprite() *ebiten.Image {
	if o.sprMapMode {
		//log.Printf("Obj: %v | Trying to draw sprite %v\n...", o.name, o.sprKey)
		if o.sprites.images[o.sprKey] == nil {
			log.Fatalf("No sprite under key \"%v\"", o.sprKey)
		}
		return o.sprites.images[o.sprKey]
	}
	if o.sprites.imagesQ[o.sprIdx] == nil {
		log.Fatalf("No sprite under index \"%v\"", o.sprIdx)
	}
	return o.sprites.imagesQ[o.sprIdx]
}

func (g *Game) DrawMatrixObjectAt(layerZ, x, y int, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	l := g.matrixLayers[layerZ]
	gx, gy := float64(x*int(l.squareLength)), float64(y*int(l.squareLength))
	op.GeoM.Translate(gx+l.xOffset, gy+l.yOffset)
	img := l.FirstObjectAt(x, y).CurrSprite()
	screen.DrawImage(img, op)
}

func (g *Game) DrawMatrixLayer(l *MatrixLayer, screen *ebiten.Image) {
	z := l.z
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			if g.matrixLayers[z].isOccupied(x, y) {
				g.DrawMatrixObjectAt(z, x, y, screen)
			}
		}
	}
}

func (g *Game) DrawAllMatrixLayers(screen *ebiten.Image) {
	for _, l := range g.matrixLayers {
		g.DrawMatrixLayer(l, screen)
	}
}

func (g *Game) DrawFreeLayer(l *FreeObjectLayer, screen *ebiten.Image) {
	for _, o := range l.objects {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.x), float64(o.y))
		img := o.CurrSprite()
		screen.DrawImage(img, op)
	}
}

func (g *Game) DrawAllFreeLayers(screen *ebiten.Image) {
	for _, l := range g.freeLayers {
		g.DrawFreeLayer(l, screen)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawAllMatrixLayers(screen)
	g.DrawAllFreeLayers(screen)
	g.DrawCursor(screen)
	g.DrawSelectedPawnInfo(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return boardWidth, boardHeight + bottomHeight
}
