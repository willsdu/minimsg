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
