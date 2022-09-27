package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"

	"github.com/spy16/moonshot/pkg/errors"
	"github.com/spy16/moonshot/pkg/httputils"
	"github.com/spy16/moonshot/pkg/log"
)

// Routes installs the user-management routes to the given router.
func (reg *Registry) Routes(r *chi.Mux) error {
	r.Get("/auth/methods", reg.listAuthMethods)
	r.Post("/auth/login", reg.doLogin)
	return nil
}

func (reg *Registry) listAuthMethods(w http.ResponseWriter, r *http.Request) {
	redirectURL := r.URL.Query().Get("redirect_url")

	type authProvider struct {
		Name    string `json:"name"`
		State   string `json:"state"`
		AuthURL string `json:"auth_url"`
	}

	var authMethods struct {
		// NOTE: this is nested to allow email_auth flags in the future.
		AuthProviders []authProvider `json:"auth_providers"`
	}

	for _, provider := range reg.Providers {
		state := randomString(16)
		authMethods.AuthProviders = append(authMethods.AuthProviders, authProvider{
			Name:    provider.Name,
			State:   state,
			AuthURL: provider.Config.AuthCodeURL(state, oauth2.SetAuthURLParam("redirect_url", redirectURL)),
		})
	}

	httputils.Respond(w, r, http.StatusOK, authMethods)
}

func (reg *Registry) doLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Method string `json:"method"`

		// when kind='oauth2'
		Code           string `json:"code"`
		Provider       string `json:"provider"`
		RedirectionURL string `json:"redirection_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputils.Respond(w, r, http.StatusBadRequest,
			errors.ErrInvalid.WithMsgf("request body is not valid json"))
		return
	}

	switch body.Method {
	case "oauth2":
		p := reg.getProvider(body.Provider)
		if p == nil {
			httputils.Respond(w, r, http.StatusBadRequest,
				errors.ErrInvalid.WithMsgf("unknown provider '%s'", body.Provider))
			return
		}

		tok, err := p.Config.Exchange(r.Context(), body.Code)
		if err != nil {
			httputils.Respond(w, r, http.StatusInternalServerError,
				errors.ErrInvalid.
					WithMsgf("login failed").
					WithCausef(err.Error()),
			)
			return
		}
		log.Infof(r.Context(), "token issued: %v", tok)

	default:
		httputils.Respond(w, r, http.StatusBadRequest,
			errors.ErrInvalid.WithMsgf("invalid login method"))
		return
	}
}

func (reg *Registry) getProvider(name string) *Provider {
	for _, provider := range reg.Providers {
		if provider.Name == name {
			return &provider
		}
	}
	return nil
}
