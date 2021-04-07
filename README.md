# CHIP8 Emulator

My first attempt at emulation. This Go project aims to emulate a fully working CHIP8 system complete with graphics, a working CPU, memory, sound, and keyboard input.

Everything works as far as I have tested, except for audio.

**TODO (In no particular order):**

- Test every CPU instruction:

    Not Tested/Implemented:
    - 0NNN
    - 00E0
    - 2NNN
    - DXYN
    - EX9E
    - EXA1
    - FX0A
    - FX18
    - FX29
        
- Comment code
- Optimize switch statement for CPU
- Test more ROMs
- Create a way to customize controls and settings (maybe using config files, maybe some GUI)
- Add more invasive features (e.g. Soft reset, Hard reset, memory manipulation, etc.)
- Implement audio (needs a better understanding of goroutines and channels)

### Credits

- Thanks to CowGod a.k.a. [Thomas P. Greene](mailto:cowgod@rockpile.com) for their extremely helpful chip8 documentation!

    [Link to Docs](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)

- Thanks to [@Skosulor](https://github.com/Skosulor) for their chip8 [testing ROM](./ROMs/TEST1.ch8) that helped me identify a mistake in the chip8 docs.

    [Link to Testing ROM](https://github.com/Skosulor/c8int/tree/master/test)


- Thanks to [@corax](https://github.com/corax89) for their chip8 [testing ROM](./ROMs/TEST.ch8) that helped me identify an error in my stack and instruction pointer.

    [Link to Testing ROM](https://github.com/corax89/chip8-test-rom)

- Thanks to [@kripod](https://github.com/kripod) for his collection of chip8 ROMs. Most of the non-testing ROMs are from here.

    [Link to Collection](https://github.com/kripod/chip8-roms)
