# ZX Spectrum Emulator

This is my ZX Spectrum 48k / 128k emulator written in golang. It uses Z80 CPU emulator I created first. I wanted to see how easy and/or difficult it would be to create an emulator.

Features implemented:
* 48k and 128k models are supported
* beeper support
* sna, szx, tap and minimal tzx file support
* work in progress on AY emulation
* memory congestion (more or less accurate)

## Screen
It is using OpenGL although this is deprecated on macOS, but I needed something simple and I was unable to find anything else to render simple 2D pixel graphics. I may migrate it to some other framework if I can find something simple.

## Memory
Memory paging for 128k model is implemented. Contended memory implemented using this page https://sinclair.wiki.zxnet.co.uk/wiki/Contended_memory rather than https://worldofspectrum.org/faq/reference/48kreference.htm.

## Keyboard
For Shift use your left shift and for Symbol Shift use your right shift. PC specific keys like backspace, cursor keys, etc are not used at the moment.

## Beeper
Using https://github.com/hajimehoshi/oto for playing sound.
Seems to be working mostly ok, but there is some issue with longer sound generation, for example BEEP 10,1 stutters occasionally. It needs some investigating, but in games beeper sounds fine.
