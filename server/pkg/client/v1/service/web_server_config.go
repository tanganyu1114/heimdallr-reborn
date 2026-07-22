package service

import (
	"context"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// WebServerConfigService defines the interface for web server config related services
type WebServerConfigService interface {
	// GetOptions retrieves option selection parameters
	GetOptions() ([]v1.BifrostGroupMeta, error)
	// GetConfigTextLines retrieves configuration file text lines
	GetConfigTextLines(opts *v1.WebServerOptions) (string, error)
	// GetContextTextLines retrieves context configuration text lines
	GetContextTextLines(opts *v1.WebServerConfigTargetContextOptions) (string, error)
	// GetConfig retrieves configuration JSON data
	GetConfig(opts *v1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprints, error)
	// GetIncludedConfigs retrieves included configuration file paths
	GetIncludedConfigs(opts *v1.WebServerConfigTargetContextOptions) ([]string, error)
	// SearchContextPositions searches context positions
	SearchContextPositions(opts *v1.WebServerConfigContextPosSearchOptions) ([]v1.ConfigContextPos, error)
	// InsertWithClone inserts configuration context with clone
	InsertWithClone(opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]) error
	// InsertWithNew inserts new configuration context
	InsertWithNew(opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]) error
	// Remove removes configuration context
	Remove(opts *v1.WebServerConfigTargetContextOptions) error
	// UpdateConfig updates configuration JSON data
	UpdateConfig(opts *v1.WebServerConfigUpdateOptions) error
	// ModifyContextValue modifies context value
	ModifyContextValue(opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]) error
	// ModifyWithClone modifies configuration context with clone
	ModifyWithClone(opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]) error
	// ChangeContextEnabledState changes context enabled state
	ChangeContextEnabledState(opts *v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]) error
	// ModifyWithNew modifies configuration context with new
	ModifyWithNew(opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]) error
	// Move moves configuration context
	Move(opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]) error
}

// webServerConfigService implements WebServerConfigService interface
type webServerConfigService struct {
	ctx context.Context
	eps epclientv1.WebServerConfigEndpoints
}

// newWebServerConfigService creates a new WebServerConfig service
func newWebServerConfigService(svc *factory) WebServerConfigService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &webServerConfigService{
		ctx: svc.ctx,
		eps: svc.eps.WebServerConfigs(),
	}
}

// GetOptions retrieves option selection parameters
func (s *webServerConfigService) GetOptions() ([]v1.BifrostGroupMeta, error) {
	req := httpclientv1.HTTPRequest[httpclientv1.NilBody]{
		Body: httpclientv1.NilBody{},
	}
	resp, err := s.eps.GetOptions()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// GetConfigTextLines retrieves configuration file text lines
func (s *webServerConfigService) GetConfigTextLines(opts *v1.WebServerOptions) (string, error) {
	if opts == nil {
		return "", errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerOptions]{
		Body: *opts,
	}
	resp, err := s.eps.GetConfigTextLines()(s.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Response()
}

// GetContextTextLines retrieves context configuration text lines
func (s *webServerConfigService) GetContextTextLines(opts *v1.WebServerConfigTargetContextOptions) (string, error) {
	if opts == nil {
		return "", errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigTargetContextOptions]{
		Body: *opts,
	}
	resp, err := s.eps.GetContextTextLines()(s.ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Response()
}

// GetConfig retrieves configuration JSON data
func (s *webServerConfigService) GetConfig(opts *v1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprints, error) {
	if opts == nil {
		return nil, utilsV3.ConfigFingerprints{}, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerOptions]{
		Body: *opts,
	}
	resp, err := s.eps.GetConfig()(s.ctx, req)
	if err != nil {
		return nil, utilsV3.ConfigFingerprints{}, err
	}
	webServerConfig, err := resp.Response()
	if err != nil {
		return nil, utilsV3.ConfigFingerprints{}, err
	}
	ngconf, err := configuration.NewNginxConfigFromJsonBytes(webServerConfig.Config)
	if err != nil {
		return nil, utilsV3.ConfigFingerprints{}, err
	}
	return ngconf, webServerConfig.OriginalFingerprints, nil
}

// GetIncludedConfigs retrieves included configuration file paths
func (s *webServerConfigService) GetIncludedConfigs(opts *v1.WebServerConfigTargetContextOptions) ([]string, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigTargetContextOptions]{
		Body: *opts,
	}
	resp, err := s.eps.GetIncludedConfigs()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// SearchContextPositions searches context positions
func (s *webServerConfigService) SearchContextPositions(opts *v1.WebServerConfigContextPosSearchOptions) ([]v1.ConfigContextPos, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextPosSearchOptions]{
		Body: *opts,
	}
	resp, err := s.eps.SearchContextPositions()(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// InsertWithClone inserts configuration context with clone
func (s *webServerConfigService) InsertWithClone(opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.InsertWithClone()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// InsertWithNew inserts new configuration context
func (s *webServerConfigService) InsertWithNew(opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.InsertWithNew()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// Remove removes configuration context
func (s *webServerConfigService) Remove(opts *v1.WebServerConfigTargetContextOptions) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigTargetContextOptions]{
		Body: *opts,
	}
	resp, err := s.eps.Remove()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// UpdateConfig updates configuration JSON data
func (s *webServerConfigService) UpdateConfig(opts *v1.WebServerConfigUpdateOptions) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[*v1.WebServerConfigUpdateOptions]{
		Body: opts,
	}
	resp, err := s.eps.UpdateConfig()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// ModifyContextValue modifies context value
func (s *webServerConfigService) ModifyContextValue(opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.ModifyContextValue()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// ModifyWithClone modifies configuration context with clone
func (s *webServerConfigService) ModifyWithClone(opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.ModifyWithClone()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// ChangeContextEnabledState changes context enabled state
func (s *webServerConfigService) ChangeContextEnabledState(opts *v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.ChangeContextEnabledState()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// ModifyWithNew modifies configuration context with new
func (s *webServerConfigService) ModifyWithNew(opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.ModifyWithNew()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}

// Move moves configuration context
func (s *webServerConfigService) Move(opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]]{
		Body: *opts,
	}
	resp, err := s.eps.Move()(s.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}
