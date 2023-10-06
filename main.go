package main

import (
	"chip8/cpu"
	setupkeys "chip8/setupkeys"
	"flag"                                     //permet de traiter les drapeaux de ligne de commande, qui sont utilisés pour spécifier un fichier ROM à charger
	"github.com/hajimehoshi/ebiten"            //package principal de la bibliothèque de jeu Ebiten, qui sera utilisé pour créer la fenêtre de jeu et gérer les entrées
	"github.com/hajimehoshi/ebiten/audio"      //packages Ebiten sont utilisés pour la lecture audio depuis un fichier MP3
	"github.com/hajimehoshi/ebiten/audio/mp3"  //
	"github.com/hajimehoshi/ebiten/ebitenutil" //  package Ebiten fournit des utilitaires pour Ebiten
	"image/color"
	"os"
)

// variables globales pour le programme
var chip8 cpu.Chip8           // instance de la structure chip8
var audioPlayer *audio.Player // lecteur audio
var square *ebiten.Image      // image utilisée pour afficher des pixels

// La fonction init est appelée automatiquement avant la main.
// on initialise la variable square en créant une image de 10x10 pixels remplie de couleur noire.
func init() {
	square, _ = ebiten.NewImage(10, 10, ebiten.FilterNearest)
	square.Fill(color.Black)
}

// La fonction getInput parcourt la carte de touches définie dans le package setupkeys et
// met à jour les touches en fonction de l'état des touches pressées dans la fenêtre de jeu.
func getInput() bool {
	for key, value := range setupkeys.KeyMap {
		if ebiten.IsKeyPressed(key) {
			chip8.Keys[value] = 0x01
			return true
		}
	}
	return false
}

/*
**La fonction update est appelée à chaque rafraîchissement de la fenêtre de jeu (frame).
Elle remplit l'écran de jeu en noir, puis exécute la boucle principale du jeu 10 fois par frame. Dans cette boucle principale :
***chip8.Run() est appelé pour exécuter un cycle de la CPU de la CHIP-8.
***getInput() est appelé pour mettre à jour l'état des touches de la CHIP-8 en fonction des touches pressées dans la fenêtre.
**Si un dessin est requis (chip8.Draw est true) ou si aucune entrée n'a été reçue, la boucle interne dessine les pixels actifs de la CHIP-8 sur l'écran.
**L'état des touches est mis à jour en fonction des touches pressées dans la fenêtre.
**Si le minuteur sonore de la CHIP-8 est actif, le son est joué.
*/

func update(screen *ebiten.Image) error {
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
	// Boucle principale du jeu (cette boucle est exécutée 10 fois par frame)
	for i := 0; i < 10; i++ {
		chip8.Draw = false
		chip8.Inputflag = false
		gotInput := true
		chip8.Run()

		if chip8.Inputflag {
			gotInput = getInput()
			if !gotInput {
				chip8.Pc = chip8.Pc - 2
			}
		}

		if chip8.Draw || !gotInput {
			// Boucle pour dessiner les pixels actifs de la CHIP-8
			for i := 0; i < 32; i++ {
				for j := 0; j < 64; j++ {
					if chip8.Display[i][j] == 0x01 {

						opts := &ebiten.DrawImageOptions{}

						opts.GeoM.Translate(float64(j*10), float64(i*10))

						screen.DrawImage(square, opts)
					}
				}
			}
		}
		// Mettre à jour l'état des touches
		for key, value := range setupkeys.KeyMap {
			if ebiten.IsKeyPressed(key) {
				chip8.Keys[value] = 0x01
			} else {
				chip8.Keys[value] = 0x00
			}
		}

		if chip8.SoundTimer > 0 {
			audioPlayer.Play()
			audioPlayer.Rewind()
		}

	}

	return nil
}

/*
La fonction start est responsable de la configuration initiale du jeu.
Elle analyse les drapeaux de ligne de commande pour obtenir le chemin de la ROM à charger (-rom),
puis crée une instance de la CPU CHIP-8 (cpu.NewCpu()) et
charge la ROM spécifiée à l'aide de chip8.LoadProgram(*rom).
Enfin, elle démarre la boucle de jeu en appelant ebiten.Run avec la fonction update
*/
func start() {
	rom := flag.String("rom", "", "")
	flag.Parse()

	if *rom == "" {
		os.Exit(0)
	} else {
		chip8 = cpu.NewCpu()
		chip8.LoadProgram(*rom)

		if err := ebiten.Run(update, 640, 320, 1, "CHIP-8 Emulator"); err != nil {
			panic(err)
		}
	}
}

/*
La fonction main point d'entrée du programme.
Elle configure le lecteur audio en utilisant les fichiers audio spécifiés, initialise les touches de jeu avec setupkeys.SetupKeys(),
puis démarre le jeu en appelant start().
C'est essentiellement comment le point d'entrée de votre jeu CHIP-8 est organisé.
N'hésitez pas à poser des questions spécifiques sur des parties de ce code si vous avez des doutes.
*/
func main() {
	audioContext, _ := audio.NewContext(48000)
	f, _ := ebitenutil.OpenFile("assets/beep.mp3")
	d, _ := mp3.Decode(audioContext, f)
	audioPlayer, _ = audio.NewPlayer(audioContext, d)
	setupkeys.SetupKeys()
	start()
}
