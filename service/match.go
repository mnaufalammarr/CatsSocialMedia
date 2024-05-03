package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/repository"
	"errors"
	"fmt"
	"strconv"
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

var ErrCatNotFound = errors.New("cat not found")

func (s *matchService) Create(userId float64, matchRequest request.MatchRequest) (model.Match, error) {
	match := model.Match{
		MatchCatID: matchRequest.MatchCatID,
		UserCatID:  matchRequest.UserCatID,
		Message:    matchRequest.Message,
		IssuedBy:   int(userId),
	}

	matchCat, matchCatError := s.catRepository.FindByID(strconv.Itoa(match.MatchCatID))
	userCat, userCatError := s.catRepository.FindByID(strconv.Itoa(match.UserCatID))

	if userCat.UserID != int(userId) {
		return match, errors.New("THE USER CAT IS NOT BELONG TO THE USER")
	}

	if (matchCat.ID == 0 && matchCatError == nil) || (userCat.ID == 0 && userCatError == nil) {
		return match, ErrCatNotFound
	}

	if matchCat.Sex == userCat.Sex {
		return match, errors.New("THE CATS GENDER ARE SAME")
	}

	if matchCat.HasMatch || userCat.HasMatch {
		return match, errors.New("THE CATS ALREADY MATCHED")
	}

	if matchCat.UserID == userCat.UserID {
		return match, errors.New("THE CATS OWNER ARE SAME")
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
