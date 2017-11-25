
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

    "testing"
)

func TestExerciser(t *testing.T) {
/*
 Exerciser, tape 097-000004-02.

 See: ftp://bitsavers.informatik.uni-stuttgart.de/pdf/dg/software/diag/097-000004-02_Exerciser_1971.pdf 

 The comments below match the labels in the program listing of the document above.

 Program execution flow:

         00002 ; start
 BEG:    00156
         00157 
         00160 
         00161 
 STA0:   01070 ; processor tests
         ...
         03115
 BEG1:   00162 ; start/stop devices
         ...
         00235
 MEMSZ:  00236 ; size memory
         ...
         00263
 STA0:   01070 ; start next pass
*/
    program := [...]uint16 {
        00001: 0000355,
        00002: 0002003,
        00003: 0000156,
        00004: 0063077,
        00005: 0000004,

        00040: 0000531, // ITAB
        00041: 0000610,
        00042: 0000510,
        00043: 0000410,
        00044: 0000644,
        00045: 0000004,
        00046: 0000004,
        00047: 0000004,
        00050: 0000775,
        00051: 0001034,
        00052: 0000754,
        00053: 0000651,

        00054: 0000000, // LSRB
        00055: 0000000, // LSRCH
        00056: 0000000, // KRAN
        00057: 0000000, // HSRB
        00060: 0000000, // HSRCH
        00061: 0000000, // READER
        00062: 0000000, // RRAN
        00063: 0000000, // PRAN
        00064: 0000000, // OOT
        00065: 0000377, // LT
        00066: 0000000, // PUNCH

        00070: 0000377, // RPDIF
        00071: 0000000, // C170
        00072: 0000000, // ACTION
        00073: 0135525, // C13552
        00074: 0000000, // RER
        00075: 0000177, // C177
        00076: 0000010, // C10

        00077: 0000565, // CXRAND
        00100: 0000000, // ISAV0
        00101: 0000000, // ISAV1
        00102: 0000000, // ISAV2
        00103: 0000000, // ISAV3
        00104: 0000000, // ISAVC
        00105: 0000000, // TRAN
        00106: 0000000, // TYPE
        00107: 0000000, // TEM
        00110: 0177773, // CM5
        00111: 0000004, // C4

        00112: 0003777, // MSIZE
        00113: 0000040, // C40
        00114: 0003140, // BUFF
        00115: 0003600, // FIN
        00116: 0003400, // FIN200
        00117: 0003340, // BUF200
        00120: 0000003, // C3
        00121: 0000005, // C5
        00122: 0000377, // C377
        00123: 0000000, // CFOO
        00124: 0177600, // C174X
        00125: 0001070, // MAIN
        00126: 0003000, // C3000

        00130: 0001000, // C1000
        00131: 0000200, // C200
        00132: 0000000, // STP
        00133: 0000000, // STT
        00134: 0000000, // TIN
        00135: 0000000, // TINEY

        00136: 0000037, // C37
        00137: 0000000, // .STP
        00140: 0000000, // .RRAN
        00141: 0000000, // .PRAN
        00142: 0000000, // .PUNCH
        00143: 0000000, // .READ
        00144: 0000000, // .HSRB
        00145: 0000000, // .HSRCH
        00146: 0000000, // .TINEY
        00147: 0000000, // .TIN
        00150: 0000000, // .LSRB
        00151: 0000000, // .LSRCH
        00152: 0000000, // .KRAN
        00153: 0000000, // .TYPE
        00154: 0000000, // .TRAN
        00155: 0000000, // .STT

        00156: 0102400, // BEG
        00157: 0040072,
        00160: 0062677,
        00161: 0002125,

        00162: 0030072, // BEG1
        00163: 0064477,
        00164: 0150000,
        00165: 0133502,
        00166: 0004276,
        00167: 0151102,
        00170: 0004307,
        00171: 0151102,
        00172: 0004315,
        00173: 0151102,
        00174: 0004267,
        00175: 0151102,
        00176: 0004264,
        00177: 0151102,
        00200: 0004340,
        00201: 0151102,
        00202: 0004351,
        00203: 0151102,
        00204: 0004324,
        00205: 0151102,
        00206: 0004331,

        00207: 0030072,
        00210: 0044072,
        00211: 0124000,
        00212: 0133502,
        00213: 0004301,
        00214: 0151102,
        00215: 0004305,
        00216: 0151102,
        00217: 0004313,
        00220: 0151102,
        00221: 0004272,
        00222: 0151102,
        00223: 0060214,
        00224: 0151102,
        00225: 0004343,
        00226: 0101001,
        00227: 0000000,
        00230: 0151102,
        00231: 0004347,
        00232: 0151102,
        00233: 0004322,
        00234: 0151102,
        00235: 0004334,

        00236: 0024126, // MEMSZ
        00237: 0020130,
        00240: 0107000,
        00241: 0044112,
        00242: 0125112,
        00243: 0000250,
        00244: 0046112,
        00245: 0032112,
        00246: 0132415,
        00247: 0000237,
        00250: 0014112,
        00251: 0020131,
        00252: 0024112,
        00253: 0106400,
        00254: 0044115,
        00255: 0106400,
        00256: 0044116,
        00257: 0020124,
        00260: 0024072,
        00261: 0107404,
        00262: 0060177,
        00263: 0002125,

        00264: 0020120, // RTCSU
        00265: 0061114,
        00266: 0001400,

        00267: 0102400, // TTOSU
        00270: 0061111,
        00271: 0102401,
        00272: 0102000, // TTOSD
        00273: 0040106,
        00274: 0040133,
        00275: 0001400,

        00276: 0102400, // PTPSU
        00277: 0061113,
        00300: 0102401,
        00301: 0102000, // PTPSD
        00302: 0040066,
        00303: 0040132,
        00304: 0001400,

        00305: 0060212, // PTRSD
        00306: 0101011,
        00307: 0060112, // PTRSU
        00310: 0102400,
        00311: 0040061,
        00312: 0001400,

        00313: 0010133, // TTISD
        00314: 0001400,
        00315: 0060110, // TTISU
        00316: 0102400,
        00317: 0040134,
        00320: 0040135,
        00321: 0001400,

        00322: 0010146, // .TISD
        00323: 0001400,
        00324: 0060150, // .TISU
        00325: 0102400,
        00326: 0040147,
        00327: 0040146,
        00330: 0001400,

        00331: 0102400, // .TOSU
        00332: 0061151,
        00333: 0102401,
        00334: 0102000, // .TOSD
        00335: 0040153,
        00336: 0040155,
        00337: 0001400,

        00340: 0102400, // .TPSU
        00341: 0061153,
        00342: 0102401,
        00343: 0102000, // .TPSD
        00344: 0040142,
        00345: 0040137,
        00346: 0001400,

        00347: 0060252, // .TRSD
        00350: 0101011,
        00351: 0060152, // .TRSU
        00352: 0102400,
        00353: 0040143,
        00354: 0001400,

        // TODO: 00355-01067

        00355: 0,
        00356: 0,
        00357: 0,
        00360: 0,
        00361: 0,
        00362: 0,
        00363: 0,
        00364: 0,
        00365: 0,
        00366: 0,
        00367: 0,
        00370: 0,
        00371: 0,
        00372: 0,
        00373: 0,
        00374: 0,
        00375: 0,
        00376: 0,
        00377: 0,

        01070: 0030114, // STA0
        01071: 0034115,
        01072: 0102400,
        01073: 0041000,

        01074: 0025000, // STA01
        01075: 0045001,
        01076: 0151400,
        01077: 0156014,
        01100: 0000774,
        01101: 0125014,
        01102: 0063077,
        01103: 0101014,
        01104: 0063077,

        01105: 0030114, // STA02
        01106: 0034115,
        01107: 0102000,
        01110: 0041000,

        01111: 0025000, // STA03
        01112: 0045001,
        01113: 0151400,
        01114: 0156014,
        01115: 0000774,

        01116: 0124014,
        01117: 0063077,
        01120: 0100014,
        01121: 0063077,

        01122: 0030114, // STA04
        01123: 0034115,
        01124: 0020420,
        01125: 0104000,
        01126: 0045000,
        01127: 0041001,

        01130: 0021000, // STA05
        01131: 0151400,
        01132: 0156015,
        01133: 0000774,

        01134: 0021000, // STA06
        01135: 0025377,
        01136: 0123020,
        01137: 0100014,
        01140: 0063077,
        01141: 0101012,
        01142: 0063077,
        01143: 0101011,
        01144: 0052525, // C52525

        01145: 0102400, // STA10
        01146: 0034114,
        01147: 0024115,
        01150: 0055400, // STA11
        01151: 0031400,
        01152: 0156414,
        01153: 0063077,
        01154: 0175400,
        01155: 0136414,
        01156: 0000772,

        01157: 0034114, // STA12
        01160: 0031400,
        01161: 0156414,
        01162: 0063077,
        01163: 0175400,
        01164: 0136414,
        01165: 0000773,

        01166: 0101014, // STA13
        01167: 0063077,

        01170: 0034114, // STA20
        01171: 0024115,

        01172: 0160005, // STA21
        01173: 0063077,
        01174: 0041400,
        01175: 0031400,
        01176: 0142414,
        01177: 0063077,
        01200: 0175400,
        01201: 0136414,
        01202: 0000770,

        01203: 0034114, // STA22
        01204: 0031400,
        01205: 0140000,
        01206: 0116414,
        01207: 0063077,
        01210: 0175400,
        01211: 0136414,
        01212: 0000772,

        01213: 0101010,

        01214: 0102520, // JSR10
        01215: 0034114,
        01216: 0030421,
        01217: 0024115,
        01220: 0051400,
        01221: 0175400,
        01222: 0166414,
        01223: 0000775,

        01224: 0030114, // JSR11
        01225: 0005000,
        01226: 0063077,
        01227: 0156014,

        01230: 0063077,
        01231: 0156014,
        01232: 0102400,
        01233: 0113000,
        01234: 0146414,
        01235: 0000770,
        01236: 0101011,
        01237: 0005401, // CJSR10

        01240: 0126520, // JMP10
        01241: 0020115,
        01242: 0030114,
        01243: 0034413,
        01244: 0001000,
        01245: 0063077,
        01246: 0156014, // JMP11
        01247: 0063077,
        01250: 0156014,
        01251: 0126400,
        01252: 0133000,
        01253: 0142414,
        01254: 0000767,
        01255: 0101011,
        01256: 0001245, // CJMP11

        01257: 0034114, // POSX
        01260: 0024115,
        01261: 0055400,
        01262: 0175400,
        01263: 0136414,
        01264: 0000775,
        01265: 0030426,
        01266: 0050411,

        01267: 0020120, // POSX1

        01270: 0024407,
        01271: 0107000,
        01272: 0044405,
        01273: 0020421,
        01274: 0106032,
        01275: 0000420,
        01276: 0034114,
        01277: 0000000, // POSX2

        01300: 0024777,
        01301: 0020122,
        01302: 0107400,
        01303: 0167000,
        01304: 0146414,
        01305: 0063077,
        01306: 0020116,
        01307: 0175400,
        01310: 0116414,
        01311: 0000766,
        01312: 0000755,

        01313: 0031400, // CLDAPX
        01314: 0031577, // POSFIN
        01315: 0101010, // POSND

        01316: 0034114, // NEGX
        01317: 0024115,
        01320: 0055400,
        01321: 0175400,
        01322: 0136414,
        01323: 0000775,
        01324: 0030431,
        01325: 0050412,

        01326: 0020120, // NEGX1
        01327: 0024410,
        01330: 0106400,
        01331: 0044406,
        01332: 0020424,
        01333: 0106033,
        01334: 0000423,
        01335: 0000401,

        01336: 0030117,
        01337: 0000000, // NEGX2
        01340: 0024777,
        01341: 0020122,
        01342: 0107400,
        01343: 0100000,
        01344: 0107000,
        01345: 0147000,
        01346: 0136414,
        01347: 0063077,
        01350: 0020115,
        01351: 0151400,
        01352: 0112414,
        01353: 0000764,
        01354: 0000752,

        01355: 0035377, // NEGFIN
        01356: 0035200, // NEGXX
        01357: 0101010, // NEGND

        01360: 0034114, // POSJ
        01361: 0024115,
        01362: 0030435,
        01363: 0051400,
        01364: 0175400,
        01365: 0166414,
        01366: 0000774,
        01367: 0020427,
        01370: 0040416,

        01371: 0020121, // POSJ1
        01372: 0024414,
        01373: 0107000,
        01374: 0044412,
        01375: 0020423,
        01376: 0106032,
        01377: 0000423,

        01400: 0030114, // POSJ2
        01401: 0020116,
        01402: 0151400,
        01403: 0142415,
        01404: 0000765,
        01405: 0034414,
        01406: 0000000, // POSJ3
        01407: 0024777,
        01410: 0020122,
        01411: 0107400,
        01412: 0147000,
        01413: 0136014,
        01414: 0063077,
        01415: 0000764,

        01416: 0001000, // POSJZ
        01417: 0005400, // POSJX
        01420: 0001177, // POSJF
        01421: 0001407, // POSJY
        01422: 0101010, // POSJND

        01423: 0030432, // NEGJ
        01424: 0050415,

        01425: 0020120, // NEGJ1
        01426: 0024413,
        01427: 0106400,
        01430: 0044411,
        01431: 0020423,
        01432: 0106052,
        01433: 0000423,

        01434: 0030117, // NEGJ2
        01435: 0020115,
        01436: 0151440,
        01437: 0142435,
        01440: 0000765,

        01441: 0000000, // NEGJ3
        01442: 0024777,
        01443: 0020122,
        01444: 0107403,
        01445: 0063077,
        01446: 0100000,
        01447: 0107000,
        01450: 0147000,
        01451: 0136014,
        01452: 0063077,
        01453: 0000762,

        01454: 0005200, // NEGJX
        01455: 0005377, // NEGJY
        01456: 0101010, // NEGJD

        01457: 0034114, // PINC
        01460: 0024115,
        01461: 0020415,
        01462: 0041400, // PINC1
        01463: 0175400,
        01464: 0136414,
        01465: 0000775,
        01466: 0020411,
        01467: 0041400,
        01470: 0030410, // PINC2
        01471: 0034114,
        01472: 0001400,

        01473: 0136414, // PINC3
        01474: 0063077,
        01475: 0000404, // CINC
        01476: 0175400, // CPINC
        01477: 0001000, // PINCR

        01500: 0001473,
        01501: 0101000,

        01502: 0102520, // ISZ0
        01503: 0024122,
        01504: 0152000,
        01505: 0176400,
        01506: 0054107, // ISZ1
        01507: 0060377,
        01510: 0010107,
        01511: 0000404,
        01512: 0020107,
        01513: 0063077,
        01514: 0102400,
        01515: 0106404,
        01516: 0000770,
        01517: 0150014,
        01520: 0063077,

        01521: 0102520, // ISZ2
        01522: 0024122,
        01523: 0152400,
        01524: 0176120,
        01525: 0054107, // ISZ3
        01526: 0060377,
        01527: 0010107,
        01530: 0000404,
        01531: 0020107,
        01532: 0063077,
        01533: 0102401,
        01534: 0106404,
        01535: 0000770,
        01536: 0151014,
        01537: 0063077,

        01540: 0102520, // ISZ4
        01541: 0126620,
        01542: 0152000,
        01543: 0050107, // ISZ5
        01544: 0060377,
        01545: 0010107,
        01546: 0000404,
        01547: 0107004,
        01550: 0000773,
        01551: 0000405,
        01552: 0020107, // ISZ6
        01553: 0063077,
        01554: 0102400,
        01555: 0000766,
        01556: 0101234, // ISZ7
        01557: 0063077,

        01560: 0152520, // ISZ10
        01561: 0102620,
        01562: 0143005,
        01563: 0000422,
        01564: 0040107,
        01565: 0060377,
        01566: 0010107,
        01567: 0101011,
        01570: 0000410,
        01571: 0024107,
        01572: 0106015,
        01573: 0000767,
        01574: 0115400, // ISZ11
        01575: 0063077,
        01576: 0152400,
        01577: 0000763,
        01600: 0115405, // ISZ12
        01601: 0000410,
        01602: 0024107,
        01603: 0063077,
        01604: 0000772,
        01605: 0102000, // ISZ13
        01606: 0024107,
        01607: 0063077,
        01610: 0000766,
        01611: 0101010, // ISZ14

        01612: 0030121, // ISZ20
        01613: 0102400,
        01614: 0034121,
        01615: 0117022,
        01616: 0000417,
        01617: 0040107, // ISZ21
        01620: 0060377,
        01621: 0010107,
        01622: 0010107,
        01623: 0010107,
        01624: 0010107,
        01625: 0010107,
        01626: 0024107,
        01627: 0136414,
        01630: 0063077,
        01631: 0136414,
        01632: 0152400,
        01633: 0143000,
        01634: 0000760,
        01635: 0101010, // ISZ22

        01636: 0102520, // ISZ30
        01637: 0024115,
        01640: 0030114,
        01641: 0176000,
        01642: 0055000, // ISZ31
        01643: 0060377,
        01644: 0011000,
        01645: 0000412,
        01646: 0035000,
        01647: 0175014,
        01650: 0000405,
        01651: 0113000,
        01652: 0132414,
        01653: 0000766,
        01654: 0000407,
        01655: 0063077, // ISZ32
        01656: 0000403,
        01657: 0035000,
        01660: 0063077,
        01661: 0102400,
        01662: 0000757,
        01663: 0101010, // ISZ33

        01664: 0126400, // ISZ40
        01665: 0102000,
        01666: 0034115,
        01667: 0030114,
        01670: 0041000,
        01671: 0151400,
        01672: 0156414,
        01673: 0000775,
        01674: 0045400,

        01675: 0030114, // ISZ41
        01676: 0060377,
        01677: 0011000,
        01700: 0000403,
        01701: 0151400,
        01702: 0000774,

        01703: 0021000, // ISZ42
        01704: 0156414,
        01705: 0126001,
        01706: 0156414,
        01707: 0063077,
        01710: 0125004,
        01711: 0000754,

        01712: 0102520, // DZS0
        01713: 0024122,
        01714: 0152000,
        01715: 0176000,
        01716: 0054107, // DZS1
        01717: 0060377,
        01720: 0014107,
        01721: 0000404,
        01722: 0020107,
        01723: 0063077,
        01724: 0102401,
        01725: 0106404,
        01726: 0000770,
        01727: 0150014,
        01730: 0063077,

        01731: 0102520, // DZS2
        01732: 0126620,
        01733: 0152520,
        01734: 0050107, // DZS3
        01735: 0060377,
        01736: 0014107,
        01737: 0000404,
        01740: 0107004,
        01741: 0000773,
        01742: 0000405,
        01743: 0020107, // DZS4
        01744: 0063077,
        01745: 0102400,
        01746: 0000766,
        01747: 0101234, // DZS5
        01750: 0063077,

        01751: 0152000, // DSZ10
        01752: 0176000,
        01753: 0102620,
        01754: 0101400,
        01755: 0143025,
        01756: 0000424,
        01757: 0040107,
        01760: 0060377,
        01761: 0014107,
        01762: 0176001,
        01763: 0000411,
        01764: 0024107,
        01765: 0122015,
        01766: 0000767,

        01767: 0117000, // DSZ11
        01770: 0063077,
        01771: 0152400,
        01772: 0176000,
        01773: 0000762,

        01774: 0101235, // DSZ12
        01775: 0000412,
        01776: 0024107,
        01777: 0117000,
        02000: 0063077,
        02001: 0000770,

        02002: 0024107, // DSZ13
        02003: 0176400,
        02004: 0102400,
        02005: 0063077,
        02006: 0000763,

        02007: 0101010, // DSZ14

        02010: 0000420,

        02030: 0030121, // DSZ20
        02031: 0102620,
        02032: 0034110,
        02033: 0117000,
        02034: 0040107, // DSZ21
        02035: 0060377,
        02036: 0014107,
        02037: 0014107,
        02040: 0014107,
        02041: 0014107,
        02042: 0014107,
        02043: 0024107,
        02044: 0136414,
        02045: 0063077,
        02046: 0136414,
        02047: 0152400,
        02050: 0143023,
        02051: 0000761,
        02052: 0101010, // DSZ22

        02053: 0102520, // DSZ30
        02054: 0024115,
        02055: 0030114,
        02056: 0176520,
        02057: 0055000, // DSZ31
        02060: 0060377,
        02061: 0015000,
        02062: 0000412,
        02063: 0035000,
        02064: 0175014,
        02065: 0000405,
        02066: 0113000,
        02067: 0132414,
        02070: 0000776,
        02071: 0000407,
        02072: 0063077, // DSZ32
        02073: 0000403,
        02074: 0035000,
        02075: 0063077,
        02076: 0102400,
        02077: 0006757, // DSZ33
        02100: 0101010,

        02101: 0126000, // DSZ40
        02102: 0102520,
        02103: 0034115,
        02104: 0030114,
        02105: 0041000,
        02106: 0151400,
        02107: 0156414,
        02110: 0000775,
        02111: 0051400,

        02112: 0102000, // DSZ41
        02113: 0034111,
        02114: 0030114,

        02115: 0173000, // DSZ42
        02116: 0060377,
        02117: 0015374,
        02120: 0000410,
        02121: 0015375,
        02122: 0000407,
        02123: 0015376,
        02124: 0000406,
        02125: 0015377,
        02126: 0000405,
        02127: 0000766,

        02130: 0113000, // DSZ43
        02131: 0113000,
        02132: 0113000,
        02133: 0113000,
        02134: 0034115,
        02135: 0021000,
        02136: 0156414,
        02137: 0126400,
        02140: 0156414,
        02141: 0063077,
        02142: 0124014,
        02143: 0000737,

        02144: 0102520, //DEF0
        02145: 0030114,
        02146: 0034115,
        02147: 0050107,
        02150: 0060377,
        02151: 0026107,
        02152: 0024107,
        02153: 0132414,
        02154: 0063077,
        02155: 0132414,
        02156: 0102400,
        02157: 0113000,
        02160: 0156414,
        02161: 0000766,
        02162: 0101010,

        02163: 0102520, //DEF10
        02164: 0034114,
        02165: 0030115,
        02166: 0055400,
        02167: 0175400,
        02170: 0156414,
        02171: 0000775,
        02172: 0034114,
        02173: 0054107, //DEF11
        02174: 0060377,
        02175: 0026107,
        02176: 0136414,
        02177: 0063077,
        02200: 0136414,
        02201: 0102400,
        02202: 0117000,
        02203: 0156414,
        02204: 0000767,
        02205: 0101010,

        02206: 0102520, // DEF14
        02207: 0030114,
        02210: 0034122,
        02211: 0050107,
        02212: 0150000,
        02213: 0060377,
        02214: 0052107,
        02215: 0150000,
        02216: 0024107,
        02217: 0132414,
        02220: 0063077,
        02221: 0146414,
        02222: 0102400,
        02223: 0113000,
        02224: 0116404,
        02225: 0000764,
        02226: 0101010, // DEF15

        02227: 0102520, // DEF20
        02230: 0030114,
        02231: 0050107,
        02232: 0154000,
        02233: 0060377,
        02234: 0056107,
        02235: 0025000,
        02236: 0136414,
        02237: 0063077,
        02240: 0136414,
        02241: 0102402,
        02242: 0113000,
        02243: 0034115,
        02244: 0156414,
        02245: 0000764,
        02246: 0101010, // DEF21

        02247: 0102520, // DEF30
        02250: 0030114,
        02251: 0034115,
        02252: 0051000,
        02253: 0151400,
        02254: 0156414,
        02255: 0000775,

        02256: 0030114, // DEF31
        02257: 0060377,
        02260: 0027000,
        02261: 0132414,
        02262: 0063077,
        02263: 0132414,
        02264: 0102400,
        02265: 0113000,
        02266: 0156414,
        02267: 0000770,
        02270: 0101010, // DEF32

        02271: 0102520, // DEF34
        02272: 0034115,
        02273: 0030114,
        02274: 0051001,
        02275: 0151400,
        02276: 0156414,
        02277: 0000775,
        02300: 0030114, // DEF35
        02301: 0060377,
        02302: 0053001,
        02303: 0025000,
        02304: 0132414,
        02305: 0063077,
        02306: 0132414,
        02307: 0102400,
        02310: 0113000,
        02311: 0156014,
        02312: 0000767,
        02313: 0101010, // DEF36

        02314: 0102520, // DEF40
        02315: 0030114,
        02316: 0034115,
        02317: 0051000,
        02320: 0151400,
        02321: 0156414,
        02322: 0000775,

        02323: 0030114, // DEF41
        02324: 0126620,
        02325: 0147000,
        02326: 0044107,
        02327: 0060377,
        02330: 0026107,
        02331: 0132414,
        02332: 0063077,
        02333: 0132414,
        02334: 0102400,
        02335: 0113000,
        02336: 0156414,
        02337: 0000765,
        02340: 0101010, // DEF42

        02341: 0030114, // JSR20
        02342: 0034115,
        02343: 0020422,
        02344: 0024420,
        02345: 0041000, // JSR21
        02346: 0045001,
        02347: 0045002,
        02350: 0151400,
        02351: 0151400,
        02352: 0151400,
        02353: 0156432,
        02354: 0000771,
        02355: 0045375,
        02356: 0034114, // JSR22
        02357: 0175400,
        02360: 0030406,
        02361: 0102400,
        02362: 0126400,
        02363: 0001777,
        02364: 0005000, // CJSR2
        02365: 0005402, // CJSR3
        02366: 0002367, // CJSR4
        02367: 0107004, // JSR23
        02370: 0063077,
        02371: 0020115,
        02372: 0162644,
        02373: 0100005,
        02374: 0101011,
        02375: 0063077,
        02376: 0101000,

        02377: 0034115, // DEF50
        02400: 0030114,
        02401: 0024114,
        02402: 0020437,
        02403: 0166640,
        02404: 0041000,
        02405: 0151400,
        02406: 0125404,
        02407: 0000775,
        02410: 0050107,
        02411: 0020114,
        02412: 0041000, // DEF51
        02413: 0101400,
        02414: 0151400,
        02415: 0156414,
        02416: 0000774,
        02417: 0030107, // DEF52
        02420: 0102620,
        02421: 0143000,
        02422: 0040107,
        02423: 0024114,
        02424: 0102520,
        02425: 0060377, // DEF53
        02426: 0006107,
        02427: 0136014,
        02430: 0063077,
        02431: 0136014,
        02432: 0102400,
        02433: 0107000,
        02434: 0101014,
        02435: 0010107,
        02436: 0132014,
        02437: 0000766,
        02440: 0101011, // DEF54
        02441: 0005400, // CJDF1

        02442: 0020114, // RANFL
        02443: 0040107,
        02444: 0034115,
        02445: 0024000,
        02446: 0121000, // RANFC
        02447: 0024433,
        02450: 0152620,
        02451: 0107222,
        02452: 0147000,
        02453: 0131000,
        02454: 0113520,
        02455: 0107000,
        02456: 0146400,
        02457: 0020112,
        02460: 0123400,
        02461: 0030114,
        02462: 0162433,
        02463: 0142033,
        02464: 0000762,
        02465: 0030107,
        02466: 0112415,
        02467: 0000757,
        02470: 0101120,
        02471: 0125100,
        02472: 0101200,
        02473: 0125200,
        02474: 0042107,
        02475: 0010107,
        02476: 0030107,
        02477: 0156014,
        02500: 0000746,
        02501: 0101011,
        02502: 0135753, // C1347

        02503: 0024114, // DEF60
        02504: 0121120,
        02505: 0101240,
        02506: 0040123,
        02507: 0102520,
        02510: 0131000,
        02511: 0101125,
        02512: 0000417,
        02513: 0031000,
        02514: 0151132,
        02515: 0000774,
        02516: 0050107,
        02517: 0031000,
        02520: 0102520, // DEF61
        02521: 0060377,
        02522: 0036123,
        02523: 0172414,
        02524: 0063077,
        02525: 0172414,
        02526: 0102401,
        02527: 0101005,
        02530: 0000771,
        02531: 0125400, // DEF62
        02532: 0034115,
        02533: 0136414,
        02534: 0000750,
        02535: 0101010, // DEF63

        02536: 0030114, // EXCH
        02537: 0034115,
        02540: 0021000,
        02541: 0025777,
        02542: 0045000,
        02543: 0041777,
        02544: 0151400,
        02545: 0174400,
        02546: 0174000,
        02547: 0172433,
        02550: 0000770,

        02551: 0024114, // DEF64
        02552: 0121120,
        02553: 0101240,
        02554: 0040123,
        02555: 0102520,
        02556: 0131000,
        02557: 0101125,
        02560: 0000422,
        02561: 0031000,
        02562: 0151132,
        02563: 0000774,
        02564: 0050107,
        02565: 0031000,

        02566: 0102520, // DEF65
        02567: 0060377,
        02570: 0012123,
        02571: 0036107,
        02572: 0016107,
        02573: 0156015,
        02574: 0000413,
        02575: 0063077,
        02576: 0102400,
        02577: 0052107,
        02600: 0101005, // DEF66
        02601: 0000766,
        02602: 0125400,
        02603: 0034115,
        02604: 0136414,
        02605: 0000745,
        02606: 0000406,

        02607: 0036107, // DEF67
        02610: 0172415,
        02611: 0000767,
        02612: 0063077,
        02613: 0000763,
        02614: 0101010,

        02615: 0102400, // ID0
        02616: 0152400,
        02617: 0126620,
        02620: 0041020,
        02621: 0151400,
        02622: 0125224,
        02623: 0000775,

        02624: 0176620, // ID1
        02625: 0021000,
        02626: 0025000,
        02627: 0106414,
        02630: 0063077,
        02631: 0151400,
        02632: 0175224,
        02633: 0000772,

        02634: 0102000, // ID2
        02635: 0151112,
        02636: 0000403,
        02637: 0152620,
        02640: 0000757,
        02641: 0101010,

        02642: 0020124, // ID4
        02643: 0024112,
        02644: 0176620,
        02645: 0175220,
        02646: 0137415,
        02647: 0000776,

        02650: 0055427, // ID5
        02651: 0027427,
        02652: 0025427,
        02653: 0136414,
        02654: 0000410,
        02655: 0024113,
        02656: 0175220,
        02657: 0137415,
        02660: 0000770,
        02661: 0101404,
        02662: 0000761,
        02663: 0000407,

        02664: 0063077, // ID6
        02665: 0055427,
        02666: 0060377,
        02667: 0027427,
        02670: 0025427,
        02671: 0000774,

        02672: 0102520, // ID10
        02673: 0176400,
        02674: 0152620,
        02675: 0151240,
        02676: 0151220,
        02677: 0051420,
        02700: 0060377,
        02701: 0027420,
        02702: 0025420,
        02703: 0146014,
        02704: 0063077,
        02705: 0146014,
        02706: 0102400,
        02707: 0113000,
        02710: 0151113,
        02711: 0000766,
        02712: 0175400,
        02713: 0024076,
        02714: 0136414,
        02715: 0000757,
        02716: 0101010, // ID11

        02717: 0102520, // ID14
        02720: 0176400,
        02721: 0152620,
        02722: 0151220,
        02723: 0151220,
        02724: 0051430,
        02725: 0060377,
        02726: 0027430,
        02727: 0025430,
        02730: 0132014,
        02731: 0063077,
        02732: 0132014,
        02733: 0102400,
        02734: 0112404,
        02735: 0000767,
        02736: 0175400,
        02737: 0024076,
        02740: 0136414,
        02741: 0000760,
        02742: 0101010, // ID15

        02743: 0102000,
        02744: 0176400, // ID20
        02745: 0030114,
        02746: 0051420,
        02747: 0060377,
        02750: 0053420,
        02751: 0025420,
        02752: 0146015, // ID21
        02753: 0000415,
        02754: 0063077,
        02755: 0102400,
        02756: 0041000, // ID22
        02757: 0112400,
        02760: 0024115,
        02761: 0132414,
        02762: 0000764,
        02763: 0175400,
        02764: 0024076,
        02765: 0136414,
        02766: 0000757,
        02767: 0000406,
        02770: 0025001,
        02771: 0132415,

        02772: 0000764, // ID23
        02773: 0063077,
        02774: 0000761,
        02775: 0101010, // ID24

        02776: 0102000, // ID30
        02777: 0176400,
        03000: 0030115,
        03001: 0051430,
        03002: 0060377,
        03003: 0053430,
        03004: 0025430,
        03005: 0132015, // ID31
        03006: 0000415,
        03007: 0063077,
        03010: 0102400,
        03011: 0041000, // ID32
        03012: 0113000,
        03013: 0024114,
        03014: 0132414,
        03015: 0000764,
        03016: 0175400,
        03017: 0024076,
        03020: 0136414,
        03021: 0000757,
        03022: 0000406,

        03023: 0025377, // ID33
        03024: 0132415,
        03025: 0000764,
        03026: 0063077,
        03027: 0000761,

        03030: 0101010, // ID34

        03031: 0102520, // ID40
        03032: 0126000,
        03033: 0030114,
        03034: 0034115,
        03035: 0050027, // ID41
        03036: 0045001,
        03037: 0060377,
        03040: 0012027,
        03041: 0102401,
        03042: 0113001,
        03043: 0063077,
        03044: 0156414,
        03045: 0000770,
        03046: 0101010, // ID42

        03047: 0024425, // ID44
        03050: 0034114,
        03051: 0030115,
        03052: 0050034,
        03053: 0156400,
        03054: 0046034,
        03055: 0175404,
        03056: 0000776,

        03057: 0102520, // ID45
        03060: 0024114,
        03061: 0044025,
        03062: 0006025,
        03063: 0174400,
        03064: 0174000,
        03065: 0136014,
        03066: 0102401,
        03067: 0107001,
        03070: 0063077,
        03071: 0132014,
        03072: 0000767,
        03073: 0101011, // ID46
        03074: 0005400, // ID40

        03075: 0024114, // ID50
        03076: 0030115,
        03077: 0102520,
        03100: 0050033,
        03101: 0006033,
        03102: 0156414,
        03103: 0102401,
        03104: 0112401,
        03105: 0063077,
        03106: 0132014,
        03107: 0000771,

        03110: 0034402, // DGCA
        03111: 0152001,
        03112: 0004000,
        03113: 0020112,
        03114: 0116432,
        03115: 0002421,
        03116: 0102620,
        03117: 0143005, // DGCX
        03120: 0000413,
        03121: 0041400,
        03122: 0060377,
        03123: 0015400,
        03124: 0101010,
        03125: 0025400,
        03126: 0122015,
        03127: 0000770,
        03130: 0152400,
        03131: 0063077,
        03132: 0000765,
        03133: 0020757, // DGCB
        03134: 0117000,
        03135: 0000754,
        03136: 0000162, // LAST
    }
    startAddr := 00002
    haltAddr := 00236

    n := NewNova()
    n.LoadMemory(0, program[:])

    // Stop test after one pass
    n.Deposit(haltAddr, 0063077)

    n.Switches(0)   // No devices yet

    if false {
        DisasmBlock(program[:], startAddr, haltAddr - startAddr)
        return
    }
    if false {
        _, err := n.Trace(startAddr)
        if err != nil {
            t.Error(err)
        }
        return
    }
    n.Start(startAddr)
    if n.WaitForHalt(time.Millisecond * 1000) != nil {
        n.Stop()
        t.Error("machine did not halt")
    }
    addr := n.Stop()
    if addr != haltAddr + 1 {
        t.Errorf("have: %05o, want: %05o", addr, haltAddr + 1)
    }
}

