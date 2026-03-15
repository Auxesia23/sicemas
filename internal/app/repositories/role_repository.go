package repositories

import (
	"context"
	"database/sql"
	"log"
	"situs-keagamaan/internal/entity"

	"github.com/jmoiron/sqlx"
)

type RoleRepository interface {
	Create(ctx context.Context, name string) (*entity.Role, error)
	ReadAll(ctx context.Context) ([]entity.Role, error)
	ReadOne(ctx context.Context, id string) (*entity.Role, error)
	Delete(ctx context.Context, name string) error
}

type roleRepositoryImpl struct {
	DB *sqlx.DB
}

func NewRoleRepo(db *sqlx.DB) RoleRepository {
	return &roleRepositoryImpl{
		DB: db,
	}
}

func (r *roleRepositoryImpl) Create(ctx context.Context, name string) (*entity.Role, error) {
	query := `
		INSERT INTO roles (name)
		VALUES ($1) RETURNING *;
	`
	var newRole entity.Role
	if err := r.DB.GetContext(ctx, &newRole, query, name); err != nil {
		return nil, err
	}
	return &newRole, nil
}

func (r *roleRepositoryImpl) ReadAll(ctx context.Context) ([]entity.Role, error) {
	query := `
		SELECT * FROM roles;
	`
	var roles []entity.Role
	if err := r.DB.SelectContext(ctx, &roles, query); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return roles, nil
}

func (r *roleRepositoryImpl) Delete(ctx context.Context, name string) error {
	query := `
		DELETE FROM roles
		WHERE id = $1;
	`
	result, err := r.DB.ExecContext(ctx, query, name)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *roleRepositoryImpl) ReadOne(ctx context.Context, id string) (*entity.Role, error) {
	query := `
		SELECT * FROM roles
		WHERE id = $1;
	`
	var role entity.Role
	if err := r.DB.GetContext(ctx, &role, query, id); err != nil {
		return nil, err
	}
	return &role, nil
}
