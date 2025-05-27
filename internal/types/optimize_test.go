package types

import (
	"testing"
)

type Item struct {
	ID   int
	Name string
}

// --- テスト関数 ---

func TestPermutationsOfStructPointers(t *testing.T) {
	items := []*Item{
		{ID: 1, Name: "Apple"},
		{ID: 2, Name: "Banana"},
		{ID: 3, Name: "Cherry"},
	}

	n := factorial(len(items))
	for k := 0; k < n; k++ {
		permutation, err := getPermutationByIndex(k, items)
		if err != nil {
			t.Fatalf("unexpected error for k=%d: %v", k, err)
		}

		// 出力を確認のためログ表示
		t.Logf("Permutation %d:", k)
		for _, v := range permutation {
			it := v.(*Item)
			t.Logf("  ID: %d, Name: %s", it.ID, it.Name)
		}
	}
}
