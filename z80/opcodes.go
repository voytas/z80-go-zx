package z80

const (
	useHL byte = 0
)

// Primary opcodes
const (
	nop        byte = 0x00 // nop
	ld_bc_nn   byte = 0x01 // ld bc,nn
	ld_bc_a    byte = 0x02 // ld (bc),a
	inc_bc     byte = 0x03 // inc bc
	inc_b      byte = 0x04 // inc b
	dec_b      byte = 0x05 // dec b
	ld_b_n     byte = 0x06 // ld b,n
	rlca       byte = 0x07 // rlca
	ex_af_af   byte = 0x08 // ex af,af'
	add_hl_bc  byte = 0x09 // add hl,bc
	ld_a_bc    byte = 0x0A // ld a,(bc)
	dec_bc     byte = 0x0B // dec bc
	inc_c      byte = 0x0C // inc c
	dec_c      byte = 0x0D // dec c
	ld_c_n     byte = 0x0E // ld c,n
	rrca       byte = 0x0F // rrca
	djnz       byte = 0x10 // djnz o
	ld_de_nn   byte = 0x11 // ld de,nn
	ld_de_a    byte = 0x12 // ld (de),a
	inc_de     byte = 0x13 // inc de
	inc_d      byte = 0x14 // inc d
	dec_d      byte = 0x15 // dec d
	ld_d_n     byte = 0x16 // ld d,n
	rla        byte = 0x17 // rla
	jr_o       byte = 0x18 // JR o
	add_hl_de  byte = 0x19 // add hl,de
	ld_a_de    byte = 0x1A // ld a,(de)
	dec_de     byte = 0x1B // dec de
	inc_e      byte = 0x1C // inc e
	dec_e      byte = 0x1D // dec e
	ld_e_n     byte = 0x1E // ld e,n
	rra        byte = 0x1F // rra
	jr_nz_o    byte = 0x20 // jr nz,o
	ld_hl_nn   byte = 0x21 // ld hl,nn
	ld_mm_hl   byte = 0x22 // ld (nn),hl
	inc_hl     byte = 0x23 // inc hl
	inc_h      byte = 0x24 // inc h
	dec_h      byte = 0x25 // dec h
	ld_h_n     byte = 0x26 // ld h,n
	daa        byte = 0x27 // daa
	jr_z_o     byte = 0x28 // jr z,o
	add_hl_hl  byte = 0x29 // add hl,hl
	ld_hl_mm   byte = 0x2A // ld hl,(nn)
	dec_hl     byte = 0x2B // dec hl
	inc_l      byte = 0x2C // inc l
	dec_l      byte = 0x2D // dec l
	ld_l_n     byte = 0x2E // ld l,n
	cpl        byte = 0x2F // cpl
	jr_nc_o    byte = 0x30 // jr nc,o
	ld_sp_nn   byte = 0x31 // ld sp,nn
	ld_mm_a    byte = 0x32 // ld (nn),A
	inc_sp     byte = 0x33 // inc sp
	inc_mhl    byte = 0x34 // inc (hl)
	dec_mhl    byte = 0x35 // dec (hl)
	ld_mhl_n   byte = 0x36 // ld (hl),n
	scf        byte = 0x37 // scf
	jr_c       byte = 0x38 // jr c,o
	add_hl_sp  byte = 0x39 // add hl,sp
	ld_a_mm    byte = 0x3A // ld a,(nn)
	dec_sp     byte = 0x3B // dec sp
	inc_a      byte = 0x3C // inc a
	dec_a      byte = 0x3D // dec a
	ld_a_n     byte = 0x3E // ld a,n
	ccf        byte = 0x3F // ccf
	ld_b_b     byte = 0x40 // ld b,b
	ld_b_c     byte = 0x41 // ld b,c
	ld_b_d     byte = 0x42 // ld b,d
	ld_b_e     byte = 0x43 // ld b,e
	ld_b_h     byte = 0x44 // ld b,h
	ld_b_l     byte = 0x45 // ld b,l
	ld_b_hl    byte = 0x46 // ld b,(hl)
	ld_b_a     byte = 0x47 // ld b,a
	ld_c_b     byte = 0x48 // ld c,b
	ld_c_c     byte = 0x49 // ld c,c
	ld_c_d     byte = 0x4A // ld c,d
	ld_c_e     byte = 0x4B // ld c,e
	ld_c_h     byte = 0x4C // ld c,h
	ld_c_l     byte = 0x4D // ld c,l
	ld_c_hl    byte = 0x4E // ld c,(hl)
	ld_c_a     byte = 0x4F // ld c,a
	ld_d_b     byte = 0x50 // ld d,b
	ld_d_c     byte = 0x51 // ld d,c
	ld_d_d     byte = 0x52 // ld d,d
	ld_d_e     byte = 0x53 // ld d,e
	ld_d_h     byte = 0x54 // ld d,h
	ld_d_l     byte = 0x55 // ld d,l
	ld_d_hl    byte = 0x56 // ld d,(hl)
	ld_d_a     byte = 0x57 // ld d,a
	ld_e_b     byte = 0x58 // ld e,b
	ld_e_c     byte = 0x59 // ld e,c
	ld_e_d     byte = 0x5A // ld e,d
	ld_e_e     byte = 0x5B // ld e,e
	ld_e_h     byte = 0x5C // ld e,h
	ld_e_l     byte = 0x5D // ld e,l
	ld_e_hl    byte = 0x5E // ld e,(hl)
	ld_e_a     byte = 0x5F // ld e,a
	ld_h_b     byte = 0x60 // ld h,b
	ld_h_c     byte = 0x61 // ld h,c
	ld_h_d     byte = 0x62 // ld h,d
	ld_h_e     byte = 0x63 // ld h,e
	ld_h_h     byte = 0x64 // ld h,h
	ld_h_l     byte = 0x65 // ld h,l
	ld_h_hl    byte = 0x66 // ld h,(hl)
	ld_h_a     byte = 0x67 // ld h,a
	ld_l_b     byte = 0x68 // ld l,b
	ld_l_c     byte = 0x69 // ld l,c
	ld_l_d     byte = 0x6A // ld l,d
	ld_l_e     byte = 0x6B // ld l,e
	ld_l_h     byte = 0x6C // ld l,h
	ld_l_l     byte = 0x6D // ld l,l
	ld_l_hl    byte = 0x6E // ld l,(hl)
	ld_l_a     byte = 0x6F // ld l,a
	ld_hl_b    byte = 0x70 // ld (hl),b
	ld_hl_c    byte = 0x71 // ld (hl),c
	ld_hl_d    byte = 0x72 // ld (hl),d
	ld_hl_e    byte = 0x73 // ld (hl),e
	ld_hl_h    byte = 0x74 // ld (hl),h
	ld_hl_l    byte = 0x75 // ld (hl),l
	halt       byte = 0x76 // halt
	ld_hl_a    byte = 0x77 // ld (hl),a
	ld_a_b     byte = 0x78 // ld a,b
	ld_a_c     byte = 0x79 // ld a,c
	ld_a_d     byte = 0x7A // ld a,d
	ld_a_e     byte = 0x7B // ld a,e
	ld_a_h     byte = 0x7C // ld a,h
	ld_a_l     byte = 0x7D // ld a,l
	ld_a_hl    byte = 0x7E // ld a,(hl)
	ld_a_a     byte = 0x7F // ld a,a
	add_a_b    byte = 0x80 // add a,b
	add_a_c    byte = 0x81 // add a,c
	add_a_d    byte = 0x82 // add a,d
	add_a_e    byte = 0x83 // add a,e
	add_a_h    byte = 0x84 // add a,h
	add_a_l    byte = 0x85 // add a,l
	add_a_hl   byte = 0x86 // add a,(hl)
	add_a_a    byte = 0x87 // add a,a
	adc_a_b    byte = 0x88 // adc a,b
	adc_a_c    byte = 0x89 // adc a,c
	adc_a_d    byte = 0x8A // adc a,d
	adc_a_e    byte = 0x8B // adc a,e
	adc_a_h    byte = 0x8C // adc a,h
	adc_a_l    byte = 0x8D // adc a,l
	adc_a_hl   byte = 0x8E // adc a,(hl)
	adc_a_a    byte = 0x8F // adc a,a
	sub_b      byte = 0x90 // sub b
	sub_c      byte = 0x91 // sub c
	sub_d      byte = 0x92 // sub d
	sub_e      byte = 0x93 // sub e
	sub_h      byte = 0x94 // sub h
	sub_l      byte = 0x95 // sub l
	sub_hl     byte = 0x96 // sub (hl)
	sub_a      byte = 0x97 // sub a
	sbc_a_b    byte = 0x98 // sbc a,b
	sbc_a_c    byte = 0x99 // sbc a,c
	sbc_a_d    byte = 0x9A // sbc a,d
	sbc_a_e    byte = 0x9B // sbc a,e
	sbc_a_h    byte = 0x9C // sbc a,h
	sbc_a_l    byte = 0x9D // sbc a,l
	sbc_a_hl   byte = 0x9E // sbc a,(hl)
	sbc_a_a    byte = 0x9F // sbc a,a
	and_b      byte = 0xA0 // and b
	and_c      byte = 0xA1 // and c
	and_d      byte = 0xA2 // and d
	and_e      byte = 0xA3 // and e
	and_h      byte = 0xA4 // and h
	and_l      byte = 0xA5 // and l
	and_hl     byte = 0xA6 // and (hl)
	and_a      byte = 0xA7 // and a
	xor_b      byte = 0xA8 // xor b
	xor_c      byte = 0xA9 // xor c
	xor_d      byte = 0xAA // xor d
	xor_e      byte = 0xAB // xor e
	xor_h      byte = 0xAC // xor h
	xor_l      byte = 0xAD // xor l
	xor_hl     byte = 0xAE // xor (hl)
	xor_a      byte = 0xAF // xor a
	or_b       byte = 0xB0 // or b
	or_c       byte = 0xB1 // or c
	or_d       byte = 0xB2 // or d
	or_e       byte = 0xB3 // or e
	or_h       byte = 0xB4 // or h
	or_l       byte = 0xB5 // or l
	or_hl      byte = 0xB6 // or (hl)
	or_a       byte = 0xB7 // or a
	cp_b       byte = 0xB8 // cp b
	cp_c       byte = 0xB9 // cp c
	cp_d       byte = 0xBA // cp d
	cp_e       byte = 0xBB // cp e
	cp_h       byte = 0xBC // cp h
	cp_l       byte = 0xBD // cp l
	cp_hl      byte = 0xBE // cp (hl)
	cp_a       byte = 0xBF // cp a
	ret_nz     byte = 0xC0 // ret nz
	pop_bc     byte = 0xC1 // pop bc
	jp_nz_nn   byte = 0xC2 // jp nz,nn
	jp_nn      byte = 0xC3 // jp nn
	call_nz_nn byte = 0xC4 // call nz,nn
	push_bc    byte = 0xC5 // push bc
	add_a_n    byte = 0xC6 // add a.n
	rst_00h    byte = 0xC7 // rst 00h
	ret_z      byte = 0xC8 // ret z
	ret        byte = 0xC9 // ret
	jp_z_nn    byte = 0xCA // jp z,nn
	prefix_cb  byte = 0xCB // bit operations etc prefix
	call_z_nn  byte = 0xCC // call z,nn
	call_nn    byte = 0xCD // call nn
	adc_a_n    byte = 0xCE // adc a,n
	rst_08h    byte = 0xCF // rst 08h
	ret_nc     byte = 0xD0 // ret nc
	pop_de     byte = 0xD1 // pop de
	jp_nc_nn   byte = 0xD2 // jp nc,nn
	out_n_a    byte = 0xD3 // out (n),a
	call_nc_nn byte = 0xD4 // call nc,nn
	push_de    byte = 0xD5 // push de
	sub_n      byte = 0xD6 // sub n
	rst_10h    byte = 0xD7 // rst 10h
	ret_c      byte = 0xD8 // ret c
	exx        byte = 0xD9 // exx
	jp_c_nn    byte = 0xDA // jp c,nn
	in_a_n     byte = 0xDB // in a,(n)
	call_c_nn  byte = 0xDC // call c,nn
	useIX      byte = 0xDD // IX instruction
	sbc_a_n    byte = 0xDE // sbc a,n
	rst_18h    byte = 0xDF // rst 18h
	ret_po     byte = 0xE0 // ret po
	pop_hl     byte = 0xE1 // pop hl
	jp_po_nn   byte = 0xE2 // jp po,nn
	ex_sp_hl   byte = 0xE3 // ex (sp),hl
	call_po_nn byte = 0xE4 // call po,nn
	push_hl    byte = 0xE5 // push hl
	and_n      byte = 0xE6 // and n
	rst_20h    byte = 0xE7 // rst 20h
	ret_pe     byte = 0xE8 // ret pe
	jp_hl      byte = 0xE9 // jp (hl)
	jp_pe_nn   byte = 0xEA // jp pe,nn
	ex_de_hl   byte = 0xEB // ex de,hl
	call_pe_nn byte = 0xEC // call pe,nn
	prefix_ed  byte = 0xED // ED prefix
	xor_n      byte = 0xEE // xor n
	rst_28h    byte = 0xEF // rst 28h
	ret_p      byte = 0xF0 // ret p
	pop_af     byte = 0xF1 // pop af
	jp_p_nn    byte = 0xF2 // jp p,nn
	di         byte = 0xF3 // di
	call_p_nn  byte = 0xF4 // call p,nn
	push_af    byte = 0xF5 // push af
	or_n       byte = 0xF6 // or n
	rst_30h    byte = 0xF7 // rst 30h
	ret_m      byte = 0xF8 // ret m
	ld_sp_hl   byte = 0xF9 // ld sp,hl
	jp_m_nn    byte = 0xFA // jp m,nn
	ei         byte = 0xFB // ei
	call_m_nn  byte = 0xFC // call m,nn
	useIY      byte = 0xFD // IY instruction
	cp_n       byte = 0xFE // cp n
	rst_38h    byte = 0xFF // rst 38h
)

