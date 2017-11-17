// MIT License
// 
// Copyright 2017 Jeremy Hall
// 
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
// 
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package nova

import (
    "fmt"
)

const (
    k32K        = 1<<15
    kAddrMask   = k32K - 1
)

// CPU flags
const (
    cpuC uint   = 1<<iota   // Carry
    cpuION                  // Interrupts enabled
)

// CPU state
type Nova struct {
    pc uint16           // Program counter
    ac [4]uint16        // Accumulators
    flags uint          // Processor flags
    memory [k32K]uint16 // 32KW memory
    sr uint16           // Switch register
}

type Status int

const (
    Run Status = iota
    Halt
    IndirectLoop
)

func NewNova() *Nova {
    return &Nova{}
}

// Step one instruction
func (p *Nova) Step(trace bool) Status {
    if trace {
        fmt.Println(p.State())
    }

    // Fetch next instruction
    ir := p.memory[p.pc&kAddrMask]
    p.pc++

    // Decode and execute instruction
    if ir&0100000 != 0 {
        // Arithmetic/logic IR<0>
        var alu uint32

        // Initilize alu with carry IR<10,11>
        switch (ir&000060) >> 4 {
        case 0:
            if p.flags&cpuC != 0 {
                alu = 1 << 16
            }
        case 1: // Z
        case 2: // O
            alu = 1 << 16
        case 3: // C
            if p.flags&cpuC == 0 {
                alu = 1 << 16
            }
        }

        acs := p.ac[(ir&060000) >> 13]  // ACS<0-15>
        acx := (ir&014000) >> 11        // ACD index IR<1,2>

        // Perform operation IR<5-7>
        switch (ir&003400) >> 8 {
        case 0: // COM
            alu += uint32(^acs)
        case 1: // NEG
            alu += uint32(^acs) + 1
        case 2: // MOV
            alu += uint32(acs)
        case 3: // INC
            alu += uint32(acs) + 1
        case 4: // ADC
            alu += uint32(p.ac[acx]) + uint32(^acs)
        case 5: // SUB
            alu += uint32(p.ac[acx]) + uint32(^acs) + 1
        case 6: // ADD
            alu += uint32(p.ac[acx]) + uint32(acs)
        case 7: // AND
            alu += uint32(p.ac[acx])&uint32(acs)
        }

        // Perform shift IR<8,9> and extract carry
        var c uint32
        switch (ir&000300) >> 6 {
        case 0:
            c = (alu >> 16)&1
            alu = alu&0177777
        case 1: // L
            c = (alu >> 15)&1
            alu = (alu&077777) << 1 | (alu >> 16)&1
        case 2: // R
            c = alu&1
            alu = (alu >> 1)&0177777
        case 3: // S
            c = (alu >> 16)&1
            alu = (alu&0377) << 8 | (alu >> 8)&0377
        }

        // Perform skip IR<13-15>
        switch (ir&000007) >> 0 {
        case 0:
        case 1: // SKP
            p.pc++
        case 2: // SZC
            if c == 0 {
                p.pc++
            }
        case 3: // SNC
            if c == 1 {
                p.pc++
            }
        case 4: // SZR
            if alu == 0 {
                p.pc++
            }
        case 5: // SNR
            if alu != 0 {
                p.pc++
            }
        case 6: // SEZ
            if c == 0 || alu == 0 {
                p.pc++
            }
        case 7: // SBN
            if c == 1 && alu != 0 {
                p.pc++
            }
        }

        // Save result IR<12>
        if ir&000010 == 0 {
            p.ac[acx] = uint16(alu)
            if c == 1 {
                p.flags |= cpuC
            } else {
                p.flags &^= cpuC
            }
        }
    } else if ir&060000 == 060000 {
        // I/O transfer IR<1,2>
        ac :=     (ir&0014000) >> 11
        op :=     (ir&0003400) >> 8
        f :=      (ir&0000300) >> 6
        device := (ir&0000077) >> 0

        if device == 077 {
            // Pseudo device CPU
            var halt bool

            switch op {
            case 0: // NIO
            case 1: // DIA; READS
                p.ac[ac] = p.sr
            case 2: // DOA
            case 3: // DIB; INTA
                p.ac[ac] = 0    // TODO: set interrupting device code
            case 4: // DOB; MSKO
                // TODO: ac to priority mask
            case 5: // DIC; IORST
                // TODO: Busy/done flags cleared; priority mask = 0
            case 6: // DOC; HALT
                halt = true
            case 7: // SKP
                switch f {
                case 0: // BN
                    if p.flags&cpuION != 0 {
                        p.pc++
                    }
                case 1: // BZ
                    if p.flags&cpuION == 0 {
                        p.pc++
                    }
                case 2: // DN
                case 3: // DZ
                    p.pc++
                }
            }

            if op != 7 {
                switch f {
                case 0:
                case 1: // S
                    p.flags |= cpuION;
                case 2: // C
                    p.flags &^= cpuION;
                case 3: // P
                }
            }

            if halt {
                return Halt
            }
        } else if device == 1 {
            // Pseudo device MDV
            if ac == 2 {
                switch op {
                case 6: // DOC
                    switch f {
                    case 1: // DOCS 2,MDV; DIV
                        if p.ac[0] >= p.ac[2] {
                            p.flags |= cpuC
                        } else {
                            dividend := uint32(p.ac[0]) << 16 | uint32(p.ac[1])
                            divisor := uint32(p.ac[2])
                            p.ac[1] = uint16(dividend/divisor)
                            p.ac[0] = uint16(dividend%divisor)
                            p.flags &^= cpuC
                        }
                    case 3: // DOCP 2,MDV; MUL
                        product := uint32(p.ac[1])*uint32(p.ac[2]) + uint32(p.ac[0])
                        p.ac[0] = uint16(product >> 16)
                        p.ac[1] = uint16(product)
                    }
                case 7: // SKP
                    switch f {
                    case 0: // BN
                    case 1: // BZ
                        p.pc++
                    case 2: // DN
                    case 3: // DZ
                        p.pc++
                    }
                }
            }
        } else {
            // All other devices
            if op == 7 {
                // SKP
                switch f {
                case 0: // BN
                case 1: // BZ
                    p.pc++
                case 2: // DN
                case 3: // DZ
                    p.pc++
                }
            }
        }
    } else {
        // Memory reference
        addr := ir&000377
        disp := addr
        if disp > 0177 {
            disp -= 0400
        }

        // Compute effective address IR<6,7>
        switch (ir&001400) >> 8 {
        case 0: // Page zero
        case 1: // PC relative
            addr = p.pc - 1 + disp
        case 2: // AC2 relative
            addr = p.ac[2] + disp
        case 3: // AC3 relative
            addr = p.ac[3] + disp
        }

        // Handle indirect reference IR<5>
        if ir&002000 != 0 {
            if addr >= 020 && addr < 030 {
                // Auto incrementing address
                p.memory[addr]++
                addr = p.memory[addr]
            } else if addr >= 030 && addr < 040 {
                // Auto decrementing address
                p.memory[addr]--
                addr = p.memory[addr]
            } else {
                // Follow indirection chain
                var i int
                for {
                    addr = p.memory[addr&kAddrMask]
                    if addr&(1 << 15) == 0 {
                        break
                    }
                    i++
                    if i == k32K {
                        return IndirectLoop
                    }
                }
            }
        }

        // Perform operation
        if ir&060000 == 0 {
            // Without accumulator IR<3,4>
            switch (ir&014000) >> 11 {
            case 0: // JMP
                p.pc = addr
            case 1: // JSR
                p.ac[3] = p.pc
                p.pc = addr
            case 2: // ISZ
                p.memory[addr&kAddrMask]++
                if p.memory[addr&kAddrMask] == 0 {
                    p.pc++
                }
            case 3: // DSZ
                p.memory[addr&kAddrMask]--
                if p.memory[addr&kAddrMask] == 0 {
                    p.pc++
                }
            }
        } else {
            // With accumulator IR<1,2>
            acx := (ir&014000) >> 11
            switch (ir&060000) >> 13 {
            case 1:   // LDA
                p.ac[acx] = p.memory[addr&kAddrMask]
            case 2:   // STA
                p.memory[addr&kAddrMask] = p.ac[acx]
            }
        }
    }

    // Handle data channel requests
    // Handle pending interrupts

    return Run
}

