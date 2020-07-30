package engine

import (
	"forbiddenwords/engine"
	"testing"
)

func TestHash_Add(t *testing.T) {
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewHash()
	for _, word := range wordList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("插入错误：" + err.Error())
		}
	}
}

func TestHash_Len(t *testing.T) {
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewHash()
	for _, word := range wordList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("插入错误：" + err.Error())
		}
	}
	if eng.Len() != len(wordList) {
		t.Fatal("总数错误")
	}
}

func TestHash_Find(t *testing.T) {
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewHash()
	for _, word := range wordList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("插入错误：" + err.Error())
		}
	}
	if eng.Len() != len(wordList) {
		t.Fatal("总数错误")
	}
	for _, word := range wordList {
		r := eng.Find(word)
		if len(r) <= 0 {
			t.Fatal("未找到词" + word)
		}
	}
	word := "你好"
	r := eng.IsExists(word)
	if r {
		t.Fatal("不存在的词" + word + "被找到了")
	}
}

func TestHash_IsExists(t *testing.T) {
	wordList := []string{"政治1", "政治2", "色情1", "色情2", "广告1", "广告2", "辱骂1", "辱骂2"}
	eng := engine.NewHash()
	for _, word := range wordList {
		err := eng.Add(word)
		if err != nil {
			t.Fatal("插入错误：" + err.Error())
		}
	}
	if eng.Len() != len(wordList) {
		t.Fatal("总数错误")
	}
	for _, word := range wordList {
		r := eng.IsExists(word)
		if !r {
			t.Fatal("词" + word + "不存在")
		}
	}
	word := "你好"
	r := eng.IsExists(word)
	if r {
		t.Fatal("不存在的词" + word + "存在了")
	}
}

//"政治": {"政治1", "政治2"},
//"色情": {"色情1", "色情2"},
//"广告": {"广告1", "广告2"},
//"辱骂": {"辱骂1", "辱骂2"},}{
