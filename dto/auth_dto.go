package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	StaffID      uint   `json:"staff_id"`
	Username     string `json:"username"`
	HospitalID   uint   `json:"hospital_id"`
	HospitalName string `json:"hospital_name"`
}
