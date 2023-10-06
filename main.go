package main

import (
	"chip8/cpu"
	setupkeys "chip8/setupkeys" // Configuration des touches
	"flag"                      //Gestion des drapeaux de ligne de commande
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten"            //Bibliothèque de jeu Ebiten, utilisé pour créer la fenêtre de jeu et gérer les entrées
	"github.com/hajimehoshi/ebiten/audio"      //Lecture audio depuis un fichier MP3
	"github.com/hajimehoshi/ebiten/audio/mp3"  //Support MP3 pour audio Ebiten
	"github.com/hajimehoshi/ebiten/ebitenutil" //Utilitaires Ebiten
)


var chip8 cpu.Chip8           // instance de la structure chip8
var audioPlayer *audio.Player // lecteur audio
var square *ebiten.Image      // image  pour afficher des pixels


// Cette fonction initialise une image de 10x10 pixels remplie de couleur blanche.
func init() {
	square, _ = ebiten.NewImage(10, 10, ebiten.FilterNearest)
	square.Fill(color.White)
}

// Cette fonction met à jour les touches en fonction des touches pressées dans la fenêtre de jeu
func getInput() bool {
	for key, value := range setupkeys.KeyMap {
		if ebiten.IsKeyPressed(key) {
			chip8.Keys[value] = 0x01
			return true
		}
	}
	return false
}

// Cette fonction met à jour l'écran du jeu en remplissant d'abord en noir. Ensuite, exécute la boucle principale du jeu
// pour gérer la logique de la CPU CHIP-8 et les entrées utilisateur.


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

// Cette fonction configure le jeu en chargeant la ROM spécifiée par les drapeaux de ligne de commande (-rom).
// Crée une instance de la CPU CHIP-8 et démarre la boucle de jeu.

func start() {
	rom := flag.String("rom", "", "")
	flag.Parse()

	if *rom == "" {
		os.Exit(1)
	} else {
		chip8 = cpu.NewCpu()
		chip8.LoadProgram(*rom)

		if err := ebiten.Run(update, 640, 320, 1, "CHIP-8 Emulator"); err != nil {
			panic(err)
		}
	}
}

// Cette fonction initialise le lecteur audio, configure les touches de jeu et lance le jeu CHIP-8.

func main() {
	audioContext, _ := audio.NewContext(48000)
	f, _ := ebitenutil.OpenFile("assets/beep.mp3")
	d, _ := mp3.Decode(audioContext, f)
	audioPlayer, _ = audio.NewPlayer(audioContext, d)
	setupkeys.SetupKeys()
	start()
}

