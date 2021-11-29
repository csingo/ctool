package response

import "net/http"

func Success(data interface{}) (int, interface{}) {
	res := &apiData{
		Code:    0,
		Message: "success",
		Data:    data,
	}

	return http.StatusOK, res
}

func Error(code int, msg string) (int, interface{}) {
	res := &apiData{
		Code:    code,
		Message: msg,
		Data:    struct{}{},
	}

	return http.StatusOK, res
}
