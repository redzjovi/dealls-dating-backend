package repository

import (
	"dealls-dating/internal/entity"
)

type MatchRepository struct {
	Repository[entity.Match]
}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{}
}
