package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

func (s *Service) CreateEmployeePatent(ctx context.Context, req *dto.CreateEmployeePatentRequest) (*dto.EmployeePatentResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeePatentArgs := sqlc.CreateEmployeePatentParams{
		EmployeeID:   req.EmployeeID,
		LanguageCode: req.LanguageCode,
		PatentTitle:  req.PatentTitle,
		Description:  req.Description,
	}
	createdEmployeePatentResult, err := s.store.Queries.CreateEmployeePatent(ctx, createdEmployeePatentArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeePatent (service) -> CreateEmployeePatent(repo) params %v: %w", createdEmployeePatentArgs, err))
	}

	return &dto.EmployeePatentResponse{
		ID:          createdEmployeePatentResult.ID,
		PatentTitle: createdEmployeePatentArgs.PatentTitle,
		Description: createdEmployeePatentArgs.Description,
		CreatedAt:   createdEmployeePatentResult.CreatedAt.Time,
		UpdatedAt:   createdEmployeePatentResult.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeePatent(ctx context.Context, req *dto.UpdateEmployeePatentRequest) (*dto.EmployeePatentResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	updatedEmployeePatentArgs := sqlc.UpdateEmployeePatentParams{
		ID:          req.ID,
		PatentTitle: req.PatentTitle,
		Description: req.Description,
	}
	updatedEmployeePatentResult, err := s.store.Queries.UpdateEmployeePatent(ctx, updatedEmployeePatentArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeePatent (service) -> CreateEmployeUpdateEmployeePatentePatent(repo) params %v: %w", updatedEmployeePatentArgs, err))
	}

	return &dto.EmployeePatentResponse{
		ID:          updatedEmployeePatentResult.ID,
		PatentTitle: updatedEmployeePatentResult.PatentTitle,
		Description: updatedEmployeePatentResult.Description,
		CreatedAt:   updatedEmployeePatentResult.CreatedAt.Time,
		UpdatedAt:   updatedEmployeePatentResult.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteEmployeePatent(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	err := s.store.Queries.DeleteEmployeePatent(ctx, id)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}

		return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeePatent (service) -> DeleteEmployeePatent (repo) params %d: %w", id, err))
	}

	return nil
}

func (s *Service) GetEmployeePatentsByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeePatentResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"emplyoeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeePatentsByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}
	employeePatentsResult, err := s.store.Queries.GetEmployeePatentsByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeePatentsByEmployeeIDAndLanguageCode(service) -> GetEmployeePatentsByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeePatentResponse, len(employeePatentsResult))
	for index, patent := range employeePatentsResult {
		result[index] = &dto.EmployeePatentResponse{
			ID:          patent.ID,
			PatentTitle: patent.PatentTitle,
			Description: patent.Description,
			CreatedAt:   patent.CreatedAt.Time,
			UpdatedAt:   patent.UpdatedAt.Time,
		}
	}

	return result, nil
}
