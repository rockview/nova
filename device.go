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
    numTTI uint16 = 010 // Teletype input
    numTTO = 011        // Teletype output
    numPTR = 012        // Paper tape reader
    numPTP = 013        // Paper type punch
    numRTC = 014        // Real time clock
    numMTA = 022        // Magnetic tape
    numDKP = 033        // Moving head disk
    numCPU = 077        // CPU
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

// Device flags
const (
    devBUSY uint = 1<<iota  // Device busy
    devDONE                 // Device done
    devINTDIS               // Interrupt disabled
)

type device struct {
    num uint16      // Device code
    flags uint      // Status flags
    mask uint16     // Priority mask
    data uint16     // Read/write data
    n *Nova         // CPU
}

func (d *device) busy() bool {
    return d.flags&devBUSY != 0;
}

func (d *device) done() bool {
    return d.flags&devDONE != 0;
}

func (d *device) reset() {
    d.flags = 0
    d.n.clrINTR(d.num)
}

func (d *device) msko(mask uint16) {
    if d.mask&mask != 0 {
        d.flags |= devINTDIS
        d.update()
    }
}

func (d *device) inta() bool {
    return d.flags&devDONE != 0 && !(d.flags&devINTDIS == 0)
}

// State change update
func (d *device) update() {
    if d.inta() {
        d.n.setINTR(d.num)
    } else {
        d.n.clrINTR(d.num)
    }
}

// Apply flag control
func (d *device) control(f uint16) {
    switch f {
    case 0:
    case 1: // STRT
        d.flags |= devBUSY
        d.flags &^= devDONE
    case 2: // CLR
        d.flags &^= devBUSY
        d.flags |= devDONE
    case 3: // IOPLS
    }
    d.update()
}
