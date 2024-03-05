package auth

type (
	ServiceClientGetRequest struct {
		ClientID string `uri:"client_id" binding:"required"`
	}
	ServiceClientGetResponse struct {
		ClientID string `json:"client_id"`
		Name     string `json:"name"`
		Scope    string `json:"scope"`
	}
)

// 認証
type (
	AuthenticationRequest struct {
		ClientID string `json:"client_id" binding:"required"`
		UserID   string `json:"user_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	AuthenticationResponse struct {
		// set clientID, userID
		JWT string `json:"jwt"`
	}
)

// 認可リクエスト(OAuth2.0)
type (
	AuthorizationRequest struct {
		JWT          string `json:"jwt" binding:"required"`
		ClientID     string `json:"client_id" binding:"required"`
		ResponseType string `json:"response_type" binding:"required"` // must 'code'
		Scope        string `json:"scope" binding:"required"`
	}
	AuthorizationResponse struct {
		// redirect URI (expected in server) preconfigured for each client with authorization_code
		Code string `json:"code"`
	}
)

// アクセストークンリクエスト(OAuth2.0)
type (
	AccessTokenRequest struct {
		GrantType    string `json:"grant_type" binding:"required"` // must 'authorization_code'
		ClientID     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
		Code         string `json:"code" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"-"`
	}
	AccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    uint   `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
)
