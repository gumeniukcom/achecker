package structs

import (
	_ "github.com/mailru/easyjson"
	_ "github.com/mailru/easyjson/gen"
)

//easyjson:json
type CheckResult struct {
	Domain     string `json:"domain"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}
