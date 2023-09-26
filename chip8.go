package main

type Chip8 struct {
	Registers       [16]uint8
	registerAddress [1]uint16
	registerProgram [1]uint8
	registerSound   [1]uint8
	registerDelay   [1]uint8
	stack           [16]uint8 // pile de 16 registres //
	keyboard        [16]bool  // clavier de 16 touche //

}
