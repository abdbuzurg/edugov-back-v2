package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

func (s *Service) CreateEmployeeResearchActivity(ctx context.Context, req *dto.CreateEmployeeResearchActivityRequest) (*dto.EmployeeResearchActivityResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.CreateEmployeeResearchActivityParams{
		EmployeeID:            req.EmployeeID,
		LanguageCode:          req.LanguageCode,
		ResearchActivityTitle: req.ResearchActivityTitle,
		EmployeeRole:          req.EmployeeRole,
	}

	created, err := s.store.Queries.CreateEmployeeResearchActivity(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("CreateEmployeeResearchActivity(service) -> CreateEmployeeResearchActivity(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeResearchActivityResponse{
		ID:                    created.ID,
		ResearchActivityTitle: created.ResearchActivityTitle,
		EmployeeRole:          created.EmployeeRole,
		CreatedAt:             created.CreatedAt.Time,
		UpdatedAt:             created.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeeResearchActivity(ctx context.Context, req *dto.UpdateEmployeeResearchActivityRequest) (*dto.EmployeeResearchActivityResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.UpdateEmployeeResearchActivityParams{
		ID:                    req.ID,
		ResearchActivityTitle: req.ResearchActivityTitle,
		EmployeeRole:          req.EmployeeRole,
	}

	updated, err := s.store.Queries.UpdateEmployeeResearchActivity(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("UpdateEmployeeResearchActivity(service) -> UpdateEmployeeResearchActivity(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeResearchActivityResponse{
		ID:                    updated.ID,
		ResearchActivityTitle: updated.ResearchActivityTitle,
		EmployeeRole:          updated.EmployeeRole,
		CreatedAt:             updated.CreatedAt.Time,
		UpdatedAt:             updated.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteEmployeeResearchActivity(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	err := s.store.Queries.DeleteEmployeeResearchActivity(ctx, id)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}
		return apperr.Internal(
			"internal error",
			fmt.Errorf("DeleteEmployeeResearchActivity(service) -> DeleteEmployeeResearchActivity(repo) params %v: %w", id, err),
		)
	}

	return nil
}

func (s *Service) GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeeResearchActivityResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}

	rows, err := s.store.Queries.GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode(service) -> GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err),
		)
	}

	result := make([]*dto.EmployeeResearchActivityResponse, len(rows))
	for i, r := range rows {
		result[i] = &dto.EmployeeResearchActivityResponse{
			ID:                    r.ID,
			ResearchActivityTitle: r.ResearchActivityTitle,
			EmployeeRole:          r.EmployeeRole,
			CreatedAt:             r.CreatedAt.Time,
			UpdatedAt:             r.UpdatedAt.Time,
		}
	}

	return result, nil
}
