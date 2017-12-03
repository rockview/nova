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

type stdWriter struct {
    controller
    w io.Writer
}

func newStdWriter(n *Nova, num, pri uint16, rate float32) *stdWriter {
    d := &stdWriter{
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

func (d *stdWriter) device(rate float32) {
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
            case ioDOA:
                // Load output register
                d.data = msg.data
                fallthrough
            case ioNIO, ioDIA, ioDIB, ioDOB, ioDIC, ioDOC:
                if msg.flags == ioS {
                    // Start device; delay until end of frame before write
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
            // Write to device
            expired = true
            if d.w != nil {
                b := []byte{byte(d.data)}
                if _, err := d.w.Write(b); err != nil {
                    panic(fmt.Sprintf("%s: %v", deviceName(d.num), err))
                }
            }
            d.complete()
        }
    }
}

func (d *stdWriter) attach(w io.Writer) {
    d.w = w
}
