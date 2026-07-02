package service

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
)

// WebServerConfigService defines the interface for web server config related services
type WebServerConfigService interface {
	// GetOptions retrieves option selection parameters
	GetOptions() ([]v1.BifrostGroupMeta, error)
	// GetConfigTextLines retrieves configuration file text lines
	GetConfigTextLines(opts *metav1.WebServerOptions) (string, error)
	// GetContextTextLines retrieves context configuration text lines
	GetContextTextLines(opts *metav1.WebServerConfigTargetContextOptions) (string, error)
	// GetConfig retrieves configuration JSON data
	GetConfig(opts *metav1.WebServerOptions) (*metav1.WebServerConfig, error)
	// GetIncludedConfigs retrieves included configuration file paths
	GetIncludedConfigs(opts *metav1.WebServerConfigTargetContextOptions) ([]string, error)
	// SearchContextPositions searches context positions
	SearchContextPositions(opts *metav1.WebServerConfigContextPosSearchOptions) ([]metav1.ConfigContextPos, error)
	// InsertWithClone inserts configuration context with clone
	InsertWithClone(opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]) error
	// InsertWithNew inserts new configuration context
	InsertWithNew(opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]) error
	// Remove removes configuration context
	Remove(opts *metav1.WebServerConfigTargetContextOptions) error
	// ModifyContextValue modifies context value
	ModifyContextValue(opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]) error
	// ModifyWithClone modifies configuration context with clone
	ModifyWithClone(opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]) error
	// ChangeContextEnabledState changes context enabled state
	ChangeContextEnabledState(opts *metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]) error
	// ModifyWithNew modifies configuration context with new
	ModifyWithNew(opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]) error
	// Move moves configuration context
	Move(opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]) error
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
	return s.eps.GetOptions()(s.ctx, req)
}

// GetConfigTextLines retrieves configuration file text lines
func (s *webServerConfigService) GetConfigTextLines(opts *metav1.WebServerOptions) (string, error) {
	if opts == nil {
		return "", errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerOptions]{
		Body: *opts,
	}
	return s.eps.GetConfigTextLines()(s.ctx, req)
}

// GetContextTextLines retrieves context configuration text lines
func (s *webServerConfigService) GetContextTextLines(opts *metav1.WebServerConfigTargetContextOptions) (string, error) {
	if opts == nil {
		return "", errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigTargetContextOptions]{
		Body: *opts,
	}
	return s.eps.GetContextTextLines()(s.ctx, req)
}

// GetConfig retrieves configuration JSON data
func (s *webServerConfigService) GetConfig(opts *metav1.WebServerOptions) (*metav1.WebServerConfig, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerOptions]{
		Body: *opts,
	}
	return s.eps.GetConfig()(s.ctx, req)
}

// GetIncludedConfigs retrieves included configuration file paths
func (s *webServerConfigService) GetIncludedConfigs(opts *metav1.WebServerConfigTargetContextOptions) ([]string, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigTargetContextOptions]{
		Body: *opts,
	}
	return s.eps.GetIncludedConfigs()(s.ctx, req)
}

// SearchContextPositions searches context positions
func (s *webServerConfigService) SearchContextPositions(opts *metav1.WebServerConfigContextPosSearchOptions) ([]metav1.ConfigContextPos, error) {
	if opts == nil {
		return nil, errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextPosSearchOptions]{
		Body: *opts,
	}
	return s.eps.SearchContextPositions()(s.ctx, req)
}

// InsertWithClone inserts configuration context with clone
func (s *webServerConfigService) InsertWithClone(opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]]{
		Body: *opts,
	}
	_, err := s.eps.InsertWithClone()(s.ctx, req)
	return err
}

// InsertWithNew inserts new configuration context
func (s *webServerConfigService) InsertWithNew(opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]]{
		Body: *opts,
	}
	_, err := s.eps.InsertWithNew()(s.ctx, req)
	return err
}

// Remove removes configuration context
func (s *webServerConfigService) Remove(opts *metav1.WebServerConfigTargetContextOptions) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigTargetContextOptions]{
		Body: *opts,
	}
	_, err := s.eps.Remove()(s.ctx, req)
	return err
}

// ModifyContextValue modifies context value
func (s *webServerConfigService) ModifyContextValue(opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]]{
		Body: *opts,
	}
	_, err := s.eps.ModifyContextValue()(s.ctx, req)
	return err
}

// ModifyWithClone modifies configuration context with clone
func (s *webServerConfigService) ModifyWithClone(opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]]{
		Body: *opts,
	}
	_, err := s.eps.ModifyWithClone()(s.ctx, req)
	return err
}

// ChangeContextEnabledState changes context enabled state
func (s *webServerConfigService) ChangeContextEnabledState(opts *metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]]{
		Body: *opts,
	}
	_, err := s.eps.ChangeContextEnabledState()(s.ctx, req)
	return err
}

// ModifyWithNew modifies configuration context with new
func (s *webServerConfigService) ModifyWithNew(opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]]{
		Body: *opts,
	}
	_, err := s.eps.ModifyWithNew()(s.ctx, req)
	return err
}

// Move moves configuration context
func (s *webServerConfigService) Move(opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]) error {
	if opts == nil {
		return errors.New("opts cannot be nil")
	}
	req := httpclientv1.HTTPRequest[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]]{
		Body: *opts,
	}
	_, err := s.eps.Move()(s.ctx, req)
	return err
}
