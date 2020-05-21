package calculations

import (
	"fmt"

	"gopkg.in/guregu/null.v4"

	"github.com/umahmood/haversine"
)

const (
	DistanceOverGroundDependancy = "distance_over_ground"
)

func RegisterDistanceOverGround(f *functionset) error {
	return f.RegisterIteratorFn(DistanceOverGroundDependancy, DistanceOverGroundFn, LocationLongitudeDependancy, LocationLatitudeDependancy)
}

func DistanceOverGroundFn(rows []*DataRow) error {
	fmt.Println("Distance over ground called") // TODO : how do you log?

	var previous *haversine.Coord

	for _, r := range rows {
		if previous == nil {
			ok, p := coordFromRow(r)
			if ok {
				previous = &p
			}
		} else {
			ok, latest := coordFromRow(r)
			if ok {
				_, km := haversine.Distance(*previous, latest)

				r.DistanceOverGround = null.FloatFrom(km)
				previous = &latest

			}
		}
	}
	return nil
}

func coordFromRow(row *DataRow) (bool, haversine.Coord) {
	if row.Location.Longitude.Valid && row.Location.Latitude.Valid {
		return true, haversine.Coord{Lon: row.Location.Longitude.Float64, Lat: row.Location.Latitude.Float64}
	}
	return false, haversine.Coord{}
}
