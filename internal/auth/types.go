package auth

type (
	ServiceClientGetRequest struct {
		ClientId string `uri:"client_id" binding:"required"`
	}
	ServiceClientGetResponse struct {
		ClientId    string `json:"client_id"`
		Name        string `json:"name"`
		Scope       string `json:"scope"`
		RedirectUri string `json:"redirect_uri"`
	}
)

// 認証
type (
	AuthenticationRequest struct {
		ClientId string `json:"client_id" binding:"required"`
		UserId   string `json:"user_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	AuthenticationResponse struct {
		// set clientId, userId
		JWT string `json:"jwt"`
	}
)

// 認可リクエスト(OAuth2.0)
type (
	AuthorizationRequest struct {
		JWT          string `json:"jwt" binding:"required"`
		ClientId     string `json:"client_id" binding:"required"`
		ResponseType string `json:"response_type" binding:"required"` // must 'code'
		Scope        string `json:"scope" binding:"required"`
	}
	AuthorizationResponse struct {
		Code string `json:"code"`
	}
)

// アクセストークンリクエスト(OAuth2.0)
type (
	AccessTokenRequest struct {
		GrantType    string `json:"grant_type" binding:"required"` // must 'authorization_code'
		ClientId     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
		Code         string `json:"code" binding:"-"`
		RefreshToken string `json:"refresh_token" binding:"-"`
	}
	AccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    uint   `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
)
