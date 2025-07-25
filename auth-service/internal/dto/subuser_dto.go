package dto

type CreateSubUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TC        string `json:"tc"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type SubUserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TC        string `json:"tc"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type UpdateSubUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TC        string `json:"tc"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Role      string `json:"role"` // "yetkili" veya "calisan"
}
