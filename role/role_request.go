package role

type RoleCreateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type RoleUpdateRequest struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=1,max=100"`
}
