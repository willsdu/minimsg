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

func TestSignature(t *testing.T) {
	s := "f4471f7d6658bec36c0c58aca4a1f499847e7ab0"
	ns := GenSignature("1593998600", "1825473595")
	if s != ns {
		t.Error("error")
		return
	}
	t.Log("matched")
}
