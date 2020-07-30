package server

import (
	"forbiddenwords/server"
	"testing"
)

func TestFilter_UpdateFilter(t *testing.T) {
	f := server.NewFilter()
	f.UpdateFilter()
	filterList := f.GetFilterStringList()
	if len(filterList) <= 0 {
		t.Fatal("len error")
	}
	for tableName, filter := range filterList {
		if len(filter) <= 0 {
			t.Fatal(tableName + " len error")
		}
	}
}

func TestFilter_UpdateFilterTwice(t *testing.T) {
	r := server.NewFilter()
	r.UpdateFilter()
	r.UpdateFilter()
}

func TestFilter_Run(t *testing.T) {
	f := server.NewFilter()
	f.UpdateFilter()
	var str, strNew string

	str = "c-A-o"
	strNew = f.Run(server.FILTER_SYMBOL, str)
	if strNew != "cao" {
		t.Fatalf("filter symbol error: %s / %s", strNew, str)
	}

	str = "cA\xF0\x9F\x98\x81o"
	strNew = f.Run(server.FILTER_EMOJI, str)
	if strNew != "cao" {
		t.Fatalf("filter emoji error: %s / %s", strNew, str)
	}

	str = "cA万岁毛主席o"
	strNew = f.Run(server.FILTER_WHITE_WORD, str)
	if strNew != "cao" {
		t.Fatalf("filter white word error: %s / %s", strNew, str)
	}

	str = "c-A万\xF0\x9F\x98\x81岁毛主席o"
	strNew = f.Run(server.FILTER_EMOJI|server.FILTER_WHITE_WORD, str)
	if strNew != "c-ao" {
		t.Fatalf("filter white word error: %s / %s", strNew, str)
	}

	str = "c-A万\xF0\x9F\x98\x81岁毛主席o"
	strNew = f.Run(server.FILTER_SYMBOL|server.FILTER_EMOJI|server.FILTER_WHITE_WORD, str)
	if strNew != "cao" {
		t.Fatalf("filter white word error: %s / %s", strNew, str)
	}
}
