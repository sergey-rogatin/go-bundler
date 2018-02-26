package htmlLoader

func printAst(n ast) string {
	switch n.t {
	case n_TEXT_CHUNK,
		n_ATTRIBUTE_VALUE:
		return n.value

	case n_TAG_CHILDREN,
		n_DOCUMENT:
		res := ""
		for _, c := range n.children {
			res += printAst(c)
		}
		return res

	case n_TAG_ATTRIBUTES:
		res := ""
		for _, c := range n.children {
			res += " " + printAst(c)
		}
		return res

	case n_ATTRIBUTE:
		val := n.children[0]
		res := n.value
		if val.t == n_ATTRIBUTE_VALUE {
			res += "=" + printAst(val)
		}
		return res

	case n_TAG:
		attrs := n.children[0]
		childs := n.children[1]
		res := "<" + n.value + printAst(attrs)
		if childs.children != nil && len(childs.children) > 0 {
			res += ">" + printAst(childs) + "</" + n.value + ">"
		} else {
			res += "/>"
		}
		return res
	}

	return "SOMETHING WENT WRONG"
}