// Run program at "addr"
func (p *Nova) Run(addr int, trace bool) {
    p.pc = uint16(addr)&kAddrMask
    for {
        switch p.Step(trace) {
        case Run:
            continue
        case Halt:
            fmt.Printf("program halt: %06o\n", (p.pc - 1) & kAddrMask)
        case IndirectLoop:
            fmt.Printf("infinite indirection\n")
        }
        break
    }
}

// Load program at "addr"
func (p *Nova) LoadMemory(addr int, words []uint16) {
    copy(p.memory[uint16(addr)&kAddrMask:], words)
}

// Simulate pressing the "RESET" switch
func (p *Nova) Reset() {
    p.flags &^= cpuION
    // IORST pulse
}

func (p *Nova) SetPC(addr uint16) {
    p.pc = addr&kAddrMask
}

func (p *Nova) SetSR(data uint16) {
    p.sr = data
}

// Simulate pressing the "PROGRAM LOAD" switch
func (p *Nova) ProgramLoad(device int, trace bool) {
    var code = []uint16{
        000: 0062677,
        001: 0060477,
        002: 0024026,
        003: 0107400,
        004: 0124000,
        005: 0010014,
        006: 0010030,
        007: 0010032,
        010: 0125404,
        011: 0000005,
        012: 0030016,
        013: 0050377,
        014: 0060077,
        015: 0101102,
        016: 0000377,
        017: 0004030,
        020: 0101065,
        021: 0000017,
        022: 0004027,
        023: 0046026,
        024: 0010100,
        025: 0000022,
        026: 0000077,
        027: 0126420,
        030: 0063577,
        031: 0000030,
        032: 0060477,
        033: 0107363,
        034: 0000030,
        035: 0125300,
        036: 0001400,
        037: 0000000 }
    addr := 0
    p.LoadMemory(addr, code)
    p.sr = uint16(device)
    p.Run(addr, trace)
}

func (p *Nova) State() string {
    var c int
    if p.flags&cpuC != 0 {
        c = 1
    }
    ir := p.memory[p.pc&kAddrMask]
    var ion int
    if p.flags&cpuION != 0 {
        ion = 1
    }
    return fmt.Sprintf("%05o %06o  %06o %06o %06o %06o  %d %d ; %s", p.pc, ir, p.ac[0], p.ac[1], p.ac[2], p.ac[3], c, ion, DisasmWord(ir))
}
