package constants

type UserRole int8

const (
	ADMIN UserRole = iota + 1
	CUSTOMER
	COMMITTEE

	PATH string = "/proto.AuthService/"

	TOKEN_EXPIRED             string = "token has been expired"
	REFRESH_TOKEN_EXPIRED     string = "refresh token has been expired"
	FAILED_TO_EXTRACT         string = "failed to extract jwt payload"
	UNEXPECTED_SIGNING_METHOD string = "unexpected signing method: %v"

	INFO  string = "INFO"
	WARN  string = "WARN"
	ERROR string = "ERROR"
	FATAL string = "FATAL"
)
