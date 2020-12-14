#

### Flags
* S - sign, typically set if result of an operation is negative
* Z - zero, typically set if result of an operation is zero
* H - half carry, typically set if result of an operation causes carry from lower nibble
* PV - parity or overflow
* C - carry, typically set if result of an operation exceeds the register size

### Half carry (H flag) calculation logic (a + b):
a7 a6 a5 a4 a3 a2 a1 a0
b7 b6 b5 b4 b3 b2 b1 b0
When there is no carry from bit 3 then bit 4 would normally be equal a4 + b4:
a4 b4 r4  a4^b4^r4
0  0   0     0
0  1   1     0
1  0   1     0
1  1   0     0
However when carry happened from bit 3 to bit 4 then result would be r4 + carry:
a4 b4 r4  a4^b4^r4
0  0   1     1
0  1   0     1
1  0   0     1
1  1   1     1
Using simple xor operation detects overflow situation. This will actually work for any bit, not just carry from bit 3 to bit 4.
