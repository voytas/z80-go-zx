package z80

const (
	NOP       byte = 0x00 // nop
	LD_BC_nn  byte = 0x01 // ld bc,nn
	LD_BC_A   byte = 0x02 // ld (bc),a
	INC_BC    byte = 0x03 // inc bc
	INC_B     byte = 0x04 // inc b
	DEC_B     byte = 0x05 // dec b
	LD_B_n    byte = 0x06 // ld b,n
	RLCA      byte = 0x07 // rlca
	EX_AF_AF  byte = 0x08 // ex af,af'
	ADD_HL_BC byte = 0x09 // add hl,bc
	LD_A_BC   byte = 0x0A // ld a,(bc)
	DEC_BC    byte = 0x0B // dec bc
	INC_C     byte = 0x0C // inc c
	DEC_C     byte = 0x0D // dec c
	LD_C_n    byte = 0x0E // ld c,n
	RRCA      byte = 0x0F // rrca
	DJNZ      byte = 0x10 // djnz o
	LD_DE_nn  byte = 0x11 // ld de,nn
	LD_DE_A   byte = 0x12 // ld (de),a
	INC_DE    byte = 0x13 // inc de
	INC_D     byte = 0x14 // inc d
	DEC_D     byte = 0x15 // dec d
	LD_D_n    byte = 0x16 // ld d,n
	RLA       byte = 0x17 // rla
	JR        byte = 0x18 // JR n
	ADD_HL_DE byte = 0x19 // add hl,de
	LD_A_DE   byte = 0x1A // ld a,(de)
	DEC_DE    byte = 0x1B // dec de
	INC_E     byte = 0x1C // inc e
	DEC_E     byte = 0x1D // dec e
	LD_E_n    byte = 0x1E // ld e,n
	RRA       byte = 0x1F // rra
	JR_NZ     byte = 0x20 // jr nz,o
	LD_HL_nn  byte = 0x21 // ld hl,nn
	LD_mm_HL  byte = 0x22 // ld (nn),hl
	INC_HL    byte = 0x23 // inc hl
	INC_H     byte = 0x24 // inc h
	DEC_H     byte = 0x25 // dec h
	LD_H_n    byte = 0x26 // ld h,n
	DAA       byte = 0x27 // daa
	JR_Z      byte = 0x28 // jr z,o
	ADD_HL_HL byte = 0x29 // add hl,hl
	LD_HL_mm  byte = 0x2A // ld hl,(nn)
	DEC_HL    byte = 0x2B // dec hl
	INC_L     byte = 0x2C // inc l
	DEC_L     byte = 0x2D // dec l
	LD_L_n    byte = 0x2E // ld l,n
	CPL       byte = 0x2F // cpl
	JR_NC     byte = 0x30 // jr nc,o
	LD_SP_nn  byte = 0x31 // ld sp,nn
	LD_mm_A   byte = 0x32 // ld (nn),A
	INC_SP    byte = 0x33 // inc sp
	INC_mHL   byte = 0x34 // inc (hl)
	DEC_mHL   byte = 0x35 // dec (hl)
	LD_mHL_n  byte = 0x36 // ld (hl),n
	SCF       byte = 0x37 // scf
	JR_C      byte = 0x38 // jr c, o
	ADD_HL_SP byte = 0x39 // add hl,sp
	LD_A_mm   byte = 0x3A // ld a,(nn)
	DEC_SP    byte = 0x3B // dec sp
	INC_A     byte = 0x3C // inc a
	DEC_A     byte = 0x3D // dec a
	LD_A_n    byte = 0x3E // ld a,n
	CCF       byte = 0x3F // ccf
	LD_B_B    byte = 0x40 // ld b,b
	LD_B_C    byte = 0x41 // ld b,c
	LD_B_D    byte = 0x42 // ld b,d
	LD_B_E    byte = 0x43 // ld b,e
	LD_B_H    byte = 0x44 // ld b,h
	LD_B_L    byte = 0x45 // ld b,l
	LD_B_A    byte = 0x47 // ld b,a
	LD_C_B    byte = 0x48 // ld c,b
	LD_C_C    byte = 0x49 // ld c,c
	LD_C_D    byte = 0x4A // ld c,d
	LD_C_E    byte = 0x4B // ld c,e
	LD_C_H    byte = 0x4C // ld c,h
	LD_C_L    byte = 0x4D // ld c,l
	LD_C_A    byte = 0x4F // ld c,a
	LD_D_B    byte = 0x50 // ld d,b
	LD_D_C    byte = 0x51 // ld d,c
	LD_D_D    byte = 0x52 // ld d,d
	LD_D_E    byte = 0x53 // ld d,e
	LD_D_H    byte = 0x54 // ld d,h
	LD_D_L    byte = 0x55 // ld d,l
	LD_D_A    byte = 0x57 // ld d,a
	LD_E_B    byte = 0x58 // ld e,b
	LD_E_C    byte = 0x59 // ld e,c
	LD_E_D    byte = 0x5A // ld e,d
	LD_E_E    byte = 0x5B // ld e,e
	LD_E_H    byte = 0x5C // ld e,h
	LD_E_L    byte = 0x5D // ld e,l
	LD_E_A    byte = 0x5F // ld e,a
	LD_H_B    byte = 0x60 // ld h,b
	LD_H_C    byte = 0x61 // ld h,c
	LD_H_D    byte = 0x62 // ld h,d
	LD_H_E    byte = 0x63 // ld h,e
	LD_H_H    byte = 0x64 // ld h,h
	LD_H_L    byte = 0x65 // ld h,l
	LD_H_A    byte = 0x67 // ld h,a
	LD_L_B    byte = 0x68 // ld l,b
	LD_L_C    byte = 0x69 // ld l,c
	LD_L_D    byte = 0x6A // ld l,d
	LD_L_E    byte = 0x6B // ld l,e
	LD_L_H    byte = 0x6C // ld l,h
	LD_L_L    byte = 0x6D // ld l,l
	LD_L_A    byte = 0x6F // ld l,a
	HALT      byte = 0x76 // halt
	LD_A_B    byte = 0x78 // ld a,b
	LD_A_C    byte = 0x79 // ld a,c
	LD_A_D    byte = 0x7A // ld a,d
	LD_A_E    byte = 0x7B // ld a,e
	LD_A_H    byte = 0x7C // ld a,h
	LD_A_L    byte = 0x7D // ld a,l
	LD_A_A    byte = 0x7F // ld a,a
	SUB_B     byte = 0x90 // sub b
	SUB_C     byte = 0x91 // sub c
	SUB_D     byte = 0x92 // sub d
	SUB_E     byte = 0x93 // sub e
	SUB_H     byte = 0x94 // sub h
	SUB_L     byte = 0x95 // sub l
	SUB_A     byte = 0x97 // sub a
	ADD_A_n   byte = 0xC6 // add a.n
)