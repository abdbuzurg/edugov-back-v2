package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) CreateEmployeeRefresherCourse(ctx context.Context, req *dto.CreateEmployeeRefresherCourseRequest) (*dto.EmployeeRefresherCourseResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeeRefresherCourseArgs := sqlc.CreateEmployeeRefresherCourseParams{
		EmployeeID:   req.EmployeeID,
		LanguageCode: req.LanguageCode,
		CourseTitle:  req.CourseTitle,
		DateStart: pgtype.Date{
			Time:  req.DateStart,
			Valid: !req.DateStart.IsZero(),
		},
		DateEnd: pgtype.Date{
			Time:  req.DateEnd,
			Valid: !req.DateEnd.IsZero(),
		},
	}
	createdEmployeeRefresherCourseResult, err := s.store.Queries.CreateEmployeeRefresherCourse(ctx, createdEmployeeRefresherCourseArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeeRefresherCourse(service) -> CreateEmployeeRefresherCourse(repo) params %v: %w", createdEmployeeRefresherCourseArgs, err))
	}

	return &dto.EmployeeRefresherCourseResponse{
		ID:          createdEmployeeRefresherCourseResult.ID,
		CourseTitle: createdEmployeeRefresherCourseResult.CourseTitle,
		DateStart:   createdEmployeeRefresherCourseResult.DateStart.Time,
		DateEnd:     createdEmployeeRefresherCourseResult.DateEnd.Time,
		CreatedAt:   createdEmployeeRefresherCourseResult.CreatedAt.Time,
		UpdatedAt:   createdEmployeeRefresherCourseResult.UpdatedAt.Time,
	}, nil
}

func (s *Service) UpdateEmployeeRefresherCourse(ctx context.Context, req *dto.UpdateEmployeeRefresherCourseRequest) (*dto.EmployeeRefresherCourseResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	updatedEmployeeRefresherCourseArgs := sqlc.UpdateEmployeeRefresherCourseParams{
		ID:          req.ID,
		CourseTitle: req.CourseTitle,
	}

	if req.DateStart != nil {
		updatedEmployeeRefresherCourseArgs.DateStart = pgtype.Date{Time: *req.DateStart, Valid: !(*req.DateStart).IsZero()}
	} else {
		updatedEmployeeRefresherCourseArgs.DateStart = pgtype.Date{Valid: false}
	}

	if req.DateEnd != nil {
		updatedEmployeeRefresherCourseArgs.DateEnd = pgtype.Date{Time: *req.DateEnd, Valid: !(*req.DateEnd).IsZero()}
	} else {
		updatedEmployeeRefresherCourseArgs.DateEnd = pgtype.Date{Valid: false}
	}

	updatedEmployeeRefresherCourseResult, err := s.store.Queries.UpdateEmployeeRefresherCourse(ctx, updatedEmployeeRefresherCourseArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeRefresherCourse(service) -> UpdateEmployeeRefresherCourse(repo) params %v: %w", updatedEmployeeRefresherCourseArgs, err))
	}

	return &dto.EmployeeRefresherCourseResponse{
		ID:          updatedEmployeeRefresherCourseResult.ID,
		CourseTitle: updatedEmployeeRefresherCourseResult.CourseTitle,
		DateStart:   updatedEmployeeRefresherCourseResult.DateStart.Time,
		DateEnd:     updatedEmployeeRefresherCourseResult.DateEnd.Time,
		CreatedAt:   updatedEmployeeRefresherCourseResult.CreatedAt.Time,
		UpdatedAt:   updatedEmployeeRefresherCourseResult.UpdatedAt.Time,
	}, nil
}

func (s *Service) DeleteEmployeeRefresherCourse(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	err := s.store.Queries.DeleteEmployeeRefresherCourse(ctx, id)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}

		return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeeRefresherCourse(service) -> DeleteEmployeeRefresherCourse(repo) params %v: %w", id, err))
	}

	return nil
}

func (s *Service) GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, langCode string) ([]*dto.EmployeeRefresherCourseResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid id", map[string]string{
			"employeeID": "id must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: langCode,
	}
	employeeRCs, err := s.store.Queries.GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode(service) -> GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeeRefresherCourseResponse, len(employeeRCs))
	for index, rc := range employeeRCs {
		result[index] = &dto.EmployeeRefresherCourseResponse{
			ID:          rc.ID,
			CourseTitle: rc.CourseTitle,
			DateStart:   rc.DateStart.Time,
			DateEnd:     rc.DateEnd.Time,
			CreatedAt:   rc.CreatedAt.Time,
			UpdatedAt:   rc.UpdatedAt.Time,
		}
	}

	return result, nil
}
