package router

type ApiResponse struct {
	Response interface{} `json:"response"`
	Message  string      `json:"message"`
	Code     string      `json:"code"`
}

const RESPONSE_CODE_SUCCESS = 200
const RESPONSE_CODE_ADDED = 201
const RESPONSE_CODE_ACCEPTED = 202
const RESPONSE_CODE_UPDATED = 204
const RESPONSE_CODE_ERROR = 500
const RESPONSE_CODE_UNAUTHORIZED = 401
const RESPONSE_CODE_NOT_FOUND = 404

const RESPONSE_CODE_STORE_ID_NOT_PRESENT = 480
const RESPONSE_CODE_BAD_REQUEST = 400
const AuthHeader = `Authorization`
const DateTimeDefaultLayout = `2006-01-02 15:04:05`
const DateDefaultLayout = `2006-01-02`

func GetSuccessResponse(response interface{}) ApiResponse {
	return ApiResponse{Response: response, Message: "Success", Code: "OK"}
}
func GetErrorResponse(message string) ApiResponse {
	return GetErrorResponseWithCode(message, "ERROR")
}

func GetErrorResponseWithCode(message string, code string) ApiResponse {
	return ApiResponse{Message: message, Code: code}
}

func GetErrorResponseWithCodeAndResponse(response interface{}, message string, code string) ApiResponse {
	return ApiResponse{Response: response, Message: message, Code: code}
}
