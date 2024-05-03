package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/repository"
	"errors"
	"fmt"
)

type MatchService interface {
	Create(userID float64, matchRequest request.MatchRequest) (model.Match, error)
	Approval(userId float64, matchId int, isAprrove bool) (int, error)
	Delete(userId float64, match int) error
}

type matchService struct {
	repository    repository.MatchRepository
	catRepository repository.CatRepository
}

func NewMatchService(repository repository.MatchRepository, catRepository repository.CatRepository) *matchService {
	return &matchService{repository, catRepository}
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

func (s *matchService) Approval(userId float64, matchId int, isAprrove bool) (int, error) {
	match, getMatchError := s.repository.MatchIsExist(matchId)

	if getMatchError != nil {
		return matchId, getMatchError
	}

	fmt.Println(match)

	if match == (model.Match{}) {
		return matchId, errors.New("MATCH DOES NOT EXIST")
	}

	if isAprrove {
		s.catRepository.UpdateHasMatch(match.MatchCatID, isAprrove)
		s.catRepository.UpdateHasMatch(match.UserCatID, isAprrove)
		fmt.Println("approved cat")
	}

	id, err := s.repository.MatchApproval(matchId, isAprrove)

	if err != nil {
		fmt.Println(err)
		return id, err
	}

	return id, nil
}

func (s *matchService) Delete(userId float64, match int) error {
	err := s.repository.Delete(match)
	return err
}
