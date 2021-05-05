package z80

// Handles opcodes with CB prefix
func (z80 *Z80) prefixCB() {
	const HL = 0b110
	var v byte
	var hl uint16

	opcode := z80.fetch()
	var offset byte
	if z80.Reg.prefix != noPrefix {
		offset = opcode
		opcode = z80.fetch()
	}
	reg := opcode & 0b00000111
	if reg == HL {
		hl = z80.getHLOffset(offset)
		v = z80.read(hl)
	} else {
		if z80.Reg.prefix == noPrefix {
			v = *z80.Reg.raw[reg]
		} else {
			hl = z80.getHLOffset(offset)
			v = z80.read(hl)
		}
	}

	var cf byte
	write := func(flags bool) {
		if reg != HL {
			*z80.Reg.raw[reg] = v
		}
		if reg == HL || z80.Reg.prefix != noPrefix {
			z80.addContention(hl, 1)
			z80.write(hl, v)
		}
		if flags {
			z80.Reg.F = (FS | FY | FX) & v
			if v == 0 {
				z80.Reg.F |= FZ
			}
			z80.Reg.F |= parity[v] | cf
		}
	}

	switch opcode & 0b11111000 {
	case rlc_r:
		cf = v >> 7
		v = v<<1 | cf
		write(true)
	case rrc_r:
		cf = v & FC
		v = v>>1 | cf<<7
		write(true)
	case rl_r:
		cf = v >> 7
		v = v<<1 | FC&z80.Reg.F
		write(true)
	case rr_r:
		cf = v & FC
		v = v>>1 | FC&z80.Reg.F<<7
		write(true)
	case sla_r:
		cf = v >> 7
		v = v << 1
		write(true)
	case sra_r:
		cf = v & FC
		v = v&0x80 | v>>1
		write(true)
	case sll_r:
		cf = v >> 7
		v = v<<1 | 0x01
		write(true)
	case srl_r:
		cf = v & FC
		v = v >> 1
		write(true)
	default:
		bit := (opcode & 0b00111000) >> 3
		switch opcode & 0b11000000 {
		case bit_b:
			z80.Reg.F = z80.Reg.F&FC | FH
			test := v & bitMask[bit]
			if test == 0 {
				z80.Reg.F |= FZ | FP
			}
			if reg == HL {
				// Might not be 100%, this undocumented behaviour is not clear, but it passses test
				z80.Reg.F |= FS&test | (FY|FX)&byte(hl>>8)
				z80.addContention(hl, 1)
			} else {
				z80.Reg.F |= FS&test | (FY|FX)&v
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
