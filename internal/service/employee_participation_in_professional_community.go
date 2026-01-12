package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

func (s *Service) CreateEmployeeParticipationInProfessionalCommunity(ctx context.Context, req *dto.CreateEmployeeParticipationInProfessionalCommunityRequest) (*dto.EmployeeParticipationInProfessionalCommunityResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeePIPCArgs := sqlc.CreateEmployeeParticipationInProfessionalCommunityParams{
		EmployeeID:                  req.EmployeeID,
		LanguageCode:                req.LanguageCode,
		ProfessionalCommunityTitle:  req.ProfessionalCommunityTitle,
		RoleInProfessionalCommunity: req.RoleInProfessionalCommunity,
	}
	createdEmployeePIPCResult, err := s.store.Queries.CreateEmployeeParticipationInProfessionalCommunity(ctx, createdEmployeePIPCArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeeParticipationInProfessionalCommunity (service) -> CreateEmployeeParticipationInProfessionalCommunity (repo) params %v: %w", createdEmployeePIPCArgs, err))
	}

	return &dto.EmployeeParticipationInProfessionalCommunityResponse{
		ID:                          createdEmployeePIPCResult.ID,
		ProfessionalCommunityTitle:  createdEmployeePIPCResult.ProfessionalCommunityTitle,
		RoleInProfessionalCommunity: createdEmployeePIPCResult.RoleInProfessionalCommunity,
		CreatedAt:                   createdEmployeePIPCResult.CreatedAt.Time,
		UpdatedAt:                   createdEmployeePIPCResult.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeeParticipationInProfessionalCommunity(ctx context.Context, req *dto.UpdateEmployeeParticipationInProfessionalCommunityRequest) (*dto.EmployeeParticipationInProfessionalCommunityResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	updatedEmployeePIPCArgs := sqlc.UpdateEmployeeParticipationInProfessionalCommunityParams{
		ID:                          req.ID,
		ProfessionalCommunityTitle:  req.ProfessionalCommunityTitle,
		RoleInProfessionalCommunity: req.RoleInProfessionalCommunity,
	}
	updatedEmployeePIPCResult, err := s.store.Queries.UpdateEmployeeParticipationInProfessionalCommunity(ctx, updatedEmployeePIPCArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeParticipationInProfessionalCommunity (service) -> UpdateEmployeeParticipationInProfessionalCommunity (repo) params %v: %w", updatedEmployeePIPCArgs, err))
	}

	return &dto.EmployeeParticipationInProfessionalCommunityResponse{
		ID:                          updatedEmployeePIPCResult.ID,
		ProfessionalCommunityTitle:  updatedEmployeePIPCResult.ProfessionalCommunityTitle,
		RoleInProfessionalCommunity: updatedEmployeePIPCResult.RoleInProfessionalCommunity,
	}, nil
}

func (s *Service) DeleteEmployeeParticipationInProfessionalCommunity(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	if err := s.store.Queries.DeleteEmployeeParticipationInProfessionalCommunity(ctx, id); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}

		return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeeParticipationInProfessionalCommunity (service) -> DeleteEmployeeParticipationInProfessionalCommunity (repo) %d: %w", id, err))
	}

	return nil
}

func (s *Service) GetEmployeeParticipationInProfessionalCommunitiesByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeeParticipationInProfessionalCommunityResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}
	employeePIPCs, err := s.store.Queries.GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCode (service) -> GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCode (repo) params %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeeParticipationInProfessionalCommunityResponse, len(employeePIPCs))
	for index, pipc := range employeePIPCs {
		result[index] = &dto.EmployeeParticipationInProfessionalCommunityResponse{
			ID:                          pipc.ID,
			ProfessionalCommunityTitle:  pipc.ProfessionalCommunityTitle,
			RoleInProfessionalCommunity: pipc.RoleInProfessionalCommunity,
			CreatedAt:                   pipc.CreatedAt.Time,
			UpdatedAt:                   pipc.UpdatedAt.Time,
		}
	}

	return result, nil
}
