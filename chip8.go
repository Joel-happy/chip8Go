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
    pc    uint16
    sp uint8

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
        // Jump to location nnn.
        address := opcode & 0x0FFF
        chip8.pc = address
		
	case 0x2000:
        address := opcode & 0x0FFF

        // Incrémentez le pointeur de pile (stack pointer).
        chip8.sp++

        // Placez l'adresse actuelle du PC (avant l'appel de sous-routine) en haut de la pile.
        chip8.stack[chip8.sp] = chip8.pc

        // Mettez à jour le PC avec l'adresse de la sous-routine.
        chip8.pc = address

    case 0x3000:
        // Skip next instruction if Vx == kk.
        x := (opcode >> 8) & 0x0F  // Extraire le numéro du registre Vx (les 4 bits du milieu).
        kk := uint8(opcode & 0x00FF)  // Extraire la valeur immédiate kk (les 8 bits de droite).

        if chip8.registers[x] == kk {
            // Si Vx est égal à kk, incrémentez le compteur de programme (pc) de 2 pour sauter l'instruction suivante.
            chip8.pc += 2
        }
    case 0x4000:
        // Skip next instruction if Vx != kk.
         x := int((opcode >> 8) & 0x0F) // Extraire le numéro du registre Vx (les 4 bits du milieu).
        kk := uint8(opcode & 0x00FF)   // Extraire la valeur immédiate kk (les 8 bits de droite).

        if chip8.registers[x] != kk {
        // Si Vx n'est pas égal à kk, incrémentez le compteur de programme (pc) de 2 pour sauter l'instruction suivante.
        chip8.pc += 2
        }
    case 0x5000:
        // Skip next instruction if Vx == Vy.
        x := int((opcode >> 8) & 0x0F) // Extraire le numéro du registre Vx (les 4 bits du milieu).
        y := int((opcode >> 4) & 0x0F) // Extraire le numéro du registre Vy (les 4 bits de droite).

        if chip8.registers[x] == chip8.registers[y] {
        // Si Vx est égal à Vy, incrémentez le compteur de programme (pc) de 2 pour sauter l'instruction suivante.
        chip8.pc += 2
        }
    case 0x6000:
        // Set Vx = kk.
        x := int((opcode >> 8) & 0x0F) // Extraire le numéro du registre Vx (les 4 bits du milieu).
        kk := uint8(opcode & 0x00FF)   // Extraire la valeur immédiate kk (les 8 bits de droite).
    
        chip8.registers[x] = kk // Mettre la valeur kk dans le registre Vx.
    case 0x7000:
        x := int((opcode >> 8) & 0x0F) // Extraire le numéro du registre Vx (les 4 bits du milieu).
        kk := uint8(opcode & 0x00FF)   // Extraire la valeur immédiate kk (les 8 bits de droite).
    
        chip8.registers[x] += kk // Ajouter la valeur kk à la valeur du registre Vx.
    
    case 0x8000:
        x := int((opcode >> 8) & 0x0F) // Numéro du registre Vx.
        y := int((opcode >> 4) & 0x0F) // Numéro du registre Vy.
    
        switch opcode & 0x000F {
        case 0x0000:
            // LD Vx, Vy : Définir Vx = Vy.
            chip8.registers[x] = chip8.registers[y]
    
        case 0x0001:
            // OR Vx, Vy : Définir Vx = Vx OR Vy.
            chip8.registers[x] |= chip8.registers[y]
    
        case 0x0002:
            // AND Vx, Vy : Définir Vx = Vx AND Vy.
            chip8.registers[x] &= chip8.registers[y]
    
        case 0x0003:
            // XOR Vx, Vy : Définir Vx = Vx XOR Vy.
            chip8.registers[x] ^= chip8.registers[y]
    
        case 0x0004:
            // ADD Vx, Vy : Définir Vx = Vx + Vy, set VF = carry.
            if chip8.registers[y] > 0xFF-chip8.registers[x] {
                chip8.registers[0xF] = 1 // La retenue dépasse 8 bits.
            } else {
                chip8.registers[0xF] = 0
            }
            chip8.registers[x] += chip8.registers[y]
    
        case 0x0005:
            // SUB Vx, Vy : Définir Vx = Vx - Vy, set VF = NOT borrow.
            if chip8.registers[x] > chip8.registers[y] {
                chip8.registers[0xF] = 1 // Pas d'emprunt.
            } else {
                chip8.registers[0xF] = 0
            }
            chip8.registers[x] -= chip8.registers[y]
    
        case 0x0006:
            // SHR Vx {, Vy} : Définir Vx = Vx SHR 1, set VF = bit le moins significatif.
            chip8.registers[0xF] = chip8.registers[x] & 0x1 // Bit le moins significatif avant le décalage.
            chip8.registers[x] >>= 1
    
        case 0x0007:
            // SUBN Vx, Vy : Définir Vx = Vy - Vx, set VF = NOT borrow.
            if chip8.registers[y] > chip8.registers[x] {
                chip8.registers[0xF] = 1 // Pas d'emprunt.
            } else {
                chip8.registers[0xF] = 0
            }
            chip8.registers[x] = chip8.registers[y] - chip8.registers[x]
    
        case 0x000E:
            // SHL Vx {, Vy} : Définir Vx = Vx SHL 1, set VF = bit le plus significatif.
            chip8.registers[0xF] = (chip8.registers[x] >> 7) & 0x1 // Bit le plus significatif avant le décalage.
            chip8.registers[x] <<= 1
        }
    case 0x9000:
        

}
}