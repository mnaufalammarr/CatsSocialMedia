package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/repository"
)

type MatchService interface {
	Create(userID float64, matchRequest request.MatchRequest) (model.Match, error)
}

type matchService struct {
	repository repository.MatchRepository
}

func NewMatchService(repository repository.MatchRepository) *matchService {
	return &matchService{repository}
}

func (s *matchService) Create(userId float64, matchRequest request.MatchRequest) (model.Match, error) {
	match := model.Match{
		MatchCatID: matchRequest.MatchCatID,
		UserCatID:  matchRequest.UserCatID,
		Message:    matchRequest.Message,
		IssuedBy:   int(userId),
	}

	newMatch, err := s.repository.Create(match)
	return newMatch, err
}
