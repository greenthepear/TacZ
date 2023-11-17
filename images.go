package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImagePack struct {
	images  map[string]*ebiten.Image
	imagesQ []*ebiten.Image
}

func initOneImagePack(paths []string) *ImagePack {
	imgMap := make(map[string]*ebiten.Image)
	imgQ := make([]*ebiten.Image, 0)
	r := regexp.MustCompile(`^[^/]*/`)
	for _, p := range paths {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("./Graphics/%v.png", p))
		if err != nil {
			log.Fatal(err)
		}
		keyName := r.ReplaceAllString(p, "")
		imgMap[keyName] = img
		imgQ = append(imgQ, img)
	}

	pack := &ImagePack{
		images:  imgMap,
		imagesQ: imgQ,
	}
	return pack
}

func initImagePacks() map[string]*ImagePack {
	packs := make(map[string]*ImagePack, 0)
	//Ground
	groundPaths := []string{
		"Terrain/terrain0", "Terrain/terrain1", "Terrain/terrain2", "Terrain/terrain3"}
	packs["Terrain"] = initOneImagePack(groundPaths)

	//Pawns
	pawnPaths := []string{
		"Pawns/playerDef"}
	packs["Pawn"] = initOneImagePack(pawnPaths)

	//UI
	uiPaths := []string{
		"UI/cursor0"}
	packs["UI"] = initOneImagePack(uiPaths)

	return packs
}
