# ZX Spectrum Emulator

This is my fun project that is basically a ZX Spectrum 48k / 128k emulator. It uses Z80 CPU emulator I created previously. I wanted to see how easy or difficult it is to create some basic emulator.

Features implemented:
* 48k and 128k models are supported
* beeper support
* sna, szx, tap and minimal tzx file support
* work in progress on AY emulation
* memory congestion

This emulator is just a proof of concept and learning exercise. Emulation may not be 100% accurate.


## Screen
I couldn't find any simple graphics library for go that will allow rendering just 2D pixels, so for now I am using OpenGL (some very old version, because it has been deprecated on macOS).

## Memory
Memory paging for 128k model is implemented. Contended memory implemented using this page https://sinclair.wiki.zxnet.co.uk/wiki/Contended_memory rather than https://worldofspectrum.org/faq/reference/48kreference.htm.

## Keyboard


## Beeper
Using https://github.com/hajimehoshi/oto for playing sound.
Seems to be working mostly ok, but there is some issue with longer sound generation, e.g. BEEP 10,1 would have hearable frequency changes for some reason. It needs some investigating, but in games is not causing issues.
