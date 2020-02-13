package newsapi

import "fmt"

func (e Error) Error() string {
	return fmt.Sprint("[%s][%s] %s", e.Status, e.Code, e.Message)
}

func ApiError(err error) bool {
	_, ok := err.(*Error)
	return ok
}