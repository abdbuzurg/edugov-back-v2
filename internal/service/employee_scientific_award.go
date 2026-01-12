package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

func (s *Service) CreateEmployeeScientificAward(ctx context.Context, req *dto.CreateEmployeeScientificAwardRequest) (*dto.EmployeeScientificAwardResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.CreateEmployeeScientificAwardParams{
		EmployeeID:           req.EmployeeID,
		LanguageCode:         req.LanguageCode,
		ScientificAwardTitle: req.ScientificAwardTitle,
		GivenBy:              req.GivenBy,
	}

	created, err := s.store.Queries.CreateEmployeeScientificAward(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("CreateEmployeeScientificAward(service) -> CreateEmployeeScientificAward(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeScientificAwardResponse{
		ID:                   created.ID,
		ScientificAwardTitle: created.ScientificAwardTitle,
		GivenBy:              created.GivenBy,
		CreatedAt:            created.CreatedAt.Time,
		UpdatedAt:            created.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeeScientificAward(ctx context.Context, req *dto.UpdateEmployeeScientificAwardRequest) (*dto.EmployeeScientificAwardResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.UpdateEmployeeScientificAwardParams{
		ID:                   req.ID,
		ScientificAwardTitle: req.ScientificAwardTitle,
		GivenBy:              req.GivenBy,
	}

	updated, err := s.store.Queries.UpdateEmployeeScientificAward(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("UpdateEmployeeScientificAward(service) -> UpdateEmployeeScientificAward(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeScientificAwardResponse{
		ID:                   updated.ID,
		ScientificAwardTitle: updated.ScientificAwardTitle,
		GivenBy:              updated.GivenBy,
		CreatedAt:            updated.CreatedAt.Time,
		UpdatedAt:            updated.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteEmployeeScientificAward(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	if err := s.store.Queries.DeleteEmployeeScientificAward(ctx, id); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}
		return apperr.Internal(
			"internal error",
			fmt.Errorf("DeleteEmployeeScientificAward(service) -> DeleteEmployeeScientificAward(repo) params %v: %w", id, err),
		)
	}

	return nil
}

func (s *Service) GetEmployeeScientificAwardsByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeeScientificAwardResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeScientificAwardsByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}

	rows, err := s.store.Queries.GetEmployeeScientificAwardsByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("GetEmployeeScientificAwardsByEmployeeIDAndLanguageCode(service) -> GetEmployeeScientificAwardsByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err),
		)
	}

	result := make([]*dto.EmployeeScientificAwardResponse, len(rows))
	for i, r := range rows {
		result[i] = &dto.EmployeeScientificAwardResponse{
			ID:                   r.ID,
			ScientificAwardTitle: r.ScientificAwardTitle,
			GivenBy:              r.GivenBy,
			CreatedAt:            r.CreatedAt.Time,
			UpdatedAt:            r.UpdatedAt.Time,
		}
	}

	return result, nil
}
