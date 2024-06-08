package errors

const ( // Business error code
	MessageSuccess = "success"
	CodeSuccess    = 200

	CodeNotAuthorized  = 1001
	CodeAuthFailed     = 1002
	CodeTokenInvalid   = 1003
	CodeTokenExpired   = 1004
	CodeTokenMissed    = 1005
	CodePermissionDeny = 1006

	CodeInvalidRequest = 2001
	CodeIdempotency    = 2002

	CodeNotFound        = 3001
	CodeOperationFailed = 3002
	CodeDuplicateKey    = 3003

	CodeServerBusy   = 4001
	CodeServiceError = 4002

	CodeUnknownError = 9999
)
