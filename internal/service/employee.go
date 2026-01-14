package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// GetEmployeeByUniqueID returns a full employee profile by unique ID and language.
func (s *Service) GetEmployeeByUniqueID(ctx context.Context, uniqueID string, langCode string) (*dto.EmployeeResponse, error) {
	if uniqueID == "" {
		return nil, apperr.Validation("invalid uniqueID", map[string]string{
			"uniqueID": "unique_id must not be empty",
		})
	}

	var result *dto.EmployeeResponse
	err := s.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		employee, err := q.GetEmployeeByUniqueIdentifier(ctx, uniqueID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}

			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeByUniqueIdentifier params %v: %w", uniqueID, err))
		}

		result = &dto.EmployeeResponse{
			ID:        employee.ID,
			UniqueID:  employee.UniqueID,
			CreatedAt: employee.CreatedAt.Time,
			UpdatedAt: employee.UpdatedAt.Time,
		}

		if employee.Gender != nil {
			result.Gender = *employee.Gender
		}

		// Degree
		employeeDetails, err := q.GetEmployeeDetailsByEmployeeID(ctx, employee.ID)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeDetailsByEmployeeID params %v: %w", employee.ID, err))
		}
		for _, ed := range employeeDetails {
			detail := &dto.EmployeeDetailsResponse{
				ID:                   ed.ID,
				LanguageCode:         ed.LanguageCode,
				Surname:              ed.Surname,
				Name:                 ed.Name,
				IsEmployeeDetailsNew: ed.IsEmployeeDetailsNew,
				CreatedAt:            ed.CreatedAt.Time,
				UpdatedAt:            ed.UpdatedAt.Time,
			}
			if ed.Middlename != nil {
				detail.Middlename = *ed.Middlename
			}
			result.Details = append(result.Details, detail)
		}

		employeeDegreeArgs := sqlc.GetEmployeeDegreesByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeeDegrees, err := q.GetEmployeeDegreesByEmployeeIDAndLanguageCode(ctx, employeeDegreeArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeDegreesByEmployeeIDAndLanguageCode params %v: %w", employeeDegreeArgs, err))
		}
		for _, ed := range employeeDegrees {
			result.Degrees = append(result.Degrees, &dto.EmployeeDegreeResponse{
				ID:                 ed.ID,
				RfInstitutionID:    ed.RfInstitutionID,
				DegreeLevel:        ed.DegreeLevel,
				InstitutionName:    ed.InstitutionName,
				Speciality:         ed.Speciality,
				DateStart:          ed.DateStart.Time,
				DateEnd:            ed.DateEnd.Time,
				GivenBy:            *ed.GivenBy,
				DateDegreeRecieved: ed.DateDegreeRecieved.Time,
				CreatedAt:          ed.CreatedAt.Time,
				UpdatedAt:          ed.UpdatedAt.Time,
			})
		}

		// Work Experience
		employeeWorkExperiencesArgs := sqlc.GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeeWorkExperiences, err := q.GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode(ctx, employeeWorkExperiencesArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode params %v: %w", employeeWorkExperiencesArgs, err))
		}
		for _, ew := range employeeWorkExperiences {
			result.WorkExperiences = append(result.WorkExperiences, &dto.EmployeeWorkExperienceResponse{
				ID:          ew.ID,
				Workplace:   ew.Workplace,
				Description: ew.Description,
				JobTitle:    ew.JobTitle,
				DateStart:   ew.DateStart.Time,
				DateEnd:     ew.DateEnd.Time,
				OnGoing:     ew.OnGoing,
				CreatedAt:   ew.CreatedAt.Time,
				UpdatedAt:   ew.UpdatedAt.Time,
			})
		}

		// Main Research Area
		employeeMRAsArgs := sqlc.GetEmployeeMainResearchAreasByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeeMRAs, err := q.GetEmployeeMainResearchAreasByEmployeeIDAndLanguageCode(ctx, employeeMRAsArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeMainResearchAreasByEmployeeIDAndLanguageCode params %v: %w", employeeMRAsArgs, err))
		}
		for _, mra := range employeeMRAs {
			mainResearchAreaKeyTopics, err := q.GetEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaIDAndLanguageCode(ctx, mra.ID)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeMainResearchAreaKeyTopicsByEmployeeMainResearchAreaIDAndLanguageCode params %v: %w", mra.ID, err))
			}
			mainResearchArea := &dto.EmployeeMainResearchAreaResponse{
				ID:         mra.ID,
				Discipline: mra.Discipline,
				Area:       mra.Area,
				CreatedAt:  mra.CreatedAt.Time,
				UpdatedAt:  mra.UpdatedAt.Time,
			}
			for _, kt := range mainResearchAreaKeyTopics {
				mainResearchArea.KeyTopics = append(mainResearchArea.KeyTopics, &dto.ResearchAreaKeyTopicResponse{
					ID:            kt.ID,
					KeyTopicTitle: kt.KeyTopicTitle,
					CreatedAt:     kt.CreatedAt.Time,
					UpdatedAt:     kt.UpdatedAt.Time,
				})
			}

			result.MainResearchAreas = append(result.MainResearchAreas, mainResearchArea)
		}

		// Publications
		employeePublicationsArgs := sqlc.GetEmployeePublicationsByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeePublications, err := q.GetEmployeePublicationsByEmployeeIDAndLanguageCode(ctx, employeePublicationsArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeePublicationsByEmployeeIDAndLanguageCode params %v: %w", employeePublicationsArgs, err))
		}
		for _, ep := range employeePublications {
			result.Publications = append(result.Publications, s.mapEmployeePublicationToResponse(ep))
		}

		// Scientific Award
		employeeScientificAwardsArgs := sqlc.GetEmployeeScientificAwardsByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeeScientificAwards, err := q.GetEmployeeScientificAwardsByEmployeeIDAndLanguageCode(ctx, employeeScientificAwardsArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeScientificAwardsByEmployeeIDAndLanguageCode params %v: %w", employeeScientificAwardsArgs, err))
		}
		for _, esa := range employeeScientificAwards {
			result.ScientificAwards = append(result.ScientificAwards, &dto.EmployeeScientificAwardResponse{
				ID:                   esa.ID,
				ScientificAwardTitle: esa.ScientificAwardTitle,
				GivenBy:              esa.GivenBy,
				CreatedAt:            esa.CreatedAt.Time,
				UpdatedAt:            esa.UpdatedAt.Time,
			})
		}

		// Patent
		employeePatentsArgs := sqlc.GetEmployeePatentsByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeePatents, err := q.GetEmployeePatentsByEmployeeIDAndLanguageCode(ctx, employeePatentsArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeePatentsByEmployeeIDAndLanguageCode params %v: %w", employeePatentsArgs, err))
		}
		for _, ep := range employeePatents {
			result.Patents = append(result.Patents, &dto.EmployeePatentResponse{
				ID:          ep.ID,
				PatentTitle: ep.PatentTitle,
				Description: ep.Description,
				CreatedAt:   ep.CreatedAt.Time,
				UpdatedAt:   ep.UpdatedAt.Time,
			})
		}

		// Participation In Events
		employeePIPCsArgs := sqlc.GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeePIPCs, err := q.GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCode(ctx, employeePIPCsArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeParticipationInProfessionalCommunitysByEmployeeIDAndLanguageCode params %v: %w", employeePIPCsArgs, err))
		}
		for _, pipc := range employeePIPCs {
			result.ParticipationInProfessionalCommunities = append(result.ParticipationInProfessionalCommunities, &dto.EmployeeParticipationInProfessionalCommunityResponse{
				ID:                          pipc.ID,
				ProfessionalCommunityTitle:  pipc.ProfessionalCommunityTitle,
				RoleInProfessionalCommunity: pipc.RoleInProfessionalCommunity,
				CreatedAt:                   pipc.CreatedAt.Time,
				UpdatedAt:                   pipc.UpdatedAt.Time,
			})
		}

		// Refresher Course
		employeeRCArgs := sqlc.GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeeRCs, err := q.GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode(ctx, employeeRCArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode params %v: %w", employeeRCArgs, err))
		}
		for _, rc := range employeeRCs {
			result.RefresherCourses = append(result.RefresherCourses, &dto.EmployeeRefresherCourseResponse{
				ID:          rc.ID,
				CourseTitle: rc.CourseTitle,
				DateStart:   rc.DateStart.Time,
				DateEnd:     rc.DateEnd.Time,
				CreatedAt:   rc.CreatedAt.Time,
				UpdatedAt:   rc.UpdatedAt.Time,
			})
		}

		// Participation in Events
		employeePIEArgs := sqlc.GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeePIEs, err := q.GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode(ctx, employeePIEArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode params %v: %w", employeePIEArgs, err))
		}
		for _, pie := range employeePIEs {
			result.ParticipationInEvents = append(result.ParticipationInEvents, &dto.EmployeeParticipationInEventResponse{
				ID:         pie.ID,
				EventTitle: pie.EventTitle,
				EventDate:  pie.EventDate.Time,
				CreatedAt:  pie.CreatedAt.Time,
				UpdatedAt:  pie.UpdatedAt.Time,
			})
		}

		// Research Activity
		employeeRAsArgs := sqlc.GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCodeParams{
			EmployeeID:   employee.ID,
			LanguageCode: langCode,
		}
		employeeRAs, err := q.GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode(ctx, employeeRAsArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode params %v: %w", employeeRAsArgs, err))
		}
		for _, ra := range employeeRAs {
			result.ResearchActivities = append(result.ResearchActivities, &dto.EmployeeResearchActivityResponse{
				ID:                    ra.ID,
				ResearchActivityTitle: ra.ResearchActivityTitle,
				EmployeeRole:          ra.EmployeeRole,
				CreatedAt:             ra.CreatedAt.Time,
				UpdatedAt:             ra.UpdatedAt.Time,
			})
		}

		// Socials
		employeeSocials, err := q.GetEmployeeSocialsByEmployeeID(ctx, employee.ID)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("GetEmployeeByUniqueID -> GetEmployeeSocialsByEmployeeID params %v: %w", employee.ID, err))
		}
		for _, social := range employeeSocials {
			result.Socials = append(result.Socials, &dto.EmployeeSocialResponse{
				ID:           social.ID,
				SocialName:   social.SocialName,
				LinkToSocial: social.LinkToSocial,
				CreatedAt:    social.CreatedAt.Time,
				UpdatedAt:    social.UpdatedAt.Time,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
