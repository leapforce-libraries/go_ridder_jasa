package ridder_jasa

type Person struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Prefix    string `json:"Prefix"`
	Title     int32  `json:"Title"`
	Initials  string `json:"Initials"`
	Gender    int32  `json:"Gender"`
	Deceased  bool   `json:"Deceased"`
}
