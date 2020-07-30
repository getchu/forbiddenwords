package server

import (
	"forbiddenwords/server"
	"testing"
)

func TestWords_hash_UpdateWords(t *testing.T) {
	w := server.NewWords("hash")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
}

func TestWords_hash_UpdateWordsTwice(t *testing.T) {
	w := server.NewWords("hash")
	w.UpdateWords()
	w.UpdateWords()
}

func TestWords_hash_Find(t *testing.T) {
	w := server.NewWords("hash")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
	//单表查 - 查找存在的
	for tableName, msg := range map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"} {
		tn, result := w.Find(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result != msg {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//单表查 - 查找存在的包含的
	resultList := map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"}
	for tableName, msg := range map[string]string{"political": "我想知道国母是什么", "gun": "我想知道短枪是什么", "sexy": "我想知道骚b是什么"} {
		tn, result := w.Find(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result != resultList[tableName] {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//单表查 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "你好", "gun": "工作", "sexy": "学习"} {
		tn, result := w.Find(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result == msg {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//所有表都查一遍 - 查找存在的
	for tableName, msg := range map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"} {
		tn, result := w.Find("", msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result != msg {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//所有表都查一遍 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "你好", "gun": "工作", "sexy": "学习"} {
		tn, result := w.Find("", msg)
		if len(tn) > 0 {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result == msg {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
	//所有表都查一遍 - 查找存在的包含的
	resultList = map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"}
	for tableName, msg := range map[string]string{"political": "我想知道国母是什么", "gun": "我想知道短枪是什么", "sexy": "我想知道骚b是什么"} {
		tn, result := w.Find("", msg)
		if len(tn) <= 0 {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result != resultList[tableName] {
			t.Fatal("result error: " + msg + "/" + result)
		}
	}
}

func TestWords_hash_IsExists(t *testing.T) {
	w := server.NewWords("hash")
	w.UpdateWords()
	engineList := w.GetEngineList()
	for tableName, eng := range engineList {
		if eng.Len() <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
	//单表查 - 查找存在的
	for tableName, msg := range map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"} {
		tn, result := w.IsExists(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if !result {
			t.Fatal("result error: " + msg)
		}
	}
	//单表查 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "国母1", "gun": "短枪1", "sexy": "骚b1"} {
		tn, result := w.IsExists(tableName, msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result {
			t.Fatal("result error: " + msg)
		}
	}
	//所有表都查一遍 - 查找存在的
	for tableName, msg := range map[string]string{"political": "国母", "gun": "短枪", "sexy": "骚b"} {
		tn, result := w.IsExists("", msg)
		if tn != tableName {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if !result {
			t.Fatal("result error: " + msg)
		}
	}
	//所有表都查一遍 - 查找不存在的
	for tableName, msg := range map[string]string{"political": "国母1", "gun": "短枪1", "sexy": "骚b1"} {
		tn, result := w.IsExists("", msg)
		if len(tn) > 0 {
			t.Fatal("tableName error: " + tableName + "/" + tn)
		}
		if result {
			t.Fatal("result error: " + msg)
		}
	}
}

func TestWords_hash_Len(t *testing.T) {
	w := server.NewWords("hash")
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
