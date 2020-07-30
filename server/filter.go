package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	//过滤字符
	FILTER_SYMBOL = 1 << iota
	//过滤EMOJI
	FILTER_EMOJI
	//过滤WhiteWordList
	FILTER_WHITE_WORD
)

type Filter struct {
	tableList   [3]string
	filterList  map[string][]string
	versionList map[string]string
}

func NewFilter() *Filter {
	filter := &Filter{}
	filter.tableList = [3]string{"white_list", "symbol", "emoji"}
	filter.filterList = make(map[string][]string)
	filter.versionList = make(map[string]string)
	return filter
}

//更新过滤词库
func (f *Filter) UpdateFilter() error {
	for _, tableName := range f.tableList {
		//获取数据
		//symbol和emoji从本地加载
		if tableName == "symbol" || tableName == "emoji" {
			f.updateFilterFromLocal(tableName)
			continue
		}
		resp, err := http.Get("http://changba.com/changbalab/forbiddenwordlist.php?secretKey=5e6928cfa7aec38387d996dae156caaf&format=txt&type=" + tableName)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if len(body) <= 0 {
			continue
		}
		wordList := bytes.Split(body, []byte("\n"))
		wordListLen := len(wordList)
		version := strings.ToLower(strings.TrimSpace(string(wordList[wordListLen-1])))
		//检测版本是否发生变化
		if version == f.versionList[tableName] {
			continue
		}
		f.filterList[tableName] = make([]string, wordListLen-1)
		for k, v := range wordList {
			//最后一行是版本号
			if k == wordListLen-1 {
				break
			}
			f.filterList[tableName][k] = strings.ToLower(strings.TrimSpace(string(v)))
		}
		f.versionList[tableName] = version
	}
	return nil
}

//更新过滤词库 从本地
func (f *Filter) updateFilterFromLocal(tableName string) {
	var wordList []string
	if tableName == "symbol" {
		wordList = symbolList
	} else if tableName == "emoji" {
		wordList = emojiList
	} else {
		return
	}
	wordListLen := len(wordList)
	version := strings.ToLower(strings.TrimSpace(wordList[wordListLen-1]))
	//检测版本是否发生变化
	if version == f.versionList[tableName] {
		return
	}
	f.filterList[tableName] = make([]string, wordListLen-1)
	for k, v := range wordList {
		//最后一行是版本号
		if k == wordListLen-1 {
			break
		}
		f.filterList[tableName][k] = strings.TrimSpace(v)
	}
	f.versionList[tableName] = version
}

func (f Filter) Run(filter int, str string) string {
	str = strings.ToLower(str)
	if filter & FILTER_SYMBOL > 0 {
		str = f.replace("symbol", str)
	}
	if filter & FILTER_EMOJI > 0 {
		str = f.replace("emoji", str)
	}
	if filter & FILTER_WHITE_WORD > 0 {
		str = f.replace("white_list", str)
	}
	return str
}

//过滤字符
func (f Filter) replace(tableName, str string) string {
	if filterList, ok := f.filterList[tableName]; ok {
		for _, filter := range filterList {
			str = strings.Replace(str, filter, "", -1)
		}
	}
	return str
}

//仅仅是为了_test文件用的，外部请勿调用
func (f *Filter) GetFilterStringList() map[string][]string {
	return f.filterList
}
