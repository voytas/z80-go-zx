# ZX Spectrum Emulator

This is my fun project that is basically a ZX Spectrum 48k / 128k emulator. It uses Z80 CPU emulator I created previously. I wanted to see how easy or difficult it would be to an emulator. Simplicity of ZX Spectrum and fact that this is my favourite 8-bit machine means it was my choice.

Features implemented:
* 48k and 128k models are supported
* beeper support
* sna, szx, tap and minimal tzx file support
* work in progress on AY emulation
* memory congestion

This emulator is just a proof of a concept and learning exercise. Emulation may not be 100% accurate. Accuracy wasn't a goal, it requires some extra work that is very much hardware specific.

## Screen
It is using OpenGL although this is deprecated on macOS, but I needed something simple and I was unable to find anything else to render simple 2D pixel graphics.

## Memory
Memory paging for 128k model is implemented. Contended memory implemented using this page https://sinclair.wiki.zxnet.co.uk/wiki/Contended_memory rather than https://worldofspectrum.org/faq/reference/48kreference.htm.

## Keyboard


## Beeper
Using https://github.com/hajimehoshi/oto for playing sound.
Seems to be working mostly ok, but there is some issue with longer sound generation, for example BEEP 10,1 stutters occasionally. It needs some investigating, but in games beeper sounds ok.
