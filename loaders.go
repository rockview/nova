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
    "io"
    "fmt"
)

// loadBootstrapLoader loads the bootstrap loader into memory and returns its
// start address.
func (n *Nova) loadBootstrapLoader() int {
    var program = [...]uint16{
        0062677,    // 00000: IORST   
        0060477,    // 00001: READS   0
        0024026,    // 00002: LDA     1,26
        0107400,    // 00003: AND     0,1
        0124000,    // 00004: COM     1,1
        0010014,    // 00005: ISZ     14
        0010030,    // 00006: ISZ     30
        0010032,    // 00007: ISZ     32
        0125404,    // 00010: INC     1,1,SZR
        0000005,    // 00011: JMP     5
        0030016,    // 00012: LDA     2,16
        0050377,    // 00013: STA     2,377
        0060077,    // 00014: ; (NOIS 0) - 1
        0101102,    // 00015: MOVL    0,0,SZC
        0000377,    // 00016: JMP     377
        0004030,    // 00017: JSR     30
        0101065,    // 00020: MOVC    0,0,SNR
        0000017,    // 00021: JMP     17
        0004027,    // 00022: JSR     27
        0046026,    // 00023: STA     1,@26
        0010100,    // 00024: ISZ     100
        0000022,    // 00025: JMP     22
        0000077,    // 00026: JMP     77
        0126420,    // 00027: SUBZ    1,1
        0063577,    // 00030: ; (SKPDN 0) - 1
        0000030,    // 00031: JMP     30
        0060477,    // 00032: ; (DIAS 0,0) - 1
        0107363,    // 00033: ADDCS   0,1,SNC
        0000030,    // 00034: JMP     30
        0125300,    // 00035: MOVS    1,1
        0001400,    // 00036: JMP     0,3
        0000000,    // 00037: JMP     0
    }
    n.LoadMemory(0, program[:])
    return 0
}

// LoadBinaryLoader loads the binary loader into memory and returns its
// start address.
func (n *Nova) LoadBinaryLoader() int {
    var program = [...]uint16{
        0177636,    // 77635:  
        0054512,    // 77636:  STA     3,.+112
        0004407,    // 77637:  JSR     .+7
        0171300,    // 77640:  MOVS    3,2
        0004405,    // 77641:  JSR     .+5
        0173300,    // 77642:  ADDS    3,2
        0143000,    // 77643:  ADD     2,0
        0002504,    // 77644:  JMP     @.+104
        0000004,    // 77645:  
        0054503,    // 77646:  STA     3,.+103
        0034503,    // 77647:  LDA     3,.+103
        0175103,    // 77650:  MOVL    3,3,SNC
        0000405,    // 77651:  JMP     .+5
        0063612,    // 77652:  SKPDN   PTR
        0000777,    // 77653:  JMP     .-1
        0074512,    // 77654:  DIAS    3,PTR
        0002474,    // 77655:  JMP     @.+74
        0063510,    // 77656:  SKPBZ   TTI
        0000777,    // 77657:  JMP     .-1
        0074510,    // 77660:  DIAS    3,TTI
        0002470,    // 77661:  JMP     @.+70
        0062677,    // 77662:  IORST   
        0060477,    // 77663:  READS   0
        0040466,    // 77664:  STA     0,.+66
        0060110,    // 77665:  NIOS    TTI
        0060112,    // 77666:  NIOS    PTR
        0004757,    // 77667:  JSR     .-21
        0171305,    // 77670:  MOVS    3,2,SNR
        0000776,    // 77671:  JMP     .-2
        0004754,    // 77672:  JSR     .-24
        0173300,    // 77673:  ADDS    3,2
        0141000,    // 77674:  MOV     2,0
        0145000,    // 77675:  MOV     2,1
        0004740,    // 77676:  JSR     .-40
        0050477,    // 77677:  STA     2,.+77
        0004736,    // 77700:  JSR     .-42
        0125113,    // 77701:  MOVL#   1,1,SNC
        0000426,    // 77702:  JMP     .+26
        0044450,    // 77703:  STA     1,.+50
        0030445,    // 77704:  LDA     2,.+45
        0034740,    // 77705:  LDA     3,.-40
        0172400,    // 77706:  SUB     3,2
        0034467,    // 77707:  LDA     3,.+67
        0136400,    // 77710:  SUB     1,3
        0172023,    // 77711:  ADCZ    3,2,SNC
        0000414,    // 77712:  JMP     .+14
        0030441,    // 77713:  LDA     2,.+41
        0147033,    // 77714:  ADDZ#   2,1,SNC
        0010436,    // 77715:  ISZ     .+36
        0147022,    // 77716:  ADDZ    2,1,SZC
        0125113,    // 77717:  MOVL#   1,1,SNC
        0004716,    // 77720:  JSR     .-62
        0052455,    // 77721:  STA     2,@.+55
        0010454,    // 77722:  ISZ     .+54
        0010430,    // 77723:  ISZ     .+30
        0000773,    // 77724:  JMP     .-5
        0101004,    // 77725:  MOV     0,0,SZR
        0063077,    // 77726:  HALT    
        0000740,    // 77727:  JMP     .-40
        0125224,    // 77730:  MOVZR   1,1,SZR
        0000411,    // 77731:  JMP     .+11
        0101004,    // 77732:  MOV     0,0,SZR
        0000773,    // 77733:  JMP     .-5
        0030442,    // 77734:  LDA     2,.+42
        0062677,    // 77735:  IORST   
        0151113,    // 77736:  MOVL#   2,2,SNC
        0001000,    // 77737:  JMP     0,2
        0063077,    // 77740:  HALT    
        0000777,    // 77741:  JMP     .-1
        0004704,    // 77742:  JSR     .-74
        0020404,    // 77743:  LDA     0,.+4
        0116404,    // 77744:  SUB     0,3,SZR
        0000775,    // 77745:  JMP     .-3
        0000721,    // 77746:  JMP     .-57
        0000377,    // 77747:  
        0000000,    // 77750:  
        0000000,    // 77751:  
        0000000,    // 77752:  
        0000000,    // 77753:  
        0000020,    // 77754:  
        0000000,    // 77755:  
        0000000,    // 77756:  
        0000000,    // 77757:  
        0000000,    // 77760:  
        0000000,    // 77761:  
        0000000,    // 77762:  
        0000000,    // 77763:  
        0000000,    // 77764:  
        0000000,    // 77765:  
        0000000,    // 77766:  
        0000000,    // 77767:  
        0000000,    // 77770:  
        0000000,    // 77771:  
        0000000,    // 77772:  
        0000000,    // 77773:  
        0000000,    // 77774:  
        0000000,    // 77775:  
        0000000,    // 77776:  
        0000663,    // 77777:  JMP     .-115
    }
    n.LoadMemory(077635, program[:])
    return 077777
}

// LoadAbsoluteBinary attaches the media to device dev and loads it into memory
// using the binary loader. The loaded program may HALT or auto start depending
// on the second word of the start block. Client programs should check if the
// processor is running after calling this function. If the program is halted,
// the program should be started at the address specified by the program
// documentation.
func (n *Nova) LoadAbsoluteBinary(dev int, media io.Reader) error {
    switch dev {
    case DevTTI:
        n.Switches(0000000)
    case DevPTR:
        n.Switches(0100000)
    default:
        return fmt.Errorf("invalid input device: %s", deviceName(uint16(dev)))
    }
    n.Attach(dev, media)

    startAddr := n.LoadBinaryLoader()
    n.Start(startAddr)

    return nil
}
