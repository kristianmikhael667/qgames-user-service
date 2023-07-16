package dto

type (
	PermissionResponse struct {
		UidPermission string `json:"uid_permission"`
		Name          string `json:"name"`
		Slug          string `json:"slug"`
	}

	PermissionRequestBody struct {
		UidPermission string `json:"uid_permission" validate:"omitempty"`
		Name          string `json:"name" validate:"omitempty"`
	}
)
