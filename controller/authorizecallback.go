package controller

import (
	"github.com/fabric8-services/fabric8-auth/app"
	"github.com/fabric8-services/fabric8-auth/errors"
	"github.com/fabric8-services/fabric8-auth/jsonapi"
	"github.com/fabric8-services/fabric8-auth/login"
	"github.com/fabric8-services/fabric8-auth/token"
	"github.com/goadesign/goa"
)

// AuthorizecallbackController implements the authorizecallback resource.
type AuthorizecallbackController struct {
	*goa.Controller
	Auth          login.KeycloakOAuthService
	TokenManager  token.Manager
	Configuration LoginConfiguration
}

// NewAuthorizecallbackController creates a authorizecallback controller.
func NewAuthorizecallbackController(service *goa.Service, auth *login.KeycloakOAuthProvider, tokenManager token.Manager, configuration LoginConfiguration) *AuthorizecallbackController {
	return &AuthorizecallbackController{Controller: service.NewController("AuthorizecallbackController"), Auth: auth, TokenManager: tokenManager, Configuration: configuration}
}

// Authorizecallback runs the authorizecallback action.
func (c *AuthorizecallbackController) Authorizecallback(ctx *app.AuthorizecallbackAuthorizecallbackContext) error {
	// AuthorizecallbackController_Authorizecallback: start_implement

	if ctx.State == nil {
		return jsonapi.JSONErrorResponse(ctx, errors.NewBadParameterError("state", "nil").Expected("State"))
	}
	if ctx.Code == nil {
		return jsonapi.JSONErrorResponse(ctx, errors.NewBadParameterError("code", "nil").Expected("Code"))
	}

	return c.Auth.PerformAuthorizeCallback(ctx)
}
