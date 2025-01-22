package parser

func (tree *ParseTree) Search_direct_children (name string) []*ParseTree {
	ret := []*ParseTree{}
	for _, c := range tree.Branches {
		if c.Leaf.Name == name {
			ret = append(ret, &c)
		}
	}
	return ret
}

func (tree *ParseTree) Search_first_child (name string) *ParseTree {
	for _, c := range tree.Branches {
		if c.Leaf.Name == name {
			return &c
		}
	}
	return nil
}

func (tree *ParseTree) Search_tree (name string) []*ParseTree {
	ret := []*ParseTree{}
	for _, br := range tree.Branches {
		if br.Leaf.Name == name {
			ret = append(ret, &br)
		}
		for _, ch := range br.Branches {
			ret = append(ret, ch.Search_tree(name)...)
		} 
	}
	return ret
}

func (tree *ParseTree) Search_top_occurences (name string) []*ParseTree {
	ret := []*ParseTree{}
	for _, br := range tree.Branches {
		if br.Leaf.Name == name {
			ret = append(ret, &br)
		} else {
			for _, ch := range br.Branches {
				ret = append(ret, ch.Search_tree(name)...)
			} 
		}
	}
	return ret
}