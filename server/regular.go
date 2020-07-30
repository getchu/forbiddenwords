package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Regular struct {
	tableList   [2]string
	regularList map[string][]string
	versionList map[string]string
}

func NewRegular() *Regular {
	regular := &Regular{}
	regular.tableList = [2]string{"combined_political", "reglar_viviantibetan"}
	regular.regularList = make(map[string][]string)
	regular.versionList = make(map[string]string)
	return regular
}

//更新正则库
func (r *Regular) UpdateRegular() error {
	for _, tableName := range r.tableList {
		//获取数据
        //url :="http://changba.com/changbalab/forbiddenwordlist.php?secretKey=5e6928cfa7aec38387d996dae156caaf&format=txt&type=" + tableName
        url :="http://api.changba.com/ktvbox_3rdparty.php?ac=getforbiddenwordlist&secretKey=426b78e00037cb3fd4ba3bd359a85687&format=txt&type=" + tableName
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
		//检测版本是否发生变化
		if version == r.versionList[tableName] {
			continue
		}
		r.regularList[tableName] = make([]string, wordListLen-1)
		for k, v := range wordList {
			//最后一行是版本号
			if k == wordListLen-1 {
				break
			}
			//库reglar_viviantibetan不小写
			if tableName == "reglar_viviantibetan"{
				r.regularList[tableName][k] = strings.TrimSpace(string(v))
			}else{
				r.regularList[tableName][k] = strings.ToLower(strings.TrimSpace(string(v)))
			}
		}
		r.versionList[tableName] = version
	}
	return nil
}

//匹配
//指定表名的话, 不指定表名就所有表都查
//返回表名和结果，结果是找到了就返回找到的字符串，找不到就返回空字符串
//这段逻辑是从PHP直接挪过来的
func (r *Regular) Match(tableNameStr, str string) (string, string) {
	regularList := make(map[string][]string)
	if len(tableNameStr) > 0 {
		tableNameList := strings.Split(tableNameStr, ",")
		for _, tableName := range tableNameList {
			if regStringList, ok := r.regularList[tableName]; ok {
				regularList[tableName] = make([]string, 0)
				regularList[tableName] = regStringList
			} else {
				return tableName, ""
			}
		}
	} else {
		regularList = r.regularList
	}
	for tn, regStringList := range regularList {
		for _, regString := range regStringList {
			regList := strings.Split(regString, "#")
			isMatch := true
			for _, reg := range regList {
				result := regexp.MustCompile(reg).FindAllString(str, -1)
				if result == nil || len(result) <= 0 {
					isMatch = false
					break
				}
			}
			if isMatch {
				tableName := ""
				if len(tableNameStr) <= 0 {
					tableName = tn
				}else{
					tableName = tableNameStr
				}
				return tableName, regString
			}
		}
	}
	return "", ""
}

//仅仅是为了_test文件用的，外部请勿调用
func (r *Regular) GetRegStringList() map[string][]string {
	return r.regularList
}
