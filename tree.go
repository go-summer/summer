package rayRoute

import (
	"fmt"
	"net/http"
	"strings"
)

type Value http.HandlerFunc

// Node 字典树节点
type Node struct {
	label    byte
	prefix   string
	children []*Node
	value    Value
}

// InsertNode 添加节点 TODO 未判断前缀与 urlStr 相等的情况
func (n *Node) InsertNode(urlStr string, value Value) {

	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			fmt.Println("urlStr:", urlStr)
		}
	}()

	// 无子节点
	if n.children == nil {
		n.children = make([]*Node, 0)
	}

	// 根节点无数据
	if n.value == nil {
		n.value = value
		n.prefix = urlStr
		n.label = urlStr[0]
		return
	}

	// 检查当前节点是否可分裂
	if n.label == urlStr[0] {
		if urlStr == n.prefix {
			// 与当前节点重复
			n.value = value
		} else if strings.HasPrefix(urlStr, n.prefix) {
			// prefix 为共同前缀 不需分裂 进入递归
			n.InsertNode(strings.TrimLeft(urlStr, n.prefix), value)
		} else {
			// prefix 不为共同前缀 分裂节点
			commonPrefix, _ := findCommonPrefix(n.prefix, urlStr)
			// 切掉共有前缀 递归传递
			n.InsertNode(strings.TrimLeft(urlStr, commonPrefix), value)
			n.InsertNode(strings.TrimLeft(n.prefix, commonPrefix), n.value)
			n.value = nil
		}
	}

	// 遍历子节点 查找挂载点（共有前缀最长）
	maxLength := 0
	maxLengthNode := -1
	for i, v := range n.children {
		if v.label == urlStr[0] {
			_, tempLength := findCommonPrefix(v.prefix, urlStr)
			if tempLength > maxLength {
				maxLengthNode = i
			}
		}
	}

	// 找到挂载点 进行挂载
	if maxLengthNode != -1 {
		commonPrefix, _ := findCommonPrefix(n.children[maxLengthNode].prefix, urlStr)
		n.children[maxLengthNode].
			InsertNode(strings.TrimLeft(urlStr, commonPrefix), value)
	}

	// 无挂载点直接新建一个子节点
	n.children = append(n.children,
		&Node{label: urlStr[0], prefix: urlStr, value: value})
}

// FindNode 查找节点
func (n *Node) FindNode(urlStr string) Value {
	// 匹配到当前节点 直接返回
	if n.prefix == urlStr {
		return n.value
	}

	// 查找子节点
	for _, v := range n.children {
		if v.label == urlStr[0] {
			if strings.HasPrefix(urlStr, v.prefix) {
				if len(urlStr) == len(v.prefix) {
					return v.value
				} else {
					return v.FindNode(urlStr)
				}
			}
		}
	}

	return nil
}

func findCommonPrefix(s1 string, s2 string) (string, int) {
	for i, _ := range s1 {
		if s1[i] != s2[i] {
			return s1[0:i], i
		}
	}
	return s1, len(s1)
}

func PrintTree(node Node) {
	nodeMap := make(map[int][]string)

	sortNode(node, 0, nodeMap)

	for i := 0; i < len(nodeMap); i++ {
		fmt.Println(nodeMap[i])
	}
}

func sortNode(node Node, level int, nodeMap map[int][]string) {

	if nodeMap[level] == nil {
		nodeMap[level] = make([]string, 0)
	}

	for _, v := range node.children {
		nodeMap[level] = append(nodeMap[level], v.prefix)
		sortNode(*v, level+1, nodeMap)
	}
}
