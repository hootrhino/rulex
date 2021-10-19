package test

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/ngaut/log"
)

func TestRegx(t *testing.T) {

	t.Log(Match("<a:1 b:2 c:3 d:5 e:9"))
}

type Kl struct {
	K string
	L int
}

func Match(s string) []Kl {
	p := `[a-z]+:[1-9]+`
	endian := s[0]
	// 大端
	if endian == '>' {

	}
	// 小端
	if endian == '<' {

	}
	result := []Kl{}
	matched, err0 := regexp.MatchString(p, s[0:])
	if matched {
		// [a:12 b:12]
		for _, v := range regexp.MustCompile(p).FindAllString(s[0:], -1) {
			kl := strings.Split(v, ":")
			if l, err1 := strconv.Atoi(kl[1]); err1 == nil {
				result = append(result, Kl{kl[0], l})
			}
		}
	} else {
		log.Error(err0)
	}
	return result
}
