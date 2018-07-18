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

func (s Suite) IDHashKeyLen() int {
	len, ok := suiteIDHashKeyLengths[s]
	if !ok {
		panic(fmt.Sprintf("Unknown suite: %d", s))
	}
	return len
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
	return s.IDHashKeyLen() + s.CipherKeyLen() + s.HMACKeyLen()
}

const (
	AES256CBCHMACSHA256 Suite = 0
)

var suites = map[Suite]string{
	AES256CBCHMACSHA256: "AES256-CBC-HMAC-SHA256",
}

var suiteIDHashKeyLengths = map[Suite]int{
	AES256CBCHMACSHA256: 32,
}

var suiteCipherKeyLengths = map[Suite]int{
	AES256CBCHMACSHA256: 32,
}

var suiteHMACKeyLengths = map[Suite]int{
	AES256CBCHMACSHA256: 32,
}
