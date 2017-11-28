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

// Device codes
const (
    numMDV = 001    // Multiply/divide
    numTTI = 010    // Teletype input
    numTTO = 011    // Teletype output
    numPTR = 012    // Paper tape reader
    numPTP = 013    // Paper type punch
    numRTC = 014    // Real time clock
    numMTA = 022    // Magnetic tape
    numDKP = 033    // Moving head disk
    numCPU = 077    // CPU
)

// Device priorities
const (
    priDKP = 7
    priMTA = 10
    priPTR = 11
    priRTC = 13
    priPTP = 13
    priTTI = 14
    priTTO = 15
)

// I/O instruction op codes
const (
    ioNIO = iota
    ioDIA
    ioDOA
    ioDIB
    ioDOB
    ioDIC
    ioDOC
    ioSKP
)

// I/O flags
const (
    _ = iota
    ioS // Start
    ioC // Idle
    ioP // Pulse
)

// I/O test
const (
    ioBN = iota // Test busy non-zero
    ioBZ        // Test busy zero
    ioDN        // Test done non-zero
    ioDZ        // Test done zero
)

const (
    _ = iota + ioSKP
    ioRST
)

type device interface {
    code() uint16
    priority() uint16
    reset()
    test(t uint16) bool
    read(op, f uint16) uint16
    write(op, f uint16, data uint16)
}

type devmsg struct {
    typ uint16
    flags uint16
    data uint16
}

// Device state
const (
    devIdle = iota
    devBusy
    devDone
)
type controller struct {
    num uint16      // Device code
    pri uint16      // Priority
    state int       // State
    dev chan devmsg // Device channel
    n *Nova         // CPU
}

func (c *controller) code() uint16 {
    return c.num
}

func (c *controller) priority() uint16 {
    return c.pri
}

// Reset device.
func (c *controller) reset() {
    c.dev <- devmsg{ioRST, 0, 0}
    <-c.dev
}

// Test device state.
func (c *controller) test(t uint16) bool {
    c.dev <- devmsg{ioSKP, t, 0}
    ack := <-c.dev
    return ack.data == 1
}

// Read word from device.
func (c *controller) read(op, f uint16) uint16 {
    c.dev <- devmsg{op, f, 0}
    ack := <-c.dev
    return ack.data
}

// Write word to device.
func (c *controller) write(op, f uint16, data uint16) {
    c.dev <- devmsg{op, f, data}
    <-c.dev
}

// Set skip condition.
func (c *controller) skip(msg devmsg) uint16 {
    var result uint16
    switch msg.flags {
    case ioBN:
        if c.state == devBusy {
            result = 1
        }
    case ioBZ:
        if c.state != devBusy {
            result = 1
        }
    case ioDN:
        if c.state == devDone {
            result = 1
        }
    case ioDZ:
        if c.state != devDone {
            result = 1
        }
    }
    return result
}

// Idle device.
func (c *controller) idle() {
    c.state = devIdle
    c.n.clearInt(c.num)
}

// Start I/O operation.
func (c *controller) start() {
    c.state = devBusy
    c.n.clearInt(c.num)
}

// Complete I/O operation.
func (c *controller) complete() {
    if c.state == devBusy {
        c.state = devDone
        c.n.setInt(c.num)
    } 
}

// Apply control flags.
func (c *controller) flags(msg devmsg) {
    switch msg.flags {
    case ioS:
        c.start()
    case ioC:
        c.idle()
    }
}
