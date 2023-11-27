package endpoints

import (
	"net/http"

	"github.com/maxgarvey/cq_server/admin"
)

func ValidateAndExtendSession(
	r *http.Request, admin admin.Adminer,
) (bool, error) {
	session := r.Header.Get("SESSION")
	valid, err := admin.ValidateSession(session)
	if err != nil || !valid {
		return false, err
	}

	err = admin.ExtendSession(session)
	if err != nil {
		return false, err
	}
	return true, nil
}
