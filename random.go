// SPDX-License-Identifier: MIT

package gromit

import (
	"io"
	"unicode"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
	rand "math/rand"
)

func PickString(begin, end string) (string, error) {
	br := []rune(begin)
	er := []rune(end)
	if len(br) != len(er) {
		return "", ErrBadRange
	}
	ret := make([]rune, len(br))
	for i := range ret {
		if int32(br[i]) > int32(er[i]) {
			return "", ErrBadRange
		}
		ret[i] = PickRune(br[i], er[i])
	}
	return string(ret), nil
}

func PickRune(begin, end rune) rune {
	return rune(PickInt32(int32(begin), int32(end)))
}

func PickInt32(begin, end int32) int32 {
	if begin > end {
		log.Fatal("PickInt32: invalid arguments: begin > end", begin, end)
	}
	diff := int64(end) - int64(begin)
	return int32(int64(begin) + rand.Int63n(diff+1))
}

func PickBool() bool {
	if rand.Int63()&1 == 1 {
		return true
	}
	return false
}

func IsCapital(s string) bool {
	ch, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(ch)
}

func pad(dst io.Writer) error {
	var padding string
	padding = " " // FIXME
	runes := []rune(padding)
	if len(runes) == 0 {
		return nil
	}
	r := runes[rand.Intn(len(runes))]
	_, err := io.WriteString(dst, string([]rune{r}))
	return err
}
