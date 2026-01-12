package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

func (s *Service) CreateEmployeePublication(ctx context.Context, req *dto.CreateEmployeePublicationRequest) (*dto.EmployeePublicationResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeePublicationArgs := sqlc.CreateEmployeePublicationParams{
		EmployeeID:        req.EmployeeID,
		LanguageCode:      req.LanguageCode,
		PublicationTitle:  req.PublicationTitle,
		LinkToPublication: req.LinkToPublication,
	}
	createdEmployeePublication, err := s.store.Queries.CreateEmployeePublication(ctx, createdEmployeePublicationArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeePublication(service) -> CreateEmployeePublication(repo) params %v: %w", createdEmployeePublicationArgs, err))
	}

	return &dto.EmployeePublicationResponse{
		ID:                createdEmployeePublication.ID,
		PublicationTitle:  createdEmployeePublication.PublicationTitle,
		LinkToPublication: createdEmployeePublication.LinkToPublication,
		CreatedAt:         createdEmployeePublication.CreatedAt.Time,
		UpdatedAt:         createdEmployeePublication.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeePublication(ctx context.Context, req *dto.UpdateEmployeePublicationRequest) (*dto.EmployeePublicationResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	updatedEmployeePublicationArgs := sqlc.UpdateEmployeePublicationParams{
		ID:                req.ID,
		PublicationTitle:  req.PublicationTitle,
		LinkToPublication: req.LinkToPublication,
	}
	updatedEmployeePublicationResult, err := s.store.Queries.UpdateEmployeePublication(ctx, updatedEmployeePublicationArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeePublication(service) -> UpdateEmployeePublication(repo) params %v: %w", updatedEmployeePublicationArgs, err))
	}

	return &dto.EmployeePublicationResponse{
		ID:                updatedEmployeePublicationResult.ID,
		PublicationTitle:  updatedEmployeePublicationResult.PublicationTitle,
		LinkToPublication: updatedEmployeePublicationResult.LinkToPublication,
		CreatedAt:         updatedEmployeePublicationResult.CreatedAt.Time,
		UpdatedAt:         updatedEmployeePublicationResult.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteEmployeePublication(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	err := s.store.Queries.DeleteEmployeePublication(ctx, id)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}

		return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeePublication(service) -> DeleteEmployeePublication(repo) params %v: %w", id, err))
	}

	return nil
}

func (s *Service) GetEmployeePublicationByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeePublicationResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeePublicationsByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}
	employeePublications, err := s.store.Queries.GetEmployeePublicationsByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeePublicationByEmployeeIDAndLanguageCode(Service) -> GetEmployeePublicationsByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeePublicationResponse, len(employeePublications))
	for index, publication := range employeePublications {
		result[index] = &dto.EmployeePublicationResponse{
			ID:                publication.ID,
			PublicationTitle:  publication.PublicationTitle,
			LinkToPublication: publication.LinkToPublication,
			CreatedAt:         publication.CreatedAt.Time,
			UpdatedAt:         publication.UpdatedAt.Time,
		}
	}

	return result, nil
}
