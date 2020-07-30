package server

import (
	"strconv"
	"testing"
)

//go test -v filter.go handle.go httphandle.go words.go regular.go httphandle_test.go

var h = NewHTTPHandle()

func TestHTTPHandle_find_symbol(t *testing.T) {
	msg := "ca-o"
	filterInt, _ := strconv.Atoi("1")
	filterMsg := h.Handle.filter.Run(filterInt, string(msg))
	//词库查找
	tableNameHTTP := "sexy"
	tableName, result := h.Handle.words.Find(tableNameHTTP, filterMsg)
	if len(result) <= 0 {
		//正则匹配
		tableName, result = h.Handle.regular.Match(tableNameHTTP, filterMsg)
	}
	if tableName != "sexy" {
		t.Fatalf("tableName error: %s / %s", tableName, tableNameHTTP)
	}
	if result != "cao" {
		t.Fatalf("tableName error: %s / %s", tableName, tableNameHTTP)
	}
}

func TestHTTPHandle_find_upper(t *testing.T) {
	msg := "CAO"
	//词库查找
	tableNameHTTP := "sexy"
	tableName, result := h.Handle.words.Find(tableNameHTTP, msg)
	if tableName != "sexy" {
		t.Fatalf("tableName error: %s / %s", tableName, tableNameHTTP)
	}
	if result != "cao" {
		t.Fatalf("result error: CAO / %s", result)
	}
}

func TestHTTPHandle_find_regular(t *testing.T) {
	msg := "羡慕北京的学生可以经常去天安门广场看升旗"
	//词库查找
	tableNameHTTP := "combined_political"
	h := NewHTTPHandle()
	tableName, result := h.Handle.regular.Match(tableNameHTTP, msg)
	if tableName != "combined_political" {
		t.Fatal("tableName: " + tableName + " error")
	}
	if result != "学生#广场#北京" {
		t.Fatal("result: " + result + " error")
	}
}
