// Scramble functions.
// Input `s []byte` is required to be not nil.
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
)

var Salt []byte

const emailDomain = "@example.com"
const bindUrl = "https://example.com?quote_gid="
const DomainLen = 10

func GenScrambleBytes(maxLength uint) func([]byte) []byte {
	return func(s []byte) []byte {
		scrambled := ScrambleBytes(s)
		if maxLength > uint(len(scrambled)) {
			return scrambled
		} else {
			return scrambled[:maxLength]
		}
	}
}

var bytesOutputAlphabetLength = byte(len(bytesOutputAlphabet))
var bytesSafeAlphabetLength = byte(len(bytesSafeAlphabet))
var bytesKeep = []byte("',\\{}")
var bytesOutputAlphabet = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+-_")
var bytesSafeAlphabet = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
var DomainPostfix = []byte(".example")

// Modifies `s` in-place.
func ScrambleBytes(s []byte) []byte {
	isArray := len(s) >= 2 && s[0] == '{' && s[len(s)-1] == '}'

	hash := sha256.New()
	// Hard-coding this constant wins less than 3% in BenchmarkScrambleBytes
	const sumLength = 32 // SHA256/8
	hash.Write(Salt)
	hash.Write(s)
	sumBytes := hash.Sum(nil)

	reader := bytes.NewReader(s)
	var r rune
	var err error
	for i := 0; ; i++ {
		r, _, err = reader.ReadRune()
		if err != nil {
			s = s[:i]
			break
		}
		if !isArray || bytes.IndexRune(bytesKeep, r) == -1 {
			// Do not insert, so should not obstruct reader.
			s[i] = bytesOutputAlphabet[(sumBytes[i%sumLength]+byte(r))%bytesOutputAlphabetLength]
		} else {
			// Possibly shift bytes to beginning of s.
			s[i] = byte(r)
		}
	}
	return s
}

// Modifies `s` in-place.
func ScrambleSafeBytes(s []byte) []byte {
	isArray := len(s) >= 2 && s[0] == '{' && s[len(s)-1] == '}'

	hash := sha256.New()
	// Hard-coding this constant wins less than 3% in BenchmarkScrambleBytes
	const sumLength = 32 // SHA256/8
	hash.Write(Salt)
	hash.Write(s)
	sumBytes := hash.Sum(nil)

	reader := bytes.NewReader(s)
	var r rune
	var err error
	for i := 0; ; i++ {
		r, _, err = reader.ReadRune()
		if err != nil {
			s = s[:i]
			break
		}
		if !isArray || bytes.IndexRune(bytesKeep, r) == -1 {
			// Do not insert, so should not obstruct reader.
			s[i] = bytesSafeAlphabet[(sumBytes[i%sumLength]+byte(r))%bytesSafeAlphabetLength]
		} else {
			// Possibly shift bytes to beginning of s.
			s[i] = byte(r)
		}
	}
	return s
}

func ScrambleDigits(s []byte) []byte {
	hash := sha256.New()
	const sumLength = 32 // SHA256/8
	hash.Write(Salt)
	hash.Write(s)
	sumBytes := hash.Sum(nil)

	for i, b := range s {
		if b >= '0' && b <= '9' {
			s[i] = '0' + (sumBytes[i%sumLength]+b)%10
		}
	}
	return s
}

func ScrambleIBAN(s []byte) []byte {
	//return constant for now
	b := []byte("DE75512108001245126199")
	return b
}

func scrambleOneEmail(s []byte) []byte {
	atIndex := bytes.IndexRune(s, '@')
	mailbox := Salt
	if atIndex != -1 {
		mailbox = s[:atIndex]
	}
	s = make([]byte, len(mailbox)+len(emailDomain))
	copy(s, mailbox)
	// ScrambleBytes is in-place, but may return string shorter than input.
	mailbox = ScrambleBytes(s[:len(mailbox)])
	copy(s[len(mailbox):], emailDomain)
	// So final len(mailbox) may be shorter than whole allocated string.
	return s[:len(mailbox)+len(emailDomain)]
}

