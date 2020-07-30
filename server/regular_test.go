package server

import (
	"forbiddenwords/server"
	"testing"
)

func TestRegular_UpdateRegular(t *testing.T) {
	r := server.NewRegular()
	r.UpdateRegular()
	regularList := r.GetRegStringList()
	if len(regularList) <= 0 {
		t.Fatal("len error")
	}
	for tableName, regular := range regularList {
		if len(regular) <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
}

func TestRegular_UpdateRegularTwice(t *testing.T) {
	r := server.NewRegular()
	r.UpdateRegular()
	r.UpdateRegular()
}

func TestRegular_Match(t *testing.T) {
	str := "北京的学生每年都要去天安门广场"
	r := server.NewRegular()
	r.UpdateRegular()
	tableName, result := r.Match("combined_political", str)
	if tableName != "combined_political" {
		t.Fatal("tableName: " + tableName + " error")
	}
	if result != "学生#广场#北京" {
		t.Fatal("result: " + result + " error")
	}

	str = "مرحبا"
	tableName, result = r.Match("reglar_viviantibetan", str)
	if tableName != "reglar_viviantibetan" {
		t.Fatal("tableName: " + tableName + " error")
	}
	if result != `[\p{Arabic}]` {
		t.Fatal("result: " + result + " error")
	}
}
