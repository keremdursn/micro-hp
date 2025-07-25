package usecase

import (
	"errors"

	"hospital-service/internal/dto"
	"hospital-service/internal/models"
	"hospital-service/internal/repository"
)

type PolyclinicUsecase interface {
	ListAllPolyclinics() ([]dto.PolyclinicLookup, error)
	AddPolyclinicToHospital(req *dto.AddHospitalPolyclinicRequest, hospitalID uint) (*dto.HospitalPolyclinicResponse, error)
	ListHospitalPolyclinic(hospitalID uint, page, size int) (*dto.HospitalPolyclinicListResponse, error)
	RemoveHospitalPolyclinic(id uint, hospitalID uint) error
}

type polyclinicUsecase struct {
	repo repository.PolyclinicRepository
}

func NewPolyclinicUsecase(repo repository.PolyclinicRepository) PolyclinicUsecase {
	return &polyclinicUsecase{repo: repo}
}

func (u *polyclinicUsecase) ListAllPolyclinics() ([]dto.PolyclinicLookup, error) {
	polys, err := u.repo.GetAllPolyclinic()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.PolyclinicLookup, 0, len(polys))
	for _, p := range polys {
		resp = append(resp, dto.PolyclinicLookup{
			ID:   p.ID,
			Name: p.Name,
		})
	}
	return resp, nil
}

func (u *polyclinicUsecase) AddPolyclinicToHospital(req *dto.AddHospitalPolyclinicRequest, hospitalID uint) (*dto.HospitalPolyclinicResponse, error) {
	alreadyExists, err := u.repo.IsPolyclinicAlreadyAdded(hospitalID, req.PolyclinicID)
	if err != nil {
		return nil, err
	}

	if alreadyExists {
		return nil, errors.New("this polyclinic is already added to the hospital")
	}

	poly, err := u.repo.GetPolyclinicByID(req.PolyclinicID)
	if err != nil {
		return nil, err
	}

	hp := &models.HospitalPolyclinic{
		HospitalID:   hospitalID,
		PolyclinicID: req.PolyclinicID,
	}

	if err := u.repo.CreateHospitalPolyclinic(hp); err != nil {
		return nil, err
	}

	return &dto.HospitalPolyclinicResponse{
		ID:             hp.ID,
		PolyclinicID:   poly.ID,
		PolyclinicName: poly.Name,
	}, nil

}

func (u *polyclinicUsecase) ListHospitalPolyclinic(hospitalID uint, page, size int) (*dto.HospitalPolyclinicListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	total, err := u.repo.CountByHospitalID(hospitalID)
	if err != nil {
		return nil, err
	}

	hps, err := u.repo.GetPaginatedByHospitalID(hospitalID, page, size)
	if err != nil {
		return nil, err
	}

	result := make([]dto.HospitalPolyclinicDetail, 0, len(hps))
	for _, hp := range hps {
		poly, err := u.repo.GetPolyclinicByID(hp.PolyclinicID)
		if err != nil {
			return nil, err
		}

		totalPersonnel, err := u.repo.CountPersonnel(hp.ID)
		if err != nil {
			return nil, err
		}

		groupCounts, err := u.repo.GetPersonnelGroupCounts(hp.ID)
		if err != nil {
			return nil, err
		}

		personnelGroups := make([]dto.PolyclinicPersonnelGroup, 0, len(groupCounts))
		for _, g := range groupCounts {
			personnelGroups = append(personnelGroups, dto.PolyclinicPersonnelGroup{
				GroupName: g.GroupName,
				Count:     g.Count,
			})
		}

		result = append(result, dto.HospitalPolyclinicDetail{
			ID:              hp.ID,
			PolyclinicName:  poly.Name,
			TotalPersonnel:  int(totalPersonnel),
			PersonnelGroups: personnelGroups,
		})

	}

	return &dto.HospitalPolyclinicListResponse{
		Polyclinics: result,
		Total:       int(total),
		Page:        page,
		Size:        size,
	}, nil
}

func (u *polyclinicUsecase) RemoveHospitalPolyclinic(id uint, hospitalID uint) error {
	hp, err := u.repo.GetHospitalPolyclinicByID(id)
	if err != nil {
		return errors.New("hospital polyclinic not found")
	}

	if hp.HospitalID != hospitalID {
		return errors.New("forbidden: cannot remove polyclinic from another hospital")
	}
	return u.repo.Delete(hp)

}
