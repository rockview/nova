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

const (
    ioNIO uint16 = iota
    ioDIA
    ioDOA
    ioDIB
    ioDOB
    ioDIC
    ioDOC
    ioSKP
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
    pc uint16                   // Program counter
    ac [4]uint16                // Accumulators
    flags uint                  // Processor flags
    m [k32K]uint16              // 32KW memory
    devices map[uint16]*device  // Devices
    intReq uint64               // Interrupting devices
    sr uint16                   // Switch register
    cmd chan message            // Console command
    ack chan message            // Console acknowlege
    halt chan struct{}          // Signals machine HALT
}

const (
    cpuRun int = iota
    cpuHalt
)

// NewNova creates a new instance of a nova processor. The processor is stopped
// and the interrupt on flag, the 16-bit priority mask, and all busy and done
// flags are set to 0.
func NewNova() *Nova {
    nova := Nova{
        devices: make(map[uint16]*device),
        cmd: make(chan message),
        ack: make(chan message),
        halt: make(chan struct{}),
    }
    go nova.processor()
    return &nova
}

// Execute one instruction.
func (n *Nova) step() int {
    // Handle data channel requests
    // Handle pending interrupts

    // Fetch next instruction
    ir := n.m[n.pc&kAddrMask]
    n.pc++

    // Decode and execute instruction
    if ir&0100000 != 0 {
        // Arithmetic/logic IR<0>
        var alu uint32

        // Initilize alu with carry IR<10,11>
        switch (ir&000060) >> 4 {
        case 0:
            if n.flags&cpuC != 0 {
                alu = 1 << 16
            }
        case 1: // Z
        case 2: // O
            alu = 1 << 16
        case 3: // C
            if n.flags&cpuC == 0 {
                alu = 1 << 16
            }
        }

        acs := n.ac[(ir&060000) >> 13]  // ACS<0-15>
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
            alu += uint32(n.ac[acx]) + uint32(^acs)
        case 5: // SUB
            alu += uint32(n.ac[acx]) + uint32(^acs) + 1
        case 6: // ADD
            alu += uint32(n.ac[acx]) + uint32(acs)
        case 7: // AND
            alu += uint32(n.ac[acx])&uint32(acs)
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
            n.pc++
        case 2: // SZC
            if c == 0 {
                n.pc++
            }
        case 3: // SNC
            if c == 1 {
                n.pc++
            }
        case 4: // SZR
            if alu == 0 {
                n.pc++
            }
        case 5: // SNR
            if alu != 0 {
                n.pc++
            }
        case 6: // SEZ
            if c == 0 || alu == 0 {
                n.pc++
            }
        case 7: // SBN
            if c == 1 && alu != 0 {
                n.pc++
            }
        }

        // Save result IR<12>
        if ir&000010 == 0 {
            n.ac[acx] = uint16(alu)
            if c == 1 {
                n.flags |= cpuC
            } else {
                n.flags &^= cpuC
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
            case ioNIO:
            case ioDIA: // READS
                n.ac[ac] = n.sr
            case ioDOA: 
            case ioDIB: // INTA
                n.ac[ac] = n.inta()
            case ioDOB: // MSKO
                n.msko(n.ac[ac])
            case ioDIC: // IORST
                n.reset()
            case ioDOC: // HALT
                halt = true
            case ioSKP:
                switch f {
                case 0: // BN
                    if n.flags&cpuION != 0 {
                        n.pc++
                    }
                case 1: // BZ
                    if n.flags&cpuION == 0 {
                        n.pc++
                    }
                case 2: // DN
                case 3: // DZ
                    n.pc++
                }
            }

            if op != 7 {
                switch f {
                case 0:
                case 1: // S
                    n.flags |= cpuION;
                case 2: // C
                    n.flags &^= cpuION;
                case 3: // P
                }
            }

            if halt {
                return cpuHalt
            }
        } else if device == 1 {
            // Pseudo device MDV
            if ac == 2 {
                switch op {
                case 6: // DOC
                    switch f {
                    case 1: // DOCS 2,MDV; DIV
                        if n.ac[0] >= n.ac[2] {
                            n.flags |= cpuC
                        } else {
                            dividend := uint32(n.ac[0]) << 16 | uint32(n.ac[1])
                            divisor := uint32(n.ac[2])
                            n.ac[1] = uint16(dividend/divisor)
                            n.ac[0] = uint16(dividend%divisor)
                            n.flags &^= cpuC
                        }
                    case 3: // DOCP 2,MDV; MUL
                        product := uint32(n.ac[1])*uint32(n.ac[2]) + uint32(n.ac[0])
                        n.ac[0] = uint16(product >> 16)
                        n.ac[1] = uint16(product)
                    }
                case 7: // SKP
                    switch f {
                    case 0: // BN
                    case 1: // BZ
                        n.pc++
                    case 2: // DN
                    case 3: // DZ
                        n.pc++
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
                    n.pc++
                case 2: // DN
                case 3: // DZ
                    n.pc++
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
            addr = n.pc - 1 + disp
        case 2: // AC2 relative
            addr = n.ac[2] + disp
        case 3: // AC3 relative
            addr = n.ac[3] + disp
        }

        // Handle indirect reference IR<5>
        if ir&002000 != 0 {
            for {
                next := n.m[addr&kAddrMask]
                bit0 := next&(1 << 15)
                if addr >= 020 && addr < 030 {
                    // Auto incrementing address
                    next++
                    n.m[addr] = next
                } else if addr >= 030 && addr < 040 {
                    // Auto decrementing address
                    next--
                    n.m[addr] = next
                }
                addr = next
                if bit0 == 0 {
                    break
                }
            }
        }

        // Perform operation
        if ir&060000 == 0 {
            // Without accumulator IR<3,4>
            switch (ir&014000) >> 11 {
            case 0: // JMP
                n.pc = addr
            case 1: // JSR
                n.ac[3] = n.pc
                n.pc = addr
            case 2: // ISZ
                n.m[addr&kAddrMask]++
                if n.m[addr&kAddrMask] == 0 {
                    n.pc++
                }
            case 3: // DSZ
                n.m[addr&kAddrMask]--
                if n.m[addr&kAddrMask] == 0 {
                    n.pc++
                }
            }
        } else {
            // With accumulator IR<1,2>
            acx := (ir&014000) >> 11
            switch (ir&060000) >> 13 {
            case 1:   // LDA
                n.ac[acx] = n.m[addr&kAddrMask]
            case 2:   // STA
                n.m[addr&kAddrMask] = n.ac[acx]
            }
        }
    }

    return cpuRun
}

// Reset the processor
func (n *Nova) reset() {
    n.flags &^= cpuION

    // Assert IORST
    for _, d := range n.devices {
        d.reset()
    }
}

// Assert MSKO
func (n *Nova) msko(mask uint16) {
    for _, d := range n.devices {
        d.msko(mask)
    }
}

// Assert INTA
func (n *Nova) inta() uint16 {
    var num uint16
    for _, d := range n.devices {
        if d.inta() {
            num = d.num
            break
        }
    }
    return num
}

func (n *Nova) setINTR(dev uint16) {
    n.intReq |= (1 << dev)
}

func (n *Nova) clrINTR(dev uint16) {
    n.intReq &^= (1 << dev)
}

// Simulate pressing the "PROGRAM LOAD" switch
func (n *Nova) loadBootstrapLoader() {
    var loader = [...]uint16{
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
        037: 0000000,
    }
    n.LoadMemory(0, loader[:])
    n.pc = 0
}
