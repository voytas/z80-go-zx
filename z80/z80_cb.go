package z80

// Handles opcodes with CB prefix
func (z80 *Z80) prefixCB() {
	const HL = 0b110
	var v byte
	var hl uint16

	opcode := z80.readByte()
	var offset byte
	if z80.reg.prefix != noPrefix {
		offset = opcode
		opcode = z80.readByte() // for IX and IY the actual opcode comes after the offset
		z80.t = 8               // extra 8 t-states for IX and IY
	}
	reg := opcode & 0b00000111
	if reg == HL {
		hl = z80.getHLOffset(offset)
		v = z80.mem.Read(hl)
		z80.t += 15 // 15 t-states for HL with one exception for bit
	} else {
		if z80.reg.prefix == noPrefix {
			v = *z80.reg.raw[reg]
		} else {
			hl = z80.getHLOffset(offset)
			v = z80.mem.Read(hl)
		}
		z80.t += 8 // 8 t-states for registers
	}

	var cf byte
	write := func(flags bool) {
		if reg != HL {
			*z80.reg.raw[reg] = v
		}
		if reg == HL || z80.reg.prefix != noPrefix {
			z80.mem.Write(hl, v)
		}
		if flags {
			z80.reg.F = (fS | fY | fX) & v
			if v == 0 {
				z80.reg.F |= fZ
			}
			z80.reg.F |= parity[v] | cf
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
		v = v<<1 | fC&z80.reg.F
		write(true)
	case rr_r:
		cf = v & fC
		v = v>>1 | fC&z80.reg.F<<7
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
			z80.reg.F = z80.reg.F&fC | fH
			test := v & bitMask[bit]
			if test == 0 {
				z80.reg.F |= fZ | fP
			}
			if reg == HL {
				// Might not be 100%, this undocumented behaviour is not clear, but it passses test
				z80.reg.F |= fS&test | (fY|fX)&byte(hl>>8)
				z80.t -= 3 // bit and HL is 12 t-states
			} else {
				z80.reg.F |= fS&test | (fY|fX)&v
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
