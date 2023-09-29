package main

const TAILLEMEMOIRE = 4096
const ADRESSEDEBUT = 512

type CPU struct {
	Memoire      [TAILLEMEMOIRE]byte // memoire
	V            [16]uint8           // le registre
	I            uint16              // stocke une adresse mémoire ou dessinateur
	NbrSaut      uint8               // stocke le nombre de sauts effectués pour ne pas dépasser 16
	CompteurJeu  uint8               // compteur jeu
	CompteurSon  uint8               // compteur son
	Saut         [16]uint16          // nbrSaut uint8
	Keypad       [16]uint8           // keypad
	Ecran        [64][32]uint8       // screen
	VitesseTimer uint8               // timerSpeed
	Op           uint16              // op
	PC           uint16              // pour parcourir le tableau « mémoire »
	SP           uint8               // sp
}

var cpu CPU // déclaration de notre CPU

func initialiserCpu() {
	// On initialise le tout
	var i uint16

	for i = 0; i < TAILLEMEMOIRE; i++ {
		cpu.Memoire[i] = 0
	}

	for i = 0; i < 16; i++ {
		cpu.V[i] = 0
		cpu.Saut[i] = 0
	}

	cpu.PC = ADRESSEDEBUT
	cpu.NbrSaut = 0
	cpu.CompteurJeu = 0
	cpu.CompteurSon = 0
	cpu.I = 0
}

func decompter() {
	// Implémentez la logique de décompte ici
	if cpu.CompteurJeu > 0 {
		cpu.CompteurJeu--
	}

	if cpu.CompteurSon > 0 {
		cpu.CompteurSon--
	}
}

//utiliser la fonction os.readFile("nomDuFichier")
//boucle qui lit tout le fichier et le rajoute en mémoire
