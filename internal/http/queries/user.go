package queries

type UserIsActiveQuery struct {
	UserId   string `form:"user_id" binding:"required"`
	IsActive bool   `form:"is_active" binding:"required"`
}

type UserIdQuery struct {
	UserId string `form:"user_id" binding:"required"`
}

type UserQuery struct {
	UserId   string `form:"user_id" binding:"required" validate:"required"`
	Username string `form:"username" binding:"required" validate:"required"`
	TeamName string `form:"team_name" binding:"required" validate:"required"`
	IsActive bool   `form:"is_active" binding:"required" validate:"required"`
}
