package typechecker

import "compiler/lexer"

func typecheck(channel chan any){
	for true{
		i := 1 == 2 == true
		if i{
			
		}
		select {
		case token := <- channel:
			switch token.(lexer.Token).Identifier{
				
			}
		}
	}
}