package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/IvanJSBog/goland-todo-app/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(req *http.Request, dest any) error {
	if err := json.NewDecoder(req.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode request body json error: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	v, ok := dest.(validatable)

	var err error

	if ok {
		err = v.Validate()

	} else {
		err = requestValidator.Struct(dest)
	}
	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
