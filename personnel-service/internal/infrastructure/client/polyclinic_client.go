package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dt "hospital-shared/dto"
)

type PolyclinicClient interface {
	GetHospitalPolyclinicByID(id uint) (*dt.HospitalPolyclinicResponseDTO, error)
}

type polyclinicClient struct {
	baseURL string
	client  *http.Client
}

func NewPolyclinicClient(baseURL string) PolyclinicClient {
	return &polyclinicClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (p *polyclinicClient) GetHospitalPolyclinicByID(id uint) (*dt.HospitalPolyclinicResponseDTO, error) {
	url := fmt.Sprintf("%s/api/polyclinic/hospital-polyclinics/%d", p.baseURL, id)
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %d", resp.StatusCode)
	}

	var info dt.HospitalPolyclinicResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
