package response

import "github.com/tanganyu1114/heimdallr-reborn/server/model"

type ExaCustomerResponse struct {
	Customer model.ExaCustomer `json:"customer"`
}