/*
        00000: 0,
        00001: 0,
        00002: 0,
        00003: 0,
        00004: 0,
        00005: 0,
        00006: 0,
        00007: 0,
        00010: 0,
        00011: 0,
        00012: 0,
        00013: 0,
        00014: 0,
        00015: 0,
        00016: 0,
        00017: 0,
        00020: 0,
        00021: 0,
        00022: 0,
        00023: 0,
        00024: 0,
        00025: 0,
        00026: 0,
        00027: 0,
        00030: 0,
        00031: 0,
        00032: 0,
        00033: 0,
        00034: 0,
        00035: 0,
        00036: 0,
        00037: 0,
        00040: 0,
        00041: 0,
        00042: 0,
        00043: 0,
        00044: 0,
        00045: 0,
        00046: 0,
        00047: 0,
        00050: 0,
        00051: 0,
        00052: 0,
        00053: 0,
        00054: 0,
        00055: 0,
        00056: 0,
        00057: 0,
        00060: 0,
        00061: 0,
        00062: 0,
        00063: 0,
        00064: 0,
        00065: 0,
        00066: 0,
        00067: 0,
        00070: 0,
        00071: 0,
        00072: 0,
        00073: 0,
        00074: 0,
        00075: 0,
        00076: 0,
        00077: 0,

*/