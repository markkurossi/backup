//
// algorithms.go
//
// Copyright (c) 2018 Markku Rossi
//
// All rights reserved.
//

package zone

import (
	"fmt"
)

type Suite byte

func (s Suite) String() string {
	name, ok := suites[s]
	if ok {
		return name
	}
	return fmt.Sprintf("{Suite %d}", s)
}

func (s Suite) CipherKeyLen() int {
	len, ok := suiteCipherKeyLengths[s]
	if !ok {
		panic(fmt.Sprintf("Unknown suite: %d", s))
	}
	return len
}

func (s Suite) HMACKeyLen() int {
	len, ok := suiteHMACKeyLengths[s]
	if !ok {
		panic(fmt.Sprintf("Unknown suite: %d", s))
	}
	return len
}

func (s Suite) KeyLen() int {
	return s.CipherKeyLen() + s.HMACKeyLen()
}

const (
	SuiteAES128CBCHMACSHA256 Suite = 0
)

var suites = map[Suite]string{
	SuiteAES128CBCHMACSHA256: "AES128-CBC-HMAC-SHA256",
}

var suiteCipherKeyLengths = map[Suite]int{
	SuiteAES128CBCHMACSHA256: 16,
}

var suiteHMACKeyLengths = map[Suite]int{
	SuiteAES128CBCHMACSHA256: 32,
}
