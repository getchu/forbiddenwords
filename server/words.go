package server

import (
	"bytes"
	"fmt"
	"forbiddenwords/engine"
	"io/ioutil"
	"net/http"
	"strings"
)

type Words struct {
	tableList   []string
	serverType  []string
	engType     string
	engineList  map[string]engine.Engine
	versionList map[string]string
}

func NewWords(engType string) *Words {
	words := &Words{}
	words.tableList = []string{"political", "gun", "sexy","ad","abuse","song"}
	words.serverType = []string{"common","search", "nickname","comment","livingroom"}
	words.engType = engType
	words.engineList = make(map[string]engine.Engine)
	words.versionList = make(map[string]string)
	return words
}

//更新词库
func (w *Words) UpdateWords() error {
	for _, tableName := range w.tableList {
		for _,server := range w.serverType {
		serverName :=tableName+"_"+server
		//获取数据
        url :="http://api.changba.com/ktvbox_3rdparty.php?ac=getforbiddenwordlist&secretKey=426b78e00037cb3fd4ba3bd359a85687&format=txt&type=" + tableName +"&serverType="+server 
		resp, err := http.Get(url)
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
        if version == "" {
            fmt.Println(string(body))
        }
		//检测版本是否发生变化
		if version == w.versionList[serverName] {
			continue
		}
		var eng engine.Engine
		if w.engType == "hash" {
			eng = engine.NewHash()
		} else if w.engType == "trie" {
			eng = engine.NewTrie()
		} else {
			return fmt.Errorf("engine invalid")
		}
		for k, v := range wordList {
			//最后一行是版本号
			if k == wordListLen-1 {
				break
			}
			eng.Add(strings.ToLower(strings.TrimSpace(string(v))))
		}
		w.engineList[serverName] = eng
		w.versionList[serverName] = version
		}
	}
	return nil
}

//查找
//指定表名的话, 不指定表名就所有表都查
//指定业务线的话, 不指定业务线就差所有业务线
//返回表名和结果，结果是找到了就返回找到的字符串，找不到就返回空字符串
//tableName是指定在哪个表查，比如色情表，政治表等
//key是待查找的字符串
func (w *Words) Find(tableNameStr,serverType, key string) (string, string) {
	var res string
	key = strings.ToLower(key)
    if len(tableNameStr) ==0{
        tableNameStr="all"
    }
    if len(serverType)==0 {
        serverType="all"
    }
    if tableNameStr=="all" && serverType =="all" {
        for _, tableName := range w.tableList {
            for _, serverName := range w.serverType {
			    tableNeed  := tableName+"_"+serverName
			    if eng, ok := w.engineList[tableNeed]; ok {
				    res=eng.Find(key)
				    if(len(res)>0){
					    return tableNeed, res
				    }
			    }
            }
		}
		return tableNameStr, ""
    }
    //查询某一业务线server下所有table
    if tableNameStr == "all" && serverType !="all"{
        for _, tableName := range w.tableList {
            tableNeed :=tableName+ "_"+serverType
		    if eng, ok := w.engineList[tableNeed]; ok {
		        res=eng.Find(key)
			    if(len(res)>0){
				    return tableNeed, res
			    }
		    }
        }
        //查询公共库
        for _, tableName := range w.tableList {
            tableNeed :=tableName+ "_common"
		    if eng, ok := w.engineList[tableNeed]; ok {
		        res=eng.Find(key)
			    if(len(res)>0){
				    return tableNeed, res
			    }
		    }
        }
		return tableNameStr, ""
    }



    //查询某一业务线下的某以server
    if tableNameStr != "all" && serverType !="all"{
        tableNeed :=tableNameStr+ "_"+serverType
		if eng, ok := w.engineList[tableNeed]; ok {
		        res=eng.Find(key)
			    if(len(res)>0){
				    return tableNeed, res
			    }
		    }
        //查询公共库
            tableNeed =tableNameStr+ "_common"
		    if eng, ok := w.engineList[tableNeed]; ok {
		        res=eng.Find(key)
			    if(len(res)>0){
				    return tableNeed, res
			    }
		    }
		return tableNameStr, ""
    }



    //查询某一table下的某一个业务线server或者全部server
    if tableNameStr != "all" && serverType =="all"{
        tableNameList2 := strings.Split(tableNameStr, ",")
		for _, tableName := range tableNameList2 {
            for _, serverName := range w.serverType {
			    tableNeed  := tableName+"_"+serverName
			    if eng, ok := w.engineList[tableNeed]; ok {
				    res=eng.Find(key)
				    if(len(res)>0){
					    return tableNeed, res
				    }
			    }
            }
        }
		return tableNameStr, ""
    }

	return "", ""
}

//是否存在
//指定表名的话, 不指定表名就所有表都查
//返回表名和结果
func (w *Words) IsExists(tableName, key string) (string, bool) {
	key = strings.ToLower(key)
	if len(tableName) > 0 {
		if eng, ok := w.engineList[tableName]; ok {
			return tableName, eng.IsExists(key)
		}
		return tableName, false
	}
	for tableName, eng := range w.engineList {
		result := eng.IsExists(key)
		if result {
			return tableName, result
		}
	}
	return "", false
}

//获取当前库容量
func (w *Words) Len() int {
	l := 0
	for _, e := range w.engineList {
		l += e.Len()
	}
	return l
}

//仅仅是为了_test文件用的，外部请勿调用
func (w *Words) GetEngineList() map[string]engine.Engine {
	return w.engineList
}
