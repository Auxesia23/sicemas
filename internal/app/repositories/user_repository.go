package repositories

import (
	"context"
	"database/sql"
	"log"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(ctx context.Context, in *dto.UserRequest, index []byte) error
	ReadAll(ctx context.Context) ([]entity.User, error)
	ReadOne(ctx context.Context, index []byte) (*entity.User, error)
	ReadById(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, in *dto.UserRequest, newIndex []byte) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepiImpl struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepiImpl{
		DB: db,
	}
}

func (r *userRepiImpl) Create(ctx context.Context, in *dto.UserRequest, index []byte) error {
	query := `
		INSERT INTO users (nip_index,nama_lengkap,nip,email,nomor_telepon)
		VALUES($1,$2,$3,$4,$5);
	`

	_, err := r.DB.ExecContext(ctx, query, index, in.NamaLengkap, in.NIP, in.Email, in.NomorTelepon)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepiImpl) ReadAll(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT *
		FROM users;
	`
	var dest []entity.User
	if err := r.DB.SelectContext(ctx, &dest, query); err != nil {
		return nil, err
	}

	return dest, nil
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
		log.Println(err.Error())
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

func (r *userRepiImpl) Update(ctx context.Context, id uuid.UUID, in *dto.UserRequest, newIndex []byte) error {
	query := `
        UPDATE users
        SET nip_index = $1,
            nama_lengkap = $2,
            nip = $3,
            email = $4,
            nomor_telepon = $5,
            updated_at = NOW()
        WHERE id = $6;
    `

	result, err := r.DB.ExecContext(ctx, query,
		newIndex,
		in.NamaLengkap,
		in.NIP,
		in.Email,
		in.NomorTelepon,
		id,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepiImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id=$1;
	`
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
