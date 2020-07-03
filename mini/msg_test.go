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

//signature=&echostr=9020835996559703542&timestamp=&nonce=247735498

func TestSignature(t *testing.T) {
	s := "a4ffbec5bbdb7906aaa4dc64fb6cf52dd8eab0ff"
	ns := GenSignature("1593775542", "247735498")
	t.Log(s == ns)

}
