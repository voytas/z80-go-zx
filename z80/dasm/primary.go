package dasm

// Table containing primary instructions
var primaryInstructions = []*instruction{
	0x00: {mnemonic: "NOP"},
	0x01: {mnemonic: "LD   BC,$1$2", args: 2},
	0x02: {mnemonic: "LD   (BC),A"},
	0x03: {mnemonic: "INC  BC"},
	0x04: {mnemonic: "INC  B"},
	0x05: {mnemonic: "DEC  B"},
	0x06: {mnemonic: "LD   B,$1", args: 1},
	0x07: {mnemonic: "RLCA"},
	0x08: {mnemonic: "EX   AF,AF'"},
	0x09: {mnemonic: "ADD  HL,BC"},
	0x0A: {mnemonic: "LD   A,(BC)"},
	0x0B: {mnemonic: "DEC  BC"},
	0x0C: {mnemonic: "INC  C"},
	0x0D: {mnemonic: "DEC  C"},
	0x0E: {mnemonic: "LD   C,$1", args: 1},
	0x0F: {mnemonic: "RRCA"},
	0x10: {mnemonic: "DJNZ $1", args: 1},
	0x11: {mnemonic: "LD   DE,$1$2", args: 2},
	0x12: {mnemonic: "LD   (DE),A"},
	0x13: {mnemonic: "INC  DE"},
	0x14: {mnemonic: "INC  D"},
	0x15: {mnemonic: "DEC  D"},
	0x16: {mnemonic: "LD   D,$1", args: 1},
	0x17: {mnemonic: "RLA"},
	0x18: {mnemonic: "JR   $1", args: 1},
	0x19: {mnemonic: "ADD  HL,DE"},
	0x1A: {mnemonic: "LD   A,(BC)"},
	0x1B: {mnemonic: "DEC  DE"},
	0x1C: {mnemonic: "INC  C"},
	0x1D: {mnemonic: "DEC  C"},
	0x1E: {mnemonic: "LD   C,$1", args: 1},
	0x1F: {mnemonic: "RRA"},
	0x20: {mnemonic: "JR NZ,$1", args: 1},
	0x21: {mnemonic: "LD   HL,$1$2", args: 2},
	0x22: {mnemonic: "LD   ($1$2),HL", args: 2},
	0x23: {mnemonic: "INC  HL"},
	0x24: {mnemonic: "INC  H"},
	0x25: {mnemonic: "DEC  H"},
	0x26: {mnemonic: "LD   H,$1", args: 1},
	0x27: {mnemonic: "DAA"},
	0x28: {mnemonic: "JR   Z,$1", args: 1},
	0x29: {mnemonic: "ADD  HL,HL"},
	0x2A: {mnemonic: "LD   HL,($1$2)", args: 2},
	0x2B: {mnemonic: "DEC  HL"},
	0x2C: {mnemonic: "INC  L"},
	0x2D: {mnemonic: "DEC  L"},
	0x2E: {mnemonic: "LD   L,$1", args: 1},
	0x2F: {mnemonic: "CPL"},
	0x30: {mnemonic: "JR   NC,$1", args: 1},
	0x31: {mnemonic: "LD   SP,$1$2", args: 2},
	0x32: {mnemonic: "LD   ($1$2),A", args: 2},
	0x33: {mnemonic: "INC  SP"},
	0x34: {mnemonic: "INC  (HL)"},
	0x35: {mnemonic: "DEC  (HL)"},
	0x36: {mnemonic: "LD   (HL),$1", args: 1},
	0x37: {mnemonic: "SCF"},
	0x38: {mnemonic: "JR   C,$1", args: 1},
	0x39: {mnemonic: "ADD  HL,SP"},
	0x3A: {mnemonic: "LD   A,($1$2)", args: 2},
	0x3B: {mnemonic: "DEC  SP"},
	0x3C: {mnemonic: "INC  A"},
	0x3D: {mnemonic: "DEC  A"},
	0x3E: {mnemonic: "LD   A,$1", args: 1},
	0x3F: {mnemonic: "CCF"},
	0x40: {mnemonic: "LD   B,B"},
	0x41: {mnemonic: "LD   B,C"},
	0x42: {mnemonic: "LD   B,D"},
	0x43: {mnemonic: "LD   B,E"},
	0x44: {mnemonic: "LD   B,H"},
	0x45: {mnemonic: "LD   B,L"},
	0x46: {mnemonic: "LD   B,(HL)"},
	0x47: {mnemonic: "LD   B,A"},
	0x48: {mnemonic: "LD   C,B"},
	0x49: {mnemonic: "LD   C,C"},
	0x4A: {mnemonic: "LD   C,D"},
	0x4B: {mnemonic: "LD   C,E"},
	0x4C: {mnemonic: "LD   C,H"},
	0x4D: {mnemonic: "LD   C,L"},
	0x4E: {mnemonic: "LD   C,(HL)"},
	0x4F: {mnemonic: "LD   C,A"},
	0x50: {mnemonic: "LD   D,B"},
	0x51: {mnemonic: "LD   D,C"},
	0x52: {mnemonic: "LD   D,D"},
	0x53: {mnemonic: "LD   D,E"},
	0x54: {mnemonic: "LD   D,H"},
	0x55: {mnemonic: "LD   D,L"},
	0x56: {mnemonic: "LD   D,(HL)"},
	0x57: {mnemonic: "LD   D,A"},
	0x58: {mnemonic: "LD   E,B"},
	0x59: {mnemonic: "LD   E,C"},
	0x5A: {mnemonic: "LD   E,D"},
	0x5B: {mnemonic: "LD   E,E"},
	0x5C: {mnemonic: "LD   E,H"},
	0x5D: {mnemonic: "LD   E,L"},
	0x5E: {mnemonic: "LD   E,(HL)"},
	0x5F: {mnemonic: "LD   E,A"},
	0x60: {mnemonic: "LD   H,B"},
	0x61: {mnemonic: "LD   H,C"},
	0x62: {mnemonic: "LD   H,D"},
	0x63: {mnemonic: "LD   H,E"},
	0x64: {mnemonic: "LD   H,H"},
	0x65: {mnemonic: "LD   H,L"},
	0x66: {mnemonic: "LD   H,(HL)"},
	0x67: {mnemonic: "LD   H,A"},
	0x68: {mnemonic: "LD   L,B"},
	0x69: {mnemonic: "LD   L,C"},
	0x6A: {mnemonic: "LD   L,D"},
	0x6B: {mnemonic: "LD   L,E"},
	0x6C: {mnemonic: "LD   L,H"},
	0x6D: {mnemonic: "LD   L,L"},
	0x6E: {mnemonic: "LD   L,(HL)"},
	0x6F: {mnemonic: "LD   L,A"},
	0x70: {mnemonic: "LD   (HL),B"},
	0x71: {mnemonic: "LD   (HL),C"},
	0x72: {mnemonic: "LD   (HL),D"},
	0x73: {mnemonic: "LD   (HL),E"},
	0x74: {mnemonic: "LD   (HL),H"},
	0x75: {mnemonic: "LD   (HL),L"},
	0x76: {mnemonic: "HALT"},
	0x77: {mnemonic: "LD   (HL),A"},
	0x78: {mnemonic: "LD   A,B"},
	0x79: {mnemonic: "LD   A,C"},
	0x7A: {mnemonic: "LD   A,D"},
	0x7B: {mnemonic: "LD   A,E"},
	0x7C: {mnemonic: "LD   A,H"},
	0x7D: {mnemonic: "LD   A,L"},
	0x7E: {mnemonic: "LD   A,(HL)"},
	0x7F: {mnemonic: "LD   A,A"},
	0x80: {mnemonic: "ADD  A,B"},
	0x81: {mnemonic: "ADD  A,C"},
	0x82: {mnemonic: "ADD  A,D"},
	0x83: {mnemonic: "ADD  A,E"},
	0x84: {mnemonic: "ADD  A,H"},
	0x85: {mnemonic: "ADD  A,L"},
	0x86: {mnemonic: "ADD  A,(HL)"},
	0x87: {mnemonic: "ADD  A,A"},
	0x88: {mnemonic: "ADC  A,B"},
	0x89: {mnemonic: "ADC  A,C"},
	0x8A: {mnemonic: "ADC  A,D"},
	0x8B: {mnemonic: "ADC  A,E"},
	0x8C: {mnemonic: "ADC  A,H"},
	0x8D: {mnemonic: "ADC  A,L"},
	0x8E: {mnemonic: "ADC  A,(HL)"},
	0x8F: {mnemonic: "ADC  A,A"},
	0x90: {mnemonic: "SUB  B"},
	0x91: {mnemonic: "SUB  C"},
	0x92: {mnemonic: "SUB  D"},
	0x93: {mnemonic: "SUB  E"},
	0x94: {mnemonic: "SUB  H"},
	0x95: {mnemonic: "SUB  L"},
	0x96: {mnemonic: "SUB  (HL)"},
	0x97: {mnemonic: "SUB  A"},
	0x98: {mnemonic: "SBC  A,B"},
	0x99: {mnemonic: "SBC  A,C"},
	0x9A: {mnemonic: "SBC  A,D"},
	0x9B: {mnemonic: "SBC  A,E"},
	0x9C: {mnemonic: "SBC  A,H"},
	0x9D: {mnemonic: "SBC  A,L"},
	0x9E: {mnemonic: "SBC  A,(HL)"},
	0x9F: {mnemonic: "SBC  A,A"},
	0xA0: {mnemonic: "AND  B"},
	0xA1: {mnemonic: "AND  C"},
	0xA2: {mnemonic: "AND  D"},
	0xA3: {mnemonic: "AND  E"},
	0xA4: {mnemonic: "AND  H"},
	0xA5: {mnemonic: "AND  L"},
	0xA6: {mnemonic: "AND  (HL)"},
	0xA7: {mnemonic: "AND  A"},
	0xA8: {mnemonic: "XOR  B"},
	0xA9: {mnemonic: "XOR  C"},
	0xAA: {mnemonic: "XOR  D"},
	0xAB: {mnemonic: "XOR  E"},
	0xAC: {mnemonic: "XOR  H"},
	0xAD: {mnemonic: "XOR  L"},
	0xAE: {mnemonic: "XOR  (HL)"},
	0xAF: {mnemonic: "XOR  A"},
	0xB0: {mnemonic: "OR   B"},
	0xB1: {mnemonic: "OR   C"},
	0xB2: {mnemonic: "OR   D"},
	0xB3: {mnemonic: "OR   E"},
	0xB4: {mnemonic: "OR   H"},
	0xB5: {mnemonic: "OR   L"},
	0xB6: {mnemonic: "OR   (HL)"},
	0xB7: {mnemonic: "OR   A"},
	0xB8: {mnemonic: "CP   B"},
	0xB9: {mnemonic: "CP   C"},
	0xBA: {mnemonic: "CP   D"},
	0xBB: {mnemonic: "CP   E"},
	0xBC: {mnemonic: "CP   H"},
	0xBD: {mnemonic: "CP   L"},
	0xBE: {mnemonic: "CP   (HL)"},
	0xBF: {mnemonic: "CP   A"},
	0xC0: {mnemonic: "RET  NZ"},
	0xC1: {mnemonic: "POP  BC"},
	0xC2: {mnemonic: "JP   NZ,$1$2", args: 2},
	0xC3: {mnemonic: "JP   $1$2", args: 2},
	0xC4: {mnemonic: "CALL NZ,$1$2", args: 2},
	0xC5: {mnemonic: "PUSH BC"},
	0xC6: {mnemonic: "ADD  A,$1", args: 1},
	0xC7: {mnemonic: "RST  0h"},
	0xC8: {mnemonic: "RET  Z"},
	0xC9: {mnemonic: "RET"},
	0xCA: {mnemonic: "JP   Z,$1$2", args: 2},
	0xCC: {mnemonic: "CALL Z,$1$2", args: 2},
	0xCD: {mnemonic: "CALL $1$2", args: 2},
	0xCE: {mnemonic: "ADC  A,$1", args: 1},
	0xCF: {mnemonic: "RST  08h"},
	0xD0: {mnemonic: "RET  NC"},
	0xD1: {mnemonic: "POP  DE"},
	0xD2: {mnemonic: "JP   NC,$1$2", args: 2},
	0xD3: {mnemonic: "OUT  ($1),A", args: 1},
	0xD4: {mnemonic: "CALL NC,$1$2", args: 2},
	0xD5: {mnemonic: "PUSH DE"},
	0xD6: {mnemonic: "SUB  $1", args: 1},
	0xD7: {mnemonic: "RST  10h"},
	0xD8: {mnemonic: "RET  C"},
	0xD9: {mnemonic: "EXX"},
	0xDA: {mnemonic: "JP   C,$1$2", args: 2},
	0xDB: {mnemonic: "IN   A,($1)", args: 1},
	0xDC: {mnemonic: "CALL C,$1$2", args: 2},
	0xDE: {mnemonic: "SBC  A,$1", args: 1},
	0xDF: {mnemonic: "RST  18h"},
	0xE0: {mnemonic: "RET  PO"},
	0xE1: {mnemonic: "POP  HL"},
	0xE2: {mnemonic: "JP   PO,$1$2", args: 2},
	0xE3: {mnemonic: "EX   (SP),HL"},
	0xE4: {mnemonic: "CALL PO,$1$2", args: 2},
	0xE5: {mnemonic: "PUSH HL"},
	0xE6: {mnemonic: "AND  $1", args: 1},
	0xE7: {mnemonic: "RST  20h"},
	0xE8: {mnemonic: "RET  PE"},
	0xE9: {mnemonic: "JP   (HL)"},
	0xEA: {mnemonic: "JP   PE,$1$2", args: 2},
	0xEB: {mnemonic: "EX   DE,HL"},
	0xEC: {mnemonic: "CALL PE,$1$2", args: 2},
	0xEE: {mnemonic: "XOR  $1", args: 1},
	0xEF: {mnemonic: "RST  28h"},
	0xF0: {mnemonic: "RET  P"},
	0xF1: {mnemonic: "POP  AF"},
	0xF2: {mnemonic: "JP   P,$1$2", args: 2},
	0xF3: {mnemonic: "DI"},
	0xF4: {mnemonic: "CALL P,$1$2", args: 2},
	0xF5: {mnemonic: "PUSH AF"},
	0xF6: {mnemonic: "OR   $1", args: 1},
	0xF7: {mnemonic: "RST  30h"},
	0xF8: {mnemonic: "RET  M"},
	0xF9: {mnemonic: "LD   SP,HL"},
	0xFA: {mnemonic: "JP   M,$1$2", args: 2},
	0xFB: {mnemonic: "EI"},
	0xFC: {mnemonic: "CALL M,$1$2", args: 2},
	0xFE: {mnemonic: "CP   $1", args: 1},
	0xFF: {mnemonic: "RST  38h"},
}