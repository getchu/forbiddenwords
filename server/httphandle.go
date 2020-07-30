package server

import (
	"encoding/json"
	"forbiddenwords/lib"
	"net/http"
	"strconv"
	"time"
)

type HTTPHandle struct {
	Handle *Handle
	cache  *lib.Cache
}

type httpContext struct{
	startTime int64
	endTime int64
	ac string
}

func newHttpContext(ac string)*httpContext{
	return &httpContext{
		startTime: time.Now().Unix(),
		ac: ac,
	}
}

func NewHTTPHandle(config *lib.Config) *HTTPHandle {
	h := &HTTPHandle{
		Handle: NewHandle(config),
		cache:  lib.NewCache(),
	}
	return h
}

//HTTP请求执行前会先调用
func (h *HTTPHandle)before(hc *httpContext){
	return
}

//HTTP请求执行后会再调用
func (h *HTTPHandle)after(hc *httpContext, msg string, resp map[string]interface{}){
	hc.endTime = time.Now().Unix()
	//大于两秒的打日志
	handleTime := hc.endTime - hc.startTime
	if handleTime >= 1{
		h.Handle.config.Logger.Output(lib.LOG_WARNING, "[SLOWLOG] ac: " + hc.ac + " msg: " + msg + " result: " + resp["result"].(string) + " filterMsg: " + resp["filterMsg"].(string) + " tableName: " + resp["tableName"].(string) + " time: " + strconv.Itoa(int(handleTime)))
	}
	if hc.ac == "find" && len(resp["result"].(string)) > 0{
		h.Handle.config.Logger.Output(lib.LOG_DEBUG, "[MATCH] msg: " + msg + " result: " + resp["result"].(string) + " filterMsg: " + resp["filterMsg"].(string) + " tableName: " + resp["tableName"].(string) + "")
	}
	return
}

func (h *HTTPHandle) output(w http.ResponseWriter, result interface{}, errorcode string) {
	resp := make(map[string]interface{})
	resp["result"] = result
	resp["errorcode"] = errorcode
	respJson, err := json.Marshal(resp)
	if err != nil {
		resp["errorcode"] = "json失败"
	}
	w.Write(respJson)
}

func (h *HTTPHandle) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello forbidden words service"))
}

//查找是否包含敏感词
func (h *HTTPHandle) find(w http.ResponseWriter, r *http.Request) {
	hc := newHttpContext("find")
	h.before(hc)
	msg := r.FormValue("msg")
	filterString := r.FormValue("filter")
	matchWayString := r.FormValue("matchWay")
	if len(msg) <= 0 {
		w.WriteHeader(400)
		h.output(w, "", "msg不存在")
		return
	}
	if len(matchWayString) <= 0 {
		w.WriteHeader(400)
		h.output(w, "", "matchWay不存在")
		return
	}
	//是否过滤emoji、符号等
	filterInt, _ := strconv.Atoi(filterString)
	filterMsg := h.Handle.filter.Run(filterInt, string(msg))
	//匹配方式：词库匹配，正则匹配
	matchWayInt, _ := strconv.Atoi(matchWayString)
	var tableName, result string
	//词库匹配
	if len(result) <= 0 && matchWayInt & MATCH_WAY_WORDS > 0 {
		wordsTableName := r.FormValue("wordsTableName")
		serverType := r.FormValue("serverType")
		tableName, result = h.Handle.words.Find(wordsTableName, serverType,filterMsg)
	}
	//正则匹配
	if len(result) <= 0 && matchWayInt & MATCH_WAY_REGULAR > 0 {
		regularTableName := r.FormValue("regularTableName")
		tableName, result = h.Handle.regular.Match(regularTableName, filterMsg)
	}
	resp := make(map[string]interface{})
	resp["result"] = result
	resp["filterMsg"] = filterMsg
	resp["tableName"] = tableName
	h.output(w, resp, "ok")
	h.after(hc, msg, resp)
}

//获取敏感词等级
func (h *HTTPHandle) forbidLevel(w http.ResponseWriter, r *http.Request) {
	hc := newHttpContext("forbidLevel")
	h.before(hc)
	msg := r.FormValue("msg")
	filterString := r.FormValue("filter")
	serverType := r.FormValue("serverType")
    wordsTableName := r.FormValue("wordsTableName")
    if len(serverType) == 0{
        serverType = "all"
    }
    if len(wordsTableName) == 0{
        wordsTableName = "all"
    }

	if len(msg) <= 0 {
		w.WriteHeader(400)
		h.output(w, "", "msg不存在")
		return
	}
	//是否过滤emoji、符号等
	filterInt, _ := strconv.Atoi(filterString)
	filterMsg := h.Handle.filter.Run(filterInt, string(msg))
	//1001是安全
	level := 1001
    num :=0
	tableGroupList := map[string]map[string]int{
		"words": map[string]int{"political": 1004, "gun": 1004, "sexy": 1002,"ad":1002,"abuse":1002,"song":1002},
		"regular": map[string]int{"combined_political": 1004, "reglar_viviantibetan": 1003},
	}
	var tableName, result string
    for group, tableList := range tableGroupList{
        for tn, lev := range tableList{
            if group == "words"{
                //类别或者业务线查找,table的只查找一次
                if len(wordsTableName)>0&&wordsTableName!="all"{
                    tn=wordsTableName
                    if num>0{
                        continue
                    }
                    num = 1
                }
			    tableName, result = h.Handle.words.Find(tn,serverType, filterMsg)
			    if len(result) > 0{
				    level = lev
				    break
			    }
            }
            if group == "regular"{
		        tableName, result = h.Handle.regular.Match(tn, filterMsg)
			    if len(result) > 0{
				    level = lev
				    break
			    }
		    }
        }
		if level != 1001{
			break
        }
	}
	resp := make(map[string]interface{})
	resp["result"] = result
	resp["filterMsg"] = filterMsg
	resp["tableName"] = tableName
	resp["level"] = level
	h.output(w, resp, "ok")
	h.after(hc, msg, resp)
}

func (h *HTTPHandle) update(w http.ResponseWriter, r *http.Request) {
	h.Handle.updateChan <- 1
	h.output(w, "success", "ok")
}
