package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) mapSQLCEmployeeDegreeModelToDTOEmployeeDegreeResponse(dbResult sqlc.EmployeeDegree) dto.EmployeeDegreeResponse {
	result := dto.EmployeeDegreeResponse{
		ID:                 dbResult.ID,
		RfInstitutionID:    dbResult.RfInstitutionID,
		DegreeLevel:        dbResult.DegreeLevel,
		InstitutionName:    dbResult.InstitutionName,
		Speciality:         dbResult.Speciality,
		DateStart:          dbResult.DateStart.Time,
		DateEnd:            dbResult.DateEnd.Time,
		DateDegreeRecieved: dbResult.DateDegreeRecieved.Time,
		CreatedAt:          dbResult.CreatedAt.Time,
		UpdatedAt:          dbResult.UpdatedAt.Time,
	}

	if dbResult.GivenBy != nil {
		result.GivenBy = *dbResult.GivenBy
	}

	return result
}

// CreateEmployeeDegree validates and stores a new degree.
func (s *Service) CreateEmployeeDegree(ctx context.Context, req *dto.CreateEmployeeDegreeRequest) (*dto.EmployeeDegreeResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	createdEmployeeDegreeArgs := sqlc.CreateEmployeeDegreeParams{
		EmployeeID:         req.EmployeeID,
		LanguageCode:       req.LanguageCode,
		RfInstitutionID:    req.RfInstitutionID,
		DegreeLevel:        req.DegreeLevel,
		InstitutionName:    req.InstitutionName,
		Speciality:         req.Speciality,
		DateStart: pgtype.Date{
			Time:  req.DateStart,
			Valid: !req.DateStart.IsZero(),
		},
		DateEnd: pgtype.Date{
			Time:  req.DateEnd,
			Valid: !req.DateEnd.IsZero(),
		},
		GivenBy: &req.GivenBy,
		DateDegreeRecieved: pgtype.Date{
			Time:  req.DateDegreeRecieved,
			Valid: !req.DateDegreeRecieved.IsZero(),
		},
	}
	createdEmployeeDegreeResult, err := s.store.Queries.CreateEmployeeDegree(ctx, createdEmployeeDegreeArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("CreateEmployeeDegree(service) -> CreateEmployeeDegree(repo) - param %v: %w", createdEmployeeDegreeArgs, err))
	}

	result := s.mapSQLCEmployeeDegreeModelToDTOEmployeeDegreeResponse(createdEmployeeDegreeResult)

	return &result, nil
}

// UpdateEmployeeDegree applies updates to an existing degree.
func (s *Service) UpdateEmployeeDegree(ctx context.Context, req *dto.UpdateEmployeeDegreeRequest) (*dto.EmployeeDegreeResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	var dateStart pgtype.Date
	if req.DateStart != nil && !req.DateStart.IsZero() {
		dateStart = pgtype.Date{Time: *req.DateStart, Valid: true}
	} else {
		dateStart = pgtype.Date{Valid: false}
	}

	var dateEnd pgtype.Date
	if req.DateEnd != nil && !req.DateEnd.IsZero() {
		dateEnd = pgtype.Date{Time: *req.DateEnd, Valid: true}
	} else {
		dateEnd = pgtype.Date{Valid: false}
	}

	var dateRec pgtype.Date
	if req.DateDegreeRecieved != nil && !req.DateDegreeRecieved.IsZero() {
		dateRec = pgtype.Date{Time: *req.DateDegreeRecieved, Valid: true}
	} else {
		dateRec = pgtype.Date{Valid: false}
	}

	updateEmployeeDegreeArgs := sqlc.UpdateEmployeeDegreeParams{
		ID:                 req.ID,
		RfInstitutionID:    req.RfInstitutionID,
		DegreeLevel:        req.DegreeLevel,
		InstitutionName:    req.InstitutionName,
		Speciality:         req.Speciality,
		DateStart:          dateStart,
		DateEnd:            dateEnd,
		GivenBy:            req.GivenBy,
		DateDegreeRecieved: dateRec,
	}
	updatedEmployeeDegreeResult, err := s.store.Queries.UpdateEmployeeDegree(ctx, updateEmployeeDegreeArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeDegree(service) -> UpdateEmployeeDegree(repo) params %v: %w", updateEmployeeDegreeArgs, err))
	}

	result := s.mapSQLCEmployeeDegreeModelToDTOEmployeeDegreeResponse(updatedEmployeeDegreeResult)

	return &result, nil
}

// DeleteEmployeeDegree removes a degree by ID.
func (s *Service) DeleteEmployeeDegree(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	err := s.store.Queries.DeleteEmployeeDegree(ctx, id)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}

		return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeeDetails(service) -> DeleteEmployeeDegree(repo) id %d: %w", id, err))
	}

	return nil
}

// GetEmployeeDegreesByEmployeeIDAndLanguageCode lists degrees for an employee in a language.
func (s *Service) GetEmployeeDegreesByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, languageCode string) ([]*dto.EmployeeDegreeResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid employee id", map[string]string{
			"employeeID": "employeeID must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeDegreesByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: languageCode,
	}
	employeeDegreesResult, err := s.store.Queries.GetEmployeeDegreesByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeDegreesByEmployeeIDAndLanguageCode(service) -> GetEmployeeDegreesByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeeDegreeResponse, len(employeeDegreesResult))
	for index := range result {
		mappedData := s.mapSQLCEmployeeDegreeModelToDTOEmployeeDegreeResponse(employeeDegreesResult[index])
		result[index] = &mappedData
	}

	return result, nil
}
