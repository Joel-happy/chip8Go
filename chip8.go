package main

import (
    "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
func main{
	switch opcode{
	case opcode[0]==0 :
			if opcode[1]!=0 {
				if 
			}
				
		}
	}
	



