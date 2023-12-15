package main

import (
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
	sx, sy := snapXYtoGridf(float64(generalGridSize), float64(cx), float64(cy))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(g.imagePacks["UI"].images["cursor0"], op)
}

func (o *GameObject) CurrSprite() *ebiten.Image {
	if o.sprMapMode {
		return o.sprites.images[o.sprKey]
	}
	return o.sprites.imagesQ[o.sprIdx]
}

func (g *Game) DrawMatrixObjectAt(layerZ, x, y int, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	l := g.matrixLayers[layerZ]
	gx, gy := float64(x*int(l.squareLength)), float64(y*int(l.squareLength))
	op.GeoM.Translate(gx, gy)
	img := l.FirstObjectAt(x, y).CurrSprite()
	screen.DrawImage(img, op)
}

func (g *Game) DrawMatrixLayer(z int, screen *ebiten.Image) {
	l := g.matrixLayers[z]
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			if g.matrixLayers[z].isOccupied(x, y) {
				g.DrawMatrixObjectAt(z, x, y, screen)
			}
		}
	}
}

func (g *Game) DrawAllMatrixLayers(screen *ebiten.Image) {
	for i := 0; i < g.matrixLayerNum; i++ {
		g.DrawMatrixLayer(i, screen)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawAllMatrixLayers(screen)
	g.DrawCursor(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
