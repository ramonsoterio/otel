package weather

type Location struct {
	City string
}

type Weather struct {
	TemperatureCelsius float64
}
type TemperatureOutput struct {
	City          string  `json:"city"`
	TempCelsius   float64 `json:"temp_C"`
	TempFarenheit float64 `json:"temp_F"`
	TempKelvin    float64 `json:"temp_K"`
}

func NewTemperatureOutput(tempCelsius float64) TemperatureOutput {
	return TemperatureOutput{
		TempCelsius:   tempCelsius,
		TempFarenheit: tempCelsius*1.8 + 32,
		TempKelvin:    tempCelsius + 273,
	}
}

func (t TemperatureOutput) WithCity(city string) TemperatureOutput {
	t.City = city
	return t
}
