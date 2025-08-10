package place

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Place struct {
	Tourapi_place_id string `json:"tourapi_place_id"`
	Category         int64  `json:"category"`
	Title            string `json:"title"`
	Address          string `json:"address"`
	Tel              string `json:"tel"`
	Longitude        string `json:"longitude"`
	Latitude         string `json:"latitude"`
	PlaceImage       string `json:"place_image"`
}
