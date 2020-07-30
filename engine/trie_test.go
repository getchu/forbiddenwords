package engine

import (
	"forbiddenwords/engine"
	"testing"
)

func TestTrie_Add(t *testing.T) {
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewTrie()
	for _, word := range wordList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("add error：" + err.Error())
		}
	}
}

func TestTrie_Len(t *testing.T) {
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewTrie()
	for _, word := range wordList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("add error：" + err.Error())
		}
	}
	l := eng.Len()
	if l != len(wordList) {
		t.Fatalf("Len error: %d / %d", l, len(wordList))
	}
}

func TestTrie_Find(t *testing.T) {
	forbiddenList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewTrie()
	for _, word := range forbiddenList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("add error：" + err.Error())
		}
	}
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	//匹配存在的
	for _, word := range wordList {
		result := eng.Find(word)
		if result != word {
			t.Fatalf("find error1：%s / %s", result, word)
		}
	}
	//匹配不存在的
	wordList = []string{"政1", "政2", "色1", "色2", "广1", "广2", "辱1", "辱2"}
	for _, word := range wordList {
		result := eng.Find(word)
		if len(result) > 0 {
			t.Fatalf("find error2：%s / %s", result, word)
		}
	}
	//匹配包含的
	wordList = []string{"哈哈政哈哈治哈哈1哈哈", "哈哈政哈哈治哈哈2哈哈", "哈哈色哈哈情哈哈1哈哈", "哈哈色哈哈情哈哈2哈哈", "哈哈广哈哈告哈哈1哈哈", "哈哈广哈哈告哈哈2哈哈", "哈哈辱哈哈骂哈哈1哈哈", "哈哈辱哈哈骂哈哈2哈哈"}
	wordList = []string{"哈哈政哈哈治哈哈1哈哈"}
	for k, word := range wordList {
		result := eng.Find(word)
		if result != forbiddenList[k] {
			t.Fatalf("find error2：%s / %s", result, word)
		}
	}
}
