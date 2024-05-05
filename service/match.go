package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/model/dto/response"
	"CatsSocialMedia/repository"
	"errors"
	"fmt"
	"strconv"
)

type MatchService interface {
	GetMatches(userId int) ([]response.MatchResponse, error)
	Create(userID int, matchRequest request.MatchRequest) (model.Match, error)
	Approval(userId int, matchId int, isAprrove bool) (int, error)
	Delete(userId int, match string) (int, error)
	LatestMatch(catID string) (model.Match, error)
}

type matchService struct {
	repository    repository.MatchRepository
	catRepository repository.CatRepository
}

func NewMatchService(repository repository.MatchRepository, catRepository repository.CatRepository) *matchService {
	return &matchService{repository, catRepository}
}

var ErrCatNotFound = errors.New("cat not found")

func (s *matchService) GetMatches(userId int) ([]response.MatchResponse, error) {
	matches, err := s.repository.GetMatches(userId)

	if err != nil {
		return []response.MatchResponse{}, err
	}

	return matches, err
}

func (s *matchService) Create(userId int, matchRequest request.MatchRequest) (model.Match, error) {
	matchCat, matchCatError := s.catRepository.FindByID(strconv.Itoa(matchRequest.MatchCatID))
	userCat, userCatError := s.catRepository.FindByID(strconv.Itoa(matchRequest.UserCatID))

	match := model.Match{
		MatchCatID: matchRequest.MatchCatID,
		UserCatID:  matchRequest.UserCatID,
		Message:    matchRequest.Message,
		IssuedBy:   userId,
		AcceptedBy: matchCat.UserID,
	}

	fmt.Println(userCat.UserID)
	fmt.Println(userId)
	if userCat.UserID != userId {
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

	// if matchCat.UserID == userCat.UserID {
	// 	return match, errors.New("THE CATS OWNER ARE SAME")
	// }

	newMatch, err := s.repository.Create(match)
	return newMatch, err
}

func (s *matchService) Approval(userId int, matchId int, isAprrove bool) (int, error) {
	match, getMatchError := s.repository.MatchIsExist(matchId)

	if getMatchError != nil {
		return matchId, getMatchError
	}

	fmt.Println(match)

	if match == (model.Match{}) {
		return matchId, errors.New("MATCH IS NOT EXIST")
	}

	fmt.Println(match.IsMatched)
	if match.IsMatched {
		return matchId, errors.New("MATCHID IS NO LONGER VALID")
	}

	matchCat, _ := s.catRepository.FindByID(strconv.Itoa(match.MatchCatID))
	if matchCat.UserID != userId {
		return matchId, errors.New("THE MATCH CAT OWNER IS NOT SAME")
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

func (s *matchService) Delete(userId int, matchID string) (int, error) {
	matchId, _ := strconv.Atoi(matchID)
	match, getMatchError := s.repository.MatchIsExist(matchId)

	if getMatchError != nil {
		return matchId, getMatchError
	}

	fmt.Println(match)

	if match == (model.Match{}) {
		return matchId, errors.New("MATCH IS NOT EXIST")
	}
	fmt.Println(match.IssuedBy)
	fmt.Println(userId)
	if match.IssuedBy != userId {
		return matchId, errors.New("UNAUTHORIZED DELETE THIS MATCH")
	}

	if match.IsMatched {
		return matchId, errors.New("MATCHID IS ALREADY APPROVED / REJECT")
	}

	err := s.repository.Delete(matchId)
	return matchId, err
}

func (s *matchService) LatestMatch(catID string) (model.Match, error) {
	match, err := s.repository.LatestMatch(catID)
	if err != nil {
		return model.Match{}, err
	}

	// Jika kucing tidak ditemukan, kembalikan error
	if match.ID == 0 {
		return model.Match{}, errors.New("match not found")
	}

	return match, nil
}
