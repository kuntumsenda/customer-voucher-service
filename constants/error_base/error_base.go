package error_base

import "net/http"

type AppError struct {
	HttpCode int
	Code     string
	Message  string
}

var (
	ErrInvalidCredentials = AppError{
		HttpCode: http.StatusUnauthorized,
		Code:     "4010",
		Message:  "Username or password is incorrect",
	}

	ErrUserNotFound = AppError{
		HttpCode: http.StatusUnauthorized,
		Code:     "4011",
		Message:  "User not found",
	}

	ErrValidationFailed = AppError{
		HttpCode: http.StatusBadRequest,
		Code:     "4001",
		Message:  "Invalid",
	}

	ErrInternalServer = AppError{
		HttpCode: http.StatusInternalServerError,
		Code:     "5001",
		Message:  "Something went wrong",
	}

	ErrDB = AppError{
		HttpCode: http.StatusInternalServerError,
		Code:     "5002",
		Message:  "Something went wrong",
	}
)
