package structs

import (
	_ "github.com/mailru/easyjson"
	_ "github.com/mailru/easyjson/gen"
)

//easyjson:json
type Task struct {
	Domain string `json:"domain"`
}
