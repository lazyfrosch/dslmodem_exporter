package units

// Mbps is Megabits per second.
//
// https://en.wikipedia.org/wiki/Data-rate_units#Megabit_per_second
type Mbps float64

func (t *Mbps) UnmarshalText(data []byte) error {
	f, err := UnmarshalFloatWithUnits(data, "Mbps")
	if err != nil {
		return err
	}

	*t = Mbps(f)

	return nil
}

func (t Mbps) BitsPerSec() float64 {
	return float64(t) * 1000 * 1000
}

// Decibel (dB) is a relative unit of measurement equal to one tenth of a bel (B).
//
// https://en.wikipedia.org/wiki/Decibel
type Decibel float64

func (t *Decibel) UnmarshalText(data []byte) error {
	f, err := UnmarshalFloatWithUnits(data, "dB")
	if err != nil {
		return err
	}

	*t = Decibel(f)

	return nil
}

// DecibelMilliwatt (dBm) is a unit of level used to indicate that a power
// level is expressed in decibels (dB) with reference to one milliwatt (mW).
//
// https://en.wikipedia.org/wiki/DBm
type DecibelMilliwatt float64

func (t *DecibelMilliwatt) UnmarshalText(data []byte) error {
	f, err := UnmarshalFloatWithUnits(data, "dBm")
	if err != nil {
		return err
	}

	*t = DecibelMilliwatt(f)

	return nil
}
