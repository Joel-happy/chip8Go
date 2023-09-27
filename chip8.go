package main

import (
   "fmt"
)

type Chip8 struct {
    memory        [4096]byte
    registers     [16]uint8
    indexRegister uint16
    programCounter uint16
    soundTimer    uint8
    delayTimer    uint8
    stack         [16]uint16
    keypad        [16]uint8
    screen        [64][32]uint8
    timerSpeed    uint8
	op  uint16

}
func (chip8 *Chip8) loadProgram(program []byte, startAddress uint16) {
    // le programme ne dépasse pas la taille de la mémoire
    if int(startAddress)+len(program) > len(chip8.memory) {
        panic("program too big for CHIP-8")
    }

    // Copiez les opcodes dans la mémoire à l'adresse de départ
    for i, opcode := range program {
        chip8.memory[startAddress+uint16(i)] = opcode
    }
}

func (chip8 *Chip8) decodeOpcode(opcode uint16) {
    switch opcode & 0xF000 {
    case 0x0000:
		switch opcode & 0x000F{
		case 0x0000:
			  // Efface l'écran en mettant toutes les valeurs à zéro
			for i := 0; i < len(chip8.screen); i++ {
                for j := 0; j < len(chip8.screen[i]); j++ {
                    chip8.screen[i][j] = 0
				}
			} 
		case 0x000E:
			 // Récupérer l'adresse en haut de la pile
			 returnAddress := chip8.stack[chip8.sp]

			 // Décrémenter le pointeur de pile (stack pointer)
			 chip8.sp--
 
			 // Régler le compteur de programme (program counter) sur l'adresse récupérée
			 chip8.pc = returnAddress
		 }
	case 0x1000:
		
	case 0x2000:


}
}