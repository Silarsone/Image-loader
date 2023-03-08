package response

import "encoding/json"

type Response struct {
	Data  any  `json:"data"`
	Error bool `json:"error"`
}

func ParseResponse(data any, isError bool) ([]byte, error) {
	resp := Response{
		Data:  data,
		Error: isError,
	}

	return json.Marshal(&resp)
}
