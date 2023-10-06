//logique de l'émulateur CHIP-8
//Il s'agit d'une grande partie du cœur de l'émulateur où les instructions CHIP-8 sont interprétées et exécutées

package cpu

/*
Ces instructions importent les packages nécessaires pour votre émulateur.
fontset  contient le jeu de caractères de la CHIP-8,
io pour l'entrée/sortie,
math/rand pour la génération de nombres aléatoires,
os pour les opérations système,
time pour la gestion du temps
*/
import (
	"chip8/fontset"
	"io"
	"math/rand"
	"os"
	"time"
)

// déclaration de deux constantes height et width pour définir la hauteur et la largeur de l'écran de la CHIP-8.
// Ce sont des valeurs spécifiques à la CHIP-8
const (
	height = byte(0x20)
	width  = byte(0x40)
)

/*
Ici, vous déclarez la structure Chip8 qui représente l'état de la machine CHIP-8.
Elle contient des champs pour le compteur de programme (Pc),
la mémoire, la pile,
les registres V0 à VF,
les timers de délai et de son, l'écran, l'état des touches, des drapeaux pour le dessin et l'entrée, et un registre pour l'entrée
*/
type Chip8 struct {
	Pc            uint16
	Memory        [4096]byte
	Stack         [16]uint16
	Sp            uint16
	V             [16]byte
	I             uint16
	DelayTimer    byte
	SoundTimer    byte
	Display       [height][width]byte
	Keys          [16]byte
	Draw          bool
	Inputflag     bool
	InputRegister byte
}

/*
La fonction NewCpu crée une nouvelle instance de Chip8,
initialise le compteur de programme (Pc) à l'adresse de démarrage typique de la CHIP-8 (0x200) et charge le jeu de caractères (LoadFontSet).
*/
func NewCpu() Chip8 {
	c := Chip8{Pc: 0x200}
	c.LoadFontSet()
	return c
}

// La fonction LoadFontSet charge le jeu de caractères de la CHIP-8 dans la mémoire de l'émulateur.
func (c *Chip8) LoadFontSet() {
	for i := 0x00; i < 0x50; i++ {
		c.Memory[i] = fontset.Fontset[i]
	}
}

// La fonction ClearDisplay efface l'écran en initialisant tous les pixels à 0.
func (c *Chip8) ClearDisplay() {
	for x := 0x00; x < 0x20; x++ {
		for y := 0x00; y < 0x40; y++ {
			c.Display[x][y] = 0x00
		}
	}
}

/*
La fonction LoadProgram charge un fichier ROM (jeu) dans la mémoire de la CHIP-8 à partir du chemin spécifié.
Elle ouvre le fichier, lit son contenu,
puis le copie dans la mémoire de l'émulateur à partir de l'adresse 0x200.
*/
func (c *Chip8) LoadProgram(rom string) int {
	f, err := os.Open(rom)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	memory := make([]byte, 3584)
	n, err := f.Read(memory)
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	for index, b := range memory {
		c.Memory[index+0x200] = b
	}
	return n
}

// La fonction Reset réinitialise complètement l'état de la machine CHIP-8, remettant tous les registres et la mémoire à zéro
func (c *Chip8) Reset() {
	c.Pc = 0x200
	c.DelayTimer = 0
	c.SoundTimer = 0
	c.I = 0
	c.Sp = 0
	for i := 0; i < len(c.Memory); i++ {
		c.Memory[i] = 0
	}

	for i := 0; i < len(c.Stack); i++ {
		c.Stack[i] = 0
	}

	for i := 0; i < len(c.V); i++ {
		c.V[i] = 0
	}

	for i := 0; i < len(c.Keys); i++ {
		c.Keys[i] = 0
	}
	c.LoadFontSet()
	c.ClearDisplay()
}

//La fonction Run est appelée à chaque rafraîchissement de la fenêtre de jeu.
//Elle exécute un cycle de la machine CHIP-8 (RunCpuCycle) et gère les timers de délai et de son.

func (c *Chip8) Run() {
	c.RunCpuCycle()

	if c.DelayTimer > 0 {
		c.DelayTimer = c.DelayTimer - 1
	}

	if c.SoundTimer > 0 {
		c.SoundTimer = c.SoundTimer - 1
	}
}

//La fonction RunCpuCycle est l'endroit où toutes les instructions de la CHIP-8 sont interprétées et exécutées.
//Elle lit un opcode depuis la mémoire, décode l'instruction et effectue les actions appropriées.

func (c *Chip8) RunCpuCycle() {
	opcode := uint16(c.Memory[c.Pc])<<8 | uint16(c.Memory[c.Pc+1])
	c.Pc = c.Pc + 2
	switch opcode & 0xF000 {

	case 0x0000:
		switch opcode & 0x000F {
		case 0x0000:
			c.ClearDisplay()
		case 0x000E:
			c.Pc = c.Stack[c.Sp-1]
			c.Sp = c.Sp - 1
		}

	case 0x1000:
		c.Pc = opcode & 0x0FFF

	case 0x2000:
		c.Stack[c.Sp] = c.Pc
		c.Sp = c.Sp + 1
		c.Pc = opcode & 0x0FFF

	case 0x3000:
		compareTo := byte(opcode & 0x00FF)
		register := (opcode & 0x0F00) >> 8
		if c.V[register] == compareTo {
			c.Pc = c.Pc + 2
		}

	case 0x4000:
		compareTo := byte(opcode & 0x00FF)
		register := (opcode & 0x0F00) >> 8
		if c.V[register] != compareTo {
			c.Pc = c.Pc + 2
		}

	case 0x5000:
		registerX := (opcode & 0x0F00) >> 8
		registerY := (opcode & 0x00F0) >> 4
		if c.V[registerX] == c.V[registerY] {
			c.Pc = c.Pc + 2
		}

	case 0x6000:
		register := byte((opcode & 0x0F00) >> 8)
		c.V[register] = byte(opcode & 0x00FF)

	case 0x7000:
		register := byte((opcode & 0x0F00) >> 8)
		value := byte(opcode & 0x00FF)
		c.V[register] = c.V[register] + value

	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerY]

		case 0x0001:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerX] | c.V[registerY]

		case 0x0002:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerX] & c.V[registerY]

		case 0x0003:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4
			c.V[registerX] = c.V[registerX] ^ c.V[registerY]

		case 0x0004:
			registerX := byte((opcode & 0x0F00) >> 8)
			registerY := byte((opcode & 0x00F0) >> 4)
			c.V[registerX] = c.V[registerX] + c.V[registerY]

			if uint16(c.V[registerX])+uint16(c.V[registerY]) > 0xFF {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}

		case 0x0005:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4

			if c.V[registerX] > c.V[registerY] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] - c.V[registerY]

		case 0x0006:
			registerX := (opcode & 0x0F00) >> 8

			if c.V[registerX]&0x1 == 1 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] >> 1

		case 0x0007:
			registerX := (opcode & 0x0F00) >> 8
			registerY := (opcode & 0x00F0) >> 4

			if c.V[registerY] > c.V[registerX] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerY] - c.V[registerX]

		case 0x000E:
			registerX := (opcode & 0x0F00) >> 8

			if c.V[registerX]&0x80 == 0x80 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[registerX] = c.V[registerX] << 1
		}

	case 0x9000:
		registerX := (opcode & 0x0F00) >> 8
		registerY := (opcode & 0x00F0) >> 4

		if c.V[registerX] != c.V[registerY] {
			c.Pc = c.Pc + 2
		}

	case 0xA000:
		c.I = (opcode & 0x0FFF)

	case 0xB000:
		c.Pc = (opcode & 0x0FFF) + uint16(c.V[0x0])

	case 0xC000:
		registerX := (opcode & 0x0F00) >> 8
		value := byte(opcode & 0x00FF)
		rand.Seed(time.Now().Unix())
		c.V[registerX] = byte(rand.Intn(256)) & value

	case 0xD000:
		registerX := (opcode & 0x0F00) >> 8
		registerY := (opcode & 0x00F0) >> 4
		nibble := byte(opcode & 0x000F)
		x := c.V[registerX]
		y := c.V[registerY]
		c.V[0xF] = 0x00

		for i := y; i < y+nibble; i++ {
			for j := x; j < x+8; j++ {
				bit := (c.Memory[c.I+uint16(i-y)] >> (7 - j + x)) & 0x01
				xIndex, yIndex := j, i
				if j >= width {
					xIndex = j - width
				}

				if i >= height {
					yIndex = i - height
				}

				if bit == 0x01 && c.Display[yIndex][xIndex] == 0x01 {
					c.V[0xF] = 0x01
				}
				c.Display[yIndex][xIndex] = c.Display[yIndex][xIndex] ^ bit
			}
		}
		c.Draw = true

	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			register := (opcode & 0x0F00) >> 8

			if c.Keys[c.V[register]] == 0x01 {
				c.Pc = c.Pc + 2
			}

		case 0x00A1:
			register := (opcode & 0x0F00) >> 8

			if c.Keys[c.V[register]] == 0x00 {
				c.Pc = c.Pc + 2
			}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x007:
			register := (opcode & 0x0F00) >> 8
			c.V[register] = c.DelayTimer

		case 0x0015:
			register := (opcode & 0x0F00) >> 8
			c.DelayTimer = c.V[register]

		case 0x0018:
			register := (opcode & 0x0F00) >> 8
			c.SoundTimer = c.V[register]

		case 0x000A:
			register := (opcode & 0x0F00) >> 8
			c.Inputflag = true
			c.InputRegister = byte(register)

		case 0x001E:
			register := (opcode & 0x0F00) >> 8
			c.I = c.I + uint16(c.V[register])

		case 0x0029:
			register := (opcode & 0x0F00) >> 8
			c.I = uint16(c.V[register] * 0x5)

		case 0x0033:
			register := (opcode & 0x0F00) >> 8
			number := c.V[register]
			c.Memory[c.I] = (number / 100) % 10
			c.Memory[c.I+1] = (number / 10) % 10
			c.Memory[c.I+2] = number % 10

		case 0x0055:
			register := (opcode & 0x0F00) >> 8

			for i := uint16(0x00); i <= register; i++ {
				c.Memory[c.I+i] = c.V[i]
			}

		case 0x0065:
			register := (opcode & 0x0F00) >> 8

			for i := uint16(0x00); i <= register; i++ {
				c.V[i] = c.Memory[c.I+i]
			}
		}
	}
}
