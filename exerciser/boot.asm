// Helper bootstrapper to load test .com file which uses couple CP/M functions
// to output the result to the screen. It is using OUT instruction to capture the
// output. This way zexdoc/zexall can be easily executed.

        device NOSLOT64K

        org $0000
        jp boot
        org $0005       // fixed CP/M bdos entry point
        jp bdos


        org $0100       // *.com programs will load at this address
prog:
        halt            // will be replaced with the *.com file


        org $F000       // Boot and bdos routies implementation
boot:
        nop             // will be modified to halt instruction to prevent infinite execution in case of jp 0
        ld hl,boot
        ld (hl),HALT    // modify nop to halt
        ld sp,$ffff     // initiate stack to top memory address
        jp prog         // and execute loaded program
bdos:
        ld a,c
        cp C_WRITE
        jr nz,is_write_str
        ld a,e
        out(PORT),a     // will be captured by the handler
        ret
is_write_str:
        cp C_WRITESTR
        ret nz
write_str:
        ld a,(de)
        cp '$'          // end of string?
        ret z
        out(PORT),a     // will be captured by the handler
        inc de
        jr write_str
        ret

C_WRITE = 2         // BDOS function 2 (C_WRITE) - Console output
C_WRITESTR = 9      // BDOS function 9 (C_WRITESTR) - Output string
PORT = 5            // port used to output text
HALT = $76          // halt instruction opcode

        // create boot.com file
        savebin "boot.com", $0000, $FFFF