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
	index := utils.HashIndex(in.NIP)
	in.Email = utils.Encrypt(in.Email)
	in.NomorTelepon = utils.Encrypt(in.NomorTelepon)
	in.NIP = utils.Encrypt(in.NIP)

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

	var response []dto.UserResponse
	for _, u := range users {
		roles, _ := s.enforcer.GetRolesForUser(u.ID.String())
		response = append(response, dto.UserResponse{
			ID:           u.ID,
			NIP:          utils.Decrypt(u.NIP),
			NamaLengkap:  u.NamaLengkap,
			Jabatan:      u.Jabatan,
			Peran:        roles[0],
			UnitKerja:    u.UnitKerja,
			Email:        utils.Decrypt(u.Email),
			NomorTelepon: utils.Decrypt(u.NomorTelepon),
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
		})
	}
	return response, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, in *dto.UserRequest) error {
	index := utils.HashIndex(in.NIP)
	in.Email = utils.Encrypt(in.Email)
	in.NomorTelepon = utils.Encrypt(in.NomorTelepon)
	in.NIP = utils.Encrypt(in.NIP)
	if err := s.userRepo.Update(ctx, id, in, index); err != nil {
		if err == sql.ErrNoRows {
			return apperror.NewNotFound("User tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}
	_, err := s.enforcer.RemoveFilteredGroupingPolicy(0, id.String())
	if err != nil {
		return apperror.NewInternal("Gagal membersihkan peran lama")
	}

	if in.Peran != "" {
		_, err = s.enforcer.AddGroupingPolicy(id.String(), in.Peran)
		if err != nil {
			return apperror.NewInternal("Gagal memperbarui peran baru")
		}
	}

	// 4. Refresh Policy agar perubahan langsung aktif di memori
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
	response := &dto.UserResponse{
		ID:           user.ID,
		NIP:          utils.Decrypt(user.NIP),
		NamaLengkap:  user.NamaLengkap,
		Jabatan:      user.Jabatan,
		Peran:        roles[0],
		UnitKerja:    user.UnitKerja,
		Email:        utils.Decrypt(user.Email),
		NomorTelepon: utils.Decrypt(user.NomorTelepon),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
	return response, nil
}
