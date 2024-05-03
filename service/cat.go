package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/model/dto/response"
	"CatsSocialMedia/repository"
	"errors"
)

type CatService interface {
	FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error)
	FindByUserID(i int) (interface{}, interface{})
	FindByID(catID string) (model.Cat, error)
	Create(catRequest request.CatRequest) (model.Cat, error)
	Update(catID string, catRequest request.CatRequest) (model.Cat, error)
	Delete(catID string) error
}

type catService struct {
	repository repository.CatRepository
}

func NewCatService(repository repository.CatRepository) *catService {
	return &catService{repository}
}

func (s *catService) FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error) {
	cats, err := s.repository.FindAll(filterParams)
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *catService) FindByID(catID string) (model.Cat, error) {
	// Cari kucing berdasarkan ID
	cat, err := s.repository.FindByID(catID)
	if err != nil {
		return model.Cat{}, err
	}

	// Jika kucing tidak ditemukan, kembalikan error
	if cat.ID == 0 {
		return model.Cat{}, errors.New("cat not found")
	}

	return cat, nil
}

func (s *catService) FindByUserID(i int) (interface{}, interface{}) {
	cat, err := s.repository.FindByUserID(i)
	if err != nil {
		return model.Cat{}, err
	}

	// Jika kucing tidak ditemukan, kembalikan error
	if cat.ID == 0 {
		return model.Cat{}, errors.New("cat not found")
	}

	return cat, nil
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
		UserID:      catRequest.UserId,
	}
	newCat, err := s.repository.Create(cat)
	return newCat, err
}

func (s *catService) Update(catID string, catRequest request.CatRequest) (model.Cat, error) {
	// Cek apakah kucing dengan ID yang diberikan ada dalam database
	existingCat, err := s.repository.FindByID(catID)
	if err != nil {
		return model.Cat{}, err
	}
	if existingCat.ID == 0 {
		return model.Cat{}, errors.New("cat not found")
	}

	// Update data kucing dengan data yang diberikan dalam request
	existingCat.Name = catRequest.Name
	existingCat.Race = catRequest.Race
	existingCat.Sex = catRequest.Sex
	existingCat.AgeInMonths = catRequest.AgeInMonths
	existingCat.Description = catRequest.Description
	existingCat.ImageUrls = catRequest.ImageUrls

	// Simpan perubahan pada database
	updatedCat, err := s.repository.Update(existingCat)
	if err != nil {
		return model.Cat{}, err
	}

	return updatedCat, nil
}

func (s *catService) Delete(catID string) error {
	err := s.repository.Delete(catID)
	return err
}
