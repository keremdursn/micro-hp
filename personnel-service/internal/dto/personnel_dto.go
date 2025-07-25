package dto

type JobGroupLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TitleLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AddStaffRequest struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	TC                   string `json:"tc"`
	Phone                string `json:"phone"`
	JobGroupID           uint   `json:"job_group_id"`
	TitleID              uint   `json:"title_id"`
	HospitalPolyclinicID *uint  `json:"hospital_polyclinic_id"`
	WorkingDays          string `json:"working_days"`
}

type StaffResponse struct {
	ID                   uint    `json:"id"`
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	TC                   string  `json:"tc"`
	Phone                string  `json:"phone"`
	JobGroupID           uint    `json:"job_group_id"`
	JobGroupName         string  `json:"job_group_name"`
	TitleID              uint    `json:"title_id"`
	TitleName            string  `json:"title_name"`
	HospitalPolyclinicID *uint   `json:"hospital_polyclinic_id"`
	PolyclinicName       *string `json:"polyclinic_name"`
	WorkingDays          string  `json:"working_days"`
}

type UpdateStaffRequest struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	TC                   string `json:"tc"`
	Phone                string `json:"phone"`
	JobGroupID           uint   `json:"job_group_id"`
	TitleID              uint   `json:"title_id"`
	HospitalPolyclinicID *uint  `json:"hospital_polyclinic_id"`
	WorkingDays          string `json:"working_days"`
}

type StaffListFilter struct {
	FirstName  string
	LastName   string
	TC         string
	JobGroupID *uint
	TitleID    *uint
}

type StaffListResponse struct {
	Staff []StaffResponse `json:"staff"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}
