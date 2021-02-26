package z80

// Handles opcodes with CB prefix
func (cpu *CPU) prefixCB() {
	var v byte
	var hl uint16

	opcode := cpu.readByte()
	var offset byte
	if cpu.reg.prefix != noPrefix {
		// for IX and IY the actual opcode comes after the offset
		offset = opcode
		opcode = cpu.readByte()
		cpu.t = 8
	}
	reg := opcode & 0b00000111
	if reg == 0b110 {
		hl = cpu.getHLOffset(offset)
		v = cpu.mem.Read(hl)
		cpu.t += 15 // the only exception is bit operation that takes 12 t-states
	} else {
		if cpu.reg.prefix == noPrefix {
			v = *cpu.reg.raw[reg]
		} else {
			hl = cpu.getHLOffset(offset)
			v = cpu.mem.Read(hl)
		}
		cpu.t += 8
	}

	var cy byte
	write := func(flags bool) {
		if reg != 0b110 {
			*cpu.reg.raw[reg] = v
		}
		if reg == 0b110 || cpu.reg.prefix != noPrefix {
			cpu.mem.Write(hl, v)
		}
		if flags {
			cpu.reg.F = f_NONE
			cpu.reg.F |= f_S & v
			if v == 0 {
				cpu.reg.F |= f_Z
			}
			cpu.reg.F |= parity[v] | cy
		}
	}

	switch opcode & 0b11111000 {
	case rlc_r:
		cy = v >> 7
		v = v<<1 | cy
		write(true)
	case rrc_r:
		cy = v & f_C
		v = v>>1 | cy<<7
		write(true)
	case rl_r:
		cy = v >> 7
		v = v<<1 | f_C&cpu.reg.F
		write(true)
	case rr_r:
		cy = v & f_C
		v = v>>1 | f_C&cpu.reg.F<<7
		write(true)
	case sla_r:
		cy = v >> 7
		v = v << 1
		write(true)
	case sra_r:
		cy = v & f_C
		v = v&0x80 | v>>1
		write(true)
	case sll_r:
		cy = v >> 7
		v = v<<1 | 0x01
		write(true)
	case srl_r:
		cy = v & f_C
		v = v >> 1
		write(true)
	default:
		bit := (opcode & 0b00111000) >> 3
		switch opcode & 0b11000000 {
		case bit_b:
			if reg == 0b110 {
				cpu.t -= 3 // for bit operation it is 12 t-states, not usual 15
			}
			cpu.reg.F &= ^(f_Z | f_N)
			cpu.reg.F |= f_H
			if v&bit_mask[bit] == 0 {
				cpu.reg.F |= f_Z
			}
		case res_b:
			v &= ^bit_mask[bit]
			write(false)
		case set_b:
			v |= bit_mask[bit]
			write(false)
		}
	}
}
