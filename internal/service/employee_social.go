package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"
)

// CreateEmployeeSocial creates a social link entry.
func (s *Service) CreateEmployeeSocial(
	ctx context.Context,
	req *dto.CreateEmployeeSocialRequest,
) (*dto.EmployeeSocialResponse, error) {

	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.CreateEmployeeSocialParams{
		EmployeeID:   req.EmployeeID,
		SocialName:   req.SocialName,
		LinkToSocial: req.LinkToSocial,
	}

	created, err := s.store.Queries.CreateEmployeeSocial(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("CreateEmployeeSocial(service) -> CreateEmployeeSocial(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeSocialResponse{
		ID:           created.ID,
		SocialName:   created.SocialName,
		LinkToSocial: created.LinkToSocial,
		CreatedAt:    created.CreatedAt.Time,
		UpdatedAt:    created.UpdatedAt.Time,
	}, nil
}

// UpdateEmployeeSocial updates an existing social link.
func (s *Service) UpdateEmployeeSocial(
	ctx context.Context,
	req *dto.UpdateEmployeeSocialRequest,
) (*dto.EmployeeSocialResponse, error) {

	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.UpdateEmployeeSocialParams{
		ID:           req.ID,
		SocialName:   req.SocialName,
		LinkToSocial: req.LinkToSocial,
	}

	updated, err := s.store.Queries.UpdateEmployeeSocial(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("UpdateEmployeeSocial(service) -> UpdateEmployeeSocial(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeSocialResponse{
		ID:           updated.ID,
		SocialName:   updated.SocialName,
		LinkToSocial: updated.LinkToSocial,
		CreatedAt:    updated.CreatedAt.Time,
		UpdatedAt:    updated.UpdatedAt.Time,
	}, nil
}

// DeleteEmployeeSocial removes a social link entry.
func (s *Service) DeleteEmployeeSocial(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	if err := s.store.Queries.DeleteEmployeeSocial(ctx, id); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}
		return apperr.Internal(
			"internal error",
			fmt.Errorf("DeleteEmployeeSocial(service) -> DeleteEmployeeSocial(repo) params %v: %w", id, err),
		)
	}

	return nil
}

// GetEmployeeSocialByEmployeeID lists social links for an employee.
func (s *Service) GetEmployeeSocialByEmployeeID(
	ctx context.Context,
	employeeID int64,
) ([]*dto.EmployeeSocialResponse, error) {

	if employeeID <= 0 {
		return nil, apperr.Validation("invalid employeeID", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	rows, err := s.store.Queries.GetEmployeeSocialsByEmployeeID(ctx, employeeID)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("GetByEmployeeID(service) -> GetEmployeeSocialsByEmployeeID(repo) params %v: %w", employeeID, err),
		)
	}

	result := make([]*dto.EmployeeSocialResponse, len(rows))
	for i, r := range rows {
		result[i] = &dto.EmployeeSocialResponse{
			ID:           r.ID,
			SocialName:   r.SocialName,
			LinkToSocial: r.LinkToSocial,
			CreatedAt:    r.CreatedAt.Time,
			UpdatedAt:    r.UpdatedAt.Time,
		}
	}

	return result, nil
}
