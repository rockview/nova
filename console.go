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
    n.cmd <- message{typ:cmdReset}
    ack := <-n.ack
    return int(ack.addr)
}

// Stop implements the console STOP function. The processor is stopped at the
// end of the current instruction. Stop has no effect if the processor is
// stopped. The current value of the program counter is returned.
func (n *Nova) Stop() int {
    n.cmd <- message{typ:cmdStop}
    ack := <-n.ack
    return int(ack.addr)
}

// Start implements the console START function. "addr" is loaded into the
// program counter and execution begins at that address. Start has no effect if
// the processor is running.
func (n *Nova) Start(addr int) {
    n.cmd <- message{typ:cmdStart, addr:uint16(addr)}
    <-n.ack
}

// Continue implements the console CONTINUE function. Execution resumes from the
// current machine state. Continue has no effect if the processor is running.
func (n *Nova) Continue() {
    n.cmd <- message{typ:cmdContinue}
    <-n.ack
}

// InstStep implements the console INST STEP function. The current instruction
// is executed and the processor is stopped. The current value of the program
// counter and an indication of whether a HALT instruction was executed are
// returned (0: no halt, 1: halt). If the processor is running, no instruction
// is executed and an error is returned.
func (n *Nova) InstStep() (pc, halt int, err error) {
    n.cmd <- message{typ:cmdInstStep}
    ack := <-n.ack
    if ack.typ == ackRunning {
        err = fmt.Errorf("processor running")
        return
    }
    pc = int(ack.addr)
    halt = int(ack.data)
    return
}

// TODO: implement and document
func (n *Nova) ProgramLoad() error {
    n.cmd <- message{typ:cmdProgramLoad}
    ack := <-n.ack
    if ack.typ == ackRunning {
        return fmt.Errorf("processor running")
    }
    return nil
}

// Deposit implements the console DEPOSIT function. The program counter is
// loaded with "addr", and "data" is stored in memory at the address specified
// by the program counter. If the processor is running, no store occurs and an
// error is returned.
func (n *Nova) Deposit(addr, data int) error {
    n.cmd <- message{typ:cmdDeposit, addr:uint16(addr), data:uint16(data)}
    ack := <-n.ack
    if ack.typ == ackRunning {
        return fmt.Errorf("processor running")
    }
    return nil
}

// DepositNext implements the console DEPOSIT NEXT function. The program counter
// is incremented and "data" is stored in memory at the address specified by the
// program counter. If the processor is running, no store occurs and an error is
// returned.
func (n *Nova) DepositNext(data int) error {
    n.cmd <- message{typ:cmdDepositNext, data:uint16(data)}
    ack := <-n.ack
    if ack.typ == ackRunning {
        return fmt.Errorf("processor running")
    }
    return nil
}

// Examine implements the console EXAMINE function. The program counter is
// loaded with "addr" and the contents of memory at the address specified by the
// program counter is returned. If the processor is running, the program counter
// is not modified and an error is returned.
func (n *Nova) Examine(addr int) (int, error) {
    n.cmd <- message{typ:cmdExamine, addr:uint16(addr)}
    ack := <-n.ack
    if ack.typ == ackRunning {
        return 0, fmt.Errorf("processor running")
    }
    return int(ack.data), nil

}

// ExamineNext implements the console EXAMINE NEXT function. The program counter
// is incremented and the contents of memory at the address specified by the
// program counter is returned. If the processor is running, the program counter
// is not modified and an error is returned.
func (n *Nova) ExamineNext() (int, error) {
    n.cmd <- message{typ:cmdExamineNext}
    ack := <-n.ack
    if ack.typ == ackRunning {
        return 0, fmt.Errorf("processor running")
    }
    return int(ack.data), nil
}

// Switches implements the console data switches function. The switch register
// is loaded with "data".
func (n *Nova) Switches(data int) {
    n.cmd <- message{typ:cmdSwitches, data:uint16(data)}
    <-n.ack
}

// IsRunning indicates whether to processor is currently running.
func (n *Nova) IsRunning() bool {
    n.cmd <- message{typ:cmdStatus}
    ack := <-n.ack
    return ack.typ == ackRunning
}

