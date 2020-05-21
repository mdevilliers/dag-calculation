package main

import (
	"gopkg.in/guregu/null.v4"

	"github.com/davecgh/go-spew/spew"
	"github.com/mdevilliers/dag-calculation/calculations"
)

func main() {

	rows := []*calculations.DataRow{
		{Location: calculations.Location{Longitude: null.FloatFrom(12.34), Latitude: null.FloatFrom(12.34)}, Weather: calculations.Weather{SwellHeight: null.FloatFrom(0.48)}},
		{Location: calculations.Location{Longitude: null.FloatFrom(12.35), Latitude: null.FloatFrom(12.35)}},
		{Location: calculations.Location{Longitude: null.FloatFrom(12.36), Latitude: null.FloatFrom(12.36)}},
		{Location: calculations.Location{Longitude: null.FloatFrom(12.37), Latitude: null.FloatFrom(12.37)}},
	}

	ds := calculations.NewResolver()

	// we have already resolved these
	// TODO : make this go away
	ds.RegisterResolved(calculations.LocationLongitudeDependancy)
	ds.RegisterResolved(calculations.LocationLatitudeDependancy)
	ds.RegisterResolved(calculations.VesselUUIDDependancy)
	ds.RegisterResolved(calculations.AccountUUIDDependancy)
	ds.RegisterResolved(calculations.EventTimeDependancy)
	ds.RegisterResolved(calculations.WeatherSwellHeightDependancy)

	funcs := calculations.NewFunctionSet()

	calculations.RegisterDistanceOverGround(funcs)
	calculations.RegisterSpeedOverGround(funcs)
	calculations.RegisterDouglasSeaState(funcs)

	err := ds.CollapseAll(rows, funcs)
	spew.Dump(rows, err)

}
