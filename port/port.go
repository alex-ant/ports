package port

// Info contains single port info data.
type Info struct {
	Name        string     `json:"name"`
	City        string     `json:"city"`
	Country     string     `json:"country"`
	Alias       []string   `json:"alias"`
	Regions     []string   `json:"regions"`
	Coordinates [2]float64 `json:"coordinates"`
	Province    string     `json:"province"`
	Timezone    string     `json:"timezone"`
	Unlocs      []string   `json:"unlocs"`
	Code        string     `json:"code"`
}
