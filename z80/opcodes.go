package z80

const (
	NOP       = 0x00 // nop
	LD_BC_nn  = 0x01 // ld bc,nn
	LD_BC_A   = 0x02 // ld (bc),a
	INC_BC    = 0x03 // inc bc
	INC_B     = 0x04 // inc b
	DEC_B     = 0x05 // dec b
	LD_B_n    = 0x06 // ld b,n
	RLCA      = 0x07 // rlca
	EX_AF_AF  = 0x08 // ex af,af'
	ADD_HL_BC = 0x09 // add hl,bc
	LD_A_BC   = 0x0A // ld a,(bc)
	DEC_BC    = 0x0B // dec bc
	INC_C     = 0x0C // inc c
	DEC_C     = 0x0D // dec c
	LD_C_n    = 0x0E // ld c,n
	RRCA      = 0x0F // rrca
	DJNZ      = 0x10 // djnz o
	LD_DE_nn  = 0x11 // ld de,nn
	INC_D     = 0x14 // inc d
	LD_D_n    = 0x16 // ld d,n
	JR        = 0x18 // JR n
	ADD_HL_DE = 0x19 // add hl,de
	DEC_DE    = 0x1B // dec de
	INC_E     = 0x1C // inc e
	LD_E_n    = 0x1E // ld e,n
	LD_HL_nn  = 0x21 // ld hl,nn
	INC_H     = 0x24 // inc h
	LD_H_n    = 0x26 // ld h,n
	ADD_HL_HL = 0x29 // add hl,hl
	DEC_HL    = 0x2B // dec hl
	INC_L     = 0x2C // inc l
	LD_L_n    = 0x2E // ld l,n
	LD_SP_nn  = 0x31 // ld sp,nn
	ADD_HL_SP = 0x39 // add hl,sp
	DEC_SP    = 0x3B // dec sp
	INC_A     = 0x3C // inc a
	LD_A_n    = 0x3E // ld a,n
	HALT      = 0x76 // halt
)
