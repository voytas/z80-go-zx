package z80

// Handles opcodes with ED prefix
func (z80 *Z80) prefixED(opcode byte) {
	switch opcode {
	case neg, 0x54, 0x64, 0x74, 0x4C, 0x5C, 0x6C, 0x7C:
		a := z80.Reg.A
		z80.Reg.A = ^a + 1
		z80.Reg.F = (FS|FY|FX)&z80.Reg.A | FN | byte(a^z80.Reg.A)&FH
		if z80.Reg.A == 0 {
			z80.Reg.F |= FZ
		}
		if z80.Reg.A == 0x80 {
			z80.Reg.F |= FP
		}
		if a != 0 {
			z80.Reg.F |= FC
		}
	case adc_hl_bc, adc_hl_de, adc_hl_hl, adc_hl_sp:
		z80.addContention(z80.Reg.IR(), 7)
		hl := z80.Reg.HL()
		nn := z80.Reg.rr(opcode & 0b00110000 >> 4)
		sum := hl + nn + uint16(z80.Reg.F&FC)
		z80.Reg.F = byte((hl^nn^sum)>>8)&FH | byte(sum>>8)&(FY|FX)
		if sum > 0x7FFF {
			z80.Reg.F |= FS
		}
		if sum == 0 {
			z80.Reg.F |= FZ
		}
		if (hl^nn)&0x8000 == 0 && (hl^sum)&0x8000 != 0 {
			z80.Reg.F |= FP
		}
		if sum < hl {
			z80.Reg.F |= FC
		}

		z80.Reg.setHL(sum)
	case sbc_hl_bc, sbc_hl_de, sbc_hl_hl, sbc_hl_sp:
		z80.addContention(z80.Reg.IR(), 7)
		hl := z80.Reg.HL()
		nn := z80.Reg.rr(opcode & 0b00110000 >> 4)
		sub := hl - nn - uint16(z80.Reg.F&FC)
		z80.Reg.F = FN | byte((hl^nn^sub)>>8)&FH | byte(sub>>8)&(FY|FX)
		if sub > 0x7FFF {
			z80.Reg.F |= FS
		}
		if sub == 0 {
			z80.Reg.F |= FZ
		}
		if (hl^nn)&0x8000 != 0 && (hl^sub)&0x8000 != 0 {
			z80.Reg.F |= FP
		}
		if sub > hl {
			z80.Reg.F |= FC
		}
		z80.Reg.setHL(sub)
	case rld:
		hl := z80.Reg.HL()
		w := (uint16(z80.Reg.A)<<8 | uint16(z80.read(hl))) << 4
		z80.addContention(hl, 4)
		z80.write(hl, byte(w)|z80.Reg.A&0x0F)
		z80.Reg.A = z80.Reg.A&0xF0 | byte(w>>8)&0x0F
		z80.Reg.F = z80.Reg.F&FC | z80.Reg.A&(FS|FY|FX) | parity[z80.Reg.A]
		if z80.Reg.A == 0 {
			z80.Reg.F |= FZ
		}
	case rrd:
		hl := z80.Reg.HL()
		w := (uint16(z80.Reg.A)<<8 | uint16(z80.read(hl)))
		z80.addContention(hl, 4)
		z80.write(hl, byte(w>>4))
		z80.Reg.A = z80.Reg.A&0xF0 | byte(w)&0x0F
		z80.Reg.F = z80.Reg.F&FC | z80.Reg.A&(FS|FY|FX) | parity[z80.Reg.A]
		if z80.Reg.A == 0 {
			z80.Reg.F |= FZ
		}
	case in_a_c, in_b_c, in_c_c, in_d_c, in_e_c, in_f_c, in_h_c, in_l_c:
		r := z80.Reg.r(opcode & 0b00111000 >> 3)
		*r = z80.readBus(z80.Reg.B, z80.Reg.C)
		z80.Reg.F = z80.Reg.F&FC | *r&FS | parity[*r]
		if *r == 0 {
			z80.Reg.F |= FZ
		}
	case out_c_a, out_c_b, out_c_c, out_c_d, out_c_e, out_c_f, out_c_h, out_c_l:
		z80.writeBus(z80.Reg.B, z80.Reg.C, *z80.Reg.r(opcode & 0b00111000 >> 3))
	case im0:
		z80.im = 0
	case im1:
		z80.im = 1
	case im2:
		z80.im = 2
	case retn, 0x55, 0x65, 0x75, 0x5D, 0x6D, reti, 0x7D:
		z80.iff1 = z80.iff2
		z80.Reg.PC = uint16(z80.read(z80.Reg.SP+1))<<8 | uint16(z80.read(z80.Reg.SP))
		z80.Reg.SP += 2
	case ld_mm_bc, ld_mm_hl_ed, ld_mm_de, ld_mm_sp:
		addr := z80.nextWord()
		rr := z80.Reg.rr(opcode & 0b00110000 >> 4)
		z80.write(addr, byte(rr))
		z80.write(addr+1, byte(rr>>8))
	case ld_bc_mm, ld_de_mm, ld_hl_mm_ed, ld_sp_mm:
		addr := z80.nextWord()
		z80.Reg.setRR(opcode&0b00110000>>4, uint16(z80.read(addr))|uint16(z80.read(addr+1))<<8)
	case ld_a_r:
		z80.addContention(z80.Reg.IR(), 1)
		z80.Reg.A = z80.Reg.R
		z80.Reg.F = z80.Reg.F&FC | z80.Reg.A&FS
		if z80.Reg.A == 0 {
			z80.Reg.F |= FZ
		}
		if z80.iff2 {
			z80.Reg.F |= FP
		}
	case ld_r_a:
		z80.addContention(z80.Reg.IR(), 1)
		z80.Reg.R = z80.Reg.A
	case ld_a_i:
		z80.addContention(z80.Reg.IR(), 1)
		z80.Reg.A = z80.Reg.I
		z80.Reg.F = z80.Reg.F&FC | z80.Reg.A&FS
		if z80.Reg.A == 0 {
			z80.Reg.F |= FZ
		}
		if z80.iff2 {
			z80.Reg.F |= FP
		}
	case ld_i_a:
		z80.addContention(z80.Reg.IR(), 1)
		z80.Reg.I = z80.Reg.A
	case ldi, ldir, ldd, lddr:
		hl := z80.Reg.HL()
		de := z80.Reg.DE()
		bc := z80.Reg.BC() - 1
		n := z80.read(hl)
		z80.write(de, n)
		z80.addContention(de, 2)
		if opcode == ldi || opcode == ldir {
			z80.Reg.setHL(hl + 1)
			z80.Reg.setDE(de + 1)
		} else {
			z80.Reg.setHL(hl - 1)
			z80.Reg.setDE(de - 1)
		}
		z80.Reg.setBC(bc)
		z80.Reg.F = z80.Reg.F & (FS | FZ | FC)
		n += z80.Reg.A
		z80.Reg.F |= FY&(n<<4) | FX&n
		if bc != 0 {
			z80.Reg.F |= FP
			if opcode == ldir || opcode == lddr {
				z80.Reg.PC -= 2
				z80.TC.Add(5)
			}
		}
	case cpi, cpir, cpd, cpdr:
		hl := z80.Reg.HL()
		bc := z80.Reg.BC() - 1
		if opcode == cpi || opcode == cpir {
			z80.Reg.setHL(hl + 1)
		} else {
			z80.Reg.setHL(hl - 1)
		}
		z80.Reg.setBC(bc)
		n := z80.read(hl)
		z80.addContention(hl, 5)
		test := z80.Reg.A - n
		z80.Reg.F = z80.Reg.F&FC | FN | test&FS
		if test == 0 {
			z80.Reg.F |= FZ
		}
		z80.Reg.F |= byte(z80.Reg.A^n^test) & FH
		if bc != 0 {
			z80.Reg.F |= FP
		}
		n = test - (z80.Reg.F&FH)>>4
		z80.Reg.F |= FY&(n<<4) | FX&n
		if (opcode == cpir || opcode == cpdr) && bc != 0 && test != 0 {
			z80.Reg.PC -= 2
			z80.addContention(hl, 5)
		}
	case ini, inir, ind, indr:
		z80.addContention(z80.Reg.IR(), 1)
		hl := z80.Reg.HL()
		n := z80.readBus(z80.Reg.B, z80.Reg.C)
		z80.write(hl, n)
		z80.Reg.B -= 1
		if opcode == ini || opcode == inir {
			z80.Reg.setHL(hl + 1)
		} else {
			z80.Reg.setHL(hl - 1)
		}
		z80.Reg.F = z80.Reg.F & ^FZ | FN
		if z80.Reg.B == 0 {
			z80.Reg.F |= FZ
		} else if opcode == inir || opcode == indr {
			z80.Reg.PC -= 2
			z80.addContention(hl, 5)
		}
	case outi, otir, outd, otdr:
		z80.addContention(z80.Reg.IR(), 1)
		hl := z80.Reg.HL()
		z80.Reg.B -= 1
		z80.writeBus(z80.Reg.B, z80.Reg.C, z80.read(hl))
		if opcode == outi || opcode == otir {
			z80.Reg.setHL(hl + 1)
		} else {
			z80.Reg.setHL(hl - 1)
		}
		z80.Reg.F = z80.Reg.F & ^FZ | FN
		if z80.Reg.B == 0 {
			z80.Reg.F |= FZ
		} else if opcode == otir || opcode == otdr {
			z80.Reg.PC -= 2
			z80.addContention(z80.Reg.BC(), 5)
		}
	}
}
