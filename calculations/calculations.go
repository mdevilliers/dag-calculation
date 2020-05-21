package calculations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

const (
	// TODO: consider removing these and replace with something like
	// "expected" or "default"
	VesselUUIDDependancy         = "vessel_uuid"
	AccountUUIDDependancy        = "account_uuid"
	EventTimeDependancy          = "event_time"
	LocationLongitudeDependancy  = "location_longitude"
	LocationLatitudeDependancy   = "location_latitude"
	WeatherSwellHeightDependancy = "weather_swell_height_dependancy"
)

type DataRow struct {
	// should already exist
	VesselUUID  string
	AccountUUID string
	EventTime   time.Time
	Location    Location
	Weather     Weather

	// calculated
	DistanceOverGround null.Float
	SpeedOverGround    null.Float
	BeaufortForce      null.Float
	DouglasSeaState    null.Int
}

type Location struct {
	Longitude null.Float
	Latitude  null.Float
}

type Weather struct {
	SwellHeight null.Float
}

func NewResolver() *resolver {
	return &resolver{}
}

type resolver struct {
	resolved []string
}

func (r *resolver) RegisterResolved(dep string) {
	r.resolved = append(r.resolved, dep) // TODO check for duplicates
}

func (r *resolver) CollapseAll(data []*DataRow, f *functionset) error {

	for f.Len() != 0 {
		if err := r.Collapse(data, f); err != nil {
			return err
		}
	}

	return nil
}

func (r *resolver) Collapse(data []*DataRow, f *functionset) error {
	// TODO : should collapse until no more columns can be calculated
	// TODO : could invoke in parrallel on multiple go routines if a useful optimisation

	// what columns already exist
	// which new columns can be calculated
	todo := f.CanInvoke(r.resolved)

	if len(todo) == 0 && f.Len() != 0 {
		return errors.New("booyah deadlock!")
	}

	// invoke one at a time
	for _, t := range todo {

		// TODO : pass in rather than hard code
		runtime := runtime{
			Ctx:    context.Background(),
			Logger: newNoopLogger(),
		}

		err := t.Fn(runtime, data)

		if err != nil {
			return err
		}

		// no error so register function as being completed
		r.RegisterResolved(t.Satisfies)

		// remove from functionset
		f.Remove(t.Satisfies)
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

func (f *functionset) Len() int {
	return len(f.registered)
}

type registeredFunc struct {
	Requires  []string
	Satisfies string
	Fn        doIteratorFn
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
