package z80

// Handles opcodes with ED prefix
func (cpu *CPU) prefixED(opcode byte) {
	// TODO: check t-states calc
	t := tStatesED[opcode]
	if t != 0 {
		cpu.t = t
	} else {
		cpu.t = 2 * tStates[nop]
	}

	switch opcode {
	case neg, 0x54, 0x64, 0x74, 0x4C, 0x5C, 0x6C, 0x7C:
		a := cpu.reg.A
		cpu.reg.A = ^a + 1
		cpu.reg.F = (fS|fY|fX)&cpu.reg.A | fN
		if cpu.reg.A == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= byte(a^cpu.reg.A) & fH
		if cpu.reg.A == 0x80 {
			cpu.reg.F |= fP
		}
		if a != 0 {
			cpu.reg.F |= fC
		}
	case adc_hl_bc, adc_hl_de, adc_hl_hl, adc_hl_sp:
		hl := cpu.reg.HL()
		nn := cpu.reg.rr(opcode & 0b00110000 >> 4)
		sum := hl + nn + uint16(cpu.reg.F&fC)
		cpu.reg.F = fNONE
		if sum > 0x7FFF {
			cpu.reg.F |= fS
		}
		if sum == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= byte((hl^nn^sum)>>8)&fH | byte(sum>>8)&(fY|fX)
		if (hl^nn)&0x8000 == 0 && (hl^sum)&0x8000 != 0 {
			cpu.reg.F |= fP
		}
		if sum < hl {
			cpu.reg.F |= fC
		}

		cpu.reg.setHL(sum)
	case sbc_hl_bc, sbc_hl_de, sbc_hl_hl, sbc_hl_sp:
		hl := cpu.reg.HL()
		nn := cpu.reg.rr(opcode & 0b00110000 >> 4)
		sub := hl - nn - uint16(cpu.reg.F&fC)
		cpu.reg.F = fN
		if sub > 0x7FFF {
			cpu.reg.F |= fS
		}
		if sub == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= byte((hl^nn^sub)>>8)&fH | byte(sub>>8)&(fY|fX)
		if (hl^nn)&0x8000 != 0 && (hl^sub)&0x8000 != 0 {
			cpu.reg.F |= fP
		}
		if sub > hl {
			cpu.reg.F |= fC
		}
		cpu.reg.setHL(sub)
	case rld:
		hl := cpu.reg.HL()
		w := (uint16(cpu.reg.A)<<8 | uint16(cpu.mem.Read(hl))) << 4
		cpu.mem.Write(hl, byte(w)|cpu.reg.A&0x0F)
		cpu.reg.A = cpu.reg.A&0xF0 | byte(w>>8)&0x0F
		cpu.reg.F = cpu.reg.F&fC | cpu.reg.A&(fS|fY|fX)
		if cpu.reg.A == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= parity[cpu.reg.A]
	case rrd:
		hl := cpu.reg.HL()
		w := (uint16(cpu.reg.A)<<8 | uint16(cpu.mem.Read(hl)))
		cpu.mem.Write(hl, byte(w>>4))
		cpu.reg.A = cpu.reg.A&0xF0 | byte(w)&0x0F
		cpu.reg.F = cpu.reg.F&fC | cpu.reg.A&(fS|fY|fX)
		if cpu.reg.A == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= parity[cpu.reg.A]
	case in_a_c, in_b_c, in_c_c, in_d_c, in_e_c, in_f_c, in_h_c, in_l_c:
		r := cpu.reg.r(opcode & 0b00111000 >> 3)
		*r = cpu.IN(cpu.reg.B, cpu.reg.C)
		cpu.reg.F &= fC
		cpu.reg.F |= *r & fS
		if *r == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= parity[*r]
	case out_c_a, out_c_b, out_c_c, out_c_d, out_c_e, out_c_f, out_c_h, out_c_l:
		cpu.OUT(cpu.reg.B, cpu.reg.C, *cpu.reg.r(opcode & 0b00111000 >> 3))
	case im0, im1, im2:
		cpu.im = opcode
	case retn, 0x55, 0x65, 0x75, 0x5D, 0x6D, reti, 0x7D:
		cpu.iff1 = cpu.iff2
		cpu.reg.PC = uint16(cpu.mem.Read(cpu.reg.SP+1))<<8 | uint16(cpu.mem.Read(cpu.reg.SP))
		cpu.reg.SP += 2
	case ld_mm_bc, ld_mm_hl_ed, ld_mm_de, ld_mm_sp:
		addr := cpu.readWord()
		rr := cpu.reg.rr(opcode & 0b00110000 >> 4)
		cpu.mem.Write(addr, byte(rr))
		cpu.mem.Write(addr+1, byte(rr>>8))
	case ld_bc_mm, ld_de_mm, ld_hl_mm_ed, ld_sp_mm:
		addr := cpu.readWord()
		cpu.reg.setRR(opcode&0b00110000>>4, uint16(cpu.mem.Read(addr))|uint16(cpu.mem.Read(addr+1))<<8)
	case ld_a_r:
		cpu.reg.A = cpu.reg.R
		cpu.reg.F = cpu.reg.F&fC | cpu.reg.A&fS
		if cpu.reg.A == 0 {
			cpu.reg.F |= fZ
		}
		if cpu.iff2 {
			cpu.reg.F |= fP
		}
	case ld_r_a:
		cpu.reg.R = cpu.reg.A
	case ld_a_i:
		cpu.reg.A = cpu.reg.I
		cpu.reg.F = cpu.reg.F&fC | cpu.reg.A&fS
		if cpu.reg.A == 0 {
			cpu.reg.F |= fZ
		}
		if cpu.iff2 {
			cpu.reg.F |= fP
		}
	case ld_i_a:
		cpu.reg.I = cpu.reg.A
	case ldi, ldir, ldd, lddr:
		hl := cpu.reg.HL()
		de := cpu.reg.DE()
		bc := cpu.reg.BC() - 1
		n := cpu.mem.Read(hl)
		cpu.mem.Write(de, n)
		if opcode == ldi || opcode == ldir {
			cpu.reg.setHL(hl + 1)
			cpu.reg.setDE(de + 1)
		} else {
			cpu.reg.setHL(hl - 1)
			cpu.reg.setDE(de - 1)
		}
		cpu.reg.setBC(bc)
		cpu.reg.F = cpu.reg.F & (fS | fZ | fC)
		n += cpu.reg.A
		cpu.reg.F |= fY&(n<<4) | fX&n
		if bc != 0 {
			cpu.reg.F |= fP
			if opcode == ldir || opcode == lddr {
				cpu.reg.PC -= 2
				cpu.t += 5
			}
		}
	case cpi, cpir, cpd, cpdr:
		hl := cpu.reg.HL()
		bc := cpu.reg.BC() - 1
		if opcode == cpi || opcode == cpir {
			cpu.reg.setHL(hl + 1)
		} else {
			cpu.reg.setHL(hl - 1)
		}
		cpu.reg.setBC(bc)
		n := cpu.mem.Read(hl)
		test := cpu.reg.A - n
		cpu.reg.F = cpu.reg.F&fC | fN | test&fS
		if test == 0 {
			cpu.reg.F |= fZ
		}
		cpu.reg.F |= byte(cpu.reg.A^n^test) & fH
		if bc != 0 {
			cpu.reg.F |= fP
		}
		n = test - (cpu.reg.F&fH)>>4
		cpu.reg.F |= fY&(n<<4) | fX&n
		if (opcode == cpir || opcode == cpdr) && bc != 0 && test != 0 {
			cpu.reg.PC -= 2
			cpu.t += 5
		}
	case ini, inir, ind, indr:
		hl := cpu.reg.HL()
		cpu.mem.Write(hl, cpu.IN(cpu.reg.B, cpu.reg.C))
		cpu.reg.B -= 1
		if opcode == ini || opcode == inir {
			cpu.reg.setHL(hl + 1)
		} else {
			cpu.reg.setHL(hl - 1)
		}
		cpu.reg.F = cpu.reg.F & ^fZ | fN
		if cpu.reg.B == 0 {
			cpu.reg.F |= fZ
		} else if opcode == inir || opcode == indr {
			cpu.reg.PC -= 2
			cpu.t += 5
		}
	case outi, otir, outd, otdr:
		hl := cpu.reg.HL()
		cpu.reg.B -= 1
		cpu.OUT(cpu.reg.B, cpu.reg.C, cpu.mem.Read(hl))
		if opcode == outi || opcode == otir {
			cpu.reg.setHL(hl + 1)
		} else {
			cpu.reg.setHL(hl - 1)
		}
		cpu.reg.F = cpu.reg.F & ^fZ | fN
		if cpu.reg.B == 0 {
			cpu.reg.F |= fZ
		} else if opcode == otir || opcode == otdr {
			cpu.reg.PC -= 2
			cpu.t += 5
		}
	}
}
