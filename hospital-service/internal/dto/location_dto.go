package dto

type CityLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DistrictLookup struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	CityID uint   `json:"city_id"`
}
