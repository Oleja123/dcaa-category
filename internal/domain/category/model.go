package category

import (
	"database/sql"
)

type Category struct {
	Id   int            `json:"id"`
	Name string         `json:"name"`
	Info sql.NullString `json:"info"`
}
