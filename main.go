package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/guregu/null.v4"
)

/*

https://github.com/90poe/voyage-monitor-calculations-experiments/blob/master/Distance%20and%20Speed%20Over%20Ground.ipynb

need long, lat
- can calculate distance_over_ground

need distance_over_ground
- can calculate speed_over_grond
*/

func main() {

	rows := []*dataRow{
		&dataRow{Longitude: 12.34, Latitude: 12.34},
		&dataRow{Longitude: 12.34, Latitude: 12.34},
		&dataRow{Longitude: 12.34, Latitude: 12.34},
		&dataRow{Longitude: 12.34, Latitude: 12.34},
	}

	ds := NewResolver()
	ds.RegisterResolved("long")
	ds.RegisterResolved("lat")

	funcs := NewFunctionSet()

	funcs.RegisterIteratorFn("distance_over_ground", func(rows []*dataRow) error {
		fmt.Println("calculating distance over ground")
		// iterate the rows and od the maths
		for _, r := range rows {
			r.DistanceOverGround = null.FloatFrom(22)
		}
		return nil
	}, "long", "lat")

	funcs.RegisterIteratorFn("speed_over_ground", func(rows []*dataRow) error {
		fmt.Println("calculating speed over ground")
		return nil
	}, "distance_over_ground")

	funcs.RegisterIteratorFn("parrallel_example", func(rows []*dataRow) error {
		fmt.Println("calculating parrallel example")
		return nil
	}, "long", "lat")

	//funcs.RegisterIterator("beaufort_force", func(rows []*dataRow) error { return nil }, "wind_speed")
	//funcs.RegisterIterator("douglas_sea_state", func(rows []*dataRow) error { return nil }, "swell_height")
	//funcs.RegisterIterator("engine_distance", func(rows []*dataRow) error { return nil }, "revolution_per_minute_average")

	// TODO : add function to validate if all functions can be resolved - e.g. is the function set solvable or do we have orphans?

	fmt.Println("first run")

	// will run distance_over_ground and parrallel_example as long and lat are resolved
	err := ds.Collapse(rows, funcs)
	spew.Dump(err)

	fmt.Println("second run")

	// will run speed_over_ground as distance_over_ground is resolved
	err = ds.Collapse(rows, funcs)
	spew.Dump(err)

}

type dataRow struct {
	Longitude          float32
	Latitude           float32
	DistanceOverGround null.Float
	SpeedOverGround    null.Float
	BeaufortForce      null.Float
	DouglasSeaState    null.Float
}

func NewResolver() *resolver {
	return &resolver{}
}

type resolver struct {
	resolved []string
}

func (r *resolver) RegisterResolved(col string) {
	r.resolved = append(r.resolved, col) // TODO check for duplicates
}

func (r *resolver) Collapse(data []*dataRow, f *functionset) error {
	// TODO : should collapse until no more columns can be calculated
	// TODO : could invoke in parrallel on multiple go routines if a useful optimisation

	// what columns already exist
	// which new columns can be calculated
	todo := f.CanInvoke(r.resolved)

	// invoke one at a time
	for _, t := range todo {
		err := t.Fn(data)

		if err != nil {
			return err
		}

		// no error so register function as being completed
		r.RegisterResolved(t.Outputs)

		// remove from functionset
		f.Remove(t.Outputs)
	}

	return nil
}

func NewFunctionSet() *functionset {
	return &functionset{
		registered: map[string]registeredFunc{},
	}
}

type functionset struct {
	registered map[string]registeredFunc
}
type doIteratorFn func(rows []*dataRow) error

func (f *functionset) RegisterIteratorFn(outputs string, fn doIteratorFn, requires ...string) {
	f.registered[outputs] = registeredFunc{
		Requires: requires,
		Outputs:  outputs,
		Fn:       fn,
	}
}

func (f *functionset) CanInvoke(resolved []string) []registeredFunc {
	ret := []registeredFunc{}
	for _, r := range f.registered {
		if r.IsSatisfied(resolved) {
			ret = append(ret, r)
		}
	}
	return ret
}

func (f *functionset) Remove(o string) {
	delete(f.registered, o)
}

type registeredFunc struct {
	Requires []string
	Outputs  string
	Fn       doIteratorFn
}

func (r registeredFunc) IsSatisfied(columns []string) bool {
	if len(columns) < len(r.Requires) {
		return false
	}

	// bleurg !!!
	c := map[string]interface{}{}
	for _, col := range columns {
		c[col] = true
	}

	for _, r := range r.Requires {
		if _, found := c[r]; !found {
			return false
		}
	}
	return true
}
