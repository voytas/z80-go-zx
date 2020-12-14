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
	LD_DE_A   = 0x12 // ld (de),a
	INC_DE    = 0x13 // inc de
	INC_D     = 0x14 // inc d
	DEC_D     = 0x15 // dec d
	LD_D_n    = 0x16 // ld d,n
	RLA       = 0x17 // rla
	JR        = 0x18 // JR n
	ADD_HL_DE = 0x19 // add hl,de
	LD_A_DE   = 0x1A // ld a,(de)
	DEC_DE    = 0x1B // dec de
	INC_E     = 0x1C // inc e
	DEC_E     = 0x1D // dec e
	LD_E_n    = 0x1E // ld e,n
	RRA       = 0x1F // rra
	JR_NZ     = 0x20 // jr nz,o
	LD_HL_nn  = 0x21 // ld hl,nn
	INC_HL    = 0x23 // inc hl
	INC_H     = 0x24 // inc h
	DEC_H     = 0x25 // dec h
	LD_H_n    = 0x26 // ld h,n
	JR_Z      = 0x28 // jr z,o
	ADD_HL_HL = 0x29 // add hl,hl
	DEC_HL    = 0x2B // dec hl
	INC_L     = 0x2C // inc l
	DEC_L     = 0x2D // dec l
	LD_L_n    = 0x2E // ld l,n
	JR_NC     = 0x30 // jr nc,o
	LD_SP_nn  = 0x31 // ld sp,nn
	INC_SP    = 0x33 // inc sp
	JR_C      = 0x38 // jr c, o
	ADD_HL_SP = 0x39 // add hl,sp
	DEC_SP    = 0x3B // dec sp
	INC_A     = 0x3C // inc a
	DEC_A     = 0x3D // dec a
	LD_A_n    = 0x3E // ld a,n
	LD_B_B    = 0x40 // ld b,b
	LD_B_C    = 0x41 // ld b,c
	LD_B_D    = 0x42 // ld b,d
	LD_B_E    = 0x43 // ld b,e
	LD_B_H    = 0x44 // ld b,h
	LD_B_L    = 0x45 // ld b,l
	LD_B_A    = 0x47 // ld b,a
	LD_C_B    = 0x48 // ld c,b
	LD_C_C    = 0x49 // ld c,c
	LD_C_D    = 0x4A // ld c,d
	LD_C_E    = 0x4B // ld c,e
	LD_C_H    = 0x4C // ld c,h
	LD_C_L    = 0x4D // ld c,l
	LD_C_A    = 0x4F // ld c,a
	LD_D_B    = 0x50 // ld d,b
	LD_D_C    = 0x51 // ld d,c
	LD_D_D    = 0x52 // ld d,d
	LD_D_E    = 0x53 // ld d,e
	LD_D_H    = 0x54 // ld d,h
	LD_D_L    = 0x55 // ld d,l
	LD_D_A    = 0x57 // ld d,a
	LD_E_B    = 0x58 // ld e,b
	LD_E_C    = 0x59 // ld e,c
	LD_E_D    = 0x5A // ld e,d
	LD_E_E    = 0x5B // ld e,e
	LD_E_H    = 0x5C // ld e,h
	LD_E_L    = 0x5D // ld e,l
	LD_E_A    = 0x5F // ld e,a
	LD_H_B    = 0x60 // ld h,b
	LD_H_C    = 0x61 // ld h,c
	LD_H_D    = 0x62 // ld h,d
	LD_H_E    = 0x63 // ld h,e
	LD_H_H    = 0x64 // ld h,h
	LD_H_L    = 0x65 // ld h,l
	LD_H_A    = 0x67 // ld h,a
	LD_L_B    = 0x68 // ld l,b
	LD_L_C    = 0x69 // ld l,c
	LD_L_D    = 0x6A // ld l,d
	LD_L_E    = 0x6B // ld l,e
	LD_L_H    = 0x6C // ld l,h
	LD_L_L    = 0x6D // ld l,l
	LD_L_A    = 0x6F // ld l,a
	HALT      = 0x76 // halt
	LD_A_B    = 0x78 // ld a,b
	LD_A_C    = 0x79 // ld a,c
	LD_A_D    = 0x7A // ld a,d
	LD_A_E    = 0x7B // ld a,e
	LD_A_H    = 0x7C // ld a,h
	LD_A_L    = 0x7D // ld a,l
	LD_A_A    = 0x7F // ld a,a
	ADD_A_n   = 0xC6 // add a.n
)
