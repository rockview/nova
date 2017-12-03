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

import "io"

// Device codes
const (
    DevTTI = 010    // Teletype input
    DevTTO = 011    // Teletype output
    DevPTR = 012    // Paper tape reader
    DevPTP = 013    // Paper type punch
    DevMTA = 022    // Magnetic tape
    DevDKP = 033    // Moving head disk

    DevTTI1 = 050   // Second teletype input
    DevTTO1 = 051   // Second teletype output
    DevPTR1 = 052   // Second paper tape reader
    DevPTP1 = 053   // Second paper type punch

    devMDV = 001    // Multiply/divide
    devRTC = 014    // Real time clock
    devCPU = 077    // CPU
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

// I/O op codes
const (
    ioNIO = iota    // No I/O
    ioDIA           // Data In A
    ioDOA           // Data Out A
    ioDIB           // Data In B
    ioDOB           // Data Out B
    ioDIC           // Data In C
    ioDOC           // Data Out C
    ioSKP           // Skip
)

// I/O control flags
const (
    _ = iota
    ioS // Signal STRT
    ioC // Signal CLR
    ioP // Pulse
)

// I/O tests flags
const (
    ioBN = iota // Test busy non-zero
    ioBZ        // Test busy zero
    ioDN        // Test done non-zero
    ioDZ        // Test done zero
)

const (
    _ = iota + ioSKP
    ioRST       // Signal IORST
)

type driver interface {
    code() uint16
    priority() uint16
    reset()
    test(t uint16) bool
    read(op, f uint16) uint16
    write(op, f uint16, data uint16)
}

type inputDriver interface {
    driver
    attach(w io.Reader)
}

type outputDriver interface {
    driver
    attach(w io.Writer)
}

// Device state.
const (
    devIdle = iota
    devBusy
    devDone
)

// Device message.
type devmsg struct {
    typ uint16      // Message type
    flags uint16    // Control or test flags
    data uint16     // Other message data
}

// Device controller.
type controller struct {
    num uint16      // Device code
    pri uint16      // Priority
    state int       // State
    data uint16     // Transferred data
    dev chan devmsg // Device message channel
    n *Nova         // CPU
}

// code returns device code.
func (c *controller) code() uint16 {
    return c.num
}

// priority returns device priority.
func (c *controller) priority() uint16 {
    return c.pri
}

// reset performs device reset.
func (c *controller) reset() {
    c.dev <- devmsg{ioRST, 0, 0}
    <-c.dev
}

// test performs I/O SKP tests.
func (c *controller) test(t uint16) bool {
    c.dev <- devmsg{ioSKP, t, 0}
    ack := <-c.dev
    return ack.data == 1
}

// read performs the I/O read operation specified by op and sets the device
// state from the flags specified by f. The data read from the device, if any,
// is returned.
func (c *controller) read(op, f uint16) uint16 {
    c.dev <- devmsg{op, f, 0}
    ack := <-c.dev
    return ack.data
}

// write performs the I/O write operation specified by op on data and sets
// the device state from the flags specified by f.
func (c *controller) write(op, f uint16, data uint16) {
    c.dev <- devmsg{op, f, data}
    <-c.dev
}

// skip returns skip condition specified by message flags.
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

// idle puts the device into an idle state.
func (c *controller) idle() {
    c.state = devIdle
    c.n.clearInt(c.num)
}

// start commences I/O operation by putting the device into a busy state.
func (c *controller) start() {
    c.state = devBusy
    c.n.clearInt(c.num)
}

// complete completes the I/O operation by putting the device into a done state.
func (c *controller) complete() {
    if c.state == devBusy {
        c.state = devDone
        c.n.setInt(c.num)
    }
}

// flags sets the device state from the message flags.
func (c *controller) flags(msg devmsg) {
    switch msg.flags {
    case ioS:
        c.start()
    case ioC:
        c.idle()
    }
}

// deviceName returns the name of the device with the code num.
func deviceName(num uint16) string {
    return ioD[num&077]
}

// addDevices add all devices to the processor. Devices begin in an idle state
// with no media attached.
func (n *Nova) addDevices() {
    n.devices[DevTTI] = newStdReader(n, DevTTI, priTTI, 10) // ASR-33
    n.devices[DevTTO] = newStdWriter(n, DevTTO, priTTO, 10) // ASR-33
    n.devices[DevPTR] = newStdReader(n, DevPTR, priPTR, 300) // 4011B
    n.devices[DevPTP] = newStdWriter(n, DevPTP, priPTP, 63.3)
    n.devices[DevTTI1] = newStdReader(n, DevTTI1, priTTI, 10) // ASR-33
    n.devices[DevTTO1] = newStdWriter(n, DevTTO1, priTTO, 10) // ASR-33
    n.devices[DevPTR1] = newStdReader(n, DevPTR1, priPTR, 300) // 4011B
    n.devices[DevPTP1] = newStdWriter(n, DevPTP1, priPTP, 63.3)
    n.devices[devRTC] = newrtc(n, 60)
}
