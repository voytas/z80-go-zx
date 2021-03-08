# Z80 CPU Emulator
This is just my fun project to see how to write CPU emulator in go. There is no real goal here apart
from having fun and learn go language a bit.

My first computer I had was Sinclair ZX Spectrum 48k and hence my first assembler programs were
written for Z80 CPU.

The implementation passes both zexdoc and zexall tests, however this does not mean that there are no bugs
or issues. I tried to implement most of the undocumented features as possible.

Now, once I had working Z80 emulator, I decided to go a bit further and see how can I emulate
my favourite ZX Spectrum computer.

## ZX Spectrum 48k emulator
Again this is just to have some fun, there are many excellent emulators out there already.
This is just a playground for me to learn how my favourite computer can be emulated.

The initial challenge I had was to find a way of displaying the actual emulated screen. There is no native
support for that in go. There are some gaming frameworks you can find, but I needed something very basic,
just to be able to draw pixels. I ended up with using OpenGL for now, although I would prefer something
more basic. And OpenGL technology is obsolete on Mac, some old 2.1 version only available, but it works,
at least for now.

## Testing
There are two popular programs usually used to test the implementation: zexdoc and zexdall. I believe
these programs were written for CP/M platform and require a little bootstrapping in order to execute
them.

In the exerciser folder there is a small exerciser command line utility that can be run the tests:

`go run ./main.go exercise ./exerciser/exercises/zexall.com`

It implements required CP/M bdos functions to output the test results.

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

However when carry happened from bit 3 to bit 4 then result would be r4 + carry:
| a4 | b4 | r4 | a4^b4^r4 |
|:--:|:--:|:--:|:--------:|
| 0  | 0  | 1  |    1     |
| 0  | 1  | 0  |    1     |
| 1  | 0  | 0  |    1     |
| 1  | 1  | 1  |    1     |

Using simple xor operation detects overflow situation. This will actually work for any bit, not just carry from bit 3 to bit 4.
The same logic applies to subtraction since it also changes expected bit value in exactly the same way.

### Overflow condition
For addition, operands with different signs never cause overflow. When adding operands with similar signs and the result contains a different sign, the Overflow Flag is set.

For subtraction, overflow can occur for operands of unalike signs. Operands of alike signs never cause overflow. When subtracting operands with different signs and the result contains a different sign than the sign of the first operand, the Overflow Flag is set.

To check signs we can use xor operation (x ^ y) & 0x80. Value other than zero means signs of the x and y are different.

### Parity
This is simple check to count bits and check whether the number is odd or event. I used simple lookup table
for this, mostly for the performance reasons.

### DAA
This is the most complicated instruction to implement. However if you follow "The Undocumented Z80 Documented",
it should be easy and it actually passes the test.