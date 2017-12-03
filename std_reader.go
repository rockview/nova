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
    "time"
    "io"
    "fmt"
)

type stdReader struct {
    controller
    r io.Reader
}

func newStdReader(n *Nova, num, pri uint16, rate float32) *stdReader {
    d := &stdReader{
        controller: controller{
            num: num,
            pri: pri,
            dev: make(chan devmsg),
            n: n,
        },
    }
    go d.device(rate)
    return d
}

func (d *stdReader) device(rate float32) {
    period := time.Duration(float32(time.Second)/rate)
    t := time.NewTimer(time.Second * 1)
    t.Stop()
    expired := true
    for {
        select {
        case msg := <-d.dev:
            switch msg.typ {
            case ioRST:
                d.idle()
            case ioDIA:
                msg.data = d.data
                fallthrough
            case ioNIO, ioDOA, ioDIB, ioDOB, ioDIC, ioDOC:
                if msg.flags == ioS {
                    // Start device; delay until end of frame before read
                    if !t.Stop() && !expired {
                        <-t.C
                    }
                    t.Reset(period)
                    expired = false
                }
                d.flags(msg)
            case ioSKP:
                msg.data = d.skip(msg)
            default:
                panic(fmt.Sprintf("%s: invalid message type", deviceName(d.num)))
            }
            d.dev <- msg    // Ack
        case <-t.C:
            // Read from device
            expired = true
            if d.r != nil {
                b := make([]byte, 1)
                if _, err := d.r.Read(b); err != nil {
                    if err != io.EOF {
                        panic(fmt.Sprintf("%s: %v", deviceName(d.num), err))
                    }
                }
                d.data = uint16(b[0])
            }
            d.complete()
        }
    }
}

func (d *stdReader) attach(r io.Reader) {
    d.r = r
}
