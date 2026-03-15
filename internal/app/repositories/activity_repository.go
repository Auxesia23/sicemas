package repositories

import (
	"context"
	"log"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository interface {
	InsertActivity(ctx context.Context, activity *entity.Activity) error
	GetActivities(ctx context.Context) ([]dto.ActivityResponse, error)
}

type ActivityRepositoryImpl struct {
	DB *sqlx.DB
}

func NewActivityRepo(db *sqlx.DB) *ActivityRepositoryImpl {
	return &ActivityRepositoryImpl{DB: db}
}

func (r *ActivityRepositoryImpl) InsertActivity(ctx context.Context, activity *entity.Activity) error {
	q := `
	INSERT INTO activities
		(user_id, action_type, entity_type, entity_id, target_name)
		VALUES ($1, $2, $3, $4, $5);`
	if _, err := r.DB.ExecContext(ctx, q, activity.UserID, activity.ActionType, activity.EntityType, activity.EntityID, activity.TargetName); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *ActivityRepositoryImpl) GetActivities(ctx context.Context) ([]dto.ActivityResponse, error) {
	q := `
	SELECT
		a.id,
		u.nama_lengkap AS actor_name,
		a.action_type,
		a.target_name,
		a.entity_type,
		a.created_at
	FROM
		activities a
	JOIN
		users u ON a.user_id = u.id
	ORDER BY a.created_at DESC
	LIMIT 5;
	`
	var activities []dto.ActivityResponse
	if err := r.DB.SelectContext(ctx, &activities, q); err != nil {
		log.Println(err)
		return nil, err
	}
	return activities, nil
}
