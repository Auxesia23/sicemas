package repositories

import (
	"context"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(ctx context.Context, in *dto.UserRegister, index []byte) error
	ReadAll(ctx context.Context) (*[]entity.User, error)
	ReadOne(ctx context.Context, index []byte) (*entity.User, error)
	ReadById(ctx context.Context, id string) (*entity.User, error)
}

type userRepiImpl struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepiImpl{
		DB: db,
	}
}

func (r *userRepiImpl) Create(ctx context.Context, in *dto.UserRegister, index []byte) error {
	query := `
		INSERT INTO users (nip_index,nama_lengkap,peran,nip,email,nomor_telepon)
		VALUES($1,$2,$3,$4,$5,$6);
	`

	_, err := r.DB.ExecContext(ctx, query, index, in.NamaLengkap, in.Peran, in.NIP, in.Email, in.NomorTelepon)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepiImpl) ReadAll(ctx context.Context) (*[]entity.User, error) {
	query := `
		SELECT *
		FROM users;
	`
	var dest []entity.User
	if err := r.DB.SelectContext(ctx, &dest, query); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (r *userRepiImpl) ReadOne(ctx context.Context, index []byte) (*entity.User, error) {
	query := `
		SELECT *
		FROM users
		WHERE nip_index = $1
		LIMIT 1;
	`

	var dest entity.User
	if err := r.DB.GetContext(ctx, &dest, query, index); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (r *userRepiImpl) ReadById(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT *
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	var dest entity.User
	if err := r.DB.GetContext(ctx, &dest, query, id); err != nil {
		return nil, err
	}

	return &dest, nil
}
