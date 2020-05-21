package calculations

import (
	"fmt"
)

const (
	SpeedOverGroundDependancy = "speed_over_ground"
)

func RegisterSpeedOverGround(f *functionset) error {
	return f.RegisterIterator(SpeedOverGroundDependancy, &sogCalculator{}, DistanceOverGroundDependancy)
}

type sogCalculator struct{}

func (s *sogCalculator) Fn(rows []*DataRow) error {
	fmt.Println("Speed over ground called") // TODO : how do you log?
	return nil
}
