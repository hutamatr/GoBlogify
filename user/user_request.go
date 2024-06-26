package user

type UserCreateRequest struct {
	Username         string `json:"username" validate:"required,min=1,max=24"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,password=6,max=24"`
	Confirm_Password string `json:"confirm_password" validate:"required,confirm_password=Password"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password=6,max=24"`
}

type UserUpdateRequest struct {
	Id         int    `json:"id" validate:"required"`
	Role_Id    int    `json:"role_id" validate:"required"`
	Username   string `json:"username" validate:"required"`
	First_Name string `json:"first_name" validate:"required"`
	Last_Name  string `json:"last_name" validate:"required"`
}
