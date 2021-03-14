# Z80 CPU Emulator
This is just my fun project to see how to write some CPU emulator and also idea to learn go. My first computer was Sinclair ZX Spectrum 48k and hence my first assembler programs were written for the Z80 CPU which I most familiar with.

I tried to implement all documented features and as many as possible undocumented ones.

Now, once I have working Z80 emulator, I decided to go a bit further and see how can I emulate my favourite [ZX Spectrum computer](spectrum/README.md).

## Testing
The implementation passes both zexdoc and zexall tests. These tests were written for CP/M and required a little bootstrapping in order to execute them (some basic BDOS functions that output the text to screen).

In the exerciser folder there is a small command line utility that can be used to execute the tests:

`go run ./main.go exercise ./exerciser/exercises/zexall.com`

## Dasm
There is a very basic disassembler in the dasm folder. I used it during debugging and testing to output
the actual instruction being executed.

## Some implementation notes
### Flags
* S - sign, typically set if result of an operation is negative
* Z - zero, typically set if result of an operation is zero
* H - half carry, typically set if result of an operation causes carry from lower nibble
* P - parity or overflow
* C - carry, typically set if result of an operation exceeds the register size
* Y - undocumented flag, typically set to bit 5 or the result
* X - undocumented flag, typically set to bit 3 or the result

### Half carry (H flag) calculation logic:

|    |    |    |    |    |    |    |    |
|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|
| a7 | a6 | a5 | a4 | a3 | a2 | a1 | a0 |
| b7 | b6 | b5 | b4 | b3 | b2 | b1 | b0 |

When there is no carry from bit 3 then bit 4 would normally be equal a4 + b4:
| a4 | b4 | r4 | a4^b4^r4 |
|:--:|:--:|:--:|:--------:|
| 0  | 0  | 0  |    0     |
| 0  | 1  | 1  |    0     |
| 1  | 0  | 1  |    0     |
| 1  | 1  | 0  |    0     |

However when carry occurred from bit 3 to bit 4 then result would be r4 + carry:
| a4 | b4 | r4 | a4^b4^r4 |
|:--:|:--:|:--:|:--------:|
| 0  | 0  | 1  |    1     |
| 0  | 1  | 0  |    1     |
| 1  | 0  | 0  |    1     |
| 1  | 1  | 1  |    1     |

Using simple xor operation detects overflow situation. This will actually work for any bit, not just carry from bit 3 to bit 4. The same logic applies to subtraction since it also changes expected bit value in exactly the same way.

### Overflow condition
For addition, operands with different signs never cause overflow. When adding operands with similar signs and the result contains a different sign, the Overflow Flag is set.

For subtraction, overflow can occur for operands of unalike signs. Operands of alike signs never cause overflow. When subtracting operands with different signs and the result contains a different sign than the sign of the first operand, the Overflow Flag is set.

To check signs we can use xor operation (x ^ y) & 0x80. Value other than zero means signs of the x and y are different.

### Parity
This is simple check to count bits and check whether the number is odd or event. I used simple lookup table
for this, mostly for the performance reasons.

### DAA
This is the most complicated instruction to implement. However if you follow "The Undocumented Z80 Documented",
it should be easy and it actually works.