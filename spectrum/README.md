# ZX Spectrum Emulator

Once I had working Z80 CPU emulator, I decided to try to create a very simple ZX Spectrum 48k emulator. This is just to exercise programming in go and have fun.

I couldn't find any simple graphics library for go that will allow rendering just 2D pixels, so for now I am using OpenGL (some very old version, because it has been deprecated on macOS).

All that information you can find [here](https://worldofspectrum.org/faq/reference/48kreference.htm)

## Screen
In terms of ZX Spectrum 48kB, screen is refreshed 50 times per second. Basically you just need a timer that will run every 20ms (50 * 20ms = 1s). However if you want to also emulate border effects, some more work is needed. Because emulation runs much faster than the real machine, we need to keep track of alle border colour changes in relation to T state and use that information when screen is rendered.

## Memory


## Keyboard
That is quite easy to emulate, there are 40 keys divided into 8 groups of 5 keys each. One thing worth mentioning is that some programs can listen to more than one group to detect 'any key' pressed scenario, for example:
```assembly
        ld c,#02    ; any key except A-G
        in a,(c)
```

## Beeper
TODO:
