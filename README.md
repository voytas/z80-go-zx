# Z80 CPU Emulator
This is a fun project I created to see how to write a CPU emulator and also to learn golang.
My first computer was Sinclair ZX Spectrum 48k so I decided to go for the Z80 CPU which I am quite familiar with.

The Z80 CPU emulates documented features and also as many as possible undocumented ones.

Having working Z80 emulator, I have created simple [ZX Spectrum emulator](spectrum/README.md).

The goal of this project is to have some fun and learn few things, not to create perfectly working feature rich emulator.
This is work in progress, although I may not be updating it on regular basis.

## Testing
Initial test was done using zexdoc/zexall instruction set exerciser (both pass). These tests (e.g. COM programs) are designed to run on CP/M system and required basic bootstrapping in order to execute them using emulator (some basic BDOS functions that output the test results to screen). You can find it in [exerciser](exerciser) folder.

You can execute tests using command line:

`go run ./main.go exercise ./exerciser/exercises/zexall.com`

This test is not exhaustive, but good for checking the basic CPU implementation.

## Dasm
There is a very basic disassembler in the [dasm](z80/dasm) folder. I used it during debugging and testing to output
the actual instruction being executed.
