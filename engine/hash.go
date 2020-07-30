package engine

import (
	"fmt"
	"strings"
)

//以Hash的方式存储, 有2层，第一层是内容的第一个字，第二层是所有的字

type Hash struct {
	//数据
	dataList map[rune]map[string]struct{}
}

func NewHash() *Hash {
	return &Hash{
		dataList: make(map[rune]map[string]struct{}),
	}
}

func (h *Hash) Add(key string) error {
	if len(key) <= 0 {
		return fmt.Errorf("key is empty")
	}
	runeList := []rune(key)
	if len(runeList) <= 0 {
		return fmt.Errorf("runeList is empty")
	}
	if _, ok := h.dataList[runeList[0]]; !ok {
		h.dataList[runeList[0]] = make(map[string]struct{})
	}
	h.dataList[runeList[0]][key] = empty
	return nil
}

//查找
//找到后返回找的字符串，找不到是空字符串
func (h *Hash) Find(key string) string {
	if len(key) <= 0 {
		return ""
	}
	runeList := []rune(key)
	if len(runeList) <= 0 {
		return ""
	}
	for _, char := range runeList {
		if _, ok := h.dataList[char]; !ok {
			continue
		}
		for data, _ := range h.dataList[char] {
			isContains := strings.Contains(key, data)
			if isContains {
				return data
			}
		}
	}
	return ""
}

func (h *Hash) Len() int {
	sum := 0
	for _, runList := range h.dataList {
		sum += len(runList)
	}
	return sum
}

//查找
//找的字符串返回true，否则返回false
func (h *Hash) IsExists(key string) bool {
	if len(key) <= 0 {
		return false
	}
	runeList := []rune(key)
	if len(runeList) <= 0 {
		return false
	}
	if _, ok := h.dataList[runeList[0]]; !ok {
		return false
	}
	if _, ok := h.dataList[runeList[0]][key]; ok {
		return true
	}
	return false
}
