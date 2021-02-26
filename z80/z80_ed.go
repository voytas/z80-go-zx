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
		cpu.reg.F = f_N
		cpu.reg.F |= f_S & cpu.reg.A
		if cpu.reg.A == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte(a^cpu.reg.A) & f_H
		if cpu.reg.A == 0x80 {
			cpu.reg.F |= f_P
		}
		if a != 0 {
			cpu.reg.F |= f_C
		}
	case adc_hl_bc, adc_hl_de, adc_hl_hl, adc_hl_sp:
		hl := cpu.reg.HL()
		nn := cpu.reg.rr(opcode & 0b00110000 >> 4)
		sum := hl + nn + uint16(cpu.reg.F&f_C)
		cpu.reg.F = f_NONE
		if sum > 0x7FFF {
			cpu.reg.F |= f_S
		}
		if sum == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte((hl^nn^sum)>>8) & f_H
		if (hl^nn)&0x8000 == 0 && (hl^sum)&0x8000 != 0 {
			cpu.reg.F |= f_P
		}
		if sum < hl {
			cpu.reg.F |= f_C
		}
		cpu.reg.setHL(sum)
	case sbc_hl_bc, sbc_hl_de, sbc_hl_hl, sbc_hl_sp:
		hl := cpu.reg.HL()
		nn := cpu.reg.rr(opcode & 0b00110000 >> 4)
		sub := hl - nn - uint16(cpu.reg.F&f_C)
		cpu.reg.F = f_N
		if sub > 0x7FFF {
			cpu.reg.F |= f_S
		}
		if sub == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte((hl^nn^sub)>>8) & f_H
		if (hl^nn)&0x8000 != 0 && (hl^sub)&0x8000 != 0 {
			cpu.reg.F |= f_P
		}
		if sub > hl {
			cpu.reg.F |= f_C
		}
		cpu.reg.setHL(sub)
	case rld:
		hl := cpu.reg.HL()
		w := (uint16(cpu.reg.A)<<8 | uint16(cpu.mem.Read(hl))) << 4
		cpu.mem.Write(hl, byte(w)|cpu.reg.A&0x0F)
		cpu.reg.A = cpu.reg.A&0xF0 | byte(w>>8)&0x0F
		cpu.reg.F &= f_C
		cpu.reg.F |= cpu.reg.A & f_S
		if cpu.reg.A == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= parity[cpu.reg.A]
	case rrd:
		hl := cpu.reg.HL()
		w := (uint16(cpu.reg.A)<<8 | uint16(cpu.mem.Read(hl)))
		cpu.mem.Write(hl, byte(w>>4))
		cpu.reg.A = cpu.reg.A&0xF0 | byte(w)&0x0F
		cpu.reg.F &= f_C
		cpu.reg.F |= cpu.reg.A & f_S
		if cpu.reg.A == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= parity[cpu.reg.A]
	case in_a_c, in_b_c, in_c_c, in_d_c, in_e_c, in_f_c, in_h_c, in_l_c:
		r := cpu.reg.r(opcode & 0b00111000 >> 3)
		*r = cpu.IN(cpu.reg.B, cpu.reg.C)
		cpu.reg.F &= f_C
		cpu.reg.F |= *r & f_S
		if *r == 0 {
			cpu.reg.F |= f_Z
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
		cpu.reg.F = cpu.reg.F&f_C | cpu.reg.A&f_S
		if cpu.reg.A == 0 {
			cpu.reg.F |= f_Z
		}
		if cpu.iff2 {
			cpu.reg.F |= f_P
		}
	case ld_r_a:
		cpu.reg.R = cpu.reg.A
	case ld_a_i:
		cpu.reg.A = cpu.reg.I
		cpu.reg.F = cpu.reg.F&f_C | cpu.reg.A&f_S
		if cpu.reg.A == 0 {
			cpu.reg.F |= f_Z
		}
		if cpu.iff2 {
			cpu.reg.F |= f_P
		}
	case ld_i_a:
		cpu.reg.I = cpu.reg.A
	case ldi, ldir, ldd, lddr:
		hl := cpu.reg.HL()
		de := cpu.reg.DE()
		bc := cpu.reg.BC() - 1
		cpu.mem.Write(de, cpu.mem.Read(hl))
		if opcode == ldi || opcode == ldir {
			cpu.reg.setHL(hl + 1)
			cpu.reg.setDE(de + 1)
		} else {
			cpu.reg.setHL(hl - 1)
			cpu.reg.setDE(de - 1)
		}
		cpu.reg.setBC(bc)
		cpu.reg.F &= ^(f_H | f_P | f_N)
		if bc != 0 {
			cpu.reg.F |= f_P
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
		cpu.reg.F = cpu.reg.F&f_C | f_N | test&f_S
		if test == 0 {
			cpu.reg.F |= f_Z
		}
		cpu.reg.F |= byte(cpu.reg.A^n^test) & f_H
		if bc != 0 {
			cpu.reg.F |= f_P
		}
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
		cpu.reg.F = cpu.reg.F & ^f_Z | f_N
		if cpu.reg.B == 0 {
			cpu.reg.F |= f_Z
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
		cpu.reg.F = cpu.reg.F & ^f_Z | f_N
		if cpu.reg.B == 0 {
			cpu.reg.F |= f_Z
		} else if opcode == otir || opcode == otdr {
			cpu.reg.PC -= 2
			cpu.t += 5
		}
	}
}
