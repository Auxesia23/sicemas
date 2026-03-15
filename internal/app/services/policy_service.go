package services

import (
	"context"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
)

type PolicyService interface {
	AddPolicy(policy dto.PolicyRequest, actorId uuid.UUID) error
	RemovePolicy(policy dto.PolicyRequest, actorId uuid.UUID) error
	GetFilteredPolicy(filter string) ([]dto.PolicyResponse, error)
}

type policyServiceImpl struct {
	activityRepo repositories.ActivityRepository
	enforcer     *casbin.Enforcer
}

func NewPolicyService(enforcer *casbin.Enforcer, activityRepo repositories.ActivityRepository) PolicyService {
	return &policyServiceImpl{
		activityRepo: activityRepo,
		enforcer:     enforcer,
	}
}

func (s *policyServiceImpl) AddPolicy(policy dto.PolicyRequest, actorId uuid.UUID) error {
	ok, err := s.enforcer.AddPolicy(policy.Subject, policy.Object, policy.Action)
	if !ok {
		return apperror.NewBadRequest("Policy sudah ada!")
	}
	if err != nil {
		return apperror.NewInternal("Terjadi Kesalahan")
	}
	ctx := context.Background()
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menambahkan policy",
		EntityType: "POLICY",
		EntityID:   uuid.Nil.String(),
		TargetName: policy.Object + ":" + policy.Action,
	})
	return nil
}

func (s *policyServiceImpl) GetFilteredPolicy(filter string) ([]dto.PolicyResponse, error) {
	policies, err := s.enforcer.GetFilteredPolicy(0, filter)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}

	if len(policies) == 0 {
		return []dto.PolicyResponse{}, nil
	}

	var response []dto.PolicyResponse
	for _, p := range policies {
		response = append(response, dto.PolicyResponse{
			Subject: p[0],
			Object:  p[1],
			Action:  p[2],
		})
	}
	return response, nil
}

func (s *policyServiceImpl) RemovePolicy(policy dto.PolicyRequest, actorId uuid.UUID) error {
	ok, err := s.enforcer.RemovePolicy(policy.Subject, policy.Object, policy.Action)
	if !ok {
		return apperror.NewBadRequest("Policy tidak ditemukan!")
	}
	if err != nil {
		return apperror.NewInternal("Terjadi Kesalahan")
	}
	ctx := context.Background()
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus policy",
		EntityType: "POLICY",
		EntityID:   uuid.Nil.String(),
		TargetName: policy.Object + ":" + policy.Action,
	})
	return nil
}
