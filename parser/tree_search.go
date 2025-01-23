package parser

func (tree *ParseTree) Search_direct_children(name string) []*ParseTree {
	ret := []*ParseTree{}
	for _, c := range tree.Branches {
		if c.Leaf.Name == name {
			ret = append(ret, &c)
		}
	}
	return ret
}

func (tree *ParseTree) Search_first_child(name string) *ParseTree {
	for _, c := range tree.Branches {
		if c.Leaf.Name == name {
			return &c
		}
	}
	return nil
}

func (tree *ParseTree) Search_tree(name string) []*ParseTree {
	ret := []*ParseTree{}
	if tree.Leaf.Name == name {
		return []*ParseTree{tree}
	}
	for _, br := range tree.Branches {
		ret = append(ret, br.Search_top_occurences(name)...)
	}
	return ret
}

func (tree *ParseTree) Search_top_occurences(name string) []*ParseTree {
	if tree.Leaf.Name == name {
		return []*ParseTree{tree}
	}

	ret := []*ParseTree{}
	for _, br := range tree.Branches {
		ret = append(ret, br.Search_top_occurences(name)...)
	}
	return ret
}

func (tree *ParseTree) Search_first_occurence_depth(name string) *ParseTree {
	if tree.Leaf.Name == name {
		return tree
	}
	for _, br := range tree.Branches {
		search := br.Search_first_occurence_depth(name)
		if search != nil {
			return search
		}
	}
	return nil
}
