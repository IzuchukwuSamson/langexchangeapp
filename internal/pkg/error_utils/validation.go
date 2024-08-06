package error_utils

type ValidationError struct {
	Err error
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}
