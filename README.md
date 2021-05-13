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
