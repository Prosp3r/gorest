// Tideland Go REST Server Library - Handlers - JWT Authorization
//
// Copyright (C) 2009-2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"time"

	"github.com/tideland/gorest/jwt"
	"github.com/tideland/gorest/rest"
)

//--------------------
// JWT AUTHORIZATION HANDLER
//--------------------

// JWTAuthorizationConfig allows to control how the JWT authorization
// handler works. All values are optional. In this case tokens are only
// decoded without using a cache, validated for the time, and there's
// no user defined gatekeeper function running afterwards.
type JWTAuthorizationConfig struct {
	Cache      jwt.Cache
	Key        jwt.Key
	Leeway     time.Duration
	Gatekeeper func(job rest.Job, claims jwt.Claims) (bool, error)
}

// jwtAuthorizationHandler checks for a valid token and then runs
// a gatekeeper function.
type jwtAuthorizationHandler struct {
	id         string
	cache      jwt.Cache
	key        jwt.Key
	leeway     time.Duration
	gatekeeper func(job rest.Job, claims jwt.Claims) (bool, error)
}

// NewJWTAuthorizationHandler creates a handler checking for a valid JSON
// Web Token in each request.
func NewJWTAuthorizationHandler(id string, config *JWTAuthorizationConfig) rest.ResourceHandler {
	h := &jwtAuthorizationHandler{
		id:     id,
		leeway: time.Minute,
	}
	if config != nil {
		if config.Cache != nil {
			h.cache = config.Cache
		}
		if config.Key != nil {
			h.key = config.Key
		}
		if config.Leeway != 0 {
			h.leeway = config.Leeway
		}
		if config.Gatekeeper != nil {
			h.gatekeeper = config.Gatekeeper
		}
	}
	return h
}

// ID is specified on the ResourceHandler interface.
func (h *jwtAuthorizationHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *jwtAuthorizationHandler) Init(env rest.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *jwtAuthorizationHandler) Get(job rest.Job) (bool, error) {
	return h.check(job)
}

// Head is specified on the HeadResourceHandler interface.
func (h *jwtAuthorizationHandler) Head(job rest.Job) (bool, error) {
	return h.check(job)
}

// Put is specified on the PutResourceHandler interface.
func (h *jwtAuthorizationHandler) Put(job rest.Job) (bool, error) {
	return h.check(job)
}

// Post is specified on the PostResourceHandler interface.
func (h *jwtAuthorizationHandler) Post(job rest.Job) (bool, error) {
	return h.check(job)
}

// Patch is specified on the PatchResourceHandler interface.
func (h *jwtAuthorizationHandler) Patch(job rest.Job) (bool, error) {
	return h.check(job)
}

// Delete is specified on the DeleteResourceHandler interface.
func (h *jwtAuthorizationHandler) Delete(job rest.Job) (bool, error) {
	return h.check(job)
}

// Options is specified on the OptionsResourceHandler interface.
func (h *jwtAuthorizationHandler) Options(job rest.Job) (bool, error) {
	return h.check(job)
}

// check is used by all methods to check the token.
func (h *jwtAuthorizationHandler) check(job rest.Job) (bool, error) {
	var jobJWT jwt.JWT
	var err error
	switch {
	case h.cache != nil && h.key != nil:
		jobJWT, err = jwt.VerifyCachedFromJob(job, h.cache, h.key)
	case h.cache != nil && h.key == nil:
		jobJWT, err = jwt.DecodeCachedFromJob(job, h.cache)
	case h.cache == nil && h.key != nil:
		jobJWT, err = jwt.VerifyFromJob(job, h.key)
	default:
		jobJWT, err = jwt.DecodeFromJob(job)
	}
	if err != nil {
		return false, h.deny(job, err.Error())
	}
	if jobJWT == nil {
		return false, h.deny(job, "no JSON Web Token")
	}
	if !jobJWT.IsValid(h.leeway) {
		return false, h.deny(job, "invalid JSON Web Token")
	}
	if h.gatekeeper != nil {
		return h.gatekeeper(job, jobJWT.Claims())
	}
	return true, nil
}

// deny sends a negative feedback to the caller.
func (h *jwtAuthorizationHandler) deny(job rest.Job, msg string) error {
	job.ResponseWriter().WriteHeader(http.StatusUnauthorized)
	switch {
	case job.AcceptsContentType(rest.ContentTypeJSON):
		return rest.NegativeFeedback(job.JSON(true), msg)
	case job.AcceptsContentType(rest.ContentTypeXML):
		return rest.NegativeFeedback(job.XML(), msg)
	default:
		job.ResponseWriter().Header().Set("Content-Type", rest.ContentTypePlain)
		job.ResponseWriter().Write([]byte(msg))
		return nil
	}
}

// EOF