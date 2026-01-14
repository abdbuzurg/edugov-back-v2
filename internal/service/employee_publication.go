package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

func (s *Service) mapEmployeePublicationToResponse(publication sqlc.EmployeePublication) *dto.EmployeePublicationResponse {
	return &dto.EmployeePublicationResponse{
		ID:                  publication.ID,
		RfPublicationTypeID: publication.RfPublicationTypeID,
		Name:                publication.Name,
		Type:                publication.Type,
		Authors:             publication.Authors,
		JournalName:         publication.JournalName,
		Volume:              publication.Volume,
		Number:              publication.Number,
		Pages:               publication.Pages,
		Year:                publication.Year,
		Link:                publication.Link,
		CreatedAt:           publication.CreatedAt.Time,
		UpdatedAt:           publication.UpdatedAt.Time,
	}
}

// CreateEmployeePublication creates a publication entry.
func (s *Service) CreateEmployeePublication(ctx context.Context, req *dto.CreateEmployeePublicationRequest) (*dto.EmployeePublicationResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeePublicationArgs := sqlc.CreateEmployeePublicationParams{
		EmployeeID:          req.EmployeeID,
		LanguageCode:        req.LanguageCode,
		RfPublicationTypeID: req.RfPublicationTypeID,
		Name:                req.Name,
		Type:                req.Type,
		Authors:             req.Authors,
		JournalName:         req.JournalName,
		Volume:              req.Volume,
		Number:              req.Number,
		Pages:               req.Pages,
		Year:                req.Year,
		Link:                req.Link,
	}
	createdEmployeePublication, err := s.store.Queries.CreateEmployeePublication(ctx, createdEmployeePublicationArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeePublication(service) -> CreateEmployeePublication(repo) params %v: %w", createdEmployeePublicationArgs, err))
	}

	return s.mapEmployeePublicationToResponse(createdEmployeePublication), nil
}

// UpdateEmployeePublication updates an existing publication.
func (s *Service) UpdateEmployeePublication(ctx context.Context, req *dto.UpdateEmployeePublicationRequest) (*dto.EmployeePublicationResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	updatedEmployeePublicationArgs := sqlc.UpdateEmployeePublicationParams{
		ID:                  req.ID,
		RfPublicationTypeID: req.RfPublicationTypeID,
		Name:                req.Name,
		Type:                req.Type,
		Authors:             req.Authors,
		JournalName:         req.JournalName,
		Volume:              req.Volume,
		Number:              req.Number,
		Pages:               req.Pages,
		Year:                req.Year,
		Link:                req.Link,
	}
	updatedEmployeePublicationResult, err := s.store.Queries.UpdateEmployeePublication(ctx, updatedEmployeePublicationArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeePublication(service) -> UpdateEmployeePublication(repo) params %v: %w", updatedEmployeePublicationArgs, err))
	}

	return s.mapEmployeePublicationToResponse(updatedEmployeePublicationResult), nil
}

// DeleteEmployeePublication removes a publication entry.
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

// GetEmployeePublicationByEmployeeIDAndLanguageCode lists publications for an employee.
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
		result[index] = s.mapEmployeePublicationToResponse(publication)
	}

	return result, nil
}
