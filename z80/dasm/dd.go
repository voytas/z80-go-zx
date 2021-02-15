package dasm

// Table containing IX instructions (DD prefix)
var ddInstructions = []*instruction{
	0x09: {mnemonic: "ADD  IX,BC"},
	0x19: {mnemonic: "ADD  IX,DE"},
	0x21: {mnemonic: "LD   IX,$1$2", args: 2},
	0x22: {mnemonic: "LD   ($1$2),IX", args: 2},
	0x23: {mnemonic: "INC  IX"},
	0x24: {mnemonic: "INC  IXH"},
	0x25: {mnemonic: "DEC  IXH"},
	0x26: {mnemonic: "LD   IXH,$1", args: 1},
	0x29: {mnemonic: "ADD  IX,IX"},
	0x2A: {mnemonic: "LD   IX,($1$2)", args: 2},
	0x2B: {mnemonic: "DEC  IX"},
	0x2C: {mnemonic: "INC  IXL"},
	0x2D: {mnemonic: "DEC  IXL"},
	0x2E: {mnemonic: "LD   IXL,$1", args: 1},
	0x34: {mnemonic: "INC  (IX+$1)", args: 1},
	0x35: {mnemonic: "DEC  (IX)"},
	0x36: {mnemonic: "LD   (IX+$2),$1", args: 2},
	0x39: {mnemonic: "ADD  IX,SP"},
	0x44: {mnemonic: "LD   B,IXH"},
	0x45: {mnemonic: "LD   B,IXL"},
	0x46: {mnemonic: "LD   B,(IX+$1)", args: 1},
	0x4C: {mnemonic: "LD   C,IXH"},
	0x4D: {mnemonic: "LD   C,IXL"},
	0x4E: {mnemonic: "LD   C,(IX+$1)", args: 1},
	0x54: {mnemonic: "LD   D,IXH"},
	0x55: {mnemonic: "LD   D,IXL"},
	0x56: {mnemonic: "LD   D,(IX+$1)", args: 1},
	0x5C: {mnemonic: "LD   E,IXH"},
	0x5D: {mnemonic: "LD   E,IXL"},
	0x5E: {mnemonic: "LD   E,(IX+$1)", args: 1},
	0x60: {mnemonic: "LD   IXH,B"},
	0x61: {mnemonic: "LD   IXH,C"},
	0x62: {mnemonic: "LD   IXH,D"},
	0x63: {mnemonic: "LD   IXH,E"},
	0x64: {mnemonic: "LD   IXH,IXH"},
	0x65: {mnemonic: "LD   IXH,IXL"},
	0x66: {mnemonic: "LD   H,(IX+$1)", args: 1},
	0x67: {mnemonic: "LD   IXH,A"},
	0x68: {mnemonic: "LD   IXL,B"},
	0x69: {mnemonic: "LD   IXL,C"},
	0x6A: {mnemonic: "LD   IXL,D"},
	0x6B: {mnemonic: "LD   IXL,E"},
	0x6C: {mnemonic: "LD   IXL,IXH"},
	0x6D: {mnemonic: "LD   IXL,IXL"},
	0x6E: {mnemonic: "LD   L,(IX+$1)", args: 1},
	0x6F: {mnemonic: "LD   IXL,A"},
	0x70: {mnemonic: "LD   (IX+$1),B", args: 1},
	0x71: {mnemonic: "LD   (IX+$1),C", args: 1},
	0x72: {mnemonic: "LD   (IX+$1),D", args: 1},
	0x73: {mnemonic: "LD   (IX+$1),E", args: 1},
	0x74: {mnemonic: "LD   (IX+$1),H", args: 1},
	0x75: {mnemonic: "LD   (IX+$1),L", args: 1},
	0x77: {mnemonic: "LD   (IX+$1),A", args: 1},
	0x7C: {mnemonic: "LD   A,IXH"},
	0x7D: {mnemonic: "LD   A,IXL"},
	0x7E: {mnemonic: "LD   A,(IX+$1)", args: 1},
	0x84: {mnemonic: "ADD  A,IXH"},
	0x85: {mnemonic: "ADD  A,IXL"},
	0x86: {mnemonic: "ADD  A,(IX+$1)", args: 1},
	0x8C: {mnemonic: "ADC  A,IXH"},
	0x8D: {mnemonic: "ADC  A,IXL"},
	0x8E: {mnemonic: "ADC  A,(IX+$1)", args: 1},
	0x94: {mnemonic: "SUB  IXH"},
	0x95: {mnemonic: "SUB  IXL"},
	0x96: {mnemonic: "SUB  (IX+$1)", args: 1},
	0x9C: {mnemonic: "SBC  A,IXH"},
	0x9D: {mnemonic: "SBC  A,IXL"},
	0x9E: {mnemonic: "SBC  A,(IX+$1)", args: 1},
	0xA4: {mnemonic: "AND  IXH"},
	0xA5: {mnemonic: "AND  IXL"},
	0xA6: {mnemonic: "AND  (IX+$1)", args: 1},
	0xAC: {mnemonic: "XOR  IXH"},
	0xAD: {mnemonic: "XOR  IXL"},
	0xAE: {mnemonic: "XOR  (IX+$1)", args: 1},
	0xB4: {mnemonic: "OR   IXH"},
	0xB5: {mnemonic: "OR   IXL"},
	0xB6: {mnemonic: "OR   (IX+$1)", args: 1},
	0xBC: {mnemonic: "CP   IXH"},
	0xBD: {mnemonic: "CP   IXL"},
	0xBE: {mnemonic: "CP   (IX+$1)", args: 1},
	0xE1: {mnemonic: "POP  IX"},
	0xE3: {mnemonic: "EX   (SP),IX"},
	0xE5: {mnemonic: "PUSH IX"},
	0xE9: {mnemonic: "JP   (IX)"},
	0xF9: {mnemonic: "LD   SP,IX"},
}
