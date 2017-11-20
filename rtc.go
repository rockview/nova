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

type devRTC struct {
    device
    lineFreq int
    freq int
}

func newRTC(n *Nova, lineFreq int) *devRTC {
    if lineFreq != 50 || lineFreq !=  60 {
        lineFreq = 60
    }
    return &devRTC{
        device: device{
            num: numRTC,
            mask: (1<<priRTC),
            n: n,
        },
        lineFreq: lineFreq,
        freq: lineFreq,
    }
}

func (d *devRTC) write(op, f, ac uint16) {
    if op == ioDOA {
        switch ac&03 {
        case 00:
            d.freq = d.lineFreq
        case 01:
            d.freq = 10
        case 02:
            d.freq = 100
        case 03:
            d.freq = 1000
        }
    }
}
