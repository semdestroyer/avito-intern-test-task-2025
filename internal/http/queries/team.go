package queries

type TeamQuery struct {
	Name    string `form:"team_name" binding:"required, string"`
	Members []TeamMemberQuery
}
