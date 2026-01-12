package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/pkg/utils"
	"fmt"
)

// func (s *Service) mapSQLCEmployeeMainResearchAreModelToDTOEmployeeMainResearchAreaResponse(dbResult sqlc.EmployeeMain) dto.EmployeeDegreeResponse {
// 	result := dto.EmployeeDegreeResponse{
// 		ID:                 dbResult.ID,
// 		DegreeLevel:        dbResult.DegreeLevel,
// 		UniversityName:     dbResult.UniversityName,
// 		Speciality:         dbResult.Speciality,
// 		DateStart:          dbResult.DateStart.Time,
// 		DateEnd:            dbResult.DateEnd.Time,
// 		DateDegreeRecieved: dbResult.DateDegreeRecieved.Time,
// 		CreatedAt:          dbResult.CreatedAt.Time,
// 		UpdatedAt:          dbResult.UpdatedAt.Time,
// 		GivenBy:            *dbResult.GivenBy,
// 	}

// 	return result
// }

func (s *Service) CreateEmployeeMainResearchArea(ctx context.Context, req *dto.CreateEmployeeMainResearchAreaRequest) (*dto.EmployeeMainResearchAreaResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	var result *dto.EmployeeMainResearchAreaResponse
	err := s.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		createdEmployeeMRAArgs := sqlc.CreateEmployeeMainResearchAreaParams{
			EmployeeID:   req.EmployeeID,
			LanguageCode: req.LanguageCode,
			Area:         req.Area,
			Discipline:   req.Discipline,
		}
		createdEmployeeMRAResult, err := q.CreateEmployeeMainResearchArea(ctx, createdEmployeeMRAArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("CreateEmployeeMainResearchArea(service) -> CreateEmployeeMainResearchArea(repo) params %v: %w", createdEmployeeMRAArgs, err))
		}

		result = &dto.EmployeeMainResearchAreaResponse{
			ID:         createdEmployeeMRAResult.ID,
			Discipline: createdEmployeeMRAArgs.Discipline,
			Area:       createdEmployeeMRAArgs.Area,
			CreatedAt:  createdEmployeeMRAResult.CreatedAt.Time,
			UpdatedAt:  createdEmployeeMRAResult.UpdatedAt.Time,
		}

		result.KeyTopics = make([]*dto.ResearchAreaKeyTopicResponse, len(req.KeyTopics))
		for index, rakt := range req.KeyTopics {
			createdMainResearchAreaKTArgs := sqlc.CreateEmployeeMainResearchAreaKeyTopicParams{
				EmployeeMainResearchAreaID: result.ID,
				KeyTopicTitle:              rakt.KeyTopicTitle,
			}
			createdMainResearchAreaKT, err := q.CreateEmployeeMainResearchAreaKeyTopic(ctx, createdMainResearchAreaKTArgs)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("CreateEmployeeMainResearchArea (service) -> CreateEmployeeMainResearchAreaKeyTopic(repo) - params %v: %w", createdMainResearchAreaKTArgs, err))
			}

			result.KeyTopics[index] = &dto.ResearchAreaKeyTopicResponse{
				ID:            createdMainResearchAreaKT.ID,
				KeyTopicTitle: createdMainResearchAreaKT.KeyTopicTitle,
				CreatedAt:     createdMainResearchAreaKT.CreatedAt.Time,
				UpdatedAt:     createdMainResearchAreaKT.UpdatedAt.Time,
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) UpdateEmployeeMainResearchArea(ctx context.Context, req *dto.UpdateEmployeeMainResearchAreaRequest) (*dto.EmployeeMainResearchAreaResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	var result *dto.EmployeeMainResearchAreaResponse
	err := s.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		updatedEmployeeMRAArgs := sqlc.UpdateEmployeeMainResearchAreaParams{
			ID:         req.ID,
			Area:       req.Area,
			Discipline: req.Discipline,
		}
		updatedEmployeeMRAResult, err := q.UpdateEmployeeMainResearchArea(ctx, updatedEmployeeMRAArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeMainResearchArea(service) -> UpdateEmployeeMainResearchArea(repo) params %v: %w", updatedEmployeeMRAArgs, err))
		}

		result = &dto.EmployeeMainResearchAreaResponse{
			ID:         updatedEmployeeMRAResult.ID,
			Discipline: updatedEmployeeMRAResult.Discipline,
			Area:       updatedEmployeeMRAResult.Area,
			CreatedAt:  updatedEmployeeMRAResult.CreatedAt.Time,
			UpdatedAt:  updatedEmployeeMRAResult.UpdatedAt.Time,
		}

		newResearchAreaKT := make([]sqlc.EmployeeMainResearchAreaKeyTopic, len(req.KeyTopics))
		for index, kt := range req.KeyTopics {
			newResearchAreaKT[index] = sqlc.EmployeeMainResearchAreaKeyTopic{
				ID:                         kt.ID,
				EmployeeMainResearchAreaID: result.ID,
				KeyTopicTitle:              *kt.KeyTopicTitle,
			}
		}

		oldResaerchAreaKT, err := q.GetEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaIDAndLanguageCode(ctx, result.ID)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeMainResearchArea(service) -> GetEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaIDAndLanguageCode(repo) params %d: %w", result.ID, err))
		}

		updatedRAKTs, newRAKTs, deletedRAKTs := utils.CompareSlices(oldResaerchAreaKT, newResearchAreaKT, func(entity sqlc.EmployeeMainResearchAreaKeyTopic) int64 { return entity.ID })
		for _, kt := range updatedRAKTs {
			updatedRAKTArgs := sqlc.UpdateEmployeeMainResearchAreaKeyTopicParams{
				ID:            kt.ID,
				KeyTopicTitle: &kt.KeyTopicTitle,
			}
			updatedRAKTResult, err := q.UpdateEmployeeMainResearchAreaKeyTopic(ctx, updatedRAKTArgs)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeMainResearchArea(service) -> UpdateEmployeeMainResearchAreaKeyTopic params %v: %w", updatedRAKTArgs, err))
			}

			result.KeyTopics = append(result.KeyTopics, &dto.ResearchAreaKeyTopicResponse{
				ID:            updatedRAKTResult.ID,
				KeyTopicTitle: updatedRAKTResult.KeyTopicTitle,
				CreatedAt:     updatedRAKTResult.CreatedAt.Time,
				UpdatedAt:     updatedRAKTResult.UpdatedAt.Time,
			})
		}

		for _, kt := range newRAKTs {
			createdRAKTArgs := sqlc.CreateEmployeeMainResearchAreaKeyTopicParams{
				EmployeeMainResearchAreaID: result.ID,
				KeyTopicTitle:              kt.KeyTopicTitle,
			}
			createdRAKTResult, err := q.CreateEmployeeMainResearchAreaKeyTopic(ctx, createdRAKTArgs)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("UpdateEmployeeMainResearchArea(service) -> CreateEmployeeMainResearchAreaKeyTopic params %v: %w", createdRAKTArgs, err))
			}

			result.KeyTopics = append(result.KeyTopics, &dto.ResearchAreaKeyTopicResponse{
				ID:            createdRAKTResult.ID,
				KeyTopicTitle: createdRAKTResult.KeyTopicTitle,
				CreatedAt:     createdRAKTResult.CreatedAt.Time,
				UpdatedAt:     createdRAKTResult.UpdatedAt.Time,
			})
		}

		for _, kt := range deletedRAKTs {
			err := q.DeleteEmployeeMainResearchAreaKeyTopic(ctx, kt.ID)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internla error", fmt.Errorf("UpdateEmployeeMainResearchArea(service) -> DeleteEmployeeMainResearchAreaKeyTopic (repo) params %v: %w", kt.ID, err))
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) DeleteEmployeeMainResearchArea(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperr.Validation("invalid id", map[string]string{
			"id": "id must be > 0",
		})
	}

	return s.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		if err := q.DeleteEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaID(ctx, id); err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeeMainResearchArea (service) -> DeleteEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaID (repo) params %d: %w", id, err))
		}

		if err := q.DeleteEmployeeMainResearchArea(ctx, id); err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("DeleteEmployeeMainResearchArea (service) -> DeleteEmployeeMainResearchArea (repo) params %d: %w", id, err))
		}

		return nil
	})
}

