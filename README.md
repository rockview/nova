# Data General Nova Emulator
Emulator for Data General Nova 16-bit minicomputer.

The primary motivation for this project is to learn about the go programming
language, but I also have a fondness for this machine as I learned to program it
when I was an undergraduate computer studies student at Lancaster University in
the UK during the 1970s. It has been interesting the become reacquainted with
the architecture and instruction set after such a long time.

There is already a perfectly good Nova emulator that is part of the simh project
[here](http://simh.trailing-edge.com), but my project focuses on just emulating
the Nova and some of its devices using the go programming language.

Device support for a number of devices is the next phase of this project with
the ultimate goal of being able to run the RDOS operating system.

After this, I plan to teach myself Verilog to implement the Nova architecture
with an FPGA. Writing the software emulator first, seemed to be a good way
easing myself into, what is for me, is an entirely new field.
