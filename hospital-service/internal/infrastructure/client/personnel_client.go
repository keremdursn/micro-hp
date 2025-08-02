package client

import (
	"encoding/json"
	"errors"
	"fmt"
	dt "hospital-shared/dto"
	"net/http"
	"time"
)

type PersonnelClient interface {
	GetPersonnelCount(hospitalPolyclinicID uint) (int, error)
	GetPersonnelGroups(hospitalPolyclinicID uint) ([]dt.PolyclinicPersonnelGroup, error)
}

type personnelClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPersonnelClient(baseURL string) PersonnelClient {
	return &personnelClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *personnelClient) GetPersonnelCount(hpID uint) (int, error) {
	url := fmt.Sprintf("%s/staff/count?hospitalPolyclinicID=%d", c.baseURL, hpID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("personnel count request failed")
	}

	var result struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

func (c *personnelClient) GetPersonnelGroups(hpID uint) ([]dt.PolyclinicPersonnelGroup, error) {
	url := fmt.Sprintf("%s/staff/groups?hospitalPolyclinicID=%d", c.baseURL, hpID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("personnel groups request failed")
	}

	var result []dt.PolyclinicPersonnelGroup
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
