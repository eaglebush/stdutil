package stdutil

import (
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

var src rand.Source

const (
	idxbits = 6              // 6 bits to represent a letter index
	idxmsk  = 1<<idxbits - 1 // All 1-bits, as many as letterIdxBits
	idxmax  = 63 / idxbits   // # of letter indices fitting in 63 bits

	ltrl string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // upper case letters
	ltrn string = "012356789"                  // numbers
)

func init() {
	src = rand.NewSource(time.Now().UTC().UnixNano())
}

// GenerateFull generates a random text that uses all letters of the alphabet, in both cases including numbers
func GenerateFull(length int) string {
	return genRndString(length, strings.ToLower(ltrl)+ltrl+ltrn)
}

// GenerateText generates a random text that uses all letters of the alphabet, in both cases not including numbers
func GenerateText(length int) string {
	return genRndString(length, strings.ToLower(ltrl)+ltrl)
}

// GenerateSeries generates random numbers
func GenerateSeries(length int) string {
	return genRndString(length, ltrn)
}

// GenerateAlpha generates a random text that uses the 26 letters of the alphabet
func GenerateAlpha(length int, lower bool) string {
	if lower {
		return genRndString(length, strings.ToLower(ltrl))
	}
	return genRndString(length, ltrl)
}

func genRndString(length int, mask string) string {

	b := make([]byte, length)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), idxmax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), idxmax
		}
		if idx := int(cache & idxmsk); idx < len(mask) {
			b[i] = mask[idx]
			i--
		}
		cache >>= idxbits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
