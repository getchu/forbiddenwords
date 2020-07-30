package engine

import (
	"fmt"
)

//以前缀数的方式存储, 有2层，第一层是内容的第一个字，第二层是所有的字

//定义子节点
type childList map[string]*node

//定义一个节点
type node struct {
	key       string
	childList childList
	//是否到达了叶子节点，不然len(childList)效率肯定没有直接end判断高
	end bool
}

func NewNode(key string) *node {
	return &node{
		key:       key,
		childList: make(childList),
		end:       false,
	}
}

//添加节点的子节点
func (n *node) addChild(key string, nod *node) error {
	n.childList[key] = nod
	return nil
}

//获取节点的子节点
func (n *node) getChild(key string) *node {
	if _, ok := n.childList[key]; ok {
		return n.childList[key]
	}
	return nil
}

//前缀树
type Trie struct {
	//根节点
	rootNode *node
	//最多跳过几个字符，比如习哈哈哈近平，可以跳过3个哈, 跳过3个字，参数值就要设置为4
	//不想跳，就是1
	maxSkipRune int
}

func NewTrie() *Trie {
	return &Trie{
		rootNode:    NewNode("root"),
		maxSkipRune: 1,
	}
}

//给树增加叶子节点
func (t *Trie) Add(key string) error {
	if len(key) <= 0 {
		return fmt.Errorf("key is empty")
	}
	//按照汉字分割
	runeList := []rune(key)
	if len(runeList) <= 0 {
		return fmt.Errorf("runeList is empty")
	}
	//同一个有单字有多字，会导致查单字查不到，因为单字有子节点，所以把单字变为双字，结尾加个\n
	if len(runeList) == 1{
		runeList = append(runeList, '\n')
	}
	curNode := t.rootNode
	for _, charRune := range runeList {
		charString := string(charRune)
		childNode := curNode.getChild(charString)
		if childNode == nil {
			childNode = NewNode(charString)
			curNode.addChild(charString, childNode)
		}
		curNode = childNode
	}
	curNode.end = true
	return nil
}

//查找
//找到后返回找的字符串，找不到是空字符串
func (t *Trie) Find(key string) string {
	if len(key) <= 0 {
		return ""
	}
	runeList := []rune(key)
	runeListLen := len(runeList)
	if runeListLen <= 0 {
		return ""
	}
	for i, charRune := range runeList {
		curNode := t.rootNode
		charString := string(charRune)
		curNode = curNode.getChild(charString)
		if curNode == nil {
			continue
		}
		maxSkipRune := t.maxSkipRune
		//ASCII小于128的，不做跳跃查询
		if charRune <= 127{
			maxSkipRune = 1
		}
		//记录匹配到的词
		matchList := charString
		//看看这个单字是不是敏感词
		tmpNode := curNode.getChild("\n")
		if tmpNode != nil {
			return matchList
		}
		//往下查找是否匹配整个敏感词。k表示当前两个相邻字符的间隔，k不能大于maxSkip
		for j, k := i+1, 0; j < runeListLen && k < maxSkipRune; j = j+1 {
			k = k + 1
			charString := string(runeList[j])
			childNode := curNode.getChild(charString)
			if childNode == nil {
				continue
			}
			matchList += charString
			if childNode.end {
				return matchList
			}
			curNode = childNode
			k = 0
		}
	}
	return ""
}

//获取叶节点长度，递归
func (t *Trie) Len() int {
	sum := 0
	for _, childNode := range t.rootNode.childList {
		t.len(childNode, &sum)
	}
	return sum
}

func (t *Trie) len(n *node, sum *int) {
	if len(n.childList) <= 0 {
		*(sum)++
		return
	}
	for _, nod := range n.childList {
		t.len(nod, sum)
	}
}

//查找
//找的字符串返回true，否则返回false
func (t *Trie) IsExists(key string) bool {
	if len(key) <= 0 {
		return false
	}
	runeList := []rune(key)
	runeListLen := len(runeList)
	if runeListLen <= 0 {
		return false
	}
	for i, charRune := range runeList {
		curNode := t.rootNode
		charString := string(charRune)
		curNode = curNode.getChild(charString)
		if curNode == nil {
			continue
		}
		//记录匹配到的词
		matchList := charString
		//往下查找是否匹配整个敏感词。k表示当前两个相邻字符的间隔，k不能大于maxSkip
		for j, k := i+1, 1; j < runeListLen && k <= t.maxSkipRune; j, k = j+1, k+1 {
			charString := string(runeList[j])
			childNode := curNode.getChild(charString)
			if childNode == nil {
				continue
			}
			matchList += charString
			if childNode.end {
				return true
			}
			curNode = childNode
			k = 0
		}
	}
	return false
}
