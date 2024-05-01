package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/repository"
)

type CatService interface {
	Create(catRequest request.CatRequest) (model.Cat, error)
}

type catService struct {
	repository repository.CatRepository
}

func NewCatService(repository repository.CatRepository) *catService {
	return &catService{repository}
}

func (s *catService) Create(catRequest request.CatRequest) (model.Cat, error) {
	//save cat
	cat := model.Cat{
		Name:        catRequest.Name,
		Race:        catRequest.Race,
		Sex:         catRequest.Sex,
		AgeInMonths: catRequest.AgeInMonths,
		Description: catRequest.Description,
		ImageUrls:   catRequest.ImageUrls,
	}
	newCat, err := s.repository.Create(cat)
	return newCat, err
}
