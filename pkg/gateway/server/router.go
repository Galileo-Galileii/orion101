package server

import (
	_ "embed"

	"github.com/orion101-ai/orion101/pkg/api"
	"github.com/orion101-ai/orion101/pkg/api/server"
)

func (s *Server) AddRoutes(mux *server.Server) {
	wrap := func(h api.HandlerFunc) api.HandlerFunc {
		return apply(h, addRequestID, addLogger, logRequest, contentType("application/json"))
	}
	// All the routes served by the API will start with `/api`
	mux.HandleFunc("GET /api/me", wrap(s.getCurrentUser))
	mux.HandleFunc("GET /api/users", wrap(s.getUsers))
	mux.HandleFunc("GET /api/users/{username}", wrap(s.getUser))
	mux.HandleFunc("PATCH /api/users/{username}", wrap(s.updateUser))
	mux.HandleFunc("DELETE /api/users/{username}", wrap(s.deleteUser))

	mux.HandleFunc("POST /api/token-request", s.tokenRequest)
	mux.HandleFunc("GET /api/token-request/{id}", s.checkForToken)
	mux.HandleFunc("GET /api/token-request/{id}/{service}", s.redirectForTokenRequest)

	mux.HandleFunc("GET /api/tokens", wrap(s.getTokens))
	mux.HandleFunc("DELETE /api/tokens/{id}", wrap(s.deleteToken))
	mux.HandleFunc("POST /api/tokens", wrap(s.newToken))

	mux.HandleFunc("POST /api/auth-providers", wrap(s.createAuthProvider))
	mux.HandleFunc("PATCH /api/auth-providers/{slug}", wrap(s.updateAuthProvider))
	mux.HandleFunc("DELETE /api/auth-providers/{slug}", wrap(s.deleteAuthProvider))
	mux.HandleFunc("GET /api/auth-providers", s.getAuthProviders)
	mux.HandleFunc("GET /api/auth-providers/{slug}", s.getAuthProvider)
	mux.HandleFunc("POST /api/auth-providers/{slug}/disable", wrap(s.disableAuthProvider))
	mux.HandleFunc("POST /api/auth-providers/{slug}/enable", wrap(s.enableAuthProvider))

	mux.HandleFunc("GET /api/oauth/start/{id}/{service}", wrap(s.oauth))
	mux.HandleFunc("/api/oauth/redirect/{service}", wrap(s.redirect))

	// CRUD routes for OAuth Apps (integrations with other service such as Microsoft 365)
	mux.HandleFunc("GET /api/oauth-apps", wrap(s.listOAuthApps))
	mux.HandleFunc("GET /api/oauth-apps/{id}", wrap(s.oauthAppByID))
	mux.HandleFunc("POST /api/oauth-apps", wrap(s.createOAuthApp))
	mux.HandleFunc("PATCH /api/oauth-apps/{id}", wrap(s.updateOAuthApp))
	mux.HandleFunc("DELETE /api/oauth-apps/{id}", wrap(s.deleteOAuthApp))

	// Routes for OAuth authorization code flow
	mux.HandleFunc("GET /api/app-oauth/authorize/{id}", wrap(s.authorizeOAuthApp))
	mux.HandleFunc("GET /api/app-oauth/refresh/{id}", wrap(s.refreshOAuthApp))
	mux.HandleFunc("GET /api/app-oauth/callback/{id}", wrap(s.callbackOAuthApp))

	// Route for credential tools to get their OAuth tokens
	mux.HandleFunc("GET /api/app-oauth/get-token", wrap(s.getTokenOAuthApp))

	// LLM proxy
	mux.HandleFunc("POST /api/llm-proxy/{path...}", s.llmProxy)

}