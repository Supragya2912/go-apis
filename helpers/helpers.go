package helpers

func ErrResponse(statusCode int, message string, code string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}
