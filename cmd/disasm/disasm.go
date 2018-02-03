package main

import (
    "os"
    "encoding/binary"
    "fmt"

    "github.com/rockview/nova"
)

func usage() {
    fmt.Fprintln(os.Stderr, "usage: disasm FILENAME")
    os.Exit(1)
}

func fatal(err error) {
    fmt.Fprintf(os.Stderr, "disasm: %v\n", err)
    os.Exit(1)
}

func main() {
    if len(os.Args) != 2 {
        usage()
    }

    filename := os.Args[1]
    f, err := os.Open(filename)
    if err != nil {
        fatal(err)
    }

    fi, err := f.Stat()
    if err != nil {
        fatal(err)
    }
    if fi.Size()%2 != 0 {
        fatal(fmt.Errorf("%s: odd size", filename))
    }

    words := make([]uint16, fi.Size()/2)
    err = binary.Read(f, binary.LittleEndian, &words)
    if err != nil {
        fatal(err)
    }

    nova.DisasmBlock(0, words)
}
