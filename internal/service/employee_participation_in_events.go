package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateEmployeeParticipationInEvent creates an event participation entry.
func (s *Service) CreateEmployeeParticipationInEvent(ctx context.Context, req *dto.CreateEmployeeParticipationInEventRequest) (*dto.EmployeeParticipationInEventResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeePIEArgs := sqlc.CreateEmployeeParticipationInEventParams{
		EmployeeID:   req.EmployeeID,
		LanguageCode: req.LanguageCode,
		EventTitle:   req.EventTitle,
		EventDate: pgtype.Date{
			Time:  req.EventDate,
			Valid: !req.EventDate.IsZero(),
		},
	}
	createdEmployeePIEResult, err := s.store.Queries.CreateEmployeeParticipationInEvent(ctx, createdEmployeePIEArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeeParticipationInEvent(service) -> CreateEmployeeParticipationInEvent(repo) params %v: %w", createdEmployeePIEArgs, err))
	}

	return &dto.EmployeeParticipationInEventResponse{
		ID:         createdEmployeePIEResult.ID,
		EventTitle: createdEmployeePIEResult.EventTitle,
		EventDate:  createdEmployeePIEResult.EventDate.Time,
		CreatedAt:  createdEmployeePIEResult.CreatedAt.Time,
		UpdatedAt:  createdEmployeePIEResult.UpdatedAt.Time,
	}, nil
}

// UpdateEmployeeParticipationInEvent updates an event participation entry.
func (s *Service) UpdateEmployeeParticipationInEvent(ctx context.Context, req *dto.UpdateEmployeeParticipationInEventRequest) (*dto.EmployeeParticipationInEventResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	updatedEmployeePIEArgs := sqlc.UpdateEmployeeParticipationInEventParams{
		ID:         req.ID,
		EventTitle: req.EventTitle,
	}
	if req.EventDate != nil {
		updatedEmployeePIEArgs.EventDate = pgtype.Date{Time: *req.EventDate, Valid: !(*req.EventDate).IsZero()}
	} else {
		updatedEmployeePIEArgs.EventDate = pgtype.Date{Valid: false}
	}

	updatedEmployeePIEResult, err := s.store.Queries.UpdateEmployeeParticipationInEvent(ctx, updatedEmployeePIEArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeParticipationInEvent(service) -> UpdateEmployeeParticipationInEvent(repo) params %v: %w", updatedEmployeePIEArgs, err))
	}

	return &dto.EmployeeParticipationInEventResponse{
		ID:         updatedEmployeePIEResult.ID,
		EventTitle: updatedEmployeePIEResult.EventTitle,
		EventDate:  updatedEmployeePIEResult.EventDate.Time,
		CreatedAt:  updatedEmployeePIEResult.CreatedAt.Time,
		UpdatedAt:  updatedEmployeePIEResult.UpdatedAt.Time,
	}, nil
}

// DeleteEmployeeParticipationInEvent removes an event participation entry.
func (s *Service) DeleteEmployeeParticipationInEvent(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	if err := s.store.Queries.DeleteEmployeeParticipationInEvent(ctx, id); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}

		return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeeParticipationInEvent (service) -> DeleteEmployeeParticipationInEvent (repo) param %d: %w", id, err))
	}

	return nil
}

// GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode lists event participations for an employee.
func (s *Service) GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeeParticipationInEventResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid employeeID", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}
	employeePIEs, err := s.store.Queries.GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode (service) -> GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode (repo) param %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeeParticipationInEventResponse, len(employeePIEs))
	for index, pie := range employeePIEs {
		result[index] = &dto.EmployeeParticipationInEventResponse{
			ID:         pie.ID,
			EventTitle: pie.EventTitle,
			EventDate:  pie.EventDate.Time,
			CreatedAt:  pie.CreatedAt.Time,
			UpdatedAt:  pie.UpdatedAt.Time,
		}
	}

	return result, nil
}
