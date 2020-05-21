package calculations

import (
	"testing"

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
	}, "one")
	fs.RegisterIteratorFn("five", func(r runtime, rows []*DataRow) error {
		return nil
	}, "three")

	err := fs.Build()
	require.Nil(t, err)

	err = fs.Collapse(nil)
	require.Nil(t, err)
}
