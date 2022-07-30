package app

import (
	"fmt"
	"net/http"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
)

type appErr struct {
	msg  string
	code int
}

var svcErrMessagesToStatusCodes = map[domain.ErrorType]int{
	domain.InvalidInput:     http.StatusBadRequest,
	domain.ResourceNotFound: http.StatusNotFound,
	domain.SystemError:      http.StatusInternalServerError,
	domain.ResourceConflict: http.StatusConflict,
	domain.Unauthorized:     http.StatusUnauthorized,
	domain.NotImplemented:   http.StatusNotImplemented,
}

func newAppErr(errMsg string, code int) *appErr {
	return &appErr{msg: errMsg, code: code}
}

func (app *App) newAppErrFromDomainErr(dErr *domain.Error) *appErr {
	code, ok := svcErrMessagesToStatusCodes[dErr.Code]
	if !ok {
		app.logger.Warnf("unknown domain error type: %v", dErr.Code.String())
		// in this case default to system error 500
		return &appErr{
			msg:  dErr.Error(),
			code: http.StatusInternalServerError,
		}
	}

	return &appErr{
		msg:  dErr.Error(),
		code: code,
	}
}

func (apperr *appErr) Error() string {
	return fmt.Sprintf(`{"error": "%s"}`, apperr.msg)
}

func (apperr *appErr) Code() int {
	return apperr.code
}

func httpErrorFromAppErr(w http.ResponseWriter, appErr *appErr) {
	http.Error(w, appErr.Error(), appErr.Code())
}
