package optimizer_test

import (
	"testing"

	"github.com/FlamingTree/expr/ast"
	"github.com/FlamingTree/expr/checker"
	"github.com/FlamingTree/expr/internal/conf"
	"github.com/FlamingTree/expr/optimizer"
	"github.com/FlamingTree/expr/parser"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptimize_constant_folding(t *testing.T) {
	tree, err := parser.Parse(`[1,2,3][5*5-25]`)
	require.NoError(t, err)

	optimizer.Optimize(&tree.Node)

	expected := &ast.IndexNode{
		Node:  &ast.ConstantNode{Value: []int{1, 2, 3}},
		Index: &ast.IntegerNode{Value: 0},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(tree.Node))
}

func TestOptimize_in_array(t *testing.T) {
	config := conf.New(map[string]int{"v": 0})

	tree, err := parser.Parse(`v in [1,2,3]`)
	require.NoError(t, err)

	_, err = checker.Check(tree, config)
	require.NoError(t, err)

	optimizer.Optimize(&tree.Node)

	expected := &ast.BinaryNode{
		Operator: "in",
		Left:     &ast.IdentifierNode{Value: "v"},
		Right:    &ast.ConstantNode{Value: optimizer.Map{1: {}, 2: {}, 3: {}}},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(tree.Node))
}

func TestOptimize_in_range(t *testing.T) {
	tree, err := parser.Parse(`age in 18..31`)
	require.NoError(t, err)

	optimizer.Optimize(&tree.Node)

	left := &ast.IdentifierNode{
		Value: "age",
	}
	expected := &ast.BinaryNode{
		Operator: "and",
		Left: &ast.BinaryNode{
			Operator: ">=",
			Left:     left,
			Right: &ast.IntegerNode{
				Value: 18,
			},
		},
		Right: &ast.BinaryNode{
			Operator: "<=",
			Left:     left,
			Right: &ast.IntegerNode{
				Value: 31,
			},
		},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(tree.Node))
}

func TestOptimize_const_range(t *testing.T) {
	tree, err := parser.Parse(`-1..1`)
	require.NoError(t, err)

	optimizer.Optimize(&tree.Node)

	expected := &ast.ConstantNode{
		Value: []int{-1, 0, 1},
	}

	assert.Equal(t, litter.Sdump(expected), litter.Sdump(tree.Node))
}
