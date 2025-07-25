package dto

type CreateHospitalRequest struct {
	Name       string `json:"name"`
	TaxNumber  string `json:"tax_number"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	CityID     uint   `json:"city_id"`
	DistrictID uint   `json:"district_id"`
}

type HospitalResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	TaxNumber    string `json:"tax_number"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	CityID       uint   `json:"city_id"`
	CityName     string `json:"city_name"`
	DistrictID   uint   `json:"district_id"`
	DistrictName string `json:"district_name"`
}
