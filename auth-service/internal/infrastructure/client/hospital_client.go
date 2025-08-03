package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dt "hospital-shared/dto"
)

type HospitalClient interface {
	CreateHospital(req *dt.CreateHospitalRequest) (*dt.HospitalResponse, error)
}

type hospitalClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHospitalClient(baseURL string) HospitalClient {
	return &hospitalClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *hospitalClient) CreateHospital(req *dt.CreateHospitalRequest) (*dt.HospitalResponse, error) {
	url := fmt.Sprintf("%s/api/hospital", c.baseURL)

	body, _ := json.Marshal(req)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("hospital service returned status %d", resp.StatusCode)
	}

	var hospitalResp dt.HospitalResponse
	if err := json.NewDecoder(resp.Body).Decode(&hospitalResp); err != nil {
		return nil, err
	}

	return &hospitalResp, nil
}
