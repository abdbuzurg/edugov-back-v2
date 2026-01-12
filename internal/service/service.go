package service

import (
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	store     *db.Store
	validator *validator.Validate
}

func New(store *db.Store, validator *validator.Validate) *Service {
	return &Service{
		store:     store,
		validator: validator,
	}
}

func (s *Service) asPgError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}

	return nil
}

func (s *Service) pgErrToAppErr(pgErr *pgconn.PgError) error {
	detail := pgErr.Detail
	if detail == "" {
		detail = pgErr.Message
	}

	switch pgErr.Code {
	// conflicts
	case "23505": // unique_violation
		return apperr.Conflict("already exists")

	case "23P01": // exclusion_violation
		return apperr.Conflict("conflict")

	// validation-ish
	case "23503": // foreign_key_violation
		return apperr.Validation("invalid reference", map[string]string{
			"constraint": pgErr.ConstraintName,
			"detail":     detail,
		})

	case "23502": // not_null_violation
		key := pgErr.ColumnName
		if key == "" {
			key = "field"
		}
		return apperr.Validation("missing required field", map[string]string{
			key: "required",
		})

	case "23514": // check_violation
		return apperr.Validation("constraint violation", map[string]string{
			"constraint": pgErr.ConstraintName,
			"detail":     detail,
		})

	// bad input / type issues
	case "22P02", // invalid_text_representation
		"22001", // string_data_right_truncation
		"22003", // numeric_value_out_of_range
		"22007", // invalid_datetime_format
		"22008": // datetime_field_overflow
		return apperr.Validation("invalid value", map[string]string{
			"detail": detail,
		})
	default:
		return apperr.Internal("internal error", pgErr)
	}
}
