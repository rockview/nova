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
    "fmt"
    "bytes"
)

// Arithmetic/logic operation mnemonics
var alOP = [...]string{
    0: "COM",
    1: "NEG",
    2: "MOV",
    3: "INC",
    4: "ADC",
    5: "SUB",
    6: "ADD",
    7: "AND" }

// Shift mnemonics
var alSH = [...]string{
    0: "",
    1: "L",
    2: "R",
    3: "S" }

// Carry mnemonics
var alC = [...]string{
    0: "",
    1: "Z",
    2: "O",
    3: "C" }

// Skip mnemonics
var alSKIP = [...]string{
    0: "",
    1: "SKP",
    2: "SZC",
    3: "SNC",
    4: "SZR",
    5: "SNR",
    6: "SEZ",
    7: "SBN" }

// Memory reference without accumulator mnemonics
var mr0OP = [...]string{
    0: "JMP",
    1: "JSR",
    2: "ISZ",
    3: "DSZ" }

// Memory reference with accumulator mnemonics
var mr1OP = [...]string{
    0: "",
    1: "LDA",
    2: "STA",
    3: "" }

// I/O transfer operation mnemonics
var ioOP = [...]string{
    0: "NIO",
    1: "DIA",
    2: "DOA",
    3: "DIB",
    4: "DOB",
    5: "DIC",
    6: "DOC",
    7: "SKP" }

// IOT function mnemonics
var ioF = [...]string{
    0: "",
    1: "S",
    2: "C",
    3: "P" }

// IOT skip condition mnemonics
var ioT = [...]string{
    0: "BN",
    1: "BZ",
    2: "DN",
    3: "DZ" }

// IOT device code mnenmonics
var ioD = [...]string{
    000: "0",
    001: "MDV",
    002: "MMU",
    003: "MMU1",
    004: "4",
    005: "5",
    006: "MCAT",
    007: "MCAR",
    010: "TTI",
    011: "TTO",
    012: "PTR",
    013: "PTP",
    014: "RTC",
    015: "PLT",
    016: "CDR",
    017: "LPT",
    020: "DSK",
    021: "ADCV",
    022: "MTA",
    023: "DACV",
    024: "DCM",
    025: "25",
    026: "26",
    027: "27",
    030: "QTY",
    031: "IBM1",
    032: "IBM2",
    033: "DKP",
    034: "CAS",
    035: "CRC",
    036: "IPB",
    037: "IVT",
    040: "DPI",
    041: "DPO",
    042: "DIO",
    043: "DIOT",
    044: "MXM",
    045: "45",
    046: "MCAT1",
    047: "MCAR1",
    050: "TTI1",
    051: "TTO1",
    052: "PTR1",
    053: "PTP1",
    054: "RTC1",
    055: "PLT1",
    056: "CDR1",
    057: "LPT1",
    060: "DSK1",
    061: "ADCV1",
    062: "MTA1",
    063: "DACV1",
    064: "FPU1",
    065: "FPU2",
    066: "FPU",
    067: "67",
    070: "QTY1",
    071: "71",
    072: "72",
    073: "DKP1",
    074: "FPU1",
    075: "FPU2",
    076: "FPU",
    077: "CPU" }

// DisasmWord disassembles the ir instruction.
func DisasmWord(ir uint16) string {
    var operator bytes.Buffer
    var operands bytes.Buffer

    if ir&0100000 != 0 {
        // Arithmetic/Logic
        acs :=    (ir & 060000) >> 13
        acd :=    (ir & 014000) >> 11
        op :=     (ir & 003400) >> 8
        sh :=     (ir & 000300) >> 6
        c :=      (ir & 000060) >> 4
        skip :=   (ir & 000007) >> 0

        operator.WriteString(alOP[op])
        operator.WriteString(alC[c])
        operator.WriteString(alSH[sh])
        if ir&0000010 != 0 {
            // No load
            operator.WriteByte('#')
        }
        fmt.Fprintf(&operands, "%o,%o", acs, acd)
        if skip != 0 {
            fmt.Fprintf(&operands, ",%s", alSKIP[skip])
        }
    } else if ir&0060000 == 0060000 {
        // TODO: needs sorting out
        // I/O transfer
        acc :=    (ir & 0014000) >> 11
        op :=     (ir & 0003400) >> 8
        f :=      (ir & 0000300) >> 6
        device := (ir & 0000077) >> 0

        if device == 077 {
            // CPU; common abbreviations
            switch ir {
            case 0062677:
                // DICC 0,CPU
                operator.WriteString("IORST")
            case 0060177:
                // NIOS 0,CPU
                operator.WriteString("INTEN")
            case 0060277:
                // NIOC 0,CPU
                operator.WriteString("INTDS")
            case 0063077:
                // DOC 0,CPU
                operator.WriteString("HALT")
            default:
                if f == 0 {
                    switch op {
                    case 1:
                        // DIA x,CPU
                        operator.WriteString("READS")
                    case 3:
                        // DIB x,CPU 
                        operator.WriteString("INTA")
                    case 4:
                        // DOB x,CPU
                        operator.WriteString("MSKO")
                    }
                    if operator.Len() > 0 {
                        fmt.Fprintf(&operands, "%o", acc)
                    }
                }
            }
        } else if (device == 1) {
            // MDV; common abbreviations
            switch ir {
            case 0073301:
                operator.WriteString("MUL")
            case 0073101:
                operator.WriteString("DIV")
            }
        }

        if operator.Len() == 0 {
            // Not abbreviated; use long format
            operator.WriteString(ioOP[op])
            if op == 7 {
                // SKP
                operator.WriteString(ioT[f])
            } else {
                // Non-SKP
                operator.WriteString(ioF[f])
            }
            if (op == 0) || (op == 7) {
                // NIO or SKP
                fmt.Fprintf(&operands, "%s", ioD[device])
            } else {
                fmt.Fprintf(&operands, "%o,%s", acc, ioD[device])
            }
        }
    } else {
        // Memory reference
        if ir&0060000 == 0 {
            // Without accumulator
            op :=  (ir & 0014000) >> 11
            operator.WriteString(mr0OP[op])
        } else {
            // With accumulator
            op :=   (ir & 0060000) >> 13
            acc :=  (ir & 0014000) >> 11
            operator.WriteString(mr1OP[op])
            fmt.Fprintf(&operands, "%o,", acc)
        }

        index :=    (ir & 0001400) >> 8
        disp :=     (ir & 0000377) >> 0

        if ir&0002000 != 0 {
            // Indirect
            operands.WriteByte('@')
        }
        switch index {
        case 0:
            // Page zero
            fmt.Fprintf(&operands, "%o", disp)
        case 1:
            // PC relative
            fmt.Fprintf(&operands, ".%+o", int8(disp))
        case 2:
            // AC2 relative
            fmt.Fprintf(&operands, "%o,2", int8(disp))
        case 3:
            // AC3 relative
            fmt.Fprintf(&operands, "%o,3", int8(disp))
        }
    }

    return fmt.Sprintf("%-8s%s", operator.String(), operands.String())
}

// DisasmBlock disassembles the words slice from the start element up to, but excluding, the limit element.
func DisasmBlock(start, limit int, words []uint16) {
    for addr := start; addr < limit; addr++ {
        word := words[addr]
        fmt.Printf("%05o %06o  %s\n", addr, word, DisasmWord(word))
    }
}
