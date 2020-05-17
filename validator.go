package saban

import (
	validate "github.com/go-playground/validator/v10"
)

var (
	validator *validate.Validate
)

func init() {
	validator = validate.New()
}
