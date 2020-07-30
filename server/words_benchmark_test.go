package server

import (
	"forbiddenwords/server"
	"testing"
)

//测试hash下查找性能
func BenchmarkHash_Find(b *testing.B) {
	w := server.NewWords("hash")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			b.Fatal(tableName + " len error")
		}
	}
	wordList := []string{"我想知道国母是什么", "我想知道短枪是什么", "我想知道骚b是什么"}
	for i := 0; i < b.N; i++ {
		for _, word := range wordList {
			w.Find("", word)
		}
	}
}

//测试hash下查找性能 - 词组中间不含其它词
func BenchmarkTrie_Find_1(b *testing.B) {
	w := server.NewWords("trie")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			b.Fatal(tableName + " len error")
		}
	}
	wordList := []string{"我想知道国母是什么", "我想知道短枪是什么", "我想知道骚b是什么"}
	for i := 0; i < b.N; i++ {
		for _, word := range wordList {
			w.Find("", word)
		}
	}
}

//测试hash下查找性能 - 词组中间包含其它词
func BenchmarkTrie_Find_2(b *testing.B) {
	w := server.NewWords("trie")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			b.Fatal(tableName + " len error")
		}
	}
	wordList := []string{"哈哈国哈哈母哈哈", "哈哈短哈哈枪哈哈", "哈哈骚哈哈b哈哈"}
	for i := 0; i < b.N; i++ {
		for _, word := range wordList {
			w.Find("", word)
		}
	}
}
