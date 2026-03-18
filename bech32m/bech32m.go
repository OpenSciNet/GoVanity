// Package bech32m implements BIP-0173 and BIP-0350.
// This file is licensed under the MIT License.
// Copyright (c) 2026 OpenSciNet
package bech32m

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"	//	excluding the number 1 and the letters b, i, o
var bech32mConst = 0x2bc830a3
var bech32mGenerator = []int{
    0x3b6a57b2,
    0x26508e6d,
    0x1ea119fa,
    0x3d4233dd,
    0x2a1462b3,
}

func bech32mPolymod(values []byte) int {
    chk := 1
    for _, v := range values {
        top := chk >> 25
        chk = (chk&0x1ffffff)<<5 ^ int(v)
        for i := 0; i < 5; i++ {
            if (top>>uint(i))&1 == 1 {
                chk ^= bech32mGenerator[i]
            }
        }
    }
    return chk
}

func bech32mHrpExpand(hrp string) []byte {
    ret := []byte{}
    for _, c := range hrp {
        ret = append(ret, byte(c>>5))
    }
    ret = append(ret, 0)
    for _, c := range hrp {
        ret = append(ret, byte(c&31))
    }
    return ret
}

func bech32mCreateChecksum(hrp string, data []byte) []byte {
    values := append(append(bech32mHrpExpand(hrp), data...), []byte{0, 0, 0, 0, 0, 0}...)
    mod := bech32mPolymod(values) ^ bech32mConst
    ret := make([]byte, 6)
    for i := 0; i < 6; i++ {
        ret[i] = byte((mod >> uint(5*(5-i))) & 31)
    }
    return ret
}

func EncodeBech32m(hrp string, data []byte) string {
    combined := append(data, bech32mCreateChecksum(hrp, data)...)
    out := make([]byte, 0, len(hrp)+1+len(combined))
    out = append(out, hrp...)
    out = append(out, '1')
    for _, v := range combined {
        out = append(out, charset[v])
    }
    return string(out)
}
