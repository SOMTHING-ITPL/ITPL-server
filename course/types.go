package course

import (
	"github.com/SOMTHING-ITPL/ITPL-server/place"
)

type Course struct {
	Places []place.Place `json:"course"`
}
