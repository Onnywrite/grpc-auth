package ve

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorsList []ValidationError

func From(ves validator.ValidationErrors) ValidationErrorsList {
	vel := make(ValidationErrorsList, 0, len(ves))

	for _, e := range ves {
		// f, ok := e.Type().FieldByName(e.StructField())
		// if !ok {
		// 	continue
		// }

		// has, ok := f.Tag.Lookup("secret")
		// if !ok {
		// 	has = "0"
		// }

		var value any
		// if secret, _ := strconv.ParseBool(has); !secret {
		value = e.Value()
		// }

		vel = append(vel, ValidationError{
			FieldName:          e.StructField(),
			ViolatedConstraint: e.Tag(),
			Value:              value,
		})
	}

	return vel
}

func (vel ValidationErrorsList) JSON() string {
	bytes, _ := json.Marshal(&vel)
	return string(bytes)
}

func (vel ValidationErrorsList) Error() string {
	return vel.JSON()
}

type ValidationError struct {
	FieldName          string `json:"Field"`
	ViolatedConstraint string `json:"Constraint"`
	Value              any    `json:"Value"`
}

func (ve *ValidationError) JSON() string {
	bytes, _ := json.Marshal(&ve)
	return string(bytes)
}
