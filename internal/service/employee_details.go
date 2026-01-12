package service

import (
	"context"
	"database/sql"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/pkg/utils"
	"errors"
	"fmt"
)

const ()

// GetEmployeeDetailsByEmployeeID returns all details records for an employee.
func (s *Service) GetEmployeeDetailsByEmployeeID(ctx context.Context, employeeID int64) ([]*dto.EmployeeDetailsResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid employee_id", map[string]string{
			"id": "must be > 0",
		})
	}

	employeeDetails, err := s.store.Queries.GetEmployeeDetailsByEmployeeID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NotFound("employee details not found")
		}
		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeDetailsByEmployeeID employeeID = %d: %w", employeeID, err))
	}

	result := make([]*dto.EmployeeDetailsResponse, len(employeeDetails))
	for index := range result {
		result[index] = &dto.EmployeeDetailsResponse{
			ID:                   employeeDetails[index].ID,
			LanguageCode:         employeeDetails[index].LanguageCode,
			Surname:              employeeDetails[index].Surname,
			Name:                 employeeDetails[index].Name,
			IsEmployeeDetailsNew: employeeDetails[index].IsEmployeeDetailsNew,
			CreatedAt:            employeeDetails[index].CreatedAt.Time,
			UpdatedAt:            employeeDetails[index].UpdatedAt.Time,
		}

		if employeeDetails[index].Middlename != nil {
			result[index].Middlename = *employeeDetails[index].Middlename
		} else {
			result[index].Middlename = ""
		}

	}

	return result, nil
}

// UpdateEmployeeDetails applies a batch of updates for employee details.
func (s *Service) UpdateEmployeeDetails(ctx context.Context, req []dto.UpdateEmployeeDetailsRequest) ([]*dto.EmployeeDetailsResponse, error) {
	for _, details := range req {
		if err := s.validator.Struct(details); err != nil {
			return nil, apperr.ValidationFromValidator(err)
		}
	}

	resp := []*dto.EmployeeDetailsResponse{}
	err := s.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		newDetails := make([]sqlc.EmployeeDetail, len(req))
		for index, details := range req {
			newDetails[index] = sqlc.EmployeeDetail{
				ID:           details.ID,
				EmployeeID:   details.EmployeeID,
				LanguageCode: details.LanguageCode,
				Middlename:   details.Middlename,
			}

			if details.Name != nil {
				newDetails[index].Name = *details.Name
			}

			if details.Surname != nil {
				newDetails[index].Surname = *details.Surname
			}
		}

		if len(newDetails) == 0 {
			return apperr.Validation("no data for update", map[string]string{
				"data": "data for update must be provided",
			})
		}

		oldDetails, err := q.GetEmployeeDetailsByEmployeeID(ctx, newDetails[0].EmployeeID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeDetails -> GetEmployeeDetailsByEmployeeID employeeID = %d: %w", newDetails[0].EmployeeID, err))
		}

		updatedDetails, createdDetails, removedDetails := utils.CompareSlices(newDetails, oldDetails, func(tmp sqlc.EmployeeDetail) int64 { return tmp.ID })

		for _, details := range updatedDetails {
			updatedDetailsArgs := sqlc.UpdateEmployeeDetailsParams{
				ID:                   details.ID,
				Surname:              details.Surname,
				Name:                 details.Name,
				Middlename:           details.Middlename,
				IsEmployeeDetailsNew: details.IsEmployeeDetailsNew,
			}
			updatedDetailsResult, err := q.UpdateEmployeeDetails(ctx, updatedDetailsArgs)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeDetails(Service) -> UpdateEmployeeDetails(Repo) param %v: %w", updatedDetailsArgs, err))
			}

			temp := &dto.EmployeeDetailsResponse{
				ID:                   updatedDetailsResult.ID,
				LanguageCode:         details.LanguageCode,
				Surname:              details.Surname,
				Name:                 details.Name,
				Middlename:           "",
				IsEmployeeDetailsNew: details.IsEmployeeDetailsNew,
				CreatedAt:            updatedDetailsResult.CreatedAt.Time,
				UpdatedAt:            updatedDetailsResult.UpdatedAt.Time,
			}

			if details.Middlename != nil {
				temp.Middlename = *details.Middlename
			} else {
				temp.Middlename = ""
			}

			resp = append(resp, temp)
		}

		for _, details := range createdDetails {
			createdDetailsArgs := sqlc.CreateEmployeeDetailsParams{
				EmployeeID:           details.EmployeeID,
				LanguageCode:         details.LanguageCode,
				Surname:              details.Surname,
				Name:                 details.Name,
				Middlename:           details.Middlename,
				IsEmployeeDetailsNew: details.IsEmployeeDetailsNew,
			}
			createdDetailsResult, err := q.CreateEmployeeDetails(ctx, createdDetailsArgs)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeDetails(Service) -> CreateEmployeeDetails(Repo) param %v: %w", createdDetailsArgs, err))
			}

			temp := &dto.EmployeeDetailsResponse{
				ID:                   createdDetailsResult.ID,
				LanguageCode:         details.LanguageCode,
				Surname:              details.Surname,
				Name:                 details.Name,
				Middlename:           "",
				IsEmployeeDetailsNew: details.IsEmployeeDetailsNew,
				CreatedAt:            createdDetailsResult.CreatedAt.Time,
				UpdatedAt:            createdDetailsResult.UpdatedAt.Time,
			}

			if details.Middlename != nil {
				temp.Middlename = *details.Middlename
			} else {
				temp.Middlename = ""
			}

			resp = append(resp, temp)
		}

		for _, details := range removedDetails {
			err := q.DeleteEmployeeDetails(ctx, details.ID)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeDetails(Service) -> DeleteEmployeeDetails(Repo) param %v: %w", details.ID, err))
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
