package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/repositories"
	"sicemas/internal/cache"
	"sicemas/internal/dto"
	"sicemas/internal/entity"
	"sicemas/internal/utils"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserService interface {
	Register(ctx context.Context, in *dto.UserRequest, actorId uuid.UUID) error
	GetAllUser(ctx context.Context) ([]dto.UserResponse, error)
	UpdateUser(ctx context.Context, id, actorId uuid.UUID, in *dto.UserRequest) error
	DeleteUser(ctx context.Context, id, actorId uuid.UUID) error
	GetProfile(ctx context.Context, id string) (*dto.UserResponse, error)
	GetPermissions(ctx context.Context, id string) (*dto.UserPermission, error)
}

type userServiceImpl struct {
	userRepo     repositories.UserRepo
	activityRepo repositories.ActivityRepository
	enforcer     *casbin.Enforcer
	cache        cache.Cache
	logger       *slog.Logger
}

func NewUserService(
	userRepo repositories.UserRepo,
	activityRepo repositories.ActivityRepository,
	enforcer *casbin.Enforcer,
	cache cache.Cache,
	logger *slog.Logger,
) UserService {
	return &userServiceImpl{
		userRepo:     userRepo,
		activityRepo: activityRepo,
		enforcer:     enforcer,
		cache:        cache,
		logger:       logger,
	}
}

func (s *userServiceImpl) Register(ctx context.Context, in *dto.UserRequest, actorId uuid.UUID) error {
	s.logger.Info("initiating user registration", "actor_id", actorId)

	var err error
	index := utils.HashIndex(in.NIP)
	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		s.logger.Error("failed to encrypt email during registration", "error", err)
		return err
	}
	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to encrypt phone number during registration", "error", err)
		return err
	}
	in.NIP, err = utils.Encrypt(in.NIP)
	if err != nil {
		s.logger.Error("failed to encrypt NIP during registration", "error", err)
		return err
	}

	if in.UnitKerja == nil || *in.UnitKerja == "" {
		defaultUnit := "KUA Kecamatan Ciemas"
		in.UnitKerja = &defaultUnit
	}

	userID, err := s.userRepo.Create(ctx, in, index)
	if err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			s.logger.Warn("user with this NIP already exists", "user_index", index)
			return apperror.NewBadRequest("User dengan NIP ini sudah ada.")
		}
		s.logger.Error("failed to create user in database", "user_index", index, "error", err)
		return apperror.NewInternal("Terjadi Kesalahan saat menyimpan user.")
	}

	if in.Peran != "" {
		_, err = s.enforcer.AddGroupingPolicy(userID, in.Peran)
		if err != nil {
			s.logger.Error("failed to assign role to user", "user_id", userID, "role", in.Peran, "error", err)
			return apperror.NewInternal("User dibuat, tapi gagal memberikan hak akses.")
		}
		s.enforcer.LoadPolicy()
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menambahkan petugas",
		EntityType: "USER",
		EntityID:   userID,
		TargetName: in.NamaLengkap,
	})

	s.logger.Info("user registered successfully", "user_id", userID)
	return nil
}

