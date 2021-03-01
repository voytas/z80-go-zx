package z80

// Handles opcodes with CB prefix
func (cpu *CPU) prefixCB() {
	const HL = 0b110
	var v byte
	var hl uint16

	opcode := cpu.readByte()
	var offset byte
	if cpu.reg.prefix != noPrefix {
		offset = opcode
		opcode = cpu.readByte() // for IX and IY the actual opcode comes after the offset
		cpu.t = 8               // extra 8 t-states for IX and IY
	}
	reg := opcode & 0b00000111
	if reg == HL {
		hl = cpu.getHLOffset(offset)
		v = cpu.mem.Read(hl)
		cpu.t += 15 // 15 t-states for HL with one exception for bit
	} else {
		if cpu.reg.prefix == noPrefix {
			v = *cpu.reg.raw[reg]
		} else {
			hl = cpu.getHLOffset(offset)
			v = cpu.mem.Read(hl)
		}
		cpu.t += 8 // 8 t-states for registers
	}

	var cf byte
	write := func(flags bool) {
		if reg != HL {
			*cpu.reg.raw[reg] = v
		}
		if reg == HL || cpu.reg.prefix != noPrefix {
			cpu.mem.Write(hl, v)
		}
		if flags {
			cpu.reg.F = (fS | fY | fX) & v
			if v == 0 {
				cpu.reg.F |= fZ
			}
			cpu.reg.F |= parity[v] | cf
		}
	}

	switch opcode & 0b11111000 {
	case rlc_r:
		cf = v >> 7
		v = v<<1 | cf
		write(true)
	case rrc_r:
		cf = v & fC
		v = v>>1 | cf<<7
		write(true)
	case rl_r:
		cf = v >> 7
		v = v<<1 | fC&cpu.reg.F
		write(true)
	case rr_r:
		cf = v & fC
		v = v>>1 | fC&cpu.reg.F<<7
		write(true)
	case sla_r:
		cf = v >> 7
		v = v << 1
		write(true)
	case sra_r:
		cf = v & fC
		v = v&0x80 | v>>1
		write(true)
	case sll_r:
		cf = v >> 7
		v = v<<1 | 0x01
		write(true)
	case srl_r:
		cf = v & fC
		v = v >> 1
		write(true)
	default:
		bit := (opcode & 0b00111000) >> 3
		switch opcode & 0b11000000 {
		case bit_b:
			cpu.reg.F = cpu.reg.F&fC | fH
			test := v & bitMask[bit]
			if test == 0 {
				cpu.reg.F |= fZ | fP
			}
			if reg == HL {
				// Might not be 100%, this undocumented behaviour is not clear, but it passses test
				cpu.reg.F |= fS&test | (fY|fX)&byte(hl>>8)
				cpu.t -= 3 // bit and HL is 12 t-states
			} else {
				cpu.reg.F |= fS&test | (fY|fX)&v
			}
		case res_b:
			v &= ^bitMask[bit]
			write(false)
		case set_b:
			v |= bitMask[bit]
			write(false)
		}
	}
}
