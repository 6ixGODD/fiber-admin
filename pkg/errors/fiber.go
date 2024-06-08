package errors

import (
	"github.com/gofiber/fiber/v2"
)

func NotAuthorized(err error) *AppError {
	return NewAppError(CodeNotAuthorized, fiber.StatusUnauthorized, "Not authorized", err)
}

func AuthFailed(err error) *AppError {
	return NewAppError(CodeAuthFailed, fiber.StatusUnauthorized, "Authentication failed", err)
}

func TokenInvalid(err error) *AppError {
	return NewAppError(CodeTokenInvalid, fiber.StatusUnauthorized, "Token invalid", err)
}

func TokenExpired(err error) *AppError {
	return NewAppError(CodeTokenExpired, fiber.StatusUnauthorized, "Token expired", err)
}

func TokenMissed(err error) *AppError {
	return NewAppError(CodeTokenMissed, fiber.StatusUnauthorized, "Token missed", err)
}

func PermissionDeny(err error) *AppError {
	return NewAppError(CodePermissionDeny, fiber.StatusForbidden, "Permission deny", err)
}

func InvalidRequest(err error) *AppError {
	return NewAppError(CodeInvalidRequest, fiber.StatusBadRequest, "Invalid request", err)
}

func Idempotency(err error) *AppError {
	return NewAppError(CodeIdempotency, fiber.StatusBadRequest, "Idempotency check failed", err)
}

func NotFound(err error) *AppError {
	return NewAppError(CodeNotFound, fiber.StatusNotFound, "Not found", err)
}

func OperationFailed(err error) *AppError {
	return NewAppError(CodeOperationFailed, fiber.StatusInternalServerError, "Operation failed", err)
}

func DuplicateKeyError(err error) *AppError {
	return NewAppError(CodeDuplicateKey, fiber.StatusBadRequest, "Duplicate key error", err)
}

func ServerBusy(err error) *AppError {
	return NewAppError(CodeServerBusy, fiber.StatusInternalServerError, "Server busy", err)
}

func ServiceError(err error) *AppError {
	return NewAppError(CodeServiceError, fiber.StatusInternalServerError, "Service error", err)
}
