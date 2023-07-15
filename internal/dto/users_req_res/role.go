package dto

type (
	RoleResponse struct {
		UidRole string `json:"uid_role"`
		Name    string `json:"name"`
		Desc    string `json:"desc"`
		Data    string `json:"data"`
		Status  string `json:"status"`
	}

	RoleRequestBody struct {
		Name   string `json:"name" validate:"omitempty"`
		Desc   string `json:"desc" validate:"omitempty"`
		Data   string `json:"data"`
		Status string `json:"status" validate:"omitempty"`
	}
)
