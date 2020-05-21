package calculations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
)

func NewFunctionSet() *functionset {
	return &functionset{
		registered: map[string]registeredFunc{},
		graph:      &dag.AcyclicGraph{},
	}
}

type functionset struct {
	registered map[string]registeredFunc
	graph      *dag.AcyclicGraph
}

type doIteratorFn func(r runtime, rows []*DataRow) error
type doIterator interface {
	Fn(r runtime, rows []*DataRow) error
}

func (f *functionset) RegisterIteratorFn(satisfies string, fn doIteratorFn, requires ...string) error {
	return f.registerInternal(registeredFunc{
		Requires:  requires,
		Satisfies: satisfies,
		Fn:        fn,
	})
}

func (f *functionset) RegisterIterator(satisfies string, it doIterator, requires ...string) error {
	return f.registerInternal(registeredFunc{
		Requires:  requires,
		Satisfies: satisfies,
		Fn:        it.Fn,
	})
}

func (f *functionset) registerInternal(r registeredFunc) error {
	if _, exists := f.registered[r.Satisfies]; exists {
		return fmt.Errorf("iterator already registered : %s", r.Satisfies)
	}
	f.registered[r.Satisfies] = r
	return nil
}

func (f *functionset) Build() error {

	f.graph.Add("ROOT") // TODO : swap out with a real root node

	for _, r := range f.registered {
		f.graph.Add(r.Satisfies)
		f.graph.Connect(dag.BasicEdge("ROOT", r.Satisfies))

		for _, s := range r.Requires {
			f.graph.Connect(dag.BasicEdge(r.Satisfies, s))
		}
	}
	//fmt.Println(f.graph.StringWithNodeTypes())
	f.graph.TransitiveReduction()
	//fmt.Println(string(f.graph.Dot(&dag.DotOpts{Verbose: true, DrawCycles: true})))

	return f.graph.Validate()
}

func (f *functionset) Collapse(data []*DataRow) error {
	f.graph.Walk(func(v dag.Vertex) tfdiags.Diagnostics {

		registered, found := f.registered[v.(string)]

		if !found {
			// TODO : work out how to return a Diagnostic
			return nil
		}

		// TODO : pass in rather than hard code
		runtime := runtime{
			Ctx:    context.Background(),
			Logger: newNoopLogger(),
		}

		err := registered.Fn(runtime, data)

		if err != nil {
			// TODO : work out how to return a Diagnostic
			return nil
		}

		return nil
	})
	return nil
}

type registeredFunc struct {
	Requires  []string
	Satisfies string
	Fn        doIteratorFn
}