func (s *userServiceImpl) GetAllUser(ctx context.Context) ([]dto.UserResponse, error) {
	s.logger.Info("fetching all users")

	users, err := s.userRepo.ReadAll(ctx)
	if err != nil {
		s.logger.Error("failed to fetch users", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	if len(users) == 0 {
		s.logger.Info("no users found")
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
			s.logger.Error("failed to decrypt NIP", "user_id", u.ID, "error", err)
			return nil, err
		}
		decryptedEmail, err := utils.Decrypt(u.Email)
		if err != nil {
			s.logger.Error("failed to decrypt email", "user_id", u.ID, "error", err)
			return nil, err
		}
		decryptedNomorTelepon, err := utils.Decrypt(u.NomorTelepon)
		if err != nil {
			s.logger.Error("failed to decrypt phone number", "user_id", u.ID, "error", err)
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

	s.logger.Info("fetched users successfully", "count", len(response))
	return response, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id, actorId uuid.UUID, in *dto.UserRequest) error {
	s.logger.Info("initiating user update", "user_id", id, "actor_id", actorId)

	var err error
	index := utils.HashIndex(in.NIP)
	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		s.logger.Error("failed to encrypt email during update", "user_id", id, "error", err)
		return err
	}
	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to encrypt phone number during update", "user_id", id, "error", err)
		return err
	}
	in.NIP, err = utils.Encrypt(in.NIP)
	if err != nil {
		s.logger.Error("failed to encrypt NIP during update", "user_id", id, "error", err)
		return err
	}

	if err := s.userRepo.Update(ctx, id, in, index); err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == "23505" {
			s.logger.Warn("user with this NIP already exists", "user_index", index)
			return apperror.NewBadRequest("User dengan NIP ini sudah ada.")
		}
		if err == sql.ErrNoRows {
			s.logger.Warn("user not found for update", "user_id", id)
			return apperror.NewNotFound("User tidak ditemukan")
		}
		s.logger.Error("failed to update user in database", "user_id", id, "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_, err = s.enforcer.RemoveFilteredGroupingPolicy(0, id.String())
	if err != nil {
		s.logger.Error("failed to clear old roles", "user_id", id, "error", err)
		return apperror.NewInternal("Gagal membersihkan peran lama")
	}

	if in.Peran != "" {
		_, err = s.enforcer.AddGroupingPolicy(id.String(), in.Peran)
		if err != nil {
			s.logger.Error("failed to assign new role", "user_id", id, "role", in.Peran, "error", err)
			return apperror.NewInternal("Gagal memperbarui peran baru")
		}
	}

	s.enforcer.LoadPolicy()
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memperbarui data petugas",
		EntityType: "USER",
		EntityID:   id.String(),
		TargetName: in.NamaLengkap,
	})

	s.logger.Info("user updated successfully", "user_id", id)
	return nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id, actorId uuid.UUID) error {
	s.logger.Info("deleting user", "user_id", id, "actor_id", actorId)

	userTarget, err := s.userRepo.ReadById(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user not found for deletion", "user_id", id)
			return apperror.NewNotFound("User tidak ditemukan")
		}
		s.logger.Error("failed to read user for deletion", "error", err)
		return apperror.NewInternal("Terjadi kesalahan saat mencari data user")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user not found during deletion", "user_id", id)
			return apperror.NewNotFound("User tidak ditemukan")
		}
		s.logger.Error("failed to delete user", "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_, err = s.enforcer.RemoveFilteredGroupingPolicy(0, id.String())
	if err != nil {
		s.logger.Warn("failed to remove Casbin policies for user", "user_id", id, "error", err)
	}

	s.enforcer.LoadPolicy()
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus petugas",
		EntityType: "USER",
		EntityID:   id.String(),
		TargetName: userTarget.NamaLengkap,
	})

	s.logger.Info("user deleted successfully", "user_id", id)
	return nil
}

func (s *userServiceImpl) GetProfile(ctx context.Context, id string) (*dto.UserResponse, error) {
	s.logger.Info("fetching user profile", "user_id", id)

	user, err := s.userRepo.ReadById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("user not found for profile", "user_id", id)
			return nil, apperror.NewNotFound("User tidak ditemukan")
		}
		s.logger.Error("failed to read user profile", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	roles, _ := s.enforcer.GetRolesForUser(user.ID.String())
	peran := ""
	if len(roles) > 0 {
		peran = roles[0]
	}
	decryptedNIP, err := utils.Decrypt(user.NIP)
	if err != nil {
		s.logger.Error("failed to decrypt NIP for profile", "user_id", id, "error", err)
		return nil, err
	}
	decryptedEmail, err := utils.Decrypt(user.Email)
	if err != nil {
		s.logger.Error("failed to decrypt email for profile", "user_id", id, "error", err)
		return nil, err
	}
	decryptedNomorTelepon, err := utils.Decrypt(user.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to decrypt phone number for profile", "user_id", id, "error", err)
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

	s.logger.Info("user profile fetched successfully", "user_id", id)
	return response, nil
}

func (s *userServiceImpl) GetPermissions(ctx context.Context, id string) (*dto.UserPermission, error) {
	s.logger.Info("fetching user permissions", "user_id", id)

	user, err := s.userRepo.ReadById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("user not found for permissions", "user_id", id)
			return nil, apperror.NewNotFound("User tidak ditemukan")
		}
		s.logger.Error("failed to read user for permissions", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	permissions, err := s.enforcer.GetImplicitPermissionsForUser(id)
	if err != nil {
		s.logger.Error("failed to get user permissions", "user_id", id, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}

	var formattedPermissions []string
	for _, permission := range permissions {
		formattedPermissions = append(formattedPermissions, fmt.Sprintf("%s:%s", permission[1], permission[2]))
	}

	response := &dto.UserPermission{
		NamaLengkap: user.NamaLengkap,
		Permissions: formattedPermissions,
	}

	s.logger.Info("user permissions fetched successfully", "user_id", id, "permission_count", len(formattedPermissions))
	return response, nil
}