func ScrambleBindUrls(s []byte) []byte {
	paramName := []byte("gid=")

	if len(s) == 0 {
		return s
	}

	if bytes.LastIndex(s, paramName) == -1 {
		return s
	}

	arr := bytes.Split(s, paramName)
	gid := arr[1]

	// ScrambleBytes is in-place, but may return string shorter than input.
	scrambledGid := ScrambleBytes(gid)

	s = make([]byte, len(bindUrl)+len(scrambledGid))
	copy(s, bindUrl)
	copy(s[len(bindUrl):], scrambledGid)

	return s
}

func scrambleOneUniqueEmail(s []byte) []byte {
	atIndex := bytes.IndexRune(s, '@')
	mailbox := Salt
	if atIndex != -1 {
		mailbox = s[:atIndex]
	} else {
		mailbox = append([]byte(nil), Salt...)
	}
	domain := s[atIndex+1:]

	scrambledMailbox := ScrambleBytes(mailbox)
	scrambledDomain := ScrambleSafeBytes(domain)
	total := len(scrambledMailbox) + 1 + len(scrambledDomain) + len(DomainPostfix)
	s2 := make([]byte, total)
	copy(s2, scrambledMailbox)
	copy(s2[len(scrambledMailbox):], []byte("@"))
	copy(s2[len(scrambledMailbox)+1:], scrambledDomain)
	copy(s2[len(scrambledMailbox)+1+len(scrambledDomain):], DomainPostfix)
	return s2
}

// Supports array of emails in format {email1,email2}
func ScrambleEmail(s []byte) []byte {
	if len(s) < 2 {
		// panic("ScrambleEmail: input is too small: '" + string(s) + "'")
		return s
	}
	if s[0] != '{' && s[len(s)-1] != '}' {
		return scrambleOneEmail(s)
	}
	parts := bytes.Split(s[1:len(s)-1], []byte{','})
	partsNew := make([][]byte, len(parts))
	outLength := 2 + len(parts) - 1
	for i, bs := range parts {
		partsNew[i] = scrambleOneEmail(bs)
		outLength += len(partsNew[i])
	}
	s = make([]byte, outLength)
	s[0] = '{'
	s[len(s)-1] = '}'
	copy(s[1:len(s)-1], bytes.Join(partsNew, []byte{','}))
	return s
}

// Supports array of emails in format {email1,email2}
func ScrambleUniqueEmail(s []byte) []byte {
	if len(s) < 2 {
		// panic("ScrambleEmail: input is too small: '" + string(s) + "'")
		return s
	}
	if s[0] != '{' && s[len(s)-1] != '}' {
		return scrambleOneUniqueEmail(s)
	}
	parts := bytes.Split(s[1:len(s)-1], []byte{','})
	partsNew := make([][]byte, len(parts))
	outLength := 2 + len(parts) - 1
	for i, bs := range parts {
		partsNew[i] = scrambleOneUniqueEmail(bs)
		outLength += len(partsNew[i])
	}
	s = make([]byte, outLength)
	s[0] = '{'
	s[len(s)-1] = '}'
	copy(s[1:len(s)-1], bytes.Join(partsNew, []byte{','}))
	return s
}

func ScrambleInet(s []byte) []byte {
	hash := sha256.New()
	const sumLength = 32 // SHA256/8
	hash.Write(Salt)
	hash.Write(s)
	sumBytes := hash.Sum(nil)

	var ip net.IP
	// Decide to output IPv4 or IPv6 based on first bit of hash sum.
	// Gives about 50% chance to each option.
	if sumBytes[0]&0x80 != 0 {
		ip = net.IP(sumBytes[:16])
	} else {
		ip = net.IP(sumBytes[:4])
	}
	return []byte(ip.String())
}

func GetScrambleByName(value string) (func(s []byte) []byte, error) {
	switch value {
	case "bytes":
		return ScrambleBytes, nil
	case "sbytes":
		return ScrambleSafeBytes, nil
	case "digits":
		return ScrambleDigits, nil
	case "email":
		return ScrambleEmail, nil
	case "uemail":
		return ScrambleUniqueEmail, nil
	case "bindurl":
		return ScrambleBindUrls, nil
	case "inet":
		return ScrambleInet, nil
	case "iban":
		return ScrambleIBAN, nil
	}
	return nil, errors.New(fmt.Sprintf("%s is not registered scramble function", value))
}

func init() {
	Salt = make([]byte, 16)
	_, err := rand.Read(Salt)
	if err != nil {
		panic(err)
	}
}
