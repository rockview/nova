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

/*
Package nova provides a Data General Nova 16-bit minicomputer emulator. 

The architecture and programming of this machine is described here:
ftp://bitsavers.informatik.uni-stuttgart.de/pdf/dg/014-000631_NovaLinePgmg_Jul79.pdf.

This implementation supports the following devices:

TODO: list supported devices, memory, etc.

Note: when the current value of the program counter is returned by a console
function it will typically point to the next instruction to be executed.
*/
package nova

import (
    "fmt"
    "time"
)

// Reset implements the console RESET function. The processor is halted at the
// end of the current instruction. The Interrupt On flag, the 16-bit priority
// mask, and all Busy and Done flags are set to 0. Reset has no effect if the
// processor is stopped. The current value of the program counter is returned.
func (n *Nova) Reset() int {
    n.con <- conmsg{typ:conReset}
    con := <-n.con
    return int(con.addr)
}

// Stop implements the console STOP function. The processor is stopped at the
// end of the current instruction. Stop has no effect if the processor is
// stopped. The current value of the program counter is returned.
func (n *Nova) Stop() int {
    n.con <- conmsg{typ:conStop}
    con := <-n.con
    return int(con.addr)
}

// Start implements the console START function. addr is loaded into the program
// counter and execution begins at that address. Start has no effect if the
// processor is running.
func (n *Nova) Start(addr int) {
    n.con <- conmsg{typ:conStart, addr:uint16(addr)}
    <-n.con
}

// Continue implements the console CONTINUE function. Execution resumes from the
// current machine state. Continue has no effect if the processor is running.
func (n *Nova) Continue() {
    n.con <- conmsg{typ:conContinue}
    <-n.con
}

// InstStep implements the console INST STEP function. The current instruction
// is executed and the processor is stopped. The current value of the program
// counter and an indication of whether a HALT instruction was executed are
// returned. halt will have the value 1 if a halt was executed, 0 otheriwse. If
// the processor is running, no instruction is executed and an error is
// returned.
func (n *Nova) InstStep() (pc, halt int, err error) {
    n.con <- conmsg{typ:conInstStep}
    con := <-n.con
    if con.typ == conRunning {
        err = fmt.Errorf("processor running")
        return
    }
    pc = int(con.addr)
    halt = int(con.data)
    return
}

// TODO: implement and document
func (n *Nova) ProgramLoad() error {
    n.con <- conmsg{typ:conProgramLoad}
    con := <-n.con
    if con.typ == conRunning {
        return fmt.Errorf("processor running")
    }
    return nil
}

// Deposit implements the console DEPOSIT function. The program counter is
// loaded with addr, and data is stored in memory at the address specified by
// the program counter. If the processor is running, no store occurs and an
// error is returned.
func (n *Nova) Deposit(addr, data int) error {
    n.con <- conmsg{typ:conDeposit, addr:uint16(addr), data:uint16(data)}
    con := <-n.con
    if con.typ == conRunning {
        return fmt.Errorf("processor running")
    }
    return nil
}

// DepositNext implements the console DEPOSIT NEXT function. The program counter
// is incremented and data is stored in memory at the address specified by the
// program counter. If the processor is running, no store occurs and an error is
// returned.
func (n *Nova) DepositNext(data int) error {
    n.con <- conmsg{typ:conDepositNext, data:uint16(data)}
    con := <-n.con
    if con.typ == conRunning {
        return fmt.Errorf("processor running")
    }
    return nil
}

// Examine implements the console EXAMINE function. The program counter is
// loaded with addr and the contents of memory at the address specified by the
// program counter is returned. If the processor is running, the program counter
// is not modified and an error is returned.
func (n *Nova) Examine(addr int) (int, error) {
    n.con <- conmsg{typ:conExamine, addr:uint16(addr)}
    con := <-n.con
    if con.typ == conRunning {
        return 0, fmt.Errorf("processor running")
    }
    return int(con.data), nil

}

// ExamineNext implements the console EXAMINE NEXT function. The program counter
// is incremented and the contents of memory at the address specified by the
// program counter is returned. If the processor is running, the program counter
// is not modified and an error is returned.
func (n *Nova) ExamineNext() (int, error) {
    n.con <- conmsg{typ:conExamineNext}
    con := <-n.con
    if con.typ == conRunning {
        return 0, fmt.Errorf("processor running")
    }
    return int(con.data), nil
}

// Switches implements the console data switches function. The switch register
// is loaded with data.
func (n *Nova) Switches(data int) {
    n.con <- conmsg{typ:conSwitches, data:uint16(data)}
    <-n.con
}

// IsRunning indicates whether to processor is currently running.
func (n *Nova) IsRunning() bool {
    n.con <- conmsg{typ:conStatus}
    con := <-n.con
    return con.typ == conRunning
}

// LoadMemory copies the words slice to memory starting from addr.  If the
// processor is running, memory remains unchanged and an error is returned.
func (n *Nova) LoadMemory(addr int, words []uint16) error {
    if n.IsRunning() {
        return fmt.Errorf("processor running")
    }
    copy(n.m[addr&kAddrMask:], words)
    return nil
}

