package calculations

import (
	"time"

	"gopkg.in/guregu/null.v4"
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
