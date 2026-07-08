package service

import (
	"context"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// WebServerBinCMDService defines the interface for web server binary command related services
type WebServerBinCMDService interface {
	// Exec executes a binary command
	Exec(req *metav1.WebServerBinCMDExecRequest) (*metav1.WebServerBinCMDExecResponse, error)
}

// webServerBinCMDService implements WebServerBinCMDService interface
type webServerBinCMDService struct {
	ctx context.Context
	eps epclientv1.WebServerBinCMDEndpoints
}

// newWebServerBinCMDService creates a new WebServerBinCMD service
func newWebServerBinCMDService(svc *factory) WebServerBinCMDService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &webServerBinCMDService{
		ctx: svc.ctx,
		eps: svc.eps.WebServerBinCMDs(),
	}
}

// Exec executes a binary command
func (s *webServerBinCMDService) Exec(req *metav1.WebServerBinCMDExecRequest) (*metav1.WebServerBinCMDExecResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	httpReq := httpclientv1.HTTPRequest[metav1.WebServerBinCMDExecRequest]{
		Body: *req,
	}
	resp, err := s.eps.Exec()(s.ctx, httpReq)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}
