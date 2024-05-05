package service

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/model/dto/response"
	"CatsSocialMedia/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CatService interface {
	FindAll(filterParams map[string]interface{}) ([]response.CatResponse, error)
	FindByUserID(i int) (interface{}, interface{})
	FindByID(catID string) (model.Cat, error)
	FindByIDAndUserID(catID string, userID int) (model.Cat, error)
	Create(catRequest request.CatRequest) (response.CreateCatResponse, error)
	Update(catID string, catRequest request.CatRequest) (model.Cat, error)
	Delete(catID string, userID int) error
	ConverToString(cat []response.CatResponse) []response.CatResponseString
}

type catService struct {
	repository   repository.CatRepository
	matchService MatchService
}

func NewCatService(repository repository.CatRepository, matchService MatchService) *catService {
	return &catService{repository, matchService}
}

func (s *catService) FindAll(filterParams map[string]interface{}) ([]response.CatResponseString, error) {
	fmt.Println(filterParams)
	cats, err := s.repository.FindAll(filterParams)
	if err != nil {
		return nil, err
	}
	catString := s.CovertToString(cats)
	return catString, nil
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

func (s *catService) FindByIDAndUserID(catID string, userID int) (model.Cat, error) {
	cat, err := s.repository.FindByIDAndUserID(catID, userID)
	if err != nil {
		return model.Cat{}, err
	}

	// Jika kucing tidak ditemukan, kembalikan error
	if cat.ID == 0 {
		return model.Cat{}, errors.New("cat not found")
	}

	return cat, nil
}

func (s *catService) Create(catRequest request.CatRequest) (response.CreateCatResponse, error) {
	//save cat
	cat := model.Cat{
		Name:        catRequest.Name,
		Race:        catRequest.Race,
		Sex:         catRequest.Sex,
		AgeInMonth:  catRequest.AgeInMonth,
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

	match, err := s.matchService.LatestMatch(catID)
	if match.ID != 0 {
		if match.IsAproved {
			if existingCat.Sex != catRequest.Sex {
				return model.Cat{}, errors.New("cat has been matched")
			}
		}
	}

	// Update data kucing dengan data yang diberikan dalam request
	existingCat.Name = catRequest.Name
	existingCat.Race = catRequest.Race
	existingCat.Sex = catRequest.Sex
	existingCat.AgeInMonth = catRequest.AgeInMonth
	existingCat.Description = catRequest.Description
	existingCat.ImageUrls = catRequest.ImageUrls

	// Simpan perubahan pada database
	updatedCat, err := s.repository.Update(existingCat)
	if err != nil {
		return model.Cat{}, err
	}

	return updatedCat, nil
}

func (s *catService) Delete(catID string, userID int) error {
	err := s.repository.Delete(catID, userID)
	return err
}

func (s *catService) CovertToString(cat []response.CatResponse) []response.CatResponseString {
	var cats []response.CatResponseString
	for _, v := range cat {
		catString := response.CatResponseString{
			ID:          string(v.ID),
			Name:        v.Name,
			Race:        string(v.Race),
			Sex:         string(v.Sex),
			AgeInMonth:  strconv.Itoa(v.AgeInMonth),
			ImageURLs:   strings.Join(v.ImageURLs, ","),
			Description: v.Description,
			HasMatched:  strconv.FormatBool(v.HasMatched),
			CreatedAt:   v.CreatedAt.String(),
		}
		cats = append(cats, catString)
	}
	return cats
}
