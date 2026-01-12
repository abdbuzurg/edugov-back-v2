package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) CreateEmployeeWorkExperience(
	ctx context.Context,
	req *dto.CreateEmployeeWorkExperienceRequest,
) (*dto.EmployeeWorkExperienceResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.CreateEmployeeWorkExperienceParams{
		EmployeeID:   req.EmployeeID,
		LanguageCode: req.LanguageCode,
		Workplace:    req.Workplace,
		Description:  req.Description,
		JobTitle:     req.JobTitle,
		OnGoing:      req.OnGoing,
		DateStart: pgtype.Date{
			Time:  req.DateStart,
			Valid: !req.DateStart.IsZero(),
		},
	}

	if req.DateEnd != nil {
		args.DateEnd = pgtype.Date{Time: *req.DateEnd, Valid: !(*req.DateEnd).IsZero()}
	} else {
		args.DateEnd = pgtype.Date{Valid: false}
	}

	created, err := s.store.Queries.CreateEmployeeWorkExperience(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("CreateEmployeeWorkExperience(service) -> CreateEmployeeWorkExperience(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeWorkExperienceResponse{
		ID:          created.ID,
		Workplace:   created.Workplace,
		Description: created.Description,
		JobTitle:    created.JobTitle,
		DateStart:   created.DateStart.Time,
		DateEnd:     created.DateEnd.Time,
		OnGoing:     created.OnGoing,
		CreatedAt:   created.CreatedAt.Time,
		UpdatedAt:   created.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeeWorkExperience(
	ctx context.Context,
	req *dto.UpdateEmployeeWorkExperienceRequest,
) (*dto.EmployeeWorkExperienceResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	args := sqlc.UpdateEmployeeWorkExperienceParams{
		ID:          req.ID,
		JobTitle:    req.JobTitle,
		Workplace:   req.Workplace,
		Description: req.Description,
		OnGoing:     req.OnGoing,
	}

	// nullable dates (same pattern as refresher_course)
	if req.DateStart != nil {
		args.DateStart = pgtype.Date{Time: *req.DateStart, Valid: !(*req.DateStart).IsZero()}
	} else {
		args.DateStart = pgtype.Date{Valid: false}
	}

	if req.DateEnd != nil {
		args.DateEnd = pgtype.Date{Time: *req.DateEnd, Valid: !(*req.DateEnd).IsZero()}
	} else {
		args.DateEnd = pgtype.Date{Valid: false}
	}

	updated, err := s.store.Queries.UpdateEmployeeWorkExperience(ctx, args)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("UpdateEmployeeWorkExperience(service) -> UpdateEmployeeWorkExperience(repo) params %v: %w", args, err),
		)
	}

	return &dto.EmployeeWorkExperienceResponse{
		ID:          updated.ID,
		Workplace:   updated.Workplace,
		Description: updated.Description,
		JobTitle:    updated.JobTitle,
		DateStart:   updated.DateStart.Time,
		DateEnd:     updated.DateEnd.Time,
		OnGoing:     updated.OnGoing,
		CreatedAt:   updated.CreatedAt.Time,
		UpdatedAt:   updated.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteEmployeeWorkExperience(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	if err := s.store.Queries.DeleteEmployeeWorkExperience(ctx, id); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}
		return apperr.Internal(
			"internal error",
			fmt.Errorf("DeleteEmployeeWorkExperience(service) -> DeleteEmployeeWorkExperience(repo) params %v: %w", id, err),
		)
	}

	return nil
}

func (s *Service) GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode(
	ctx context.Context,
	employeeID int64,
	langCode string,
) ([]*dto.EmployeeWorkExperienceResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}

	rows, err := s.store.Queries.GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal(
			"internal error",
			fmt.Errorf("GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode(service) -> GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err),
		)
	}

	result := make([]*dto.EmployeeWorkExperienceResponse, len(rows))
	for i, r := range rows {
		result[i] = &dto.EmployeeWorkExperienceResponse{
			ID:          r.ID,
			Workplace:   r.Workplace,
			Description: r.Description,
			JobTitle:    r.JobTitle,
			DateStart:   r.DateStart.Time,
			DateEnd:     r.DateEnd.Time,
			OnGoing:     r.OnGoing,
			CreatedAt:   r.CreatedAt.Time,
			UpdatedAt:   r.UpdatedAt.Time,
		}
	}

	return result, nil
}
