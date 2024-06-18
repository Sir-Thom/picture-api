package models

import (
	"time"
)

type Pictures struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	//Data      []byte    `json:"data"`
	AddedDate time.Time `json:"added_date"`
}
