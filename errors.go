package main

import "errors"

type Error string

var (
	/**参数不全，缺少必传参数**/
	ERR_REQUIRED_PARAMETER_MISSING = errors.New("Required parameter missing")
	ERR_FIELD_IS_NIL               = errors.New("Field is nil")
	ERR_BAD_FILED_VALUE            = errors.New("Bad field value")
)