func (s *Service) GetEmployeeMainResearchAreaByEmployeeIDAndLanguageCode(ctx context.Context, employeeID int64, languageCode string) ([]*dto.EmployeeMainResearchAreaResponse, error) {
	if employeeID <= 0 {
		return nil, apperr.Validation("invalid employee id", map[string]string{
			"employeeID": "employeeID must be > 0",
		})
	}

	lookupArgs := sqlc.GetEmployeeMainResearchAreasByEmployeeIDAndLanguageCodeParams{
		EmployeeID:   employeeID,
		LanguageCode: languageCode,
	}
	employeeMRAsResult, err := s.store.Queries.GetEmployeeMainResearchAreasByEmployeeIDAndLanguageCode(ctx, lookupArgs)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeMainResearchAreaByEmployeeIDAndLanguageCode(service) -> GetEmployeeMainResearchAreasByEmployeeIDAndLanguageCode(repo) params %v: %w", lookupArgs, err))
	}

	result := make([]*dto.EmployeeMainResearchAreaResponse, len(employeeMRAsResult))
	for index := range result {
		result[index] = &dto.EmployeeMainResearchAreaResponse{
			ID:         employeeMRAsResult[index].ID,
			Discipline: employeeMRAsResult[index].Discipline,
			Area:       employeeMRAsResult[index].Area,
			CreatedAt:  employeeMRAsResult[index].CreatedAt.Time,
			UpdatedAt:  employeeMRAsResult[index].UpdatedAt.Time,
		}

		mainResearchAreaKTs, err := s.store.Queries.GetEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaIDAndLanguageCode(ctx, result[index].ID)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return nil, s.pgErrToAppErr(pgErr)
			}

			return nil, apperr.Internal("internal error", fmt.Errorf("GetEmployeeMainResearchAreaByEmployeeIDAndLanguageCode(service) -> GetEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaIDAndLanguageCode(repo) params %v: %w", result[index], err))
		}

		for _, kt := range mainResearchAreaKTs {
			result[index].KeyTopics = append(result[index].KeyTopics, &dto.ResearchAreaKeyTopicResponse{
				ID:            kt.ID,
				KeyTopicTitle: kt.KeyTopicTitle,
				CreatedAt:     kt.CreatedAt.Time,
				UpdatedAt:     kt.UpdatedAt.Time,
			})
		}
	}

	return result, nil
}
