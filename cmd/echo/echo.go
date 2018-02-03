package main

import (
    "os"
    "fmt"
    "time"
    "bytes"
    "io"

    "github.com/pkg/term"
    "github.com/rockview/nova"
)

func fatal(err error) {
    fmt.Fprintf(os.Stderr, "echo: %v\n", err)
    os.Exit(1)
}

func main() {
    var program = [...]uint16{
        0062677, // 00400:  IORST   
        0004417, // 00401:  JSR     .+17
        0020426, // 00402:  LDA     0,.+26
        0107400, // 00403:  AND     0,1
        0020425, // 00404:  LDA     0,.+25
        0106415, // 00405:  SUB#    0,1,SNR
        0063077, // 00406:  HALT    
        0004415, // 00407:  JSR     .+15
        0020422, // 00410:  LDA     0,.+22
        0106414, // 00411:  SUB#    0,1,SZR
        0000767, // 00412:  JMP     .-11
        0024420, // 00413:  LDA     1,.+20
        0004410, // 00414:  JSR     .+10
        0000764, // 00415:  JMP     .-14
        0000000, // 00416:
        0000000, // 00417:
        0063610, // 00420:  SKPDN   TTI      ; IN
        0000777, // 00421:  JMP     .-1
        0064610, // 00422:  DIAC    1,TTI
        0001400, // 00423:  JMP     0,3
        0063611, // 00424:  SKPDN   TTO      ; OUT
        0000777, // 00425:  JMP     .-1
        0064511, // 00426:  DIAS    1,TTO
        0001400, // 00427:  JMP     0,3
        0000177, // 00430:
        0000004, // 00431:
        0000015, // 00432:
        0000012, // 00433:
    }

    if true {
    t, err := term.Open("/dev/tty")
    if err != nil {
        fatal(err)
    }
    defer t.Close()

    err = term.RawMode(t)
    if err != nil {
        fatal(err)
    }
    defer t.Restore()

    b := make([]byte, 1)
    for {
        n, err := t.Read(b)
        if err == io.EOF {
            fmt.Println("got EOF")
        } else if err != nil {
            fmt.Printf("got error: %d, %v\n", n, err)
        } else {
            fmt.Printf("got char: %d, %c\n", n, b[0])
        }
       time.Sleep(time.Second * 3)
   }
}

    n := nova.NewNova()
    //n.Attach(nova.DevTTI, os.Stdin)
    //n.Attach(nova.DevTTO, os.Stdout)

    input := bytes.NewBufferString("abc\x04")
    var output bytes.Buffer

    n.Attach(nova.DevTTI, input)
    n.Attach(nova.DevTTO, output)

    n.LoadMemory(0400, program[:])
    //n.Trace(0400, 1000000)

    /*
    _, err := n.Examine(0400)   // Load PC
    if err != nil {
        fatal(err)
    }
    for i := 0; i < 10; i++ {
        state, _ := n.State()
        fmt.Println(state)
        ad, halt, err := n.InstStep()
        if err != nil {
            fatal(err)
        }
        if halt == 1 {
            return
        }
    }
    */

    addr, err := n.WaitForHalt(time.Second * 1)
    if err != nil {
        fatal(err)
    }
    if addr != 00406 {
        fatal(fmt.Errorf("wrong HALT address"))
    }

    fmt.Printf("output: %v\n", output)

    /*
    b := make([]byte, 1)
    for {
        _, err = t.Read(b)
        if err == io.EOF {
            continue
        }
        if err != nil {
            fatal(err)
        }
        if b[0] == 0x04 {   // ^D (EOT)
            break
        }
        _, err = t.Write(b)
        if err != nil {
            fatal(err)
        }
    }
    */
}
