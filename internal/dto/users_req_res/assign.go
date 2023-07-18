package dto

type (
	AssignResponse struct {
		UidAssign   string `json:"uid_assign"`
		Users       string `json:"users"`
		Roles       string `json:"roles"`
		Permissions string `json:"permissions"`
		Status      string `json:"status"`
	}
)
