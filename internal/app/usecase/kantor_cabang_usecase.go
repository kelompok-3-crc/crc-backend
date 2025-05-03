package usecase

import (
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"
)

type KantorCabangUsecase interface {
	Create(req dto.CreateKantorCabangRequest) (*dto.KantorCabangResponse, error)
	GetAll() ([]dto.KantorCabangResponse, error)
}

type kantorCabangUsecase struct {
	repo repository.KantorCabangRepository
}

func NewKantorCabangUsecase(repo repository.KantorCabangRepository) KantorCabangUsecase {
	return &kantorCabangUsecase{repo}
}

func (uc *kantorCabangUsecase) Create(req dto.CreateKantorCabangRequest) (*dto.KantorCabangResponse, error) {
	model := model.KantorCabang{
		Nama: req.Nama,
	}
	err := uc.repo.Create(&model)
	if err != nil {
		return nil, err
	}
	return &dto.KantorCabangResponse{
		ID:   model.ID,
		Nama: model.Nama,
	}, nil
}

func (uc *kantorCabangUsecase) GetAll() ([]dto.KantorCabangResponse, error) {
	cabangs, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.KantorCabangResponse, 0)
	for _, c := range cabangs {
		responses = append(responses, dto.KantorCabangResponse{
			ID:   c.ID,
			Nama: c.Nama,
		})
	}
	return responses, nil
}
