// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// Get - Account. GET /account
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-account/
func (s *AccountService) Get(ctx context.Context) (*AccountResponse, error) {
	return get[AccountResponse](ctx, s.client, "/account", nil, "get_account")
}
