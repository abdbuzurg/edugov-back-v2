package service

import (
	"context"
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/db/sqlc"
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/pkg/utils"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a user and linked employee record.
func (s *Service) Register(ctx context.Context, req *dto.RegisterRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return apperr.ValidationFromValidator(err)
	}

	return s.store.ExecTx(ctx, func(q *sqlc.Queries) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return apperr.Internal("internal error", fmt.Errorf("Register(service) -> GenerateFromPassword params %v: %w", req.Password, err))
		}

		createUserArgs := sqlc.CreateUserParams{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
		}
		createdUser, err := q.CreateUser(ctx, createUserArgs)
		if err != nil {
			if pgErr := s.asPgError(err); pgErr != nil {
				return s.pgErrToAppErr(pgErr)
			}

			return apperr.Internal("internal error", fmt.Errorf("Register(service) -> CreateUser params %v: %w", createUserArgs, err))
		}

		isUniqueExists := true
		for isUniqueExists {
			uniqueID, err := utils.GenerateNumericUniqueID()
			if err != nil {
				return apperr.Internal("internal error", fmt.Errorf("Register(service) -> GenerateNumericUniqueID: %w", err))
			}

			createEmployeeArgs := sqlc.CreateEmployeeParams{
				UniqueID: uniqueID,
				UserID:   &createdUser.ID,
				Gender:   &req.Gender,
				Tin:      &req.Tin,
			}
			createdEmployee, err := q.CreateEmployee(ctx, createEmployeeArgs)
			if err != nil {
				if pgErr := s.asPgError(err); pgErr != nil {
					return s.pgErrToAppErr(pgErr)
				}

				return apperr.Internal("internal error", fmt.Errorf("Register(service) -> CreateEmployee(repo) params %v: %w", createEmployeeArgs, err))
			}

			if len(createdEmployee) == 1 {
				isUniqueExists = false
			}
		}

		return nil
	})
}

// Login verifies credentials, issues tokens, and stores a refresh session.
func (s *Service) Login(ctx context.Context, req *dto.AuthRequest) (*dto.AuthResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, apperr.ValidationFromValidator(err)
	}

	user, err := s.store.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		pgErr := s.asPgError(err)
		if pgErr != nil && !errors.Is(pgErr, pgx.ErrNoRows) {
			return nil, s.pgErrToAppErr(pgErr)
		}
		if errors.Is(pgErr, pgx.ErrNoRows) {
			return nil, apperr.Validation("invalid credentials", map[string]string{
				"email": "no user with this email",
			})
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("Login(service) -> GetUserByEmail(repo) params %v: %w", req.Email, err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperr.Validation("invalid credentials", map[string]string{
			"password": "password is invalid",
		})
	}

	employee, err := s.store.Queries.GetEmployeeByUserID(ctx, &user.ID)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Validation("invalid credentials", map[string]string{
			"user": "user not linked to employee",
		})
	}

	signingSecret := viper.GetString("jwt.secret")
	now := time.Now()
	accessTTL := viper.GetDuration("jwt.access_ttl")
	if accessTTL == 0 {
		accessTTL = 2 * time.Hour
	}
	refreshTTL := viper.GetDuration("jwt.refresh_ttl")
	if refreshTTL == 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	accessExp := now.Add(accessTTL)
	refreshExp := now.Add(refreshTTL)

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(accessExp),
		IssuedAt:  jwt.NewNumericDate(now),
	}).SignedString([]byte(signingSecret))
	if err != nil {
		return nil, apperr.Internal("internal error", fmt.Errorf("Login(service) -> sign access token: %w", err))
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		IssuedAt:  jwt.NewNumericDate(now),
	}).SignedString([]byte(signingSecret))
	if err != nil {
		return nil, apperr.Internal("internal error", fmt.Errorf("Login(service) -> sign refresh token: %w", err))
	}

	if err := s.store.Queries.DeleteUserSessionByUserID(ctx, user.ID); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("Login(service) -> DeleteUserSessionByUserID(repo): %w", err))
	}

	_, err = s.store.Queries.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshExp, Valid: true},
	})
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}

		return nil, apperr.Internal("internal error", fmt.Errorf("Login(service) -> CreateUserSession(repo): %w", err))
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       employee.UniqueID,
	}, nil
}

