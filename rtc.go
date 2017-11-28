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

import "time"

type rtc struct {
    controller
}

func newRTC(n *Nova, lineFreq int) *rtc {
    if lineFreq != 50 || lineFreq !=  60 {
        lineFreq = 60
    }
    d := &rtc{
        controller: controller{
            num: numRTC,
            pri: priRTC,
            dev: make(chan devmsg),
            n: n,
        },
    }
    go d.device(lineFreq)
    return d
}

func (d *rtc) device(lineFreq int) {
    periods := [...]time.Duration {
        time.Second/time.Duration(lineFreq),
        time.Second/10,
        time.Second/100,
        time.Second/1000,
    }
    ticker := time.NewTicker(periods[0])
    for {
        select {
        case msg := <-d.dev:
            switch msg.typ {
            case ioRST:
                ticker.Stop()
                ticker = time.NewTicker(periods[0])
                d.idle()
                d.dev <- msg
            case ioDOA:
                ticker.Stop()
                ticker = time.NewTicker(periods[msg.data&03])
                fallthrough
            case ioNIO, ioDIA, ioDIB, ioDOB, ioDIC, ioDOC:
                d.flags(msg)
                d.dev <- msg
            case ioSKP:
                msg.data = d.skip(msg)
                d.dev <- msg
            }
        case <-ticker.C:
            d.complete()
        }
    }
}