// Trace shows the processor state before each instuction is executed beginning
// with the instruction at addr. If the processor is running, no instruction is
// executed and an error is returned.
func (n *Nova) Trace(addr int) (int, error) {
    _, err := n.Examine(addr)   // Load PC
    if err != nil {
        return 0, err
    }
    for {
        state, _ := n.State()
        if err != nil {
            return 0, err
        }
        fmt.Println(state)
        ad, halt, err := n.InstStep()
        if err != nil {
            return 0, err
        }
        if halt == 1 {
            return int(ad), nil
        }
    }
}

// State returns the processor state as a string having the following format:
// PC IR  AC[0] AC[1] AC[2] AC[3]  C ION ; <disassembled IR>. Note: the state
// prior to the execution of the indicated instruction is returned.
func (n *Nova) State() (string, error) {
    if n.IsRunning() {
        return "", fmt.Errorf("processor running")
    }
    var carry int
    if n.flags&cpuC != 0 {
        carry = 1
    }
    var ion int
    if n.flags&cpuION != 0 {
        ion = 1
    }
    ir := n.m[n.pc&kAddrMask]
    return fmt.Sprintf("%05o %06o  %06o %06o %06o %06o  %d %d ; %s",
        n.pc, ir, n.ac[0], n.ac[1], n.ac[2], n.ac[3], carry, ion, DisasmWord(ir)), nil
}

// WaitForHalt waits for the processor to halt. If the processor halted within
// the timeout period the current value of the program counter is returned. An
// error is returned if the processor fails to halt within the timeout period.
func (n *Nova) WaitForHalt(timeout time.Duration) (int, error) {
    select {
    case <- time.After(timeout):
        return 0, fmt.Errorf("timed out")
    case <- n.halt:
        return int(n.pc), nil
    }
}

const (
    // Request
    conReset int = iota
    conStop
    conStart
    conContinue
    conInstStep
    conDeposit
    conDepositNext
    conExamine
    conExamineNext
    conProgramLoad
    conSwitches
    conStatus

    // Response
    conStopped
    conRunning
)

// Console message
type conmsg struct {
    typ int
    addr uint16
    data uint16
}

// Initialize the processor prior to running.
func (n *Nova) initRun() {
    select {
    case <-n.halt:
    default:
    }
}

// Processor stopped; waiting for key
func (n *Nova) stopped() {
    for {
        msg := <-n.con
        switch msg.typ {
        case conReset:
            n.reset()
            n.con <- conmsg{typ:conStopped, addr:n.pc}
        case conStop:
            n.con <- conmsg{typ:conStopped, addr:n.pc}
        case conStart:
            n.initRun()
            n.pc = msg.addr
            n.con <- conmsg{typ:conRunning}
            return
        case conContinue:
            n.initRun()
            n.con <- conmsg{typ:conRunning}
            return
        case conInstStep:
            var halt uint16
            if n.step() == cpuHalt {
                halt = 1
            }
            n.con <- conmsg{typ:conStopped, addr:n.pc, data:halt}
        case conDeposit:
            n.pc = msg.addr
            n.m[n.pc&kAddrMask] = msg.data
            n.con <- conmsg{typ:conStopped}
        case conDepositNext:
            n.pc++
            n.m[n.pc&kAddrMask] = msg.data
            n.con <- conmsg{typ:conStopped}
        case conExamine:
            n.pc = msg.addr
            n.con <- conmsg{typ:conStopped, data:n.m[n.pc&kAddrMask]}
        case conExamineNext:
            n.pc++
            n.con <- conmsg{typ:conStopped, data:n.m[n.pc&kAddrMask]}
        case conSwitches:
            n.sr = msg.data
            n.con <- conmsg{typ:conStopped}
        case conProgramLoad:
            n.initRun()
            n.loadBootstrapLoader()
            n.con <- conmsg{typ:conRunning}
            return
        case conStatus:
            n.con <- conmsg{typ:conStopped}
        }
    }
}

// Processor running; run until key or halt
func (n *Nova) running() {
    for {
        select {
        case msg := <-n.con:
            switch msg.typ {
            case conReset:
                n.reset()
                n.con <- conmsg{typ:conStopped, addr:n.pc}
                return
            case conStop:
                n.con <- conmsg{typ:conStopped, addr:n.pc}
                return
            case conSwitches:
                n.sr = msg.data
                n.con <- conmsg{typ:conStopped}
            case conStart, conContinue, conInstStep, conDeposit, conDepositNext, conExamine, conExamineNext, conProgramLoad, conStatus:
                n.con <- conmsg{typ:conRunning}
            }
        default:
            if n.step() == cpuHalt {
                n.halt <- struct{}{}
                return
            }
        }
    }
}

func (n *Nova) processor() {
    for {
        n.stopped()
        n.running()
    }
}
