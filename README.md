#

### Flags
* S - sign, typically set if result of an operation is negative
* Z - zero, typically set if result of an operation is zero
* H - half carry, typically set if result of an operation causes carry from lower nibble
* PV - parity or overflow
* C - carry, typically set if result of an operation exceeds the register size

### Half carry (H flag) calculation logic:
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
The same logic applies to subtraction since it also changes expected bit value in exactly the same way.

### Overflow condition
For addition, operands with different signs never cause overflow. When adding operands with similar signs and the result contains a different sign, the Overflow Flag is set.

For subtraction, overflow can occur for operands of unalike signs. Operands of alike signs never cause overflow. When subtracting operands with different signs and the result contains a different sign than the sign of the first operand, the Overflow Flag is set.

To check signs we can use xor operation (x ^ y) & 0x80. Value other than zero means signs of the x and y are different.
