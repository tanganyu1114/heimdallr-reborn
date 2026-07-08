package service

import (
	"context"
	"gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// WebServerStatisticsService defines the interface for web server statistics related services
type WebServerStatisticsService interface {
	// GetProxyServiceInfo retrieves proxy service information
	GetProxyServiceInfo(opts *metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error)
	// ConnectivityCheckOfProxyService performs connectivity check of proxy service
	ConnectivityCheckOfProxyService(opts *metav1.ConnectivityCheckOfProxiedServersRequestOptions) (v1.ProxyServiceInfo, error)
	// ExportProxyServiceInfoToExcel exports proxy service info to Excel
	ExportProxyServiceInfoToExcel(opts *metav1.WebServerOptions) ([]byte, error)
}

// webServerStatisticsService implements WebServerStatisticsService interface
type webServerStatisticsService struct {
	ctx context.Context
	eps epclientv1.WebServerStatisticsEndpoints
}

// newWebServerStatisticsService creates a new WebServerStatistics service
func newWebServerStatisticsService(svc *factory) WebServerStatisticsService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &webServerStatisticsService{
		ctx: svc.ctx,
		eps: svc.eps.WebServerStatistics(),
	}
}

// GetProxyServiceInfo retrieves proxy service information
func (s *webServerStatisticsService) GetProxyServiceInfo(opts *metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerOptions]{
		Body: *opts,
	}
	resp, err := s.eps.GetProxyServiceInfo()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// ConnectivityCheckOfProxyService performs connectivity check of proxy service
func (s *webServerStatisticsService) ConnectivityCheckOfProxyService(opts *metav1.ConnectivityCheckOfProxiedServersRequestOptions) (v1.ProxyServiceInfo, error) {
	if opts == nil {
		return v1.ProxyServiceInfo{}, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.ConnectivityCheckOfProxiedServersRequestOptions]{
		Body: *opts,
	}
	resp, err := s.eps.ConnectivityCheckOfProxyService()(s.ctx, req)
	if err != nil {
		return v1.ProxyServiceInfo{}, err
	}
	return resp.Response()
}

// ExportProxyServiceInfoToExcel exports proxy service info to Excel
func (s *webServerStatisticsService) ExportProxyServiceInfoToExcel(opts *metav1.WebServerOptions) ([]byte, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerOptions]{
		Body: *opts,
	}
	resp, err := s.eps.ExportProxyServiceInfoToExcel()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}
