package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

const (
	Largeur       = 64 // Nombre de pixels en largeur
	Longueur      = 32 // Nombre de pixels en longueur
	DimPixel      = 8  // Taille d'un pixel (carré de côté 8)
	LargeurEcran  = Largeur * DimPixel
	LongueurEcran = Longueur * DimPixel
)

type Pixel struct {
	Position ebiten.Vec2 // Regroupe l'abscisse et l'ordonnée
	Couleur  color.Color // Comme son nom l'indique, c'est la couleur
}

var pixel [Largeur][Longueur]Pixel // Déclaration de la variable globale pixel

var Noir = color.RGBA{0, 0, 0, 0} //noir avec une opacité nulle

func initialiserPixel() {
	for x := 0; x < Largeur; x++ {
		for y := 0; y < Longueur; y++ {
			pixel[x][y].Position = ebiten.Vec2{
				float64(x * DimPixel),
				float64(y * DimPixel),
			}
			pixel[x][y].Couleur = Noir // On initialise par défaut les pixels en noir
		}
	}
}

type Game struct{}

func (g Game) Update() error {
	//implémenter logique de mise à jour du jeu
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	//implémenter logique de rendu du jeu
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//TODO implement me
	return LargeurEcran, LongueurEcran
}

func main() {
	// Initialisation d'Ebiten
	ebiten.SetWindowSize(LargeurEcran, LongueurEcran)
	ebiten.SetWindowTitle("Émulateur CHIP-8")

	//Initialisation des pixels
	initialiserPixel()

	// Votre code principal ici, y compris la boucle de jeu avec Ebiten
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
