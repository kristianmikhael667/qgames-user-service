package dto

type (
	AssignResponse struct {
		UidAssign   string `json:"uid_assign"`
		Users       string `json:"users"`
		Roles       string `json:"roles"`
		Permissions string `json:"permissions"`
		Status      string `json:"status"`
	}

	ReqAssign struct {
		Users      string  `json:"users"`
		PaymentFee float64 `json:"payment_fee"`
	}
)
