package calculations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/stretchr/testify/require"
)

func Test_FunctionSet_Build(t *testing.T) {

	fs := NewFunctionSet()

	fs.RegisterIteratorFn("one", func(r runtime, rows []*DataRow) error {
		return nil
	})
	fs.RegisterIteratorFn("two", func(r runtime, rows []*DataRow) error {
		return nil
	}, "one")
	fs.RegisterIteratorFn("three", func(r runtime, rows []*DataRow) error {
		return nil
	})
	/*	fs.RegisterIteratorFn("five", func(r runtime, rows []*DataRow) error {
			return nil
		}, "three")
	*/
	err := fs.Build()
	require.Nil(t, err)

	fmt.Println(string(fs.graph.Dot(&dag.DotOpts{Verbose: true, DrawCycles: true})))
	root, err := fs.graph.Root()
	require.Nil(t, err)

	fmt.Println(root)

	fs.graph.Walk(func(v dag.Vertex) tfdiags.Diagnostics {
		fmt.Println(v)
		return nil
	})

}
