package model

import "errors"

// Domain 層錯誤定義

var (
	// API Spec 錯誤
	ErrSpecNotFound       = errors.New("API specification not found")
	ErrSpecAlreadyExists  = errors.New("API specification already exists")
	ErrInvalidSpecFormat  = errors.New("invalid specification format")
	ErrInvalidSpecVersion = errors.New("invalid specification version")
	ErrSpecParseError     = errors.New("failed to parse specification")

	// Auth 錯誤
	ErrAuthConfigNotFound = errors.New("auth configuration not found")
	ErrInvalidAuthType    = errors.New("invalid auth type")
	ErrMissingCredentials = errors.New("missing credentials")
	ErrEncryptionFailed   = errors.New("encryption failed")
	ErrDecryptionFailed   = errors.New("decryption failed")

	// Proxy 錯誤
	ErrOperationNotFound  = errors.New("operation not found")
	ErrMissingParameter   = errors.New("missing required parameter")
	ErrInvalidParameter   = errors.New("invalid parameter")
	ErrMissingSpecID      = errors.New("missing spec ID")
	ErrMissingOperationID = errors.New("missing operation ID")
	ErrProxyRequestFailed = errors.New("proxy request failed")

	// Health Check 錯誤
	ErrHealthCheckFailed = errors.New("health check failed")

	// 通用錯誤
	ErrInternalError  = errors.New("internal error")
	ErrNotImplemented = errors.New("not implemented")
)
