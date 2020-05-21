package main

import (
	"gopkg.in/guregu/null.v4"

	"github.com/mdevilliers/dag-calculation/calculations"
)

func main() {

	rows := []*calculations.DataRow{
		{Location: calculations.Location{Longitude: null.FloatFrom(12.34), Latitude: null.FloatFrom(12.34)}, Weather: calculations.Weather{SwellHeight: null.FloatFrom(0.48)}},
		{Location: calculations.Location{Longitude: null.FloatFrom(12.35), Latitude: null.FloatFrom(12.35)}},
		{Location: calculations.Location{Longitude: null.FloatFrom(12.36), Latitude: null.FloatFrom(12.36)}},
		{Location: calculations.Location{Longitude: null.FloatFrom(12.37), Latitude: null.FloatFrom(12.37)}},
	}

	funcs := calculations.NewFunctionSet()

	calculations.RegisterDistanceOverGround(funcs)
	calculations.RegisterDouglasSeaState(funcs)
	calculations.RegisterSpeedOverGround(funcs)

	err := funcs.Build()

	if err != nil {
		panic(err)
	}

	err = funcs.Collapse(rows)

	if err != nil {
		panic(err)
	}

}
