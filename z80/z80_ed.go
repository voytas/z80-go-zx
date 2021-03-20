package z80

// Handles opcodes with ED prefix
func (z80 *Z80) prefixED(opcode byte) {
	switch opcode {
	case neg, 0x54, 0x64, 0x74, 0x4C, 0x5C, 0x6C, 0x7C:
		a := z80.reg.A
		z80.reg.A = ^a + 1
		z80.reg.F = (fS|fY|fX)&z80.reg.A | fN | byte(a^z80.reg.A)&fH
		if z80.reg.A == 0 {
			z80.reg.F |= fZ
		}
		if z80.reg.A == 0x80 {
			z80.reg.F |= fP
		}
		if a != 0 {
			z80.reg.F |= fC
		}
	case adc_hl_bc, adc_hl_de, adc_hl_hl, adc_hl_sp:
		z80.contention(z80.reg.IR(), 7)
		hl := z80.reg.HL()
		nn := z80.reg.rr(opcode & 0b00110000 >> 4)
		sum := hl + nn + uint16(z80.reg.F&fC)
		z80.reg.F = byte((hl^nn^sum)>>8)&fH | byte(sum>>8)&(fY|fX)
		if sum > 0x7FFF {
			z80.reg.F |= fS
		}
		if sum == 0 {
			z80.reg.F |= fZ
		}
		if (hl^nn)&0x8000 == 0 && (hl^sum)&0x8000 != 0 {
			z80.reg.F |= fP
		}
		if sum < hl {
			z80.reg.F |= fC
		}

		z80.reg.setHL(sum)
	case sbc_hl_bc, sbc_hl_de, sbc_hl_hl, sbc_hl_sp:
		z80.contention(z80.reg.IR(), 7)
		hl := z80.reg.HL()
		nn := z80.reg.rr(opcode & 0b00110000 >> 4)
		sub := hl - nn - uint16(z80.reg.F&fC)
		z80.reg.F = fN | byte((hl^nn^sub)>>8)&fH | byte(sub>>8)&(fY|fX)
		if sub > 0x7FFF {
			z80.reg.F |= fS
		}
		if sub == 0 {
			z80.reg.F |= fZ
		}
		if (hl^nn)&0x8000 != 0 && (hl^sub)&0x8000 != 0 {
			z80.reg.F |= fP
		}
		if sub > hl {
			z80.reg.F |= fC
		}
		z80.reg.setHL(sub)
	case rld:
		hl := z80.reg.HL()
		w := (uint16(z80.reg.A)<<8 | uint16(z80.read(hl))) << 4
		z80.contention(hl, 4)
		z80.write(hl, byte(w)|z80.reg.A&0x0F)
		z80.reg.A = z80.reg.A&0xF0 | byte(w>>8)&0x0F
		z80.reg.F = z80.reg.F&fC | z80.reg.A&(fS|fY|fX) | parity[z80.reg.A]
		if z80.reg.A == 0 {
			z80.reg.F |= fZ
		}
	case rrd:
		hl := z80.reg.HL()
		w := (uint16(z80.reg.A)<<8 | uint16(z80.read(hl)))
		z80.contention(hl, 4)
		z80.write(hl, byte(w>>4))
		z80.reg.A = z80.reg.A&0xF0 | byte(w)&0x0F
		z80.reg.F = z80.reg.F&fC | z80.reg.A&(fS|fY|fX) | parity[z80.reg.A]
		if z80.reg.A == 0 {
			z80.reg.F |= fZ
		}
	case in_a_c, in_b_c, in_c_c, in_d_c, in_e_c, in_f_c, in_h_c, in_l_c:
		r := z80.reg.r(opcode & 0b00111000 >> 3)
		*r = z80.readBus(z80.reg.B, z80.reg.C)
		z80.reg.F = z80.reg.F&fC | *r&fS | parity[*r]
		if *r == 0 {
			z80.reg.F |= fZ
		}
	case out_c_a, out_c_b, out_c_c, out_c_d, out_c_e, out_c_f, out_c_h, out_c_l:
		z80.writeBus(z80.reg.B, z80.reg.C, *z80.reg.r(opcode & 0b00111000 >> 3))
	case im0:
		z80.im = 0
	case im1:
		z80.im = 1
	case im2:
		z80.im = 2
	case retn, 0x55, 0x65, 0x75, 0x5D, 0x6D, reti, 0x7D:
		z80.iff1 = z80.iff2
		z80.reg.PC = uint16(z80.read(z80.reg.SP+1))<<8 | uint16(z80.read(z80.reg.SP))
		z80.reg.SP += 2
	case ld_mm_bc, ld_mm_hl_ed, ld_mm_de, ld_mm_sp:
		addr := z80.readWord()
		rr := z80.reg.rr(opcode & 0b00110000 >> 4)
		z80.write(addr, byte(rr))
		z80.write(addr+1, byte(rr>>8))
	case ld_bc_mm, ld_de_mm, ld_hl_mm_ed, ld_sp_mm:
		addr := z80.readWord()
		z80.reg.setRR(opcode&0b00110000>>4, uint16(z80.read(addr))|uint16(z80.read(addr+1))<<8)
	case ld_a_r:
		z80.contention(z80.reg.IR(), 1)
		z80.reg.A = z80.reg.R
		z80.reg.F = z80.reg.F&fC | z80.reg.A&fS
		if z80.reg.A == 0 {
			z80.reg.F |= fZ
		}
		if z80.iff2 {
			z80.reg.F |= fP
		}
	case ld_r_a:
		z80.contention(z80.reg.IR(), 1)
		z80.reg.R = z80.reg.A
	case ld_a_i:
		z80.contention(z80.reg.IR(), 1)
		z80.reg.A = z80.reg.I
		z80.reg.F = z80.reg.F&fC | z80.reg.A&fS
		if z80.reg.A == 0 {
			z80.reg.F |= fZ
		}
		if z80.iff2 {
			z80.reg.F |= fP
		}
	case ld_i_a:
		z80.contention(z80.reg.IR(), 1)
		z80.reg.I = z80.reg.A
	case ldi, ldir, ldd, lddr:
		hl := z80.reg.HL()
		de := z80.reg.DE()
		bc := z80.reg.BC() - 1
		n := z80.read(hl)
		z80.write(de, n)
		z80.contention(de, 2)
		if opcode == ldi || opcode == ldir {
			z80.reg.setHL(hl + 1)
			z80.reg.setDE(de + 1)
		} else {
			z80.reg.setHL(hl - 1)
			z80.reg.setDE(de - 1)
		}
		z80.reg.setBC(bc)
		z80.reg.F = z80.reg.F & (fS | fZ | fC)
		n += z80.reg.A
		z80.reg.F |= fY&(n<<4) | fX&n
		if bc != 0 {
			z80.reg.F |= fP
			if opcode == ldir || opcode == lddr {
				z80.reg.PC -= 2
				z80.TC.Add(5)
			}
		}
	case cpi, cpir, cpd, cpdr:
		hl := z80.reg.HL()
		bc := z80.reg.BC() - 1
		if opcode == cpi || opcode == cpir {
			z80.reg.setHL(hl + 1)
		} else {
			z80.reg.setHL(hl - 1)
		}
		z80.reg.setBC(bc)
		n := z80.read(hl)
		z80.contention(hl, 5)
		test := z80.reg.A - n
		z80.reg.F = z80.reg.F&fC | fN | test&fS
		if test == 0 {
			z80.reg.F |= fZ
		}
		z80.reg.F |= byte(z80.reg.A^n^test) & fH
		if bc != 0 {
			z80.reg.F |= fP
		}
		n = test - (z80.reg.F&fH)>>4
		z80.reg.F |= fY&(n<<4) | fX&n
		if (opcode == cpir || opcode == cpdr) && bc != 0 && test != 0 {
			z80.reg.PC -= 2
			z80.contention(hl, 5)
		}
	case ini, inir, ind, indr:
		z80.contention(z80.reg.IR(), 1)
		hl := z80.reg.HL()
		n := z80.readBus(z80.reg.B, z80.reg.C)
		z80.write(hl, n)
		z80.reg.B -= 1
		if opcode == ini || opcode == inir {
			z80.reg.setHL(hl + 1)
		} else {
			z80.reg.setHL(hl - 1)
		}
		z80.reg.F = z80.reg.F & ^fZ | fN
		if z80.reg.B == 0 {
			z80.reg.F |= fZ
		} else if opcode == inir || opcode == indr {
			z80.reg.PC -= 2
			z80.contention(hl, 5)
		}
	case outi, otir, outd, otdr:
		z80.contention(z80.reg.IR(), 1)
		hl := z80.reg.HL()
		z80.reg.B -= 1
		z80.writeBus(z80.reg.B, z80.reg.C, z80.read(hl))
		if opcode == outi || opcode == otir {
			z80.reg.setHL(hl + 1)
		} else {
			z80.reg.setHL(hl - 1)
		}
		z80.reg.F = z80.reg.F & ^fZ | fN
		if z80.reg.B == 0 {
			z80.reg.F |= fZ
		} else if opcode == otir || opcode == otdr {
			z80.reg.PC -= 2
			z80.contention(z80.reg.BC(), 5)
		}
	}
}
