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
    "testing"

    "time"
)

func TestDepositExamine(t *testing.T) {
    n := NewNova()

    for addr := 0; addr < k32K; addr++ {
        err := n.Deposit(addr, addr)
        if err != nil {
            t.Error(err)
        }
    }

    for addr := 0; addr < k32K; addr++ {
        data, err := n.Examine(addr)
        if err != nil {
            t.Error(err)
        }
        if data != addr {
            t.Errorf("have: %06o, want: %06o", addr, data)
        }
    }
}

func TestDepositNextExamineNext(t *testing.T) {
    n := NewNova()

    addr := 0
    err := n.Deposit(addr, addr)
    if err != nil {
        t.Error(err)
    }
    for addr++; addr < k32K; addr++ {
        err = n.DepositNext(addr)
        if err != nil {
            t.Error(err)
        }
    }

    // Wrap around
    for addr := 0; addr < k32K; addr++ {
        data, err := n.ExamineNext()
        if err != nil {
            t.Error(err)
        }
        if data != addr {
            t.Errorf("have: %06o, want: %06o", addr, data)
        }
    }
}

func TestLoadMemory(t *testing.T) {
    var words [k32K]uint16
    for addr := 0; addr < k32K; addr++ {
        words[addr] = uint16(addr)
    }
    n := NewNova()
    n.LoadMemory(0, words[:])

    for addr := 0; addr < k32K; addr++ {
        data, err := n.Examine(addr)
        if err != nil {
            t.Error(err)
        }
        if data != addr {
            t.Errorf("have: %06o, want: %06o", addr, data)
        }
    }
}

func TestExecution(t *testing.T) {
    n := NewNova()
    err := n.Deposit(1, 1) // 00001: JMP 1
    if err != nil {
        t.Error(err)
    }

    // Start/Stop
    n.Start(1)
    _, err = n.WaitForHalt(time.Millisecond * 10)
    if err == nil {
        t.Error("program: have: halt, want: timeout")
    }
    addr := n.Stop()
    if addr != 1 {
        t.Errorf("PC: have: %05o, want: %05o", addr, 1)
    }

    // Continue/Reset
    n.Continue()
    _, err = n.WaitForHalt(time.Millisecond * 10)
    if err == nil {
        t.Error("program: have: halt, want: timeout")
    }
    addr = n.Reset()
    if addr != 1 {
        t.Errorf("PC: have: %05o, want: %05o", addr, 1)
    }
}

func TestRunningKeys(t *testing.T) {
    program := [...]uint16 {
        0000001, // JMP 1
    }
    n := NewNova()
    startAddr := 1
    n.LoadMemory(startAddr, program[:])
    n.Start(startAddr)

    if !n.IsRunning() {
        t.Error("have: false, want: true")
    }

    _, _, err := n.InstStep()
    if err == nil {
        t.Error("have: nil, want: err")
    }
    err = n.Deposit(0, 0)
    if err == nil {
        t.Error("have: nil, want: err")
    }
    err = n.DepositNext(0)
    if err == nil {
        t.Error("have: nil, want: err")
    }
    _, err = n.Examine(0)
    if err == nil {
        t.Error("have: nil, want: err")
    }
    _, err = n.ExamineNext()
    if err == nil {
        t.Error("have: nil, want: err")
    }
    err = n.LoadMemory(startAddr, program[:])
    if err == nil {
        t.Error("have: nil, want: err")
    }

    n.Stop()
}

func TestInstStep(t *testing.T) {
    program := [...]uint16 {
        00000: 0000401, // JMP .+1
        00001: 0000401, // JMP .+1
        00002: 0063077, // HALT
    }
    n := NewNova()
    err := n.LoadMemory(0, program[:])
    if err != nil {
        t.Error(err)
    }

    n.Examine(0)    // Load PC
    for addr := 0; addr < 3; addr++ {
        pc, halt, err := n.InstStep()
        if err != nil {
            t.Error(err)
        }
        want := addr + 1
        if pc != want {
            t.Errorf("have: %05o, want: %05o", pc, want)
        }
        want = 3
        if halt == 1 && pc != want {
            t.Errorf("have: %05o, want: %05o", pc, want)
        }
    }
}
