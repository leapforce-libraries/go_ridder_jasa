package ridder_jasa

type Address struct {
	Street      string `json:"Street"`
	HouseNumber string `json:"HouseNumber"`
	ZipCode     string `json:"ZipCode"`
	City        string `json:"City"`
	Country     string `json:"Country"`
}
