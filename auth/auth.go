package auth

import (
	"context"
	"strings"

	"github.com/micro-community/micro-webui/namespace"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/server"
)

const (
	// BearerScheme used for Authorization header
	BearerScheme = "Bearer "
	// TokenCookieName is the name of the cookie which stores the auth token
	TokenCookieName = "micro-token"
)

//NewAuthHandlerWrapper for auth of web
func NewAuthHandlerWrapper() server.HandlerWrapper {

	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// Extract the token if the header is present. We will inspect the token regardless of if it's
			// present or not since noop auth will return a blank account upon Inspecting a blank token.
			var token string
			if header, ok := metadata.Get(ctx, "Authorization"); ok {
				// Ensure the correct scheme is being used
				if !strings.HasPrefix(header, BearerScheme) {
					return errors.Unauthorized(req.Service(), "invalid authorization header. expected Bearer schema")
				}

				// Strip the bearer scheme prefix
				token = strings.TrimPrefix(header, BearerScheme)
			}

			// Determine the namespace
			ns := auth.DefaultAuth.Options().Issuer

			var acc *auth.Account
			if a, err := auth.Inspect(token); err == nil && a.Issuer == ns {
				// We only use accounts issued by the same namespace as the service when verifying against
				// the rule set.
				ctx = auth.ContextWithAccount(ctx, a)
				acc = a
			} else if err == nil && ns == namespace.DefaultNamespace {
				// for the default domain, we want to inject the account into the context so that the
				// server can access it (since it's designed for multi-tenancy), however we don't want to
				// use it when verifying against the auth rules, since this will allow any user access to the
				// services running in the micro namespace
				ctx = auth.ContextWithAccount(ctx, a)
			}

			// ensure only accounts with the correct namespace can access this namespace,
			// since the auth package will verify access below, and some endpoints could
			// be public, we allow nil accounts access using the namespace.Public option.
			err := namespace.Authorize(ctx, ns, namespace.Public(ns))
			if err == namespace.ErrForbidden {
				return errors.Forbidden(req.Service(), err.Error())
			} else if err != nil {
				return errors.InternalServerError(req.Service(), err.Error())
			}

			// construct the resource
			res := &auth.Resource{
				Type:     "service",
				Name:     req.Service(),
				Endpoint: req.Endpoint(),
			}

			// Verify the caller has access to the resource.
			err = auth.Verify(acc, res, auth.VerifyNamespace(ns))
			if err == auth.ErrForbidden && acc != nil {
				return errors.Forbidden(req.Service(), "Forbidden call made to %v:%v by %v", req.Service(), req.Endpoint(), acc.ID)
			} else if err == auth.ErrForbidden {
				return errors.Unauthorized(req.Service(), "Unauthorized call made to %v:%v", req.Service(), req.Endpoint())
			} else if err != nil {
				return errors.InternalServerError(req.Service(), "Error authorizing request: %v", err)
			}

			// The user is authorized, allow the call
			return h(ctx, req, rsp)
		}
	}

}
