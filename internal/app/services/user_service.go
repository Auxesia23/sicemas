package services

import (
	"context"
	"database/sql"
	"log"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserService interface {
	Register(ctx context.Context, in *dto.UserRequest) error
	GetAllUser(ctx context.Context) ([]dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, in *dto.UserRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetProfile(ctx context.Context, id string) (*dto.UserResponse, error)
}

type userServiceImpl struct {
	userRepo repositories.UserRepo
	enforcer *casbin.Enforcer
	cache    cache.Cache
}

func NewUserService(userRepo repositories.UserRepo, enforcer *casbin.Enforcer, cache cache.Cache) UserService {
	return &userServiceImpl{
		userRepo,
		enforcer,
		cache,
	}
}

func (s *userServiceImpl) Register(ctx context.Context, in *dto.UserRequest) error {
	var err error
	index := utils.HashIndex(in.NIP)
	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		return err
	}
	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		return err
	}
	in.NIP, err = utils.Encrypt(in.NIP)
	if err != nil {
		return err
	}

	if in.UnitKerja == nil || *in.UnitKerja == "" {
		defaultUnit := "KUA Kecamatan Ciemas"
		in.UnitKerja = &defaultUnit
	}

	userID, err := s.userRepo.Create(ctx, in, index)
	if err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			return apperror.NewBadRequest("User dengan NIP ini sudah ada.")
		}
		return apperror.NewInternal("Terjadi Kesalahan saat menyimpan user.")
	}

	if in.Peran != "" {
		_, err = s.enforcer.AddGroupingPolicy(userID, in.Peran)
		if err != nil {
			return apperror.NewInternal("User dibuat, tapi gagal memberikan hak akses.")
		}
		s.enforcer.LoadPolicy()
	}

	return nil
}

func (s *userServiceImpl) GetAllUser(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.userRepo.ReadAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	if len(users) == 0 {
		return []dto.UserResponse{}, nil
	}

	var response []dto.UserResponse
	for _, u := range users {
		roles, _ := s.enforcer.GetRolesForUser(u.ID.String())
		peran := ""
		if len(roles) > 0 {
			peran = roles[0]
		}
		decryptedNIP, err := utils.Decrypt(u.NIP)
		if err != nil {
			return nil, err
		}
		decryptedEmail, err := utils.Decrypt(u.Email)
		if err != nil {
			return nil, err
		}
		decryptedNomorTelepon, err := utils.Decrypt(u.NomorTelepon)
		if err != nil {
			return nil, err
		}
		response = append(response, dto.UserResponse{
			ID:           u.ID,
			NIP:          decryptedNIP,
			NamaLengkap:  u.NamaLengkap,
			Jabatan:      u.Jabatan,
			Peran:        peran,
			UnitKerja:    u.UnitKerja,
			Email:        decryptedEmail,
			NomorTelepon: decryptedNomorTelepon,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
		})
	}
	return response, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, in *dto.UserRequest) error {
	var err error
	index := utils.HashIndex(in.NIP)
	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		return err
	}
	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		return err
	}
	in.NIP, err = utils.Encrypt(in.NIP)
	if err != nil {
		return err
	}
	if err := s.userRepo.Update(ctx, id, in, index); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("User tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	_, err = s.enforcer.RemoveFilteredGroupingPolicy(0, id.String())
	if err != nil {
		return apperror.NewInternal("Gagal membersihkan peran lama")
	}

	if in.Peran != "" {
		_, err = s.enforcer.AddGroupingPolicy(id.String(), in.Peran)
		if err != nil {
			return apperror.NewInternal("Gagal memperbarui peran baru")
		}
	}

	s.enforcer.LoadPolicy()
	return nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("User tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	_, err := s.enforcer.RemoveFilteredGroupingPolicy(0, id.String())
	if err != nil {
		// Log error tapi mungkin jangan gagalkan response karena user di DB sudah terhapus
		log.Printf("Warning: Gagal hapus policy Casbin untuk user %s", id)
	}

	s.enforcer.LoadPolicy()
	return nil
}

func (s *userServiceImpl) GetProfile(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.ReadById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFound("User tidak ditemukan")
		}
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}
	roles, _ := s.enforcer.GetRolesForUser(user.ID.String())
	peran := ""
	if len(roles) > 0 {
		peran = roles[0]
	}
	decryptedNIP, err := utils.Decrypt(user.NIP)
	if err != nil {
		return nil, err
	}
	decryptedEmail, err := utils.Decrypt(user.Email)
	if err != nil {
		return nil, err
	}
	decryptedNomorTelepon, err := utils.Decrypt(user.NomorTelepon)
	if err != nil {
		return nil, err
	}
	response := &dto.UserResponse{
		ID:           user.ID,
		NIP:          decryptedNIP,
		NamaLengkap:  user.NamaLengkap,
		Jabatan:      user.Jabatan,
		Peran:        peran,
		UnitKerja:    user.UnitKerja,
		Email:        decryptedEmail,
		NomorTelepon: decryptedNomorTelepon,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	return response, nil
}
