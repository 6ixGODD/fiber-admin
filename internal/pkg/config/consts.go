package config

import (
	"fiber-admin/pkg/zap"
)

// Context Key
const (
	UserIDKey    = zap.UserIDKey
	RequestIDKey = zap.RequestIDKey
)

// Enum Values
const (
	NoticeTypeUrgent = "URGENT"
	NoticeTypeNormal = "NORMAL"

	OperationTypeCreate = "CREATE"
	OperationTypeUpdate = "UPDATE"
	OperationTypeDelete = "DELETE"

	EntityTypeUser          = "USER"
	EntityTypeDocumentation = "DOCUMENTATION"
	EntityTypeNotice        = "NOTICE"

	OperationStatusSuccess = "SUCCESS"
	OperationStatusFailure = "FAILURE"

	UserRoleUser  = "USER"
	UserRoleAdmin = "ADMIN"
)

// MongoDB Collection Name
const (
	DocumentationCollectionName = "documentation"
	NoticeCollectionName        = "notice"
	LoginLogCollectionName      = "login_log"
	OperationLogCollectionName  = "operation_log"
	UserCollectionName          = "user"
)

// cache Prefix / Key
const (
	NoticeCachePrefix         = "dao:notice"
	UserCachePrefix           = "dao:user"
	DocumentationCachePrefix  = "dao:documentation"
	TokenBlacklistCachePrefix = "token:blacklist"
	IdempotencyCachePrefix    = "idempotency"

	LoginLogCacheKey     = "log:login"
	OperationLogCacheKey = "log:operation"

	CacheTrue = "1"
)
