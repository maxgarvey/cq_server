package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/data"
)

// AdminLogin
func AdminLogin(admin admin.Admin) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			admin.Logger.Error(
				fmt.Sprintf(
					"error reading request body: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		var request data.AdminLoginRequest
		err = json.Unmarshal(requestBody, &request)
		if err != nil {
			admin.Logger.Error(
				fmt.Sprintf(
					"error unmarshalling JSON: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		token, err := admin.Login(
			request.Username, request.Password,
		)
		if err != nil {
			admin.Logger.Error(
				fmt.Sprintf(
					"error performing login: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}

		response := data.AdminLoginResponse{
			Token: token,
		}
		json.NewEncoder(w).Encode(response)
	}
}
