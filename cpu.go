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

import "sync"

const (
    k32K        = 1<<15
    kAddrMask   = k32K - 1
)

// CPU flags
const (
    cpuC uint   = 1<<iota   // Carry
    cpuION                  // Interrupts enabled
    cpuIONPending           // Set ION at start of next instruction
)

// CPU state
type Nova struct {
    pc uint16                   // Program counter
    ac [4]uint16                // Accumulators
    flags uint                  // Processor flags
    m [k32K]uint16              // 32KW memory

    devices map[uint16]driver   // Devices
    mu sync.Mutex
    interrupts uint64           // Interrupting devices
    intdisable uint64           // Interrupt disabled devices

    sr uint16                   // Switch register
    con chan conmsg             // Console channel
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
    n := &Nova{
        devices: make(map[uint16]driver),
        con: make(chan conmsg),
        halt: make(chan struct{}),
    }
    n.addDevices()
    go n.processor()
    return n
}

// Execute one instruction.
func (n *Nova) step() int {
    if (n.flags&cpuIONPending) != 0 {
        n.flags &^= cpuIONPending
        n.flags |= cpuION
    }

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
        ac :=  (ir&0014000) >> 11
        op :=  (ir&0003400) >> 8
        f :=   (ir&0000300) >> 6
        num := (ir&0000077) >> 0

        if num == devCPU {
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
                case ioBN:
                    if n.flags&cpuION != 0 {
                        n.pc++
                    }
                case ioBZ:
                    if n.flags&cpuION == 0 {
                        n.pc++
                    }
                case ioDZ:
                    n.pc++
                }
            }

            if op != ioSKP {
                switch f {
                case ioS:
                    if (n.flags&cpuION) == 0 {
                        n.flags |= cpuIONPending;
                    }
                case ioC:
                    n.flags &^= cpuION;
                }
            }

            if halt {
                return cpuHalt
            }
        } else if num == devMDV {
            // Pseudo device MDV
            if ac == 2 {
                switch op {
                case ioDOC:
                    switch f {
                    case ioS: // DOCS 2,MDV; DIV
                        if n.ac[0] >= n.ac[2] {
                            n.flags |= cpuC
                        } else {
                            dividend := uint32(n.ac[0]) << 16 | uint32(n.ac[1])
                            divisor := uint32(n.ac[2])
                            n.ac[1] = uint16(dividend/divisor)
                            n.ac[0] = uint16(dividend%divisor)
                            n.flags &^= cpuC
                        }
                    case ioP: // DOCP 2,MDV; MUL
                        product := uint32(n.ac[1])*uint32(n.ac[2]) + uint32(n.ac[0])
                        n.ac[0] = uint16(product >> 16)
                        n.ac[1] = uint16(product)
                    }
                case ioSKP:
                    switch f {
                    case ioBZ, ioDZ:
                        n.pc++
                    }
                }
            }
        } else {
            // All other devices
            dev := n.devices[num]
            if dev == nil {
                // Device not present
                switch op {
                case ioSKP:
                    switch f {
                    case ioBZ, ioDZ:
                        n.pc++
                    }
                }
            } else {
                switch op {
                case ioNIO, ioDIA, ioDIB, ioDIC:
                    n.ac[ac] = dev.read(op, f)
                case ioDOA, ioDOB, ioDOC:
                    dev.write(op, f, n.ac[ac])
                case ioSKP:
                    if dev.test(f) {
                        n.pc++
                    }
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
            addr = n.loadAddr(addr)
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

    // Handle data channel requests

    // Handle interrupts
    if (n.flags&cpuION) != 0 {
        n.mu.Lock()
        intrs := n.interrupts&^n.intdisable
        n.mu.Unlock()
        if intrs != 0 {
            // Disable interrupts and jump to ISR
            n.flags &^= cpuION
            n.m[0] = n.pc
            n.pc = n.loadAddr(1)
        }
    }

    return cpuRun
}

func (n *Nova) loadAddr(addr uint16) uint16 {
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
    return addr
}

// Reset the processor and devices
func (n *Nova) reset() {
    // Assert IORST
    for _, d := range n.devices {
        d.reset()
    }

    n.flags &^= cpuION
    n.mu.Lock()
    n.interrupts = 0
    n.mu.Unlock()
    n.intdisable = 0
}

// Assert MSKO
func (n *Nova) msko(mask uint16) {
    var flags uint64
    for _, d := range n.devices {
        if mask&(1 << d.priority()) != 0 {
            flags |= (1 << d.code())
        }
    }
    n.intdisable = flags
}

// Assert INTA
// TODO: define the order of query
func (n *Nova) inta() uint16 {
    n.mu.Lock()
    defer n.mu.Unlock()
    for _, d := range n.devices {
        if (n.interrupts&(1 << d.code())) != 0 {
            return d.code()
        }
    }
    return 0
}

func (n *Nova) setInt(num uint16) {
    n.mu.Lock()
    n.interrupts |= (1 << num)
    n.mu.Unlock()
}

func (n *Nova) clearInt(num uint16) {
    n.mu.Lock()
    n.interrupts &^= (1 << num)
    n.mu.Unlock()
}
