package admin

type AdminCreateRequest struct {
	Username         string `json:"username" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,password=6,max=24"`
	Confirm_Password string `json:"confirm_password" validate:"required,confirm_password=Password"`
	Admin_Code       string `json:"admin_code" validate:"required"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password=6,max=24"`
}
