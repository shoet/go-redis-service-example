package errorutil

const (
	ErrorMessageNotFound            = "Not Found"
	ErrorMessageInternalServerError = "Internal Server Error"
	ErrorMessageUnauthorized        = "Unauthorized"
)

type ErrorResponse struct {
	Message string `json:"message"`
}
