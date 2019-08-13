package checker

import (
	"github.com/FlamingTree/expr/ast"
	"github.com/FlamingTree/expr/internal/conf"
	"github.com/FlamingTree/expr/parser"
)

type operatorPatcher struct {
	ops   map[string][]string
	types conf.TypesTable
}

func (p *operatorPatcher) Enter(node *ast.Node) {}
func (p *operatorPatcher) Exit(node *ast.Node) {
	binaryNode, ok := (*node).(*ast.BinaryNode)
	if !ok {
		return
	}

	fns, ok := p.ops[binaryNode.Operator]
	if !ok {
		return
	}

	leftType := binaryNode.Left.GetType()
	rightType := binaryNode.Right.GetType()

	_, fn, ok := conf.FindSuitableOperatorOverload(fns, p.types, leftType, rightType)
	if ok {
		newNode := &ast.FunctionNode{
			Name:      fn,
			Arguments: []ast.Node{binaryNode.Left, binaryNode.Right},
		}
		newNode.SetType((*node).GetType())
		newNode.SetLocation((*node).GetLocation())
		*node = newNode
	}
}

func PatchOperators(tree *parser.Tree, config *conf.Config) {
	if len(config.Operators) == 0 {
		return
	}
	patcher := &operatorPatcher{ops: config.Operators, types: config.Types}
	ast.Walk(&tree.Node, patcher)
}
