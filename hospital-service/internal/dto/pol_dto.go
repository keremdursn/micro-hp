package dto

type PolyclinicLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AddHospitalPolyclinicRequest struct {
	PolyclinicID uint `json:"polyclinic_id"`
}

type HospitalPolyclinicResponse struct {
	ID             uint   `json:"id"`
	PolyclinicID   uint   `json:"polyclinic_id"`
	PolyclinicName string `json:"polyclinic_name"`
}

type PolyclinicPersonnelGroup struct {
	GroupName string `json:"group_name"`
	Count     int    `json:"count"`
}

type HospitalPolyclinicDetail struct {
	ID              uint                      `json:"id"`
	PolyclinicName  string                    `json:"polyclinic_name"`
	TotalPersonnel  int                       `json:"total_personnel"`
	PersonnelGroups []PolyclinicPersonnelGroup `json:"personnel_groups"`
}

type HospitalPolyclinicListResponse struct {
	Polyclinics []HospitalPolyclinicDetail `json:"polyclinics"`
	Total       int                        `json:"total"`
	Page        int                        `json:"page"`
	Size        int                        `json:"size"`
}