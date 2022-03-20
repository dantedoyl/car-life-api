package utils

import "encoding/json"

type Error struct {
	HttpError int       `json:"-"`
	Message   string    `json:"message"`
}

func JSONError(error *Error) []byte {
	jsonError, err := json.Marshal(error)
	if err != nil {
		return []byte("")
	}
	return jsonError
}
