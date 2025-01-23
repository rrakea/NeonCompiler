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
	if tree.Leaf.Name == name {
		return []*ParseTree{tree}
	}
	for _, br := range tree.Branches {
		ret = append(ret, br.Search_top_occurences(name)...)
	}
	return ret
}

func (tree *ParseTree) Search_top_occurences (name string) []*ParseTree {
	if tree.Leaf.Name == name {
		return []*ParseTree{tree}
	}

	ret := []*ParseTree{}
	for _, br := range tree.Branches {
		ret = append(ret, br.Search_top_occurences(name)...)
	}
	return ret
}


func (tree *ParseTree) Search_first_occurenc_depth (name string) *ParseTree {
	for _, br := range tree.Branches {
		if br.Leaf.Name == name {
			return &br
		} else {
			for _, ch := range br.Branches {
				search := ch.Search_first_occurenc_depth(name)
				if  search != nil {
					return search
				}
			} 
		}
	}
	return nil
}