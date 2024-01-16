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

func initOneImagePack(paths ...string) *ImagePack {
	imgMap := make(map[string]*ebiten.Image)
	imgQ := make([]*ebiten.Image, 0)
	r := regexp.MustCompile(`^[^/]*/`)
	for _, p := range paths {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("./Graphics/%s.png", p))
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

func initImagePacks() map[string]*ImagePack { //TODO: just autoload
	packs := make(map[string]*ImagePack, 0)
	//Ground
	packs["Terrain"] = initOneImagePack(
		"Terrain/terrain0", "Terrain/terrain1", "Terrain/terrain2", "Terrain/terrain3")

	//Pawns
	packs["Pawn"] = initOneImagePack(
		"Pawns/playerDef", "Pawns/playerSel")

	//UI
	packs["UI"] = initOneImagePack(
		"UI/cursor0", "UI/walkable", "UI/attackable", "UI/enemyAttackable")

	//Zombies
	packs["Zombie"] = initOneImagePack(
		"Zombies/zombieDef")

	//Obstacles
	packs["Obstacles"] = initOneImagePack(
		"Obstacles/gravestone")

	packs["Obstacles/Fences"] = initOneImagePack(
		"Obstacles/Fences/fence1", "Obstacles/Fences/fence2", "Obstacles/Fences/fence3",
		"Obstacles/Fences/fenceEndMid", "Obstacles/Fences/fenceUpLeft",
		"Obstacles/Fences/fenceEndLeft",
		"Obstacles/Fences/fenceEndRight", "Obstacles/Fences/fenceUpRight")

	packs["Attacks"] = initOneImagePack(
		"Attacks/shove", "Attacks/throwRock")

	return packs
}
