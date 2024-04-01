package web

type AdminCreateRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	Admin_Code string `json:"admin_code" validate:"required"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
