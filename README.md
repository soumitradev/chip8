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
- Implement audio (needs a better understanding of goroutines and channels)