package service

import (
	"context"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
)

// AuthService defines the interface for authentication and authorization service
type AuthService interface {
	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(token string) (*response.ValidateTokenResponse, error)
	// RefreshToken refreshes a JWT token and returns the new token
	RefreshToken(token string, expiresAt int64) (*response.RefreshTokenResponse, error)
	// CheckPermission checks if a user has permission to perform an action
	CheckPermission(sub, obj, act string) (*response.CheckPermissionResponse, error)
	// LogOperation logs an operation record
	LogOperation(record *request.LogOperationRequest) error
}

// authService implements AuthService interface
type authService struct {
	ctx context.Context
	eps epclientv1.AuthEndpoints
}

// newAuthService creates a new auth service
func newAuthService(svc *factory) AuthService {
	if svc == nil || svc.eps == nil {
		return nil
	}
	return &authService{
		ctx: svc.ctx,
		eps: svc.eps.Authz(),
	}
}

// ValidateToken validates a JWT token and returns the claims
func (a *authService) ValidateToken(token string) (*response.ValidateTokenResponse, error) {
	req := httpclientv1.HTTPRequest[*request.ValidateTokenRequest]{
		Body: &request.ValidateTokenRequest{
			Token: token,
		},
	}
	resp, err := a.eps.ValidateToken()(a.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// RefreshToken refreshes a JWT token and returns the new token
func (a *authService) RefreshToken(token string, expiresAt int64) (*response.RefreshTokenResponse, error) {
	req := httpclientv1.HTTPRequest[*request.RefreshTokenRequest]{
		Body: &request.RefreshTokenRequest{
			Token:     token,
			ExpiresAt: expiresAt,
		},
	}
	resp, err := a.eps.RefreshToken()(a.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// CheckPermission checks if a user has permission to perform an action
func (a *authService) CheckPermission(sub, obj, act string) (*response.CheckPermissionResponse, error) {
	req := httpclientv1.HTTPRequest[*request.CheckPermissionRequest]{
		Body: &request.CheckPermissionRequest{
			Sub: sub,
			Obj: obj,
			Act: act,
		},
	}
	resp, err := a.eps.CheckPermission()(a.ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Response()
}

// LogOperation logs an operation record
func (a *authService) LogOperation(record *request.LogOperationRequest) error {
	req := httpclientv1.HTTPRequest[*request.LogOperationRequest]{
		Body: record,
	}
	resp, err := a.eps.LogOperation()(a.ctx, req)
	if err != nil {
		return err
	}
	_, err = resp.Response()
	return err
}
