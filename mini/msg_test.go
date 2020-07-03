package mini

import (
	"sort"
	"testing"
)

func TestDicSort(t *testing.T) {
	ps := []string{"sgf", "dgesarg", "zsdff"}
	sort.Slice(ps, func(i, j int) bool {
		return ps[i] < ps[j]
	})
	t.Log(ps)
}

//signature=&echostr=6044848836944286472&timestamp=&nonce=

func TestSignature(t *testing.T) {
	s := "9e425999690c965802d464489fc242e7d10c4716"
	ns := GenSignature("1593774856", "472719312")
	t.Log(s == ns)

}