// CB prefixed opcodes
const (
	rlc_r byte = 0b00000000
	rrc_r byte = 0b00001000
	rl_r  byte = 0b00010000
	rr_r  byte = 0b00011000
	sla_r byte = 0b00100000
	sra_r byte = 0b00101000
	sll_r byte = 0b00110000
	srl_r byte = 0b00111000
	bit_b byte = 0b01000000
	res_b byte = 0b10000000
	set_b byte = 0b11000000

	bit_0 byte = 0b00000000
	bit_1 byte = 0b00001000
	bit_2 byte = 0b00010000
	bit_3 byte = 0b00011000
	bit_4 byte = 0b00100000
	bit_5 byte = 0b00101000
	bit_6 byte = 0b00110000
	bit_7 byte = 0b00111000
)

// ED prefixed opcodes
const (
	in_b_c    byte = 0x40 // in b,(c)
	out_c_b   byte = 0x41 // out (c),b
	sbc_hl_bc byte = 0x42 // sbc hl,bc
	ld_mm_bc  byte = 0x43 // ld (nn),bc
	neg       byte = 0x44 // neg
	retn      byte = 0x45 // retn
	im0       byte = 0x46 // im 0
	ld_i_a    byte = 0x47 // ld i,a
	in_c_c    byte = 0x48 // in c,(c)
	out_c_c   byte = 0x49 // out (c),c
	adc_hl_bc byte = 0x4A // adc hl,bc
	ld_bc_mm  byte = 0x4B // ld bc,(nn)
	reti      byte = 0x4D // reti
	ld_r_a    byte = 0x4F // ld r,a
	in_d_c    byte = 0x50 // in d,(c)
	out_c_d   byte = 0x51 // out (c),d
	sbc_hl_de byte = 0x52 // sbc hl,de
	ld_mm_de  byte = 0x53 // ld (nn),de
	im1       byte = 0x56 // im 1
	ld_a_i    byte = 0x57 // ld a,i
	in_e_c    byte = 0x58 // in e,(c)
	out_c_e   byte = 0x59 // out (c),e
	adc_hl_de byte = 0x5A // adc hl,de
	ld_de_mm  byte = 0x5B // ld de,(nn)
	im2       byte = 0x5E // im 2
	ld_a_r    byte = 0x5F // ld a,r
	in_h_c    byte = 0x60 // in h,(c)
	out_c_h   byte = 0x61 // out (c),h
	sbc_hl_hl byte = 0x62 // sbc hl,hl
	ld_mm_hl2 byte = 0x63 // ld (nn),hl
	rrd       byte = 0x67 // rrd
	in_l_c    byte = 0x68 // in l,(c)
	out_c_l   byte = 0x69 // out (c),l
	adc_hl_hl byte = 0x6A // adc hl,hl
	ld_hl_mm2 byte = 0x6B // ld hl,(nn)
	rld       byte = 0x6F // rld
	in_f_c    byte = 0x70 // in f,(c)
	out_c_f   byte = 0x71 // out (c),f
	sbc_hl_sp byte = 0x72 // sbc hl,sp
	ld_mm_sp  byte = 0x73 // ld (nn),sp
	in_a_c    byte = 0x78 // in a,(c)
	out_c_a   byte = 0x79 // out (c),a
	adc_hl_sp byte = 0x7A // adc hl,sp
	ld_sp_mm  byte = 0x7B // ld sp,(nn)
	ldi       byte = 0xA0 // ldi
	cpi       byte = 0xA1 // cpi
	ini       byte = 0xA2 // ini
	outi      byte = 0xA3 // outi
	ldd       byte = 0xA8 // ldd
	cpd       byte = 0xA9 // cpd
	ind       byte = 0xAA // ind
	outd      byte = 0xAB // outd
	ldir      byte = 0xB0 // ldir
	cpir      byte = 0xB1 // cpir
	inir      byte = 0xB2 // inir
	otir      byte = 0xB3 // otir
	lddr      byte = 0xB8 // lddr
	cpdr      byte = 0xB9 // cpdr
	indr      byte = 0xBA // indr
	otdr      byte = 0xBB // otdr
)
