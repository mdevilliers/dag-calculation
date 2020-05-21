package calculations

const (
	SpeedOverGroundDependancy = "speed_over_ground"
)

func RegisterSpeedOverGround(f *functionset) error {
	return f.RegisterIterator(SpeedOverGroundDependancy, &sogCalculator{}, DistanceOverGroundDependancy)
}

type sogCalculator struct{}

func (s *sogCalculator) Fn(runtime runtime, rows []*DataRow) error {
	runtime.Logger.Infof("Speed over ground called")
	return nil
}
