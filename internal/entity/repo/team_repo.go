package repo

import (
	"avito-intern-test-task-2025/internal/entity"
	"avito-intern-test-task-2025/pkg/db"
)

type TeamRepo struct {
	db *db.DB
}

func NewTeamRepo(db *db.DB) *TeamRepo {
	return &TeamRepo{
		db: db,
	}
}

func (tr TeamRepo) CreateTeam() {

}

func (tr TeamRepo) TeamGetByName() entity.Team {

}
