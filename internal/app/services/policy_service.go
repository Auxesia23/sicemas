package services

import (
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/dto"

	"github.com/casbin/casbin/v2"
)

type PolicyService interface {
	AddPolicy(dto.PolicyRequest) error
	RemovePolicy(policy dto.PolicyRequest) error
	GetFilteredPolicy(filter string) ([]dto.PolicyResponse, error)
}

type policyServiceImpl struct {
	e *casbin.Enforcer
}

func NewPolicyService(e *casbin.Enforcer) PolicyService {
	return &policyServiceImpl{
		e,
	}
}

func (s *policyServiceImpl) AddPolicy(policy dto.PolicyRequest) error {
	ok, err := s.e.AddPolicy(policy.Subject, policy.Object, policy.Action)
	if !ok {
		return apperror.NewBadRequest("Policy sudah ada!")
	}
	if err != nil {
		return apperror.NewInternal("Terjadi Kesalahan")
	}
	return nil
}

func (s *policyServiceImpl) GetFilteredPolicy(filter string) ([]dto.PolicyResponse, error) {
	policies, err := s.e.GetFilteredPolicy(0, filter)
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
