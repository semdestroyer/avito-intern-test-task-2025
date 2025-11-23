package queries

type UserQuery struct {
	UserId   string `form:"user_id" binding:"required, string"`
	Username string `form:"username" binding:"required, string"`
	TeamName string `form:"team_name" binding:"required, string"`
	IsActive bool   `form:"is_active" binding:"required, boolean"`
}
