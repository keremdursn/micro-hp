package usecase

import (
	"errors"

	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/utils"
)

type SubUserUsecase interface {
	CreateSubUser(req *dto.CreateSubUserRequest, hospitalID uint) (*dto.SubUserResponse, error)
	ListSubUsers(hospitalID uint) ([]dto.SubUserResponse, error)
	UpdateSubUser(id uint, req *dto.UpdateSubUserRequest, hospitalID uint) (*dto.SubUserResponse, error)
	DeleteSubUser(id, hospitalID uint) error
}

type subUserUsecase struct {
	repo repository.SubUserRepository
}

func NewSubUserUsecase(repo repository.SubUserRepository) SubUserUsecase {
	return &subUserUsecase{repo: repo}
}

func (u *subUserUsecase) CreateSubUser(req *dto.CreateSubUserRequest, hospitalID uint) (*dto.SubUserResponse, error) {
	if req.Role != "yetkili" && req.Role != "calisan" {
		return nil, errors.New("role must be 'yetkili' or 'calisan'")
	}

	exists, err := u.repo.IsAuthorityExists(req.TC, req.Email, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with given TC, email, or phone already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	auth := &models.Authority{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		TC:         req.TC,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   hashedPassword,
		Role:       req.Role,
		HospitalID: hospitalID,
	}

	if err := u.repo.CreateAuthority(auth); err != nil {
		return nil, err
	}

	return &dto.SubUserResponse{
		ID:        auth.ID,
		FirstName: auth.FirstName,
		LastName:  auth.LastName,
		TC:        auth.TC,
		Email:     auth.Email,
		Phone:     auth.Phone,
		Role:      auth.Role,
	}, nil
}

func (u *subUserUsecase) ListSubUsers(hospitalID uint) ([]dto.SubUserResponse, error) {
	users, err := u.repo.GetAllSubUsersByHospitalID(hospitalID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.SubUserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.SubUserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			TC:        u.TC,
			Email:     u.Email,
			Phone:     u.Phone,
			Role:      u.Role,
		})
	}
	return resp, nil
}

func (u *subUserUsecase) UpdateSubUser(id uint, req *dto.UpdateSubUserRequest, hospitalID uint) (*dto.SubUserResponse, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.HospitalID != hospitalID {
		return nil, errors.New("forbidden: cannot update user from another hospital")
	}

	if req.Role != "yetkili" && req.Role != "calisan" {
		return nil, errors.New("role must be 'yetkili' or 'calisan'")
	}

	//isUnique, err := u.repo.IsUniqueForUpdate(id, req.TC, req.Email, req.Phone)
	//if err != nil || !isUnique {
	//	return nil, errors.New("another user with given TC, email, or phone already exists")
	//}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.TC = req.TC
	user.Email = req.Email
	user.Phone = req.Phone
	user.Role = req.Role

	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}

	if err := u.repo.Update(user); err != nil {
		return nil, err
	}

	return &dto.SubUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		TC:        user.TC,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
	}, nil
}

func (u *subUserUsecase) DeleteSubUser(id, hospitalID uint) error {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	if user.HospitalID != hospitalID {
		return errors.New("forbidden: cannot delete user from another hospital")
	}
	return u.repo.Delete(user)
}