// LoadMemory copies the "words" slice to memory starting from "addr".  If the
// processor is running, memory remains unchanged and an error is returned.
func (n *Nova) LoadMemory(addr int, words []uint16) error {
    if n.IsRunning() {
        return fmt.Errorf("processor running")
    }
    copy(n.m[uint16(addr)&kAddrMask:], words)
    return nil
}

func (n *Nova) Trace(addr int) (int, error) {
    _, err := n.Examine(addr)   // Load PC
    if err != nil {
        return 0, err
    }
    for {
        ad, halt, err := n.InstStep()
        if err != nil {
            return 0, err
        }
        if halt == 1 {
            return int(ad), nil
        }
        state, _ := n.State()
        if err != nil {
            return 0, err
        }
        fmt.Println(state)
    }
}

// State returns the processor state as a string having the following format:
// PC IR  AC[0] AC[1] AC[2] AC[3]  C ION ; <disassembled IR>
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

// WaitForHalt waits for the processor to halt. If the processor fails to halt
// within the specified "timeout" period, an error is returned.
func (n *Nova) WaitForHalt(timeout time.Duration) error {
    select {
    case <- time.After(timeout):
        return fmt.Errorf("timed out")
    case <- n.halt:
        return nil
    }
}

const (
    cmdReset int = iota
    cmdStop
    cmdStart
    cmdContinue
    cmdInstStep
    cmdDeposit
    cmdDepositNext
    cmdExamine
    cmdExamineNext
    cmdProgramLoad
    cmdSwitches
    cmdStatus

    ackStopped
    ackRunning
)

type message struct {
    typ int
    addr uint16
    data uint16
}

func (n *Nova) initRun() {
    select {
    case <-n.halt:
    default:
    }
}

// CPU stopped; waiting for key
func (n *Nova) stopped() {
    for {
        cmd := <-n.cmd
        switch cmd.typ {
        case cmdReset:
            n.reset()
            n.ack <- message{typ:ackStopped, addr:n.pc}
        case cmdStop:
            n.ack <- message{typ:ackStopped, addr:n.pc}
        case cmdStart:
            n.initRun()
            n.pc = cmd.addr
            n.ack <- message{typ:ackRunning}
            return
        case cmdContinue:
            n.initRun()
            n.ack <- message{typ:ackRunning}
            return
        case cmdInstStep:
            var halt uint16
            if n.step() == cpuHalt {
                halt = 1
            }
            n.ack <- message{typ:ackStopped, addr:n.pc, data:halt}
        case cmdDeposit:
            n.pc = cmd.addr
            n.m[n.pc&kAddrMask] = cmd.data
            n.ack <- message{typ:ackStopped}
        case cmdDepositNext:
            n.pc++
            n.m[n.pc&kAddrMask] = cmd.data
            n.ack <- message{typ:ackStopped}
        case cmdExamine:
            n.pc = cmd.addr
            n.ack <- message{typ:ackStopped, data:n.m[n.pc&kAddrMask]}
        case cmdExamineNext:
            n.pc++
            n.ack <- message{typ:ackStopped, data:n.m[n.pc&kAddrMask]}
        case cmdSwitches:
            n.sr = cmd.data
            n.ack <- message{typ:ackStopped}
        case cmdProgramLoad:
            n.initRun()
            n.loadBootstrapLoader()
            n.ack <- message{typ:ackRunning}
            return
        case cmdStatus:
            n.ack <- message{typ:ackStopped}
        }
    }
}

// CPU running; run until key or halt
func (n *Nova) running() {
    for {
        select {
        case cmd := <-n.cmd:
            switch cmd.typ {
            case cmdReset:
                n.reset()
                n.ack <- message{typ:ackStopped, addr:n.pc}
                return
            case cmdStop:
                n.ack <- message{typ:ackStopped, addr:n.pc}
                return
            case cmdSwitches:
                n.sr = cmd.data
                n.ack <- message{typ:ackStopped}
            case cmdStart, cmdContinue, cmdInstStep, cmdDeposit, cmdDepositNext, cmdExamine, cmdExamineNext, cmdProgramLoad, cmdStatus:
                n.ack <- message{typ:ackRunning}
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
