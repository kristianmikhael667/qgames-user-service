package dto

type ByUuidUsersRequest struct {
	Uid string `param:"uid" validate:"required"`
}
