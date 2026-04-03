package services

import (
	"context"
	"log/slog"
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
	logger       *slog.Logger
}

func NewPolicyService(
	enforcer *casbin.Enforcer,
	activityRepo repositories.ActivityRepository,
	logger *slog.Logger,
) PolicyService {
	return &policyServiceImpl{
		activityRepo: activityRepo,
		enforcer:     enforcer,
		logger:       logger,
	}
}

func (s *policyServiceImpl) AddPolicy(policy dto.PolicyRequest, actorId uuid.UUID) error {
	s.logger.Info("adding policy", "subject", policy.Subject, "object", policy.Object, "action", policy.Action, "actor_id", actorId)

	ok, err := s.enforcer.AddPolicy(policy.Subject, policy.Object, policy.Action)
	if !ok {
		s.logger.Warn("policy already exists", "subject", policy.Subject, "object", policy.Object, "action", policy.Action)
		return apperror.NewBadRequest("Policy sudah ada!")
	}
	if err != nil {
		s.logger.Error("failed to add policy", "error", err)
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

	s.logger.Info("policy added successfully", "subject", policy.Subject, "object", policy.Object, "action", policy.Action)
	return nil
}

func (s *policyServiceImpl) GetFilteredPolicy(filter string) ([]dto.PolicyResponse, error) {
	s.logger.Info("fetching filtered policies", "filter", filter)

	policies, err := s.enforcer.GetFilteredPolicy(0, filter)
	if err != nil {
		s.logger.Error("failed to get filtered policies", "error", err)
		return nil, apperror.NewInternal("Terjadi Kesalahan")
	}

	if len(policies) == 0 {
		s.logger.Info("no policies found for filter", "filter", filter)
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

	s.logger.Info("fetched filtered policies successfully", "count", len(response))
	return response, nil
}

func (s *policyServiceImpl) RemovePolicy(policy dto.PolicyRequest, actorId uuid.UUID) error {
	s.logger.Info("removing policy", "subject", policy.Subject, "object", policy.Object, "action", policy.Action, "actor_id", actorId)

	ok, err := s.enforcer.RemovePolicy(policy.Subject, policy.Object, policy.Action)
	if !ok {
		s.logger.Warn("policy not found for removal", "subject", policy.Subject, "object", policy.Object, "action", policy.Action)
		return apperror.NewBadRequest("Policy tidak ditemukan!")
	}
	if err != nil {
		s.logger.Error("failed to remove policy", "error", err)
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

	s.logger.Info("policy removed successfully", "subject", policy.Subject, "object", policy.Object, "action", policy.Action)
	return nil
}
