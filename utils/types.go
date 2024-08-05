package utils

type CtxKey string

// ErrMessage represents the structure of errors in the output
type ErrMessage struct {
	Error string `json:"error"`
}

// ResponseMsg represents the structure of responses in the output
type ResponseMsg struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MailResponse struct {
	Message string `json:"message"`
}

var GenericError = "there was an error"

const (
	UsersTable             = "users"
	EmailVerificationTable = "emailverification"
	PasswordReset          = "password_reset"
)

const (
	UserRole = "user"
)

const (
	MethodPost = "POST"
)

const (
	AuthHeader       = "Authorization"
	AuthHeaderPrefix = "Bearer:"
	TokenPrefix      = "token_"
	ResetPrefix      = "reset:"
)
