package web

import (
	"strings"
)

type trie struct {
	pattern string
	part string
	children []*trie
	wild bool //是否包含:或者*
}

func (n *trie) matchChild(part string) *trie {
	for _, child := range n.children {
		if child.part == part || child.wild {
			return child
		}
	}

	return nil
}

func (n *trie) matchChildren(part string) []*trie{

	nodes := make([]*trie,0)

	for _,child := range n.children{
		if child.part == part || child.wild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

//example：
//pattern : a/b/c
//parts: [a,b,c]
//height: 0 (start form  ZERO)
func (n *trie)add(pattern string, parts [] string, height int )  {
	//only the last part has pattern
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//current part
	part := parts[height]
	//search Part in the current node
	child := n.matchChild(part)
	//not have
	if child == nil {
		child = &trie{part:part, wild: part[0] == '*' || part[0] == ':'}
		//child node append to parent node
		n.children = append(n.children, child)
	}
	//recursion
	child.add(pattern, parts , height+1)
}


//TREE
//         a
//	/    /     \      \
// b   :name   :age  *all
// |     |       |
// c     h       g

//example : /a/m/g
//parts [a,m,g]
//height 0

func (n *trie) find(parts []string, height int)* trie{
	//end condition
	//last part || part begin with "*"
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		//recursion
		result := child.find(parts, height+1)
		if result !=nil {
			return result
		}
	}
	return nil
}





