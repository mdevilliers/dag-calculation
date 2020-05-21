package calculations

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

const (
	DouglasSeaStateDependancy = "douglas_sea_state"
)

func RegisterDouglasSeaState(f *functionset) error {
	return f.RegisterIterator(DouglasSeaStateDependancy, &dss{}, WeatherSwellHeightDependancy)
}

type dss struct{}

func (d *dss) Fn(rows []*DataRow) error {
	fmt.Println("Douglas Sea State called") // TODO : how do you log?

	for _, r := range rows {
		if r.Weather.SwellHeight.Valid {
			height := r.Weather.SwellHeight.Float64

			if height == 0 {
				r.DouglasSeaState = null.IntFrom(0)
			} else if height <= 0.1 {
				r.DouglasSeaState = null.IntFrom(1)
			} else if height <= 0.5 {
				r.DouglasSeaState = null.IntFrom(2)
			} else if height <= 1.25 {
				r.DouglasSeaState = null.IntFrom(3)
			} else if height <= 2.5 {
				r.DouglasSeaState = null.IntFrom(4)
			} else if height <= 4.00 {
				r.DouglasSeaState = null.IntFrom(5)
			} else if height <= 6.00 {
				r.DouglasSeaState = null.IntFrom(6)
			} else if height <= 9.00 {
				r.DouglasSeaState = null.IntFrom(7)
			} else if height <= 14.00 {
				r.DouglasSeaState = null.IntFrom(8)
			} else {
				r.DouglasSeaState = null.IntFrom(9)
			}

		}

	}

	return nil
}