// RefreshToken validates the refresh token and rotates tokens.
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	if refreshToken == "" {
		return nil, apperr.Validation("refresh token required", map[string]string{
			"refreshToken": "required",
		})
	}

	signingSecret := viper.GetString("jwt.secret")
	now := time.Now()

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(signingSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, apperr.Unauthorized("invalid refresh token")
	}

	if claims.ExpiresAt == nil || now.After(claims.ExpiresAt.Time) {
		return nil, apperr.Unauthorized("refresh token expired")
	}

	session, err := s.store.Queries.GetUserSessionByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.Unauthorized("invalid refresh token")
		}
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal("internal error", fmt.Errorf("RefreshToken(service) -> GetUserSessionByToken(repo): %w", err))
	}

	if session.ExpiresAt.Valid && now.After(session.ExpiresAt.Time) {
		return nil, apperr.Unauthorized("refresh token expired")
	}

	if claims.Subject != fmt.Sprintf("%d", session.UserID) {
		return nil, apperr.Unauthorized("invalid refresh token")
	}

	employee, err := s.store.Queries.GetEmployeeByUserID(ctx, &session.UserID)
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal("internal error", fmt.Errorf("RefreshToken(service) -> GetEmployeeByUserID(repo): %w", err))
	}

	accessTTL := viper.GetDuration("jwt.access_ttl")
	if accessTTL == 0 {
		accessTTL = 2 * time.Hour
	}
	refreshTTL := viper.GetDuration("jwt.refresh_ttl")
	if refreshTTL == 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	accessExp := now.Add(accessTTL)
	refreshExp := now.Add(refreshTTL)

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", session.UserID),
		ExpiresAt: jwt.NewNumericDate(accessExp),
		IssuedAt:  jwt.NewNumericDate(now),
	}).SignedString([]byte(signingSecret))
	if err != nil {
		return nil, apperr.Internal("internal error", fmt.Errorf("RefreshToken(service) -> sign access token: %w", err))
	}

	newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", session.UserID),
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		IssuedAt:  jwt.NewNumericDate(now),
	}).SignedString([]byte(signingSecret))
	if err != nil {
		return nil, apperr.Internal("internal error", fmt.Errorf("RefreshToken(service) -> sign refresh token: %w", err))
	}

	if err := s.store.Queries.DeleteUserSessionByUserID(ctx, session.UserID); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal("internal error", fmt.Errorf("RefreshToken(service) -> DeleteUserSessionByUserID(repo): %w", err))
	}

	_, err = s.store.Queries.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID:       session.UserID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshExp, Valid: true},
	})
	if err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return nil, s.pgErrToAppErr(pgErr)
		}
		return nil, apperr.Internal("internal error", fmt.Errorf("RefreshToken(service) -> CreateUserSession(repo): %w", err))
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		UserID:       employee.UniqueID,
	}, nil
}

// Logout removes all sessions for the user referenced by the refresh token.
func (s *Service) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return apperr.ValidationFromValidator(err)
	}

	session, err := s.store.Queries.GetUserSessionByToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperr.Unauthorized("invalid refresh token")
		}
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}
		return apperr.Internal("internal error", fmt.Errorf("Logout(service) -> GetUserSessionByToken(repo): %w", err))
	}

	if err := s.store.Queries.DeleteUserSessionByUserID(ctx, session.UserID); err != nil {
		if pgErr := s.asPgError(err); pgErr != nil {
			return s.pgErrToAppErr(pgErr)
		}
		return apperr.Internal("internal error", fmt.Errorf("Logout(service) -> DeleteUserSessionByUserID(repo): %w", err))
	}

	return nil
}
