package server

import (
	"forbiddenwords/server"
	"testing"
)

func TestWords_trie_UpdateWords(t *testing.T) {
	w := server.NewWords("trie")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
}

func TestWords_trie_UpdateWordsTwice(t *testing.T) {
	w := server.NewWords("trie")
	w.UpdateWords()
	w.UpdateWords()
}

func TestWords_trie_Find(t *testing.T) {
	w := server.NewWords("trie")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
	forbidden := map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"}
	//单表查 - 查找存在的
	for tableName, msg := range map[string]string{"political": "哈哈国哈哈母哈哈", "gun": "哈哈短哈哈枪哈哈", "sexy": "哈哈骚哈哈b哈哈"} {
		tn, result := w.Find(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result != forbidden[tableName] {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//单表查 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "你好", "gun": "工作", "sexy": "学习"} {
		tn, result := w.Find(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if len(result) > 0 {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//所有表都查一遍 - 查找存在的
	for tableName, msg := range map[string]string{"political": "哈哈国哈哈母哈哈", "gun": "哈哈短哈哈枪哈哈", "sexy": "哈哈骚哈哈b哈哈"} {
		tn, result := w.Find("", msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result != forbidden[tableName] {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//所有表都查一遍 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "你好", "gun": "工作", "sexy": "学习"} {
		tn, result := w.Find("", msg)
		if len(tn) > 0 {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if len(result) > 0 {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
}

func TestWords_trie_IsExists(t *testing.T) {
	w := server.NewWords("trie")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
	//单表查 - 查找存在的
	for tableName, msg := range map[string]string{"political": "哈哈国哈哈母哈哈", "gun": "哈哈短哈哈枪哈哈", "sexy": "哈哈骚哈哈b哈哈"} {
		tn, result := w.IsExists(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if !result {
			t.Fatal("result error")
		}
	}
	//单表查 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "你好", "gun": "工作", "sexy": "学习"} {
		tn, result := w.IsExists(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result {
			t.Fatal("result error")
		}
	}
	//所有表都查一遍 - 查找存在的
	for tableName, msg := range map[string]string{"political": "哈哈国哈哈母哈哈", "gun": "哈哈短哈哈枪哈哈", "sexy": "哈哈骚哈哈b哈哈"} {
		tn, result := w.IsExists("", msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if !result {
			t.Fatal("result error")
		}
	}
	//所有表都查一遍 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "你好", "gun": "工作", "sexy": "学习"} {
		tn, result := w.IsExists("", msg)
		if len(tn) > 0 {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result {
			t.Fatal("result error")
		}
	}
}

func TestWords_trie_Len(t *testing.T) {
	w := server.NewWords("trie")
	w.UpdateWords()
	l := 0
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
		l += eng.Len()
	}
	if l != w.Len() {
		t.Fatal("total len error")
	}
}
