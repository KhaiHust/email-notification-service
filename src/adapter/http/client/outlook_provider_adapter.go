package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/adapter/http/strategy"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/web/client"
	"golang.org/x/oauth2"
	"net/url"
)

type OutlookProviderAdapter struct {
	httpClient    client.ContextualHttpClient
	props         *properties.OutlookProviderProperties
	msOAuthConfig *oauth2.Config
}

func (o *OutlookProviderAdapter) GetType() string {
	return constant.EmailProviderOutlook
}

// --- EXCHANGE CODE FOR TOKEN ---
func (o *OutlookProviderAdapter) GetOAuthInfoByCode(ctx context.Context, code string) (*response.OAuthInfoResponseDto, error) {
	token, err := o.msOAuthConfig.Exchange(ctx, code)
	if err != nil {
		log.Error(ctx, "Error when exchange code to access token", err)
		return nil, err
	}
	email, err := o.getMicrosoftUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	return &response.OAuthInfoResponseDto{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Email:        *email,
		ExpiredAt:    token.Expiry.Unix(),
	}, nil
}

// --- REFRESH TOKEN ---
// todo: check later
func (o *OutlookProviderAdapter) GetOAuthByRefreshToken(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*response.OAuthInfoResponseDto, error) {
	token := &oauth2.Token{
		RefreshToken: emailProviderEntity.OAuthRefreshToken,
	}
	newToken, err := o.msOAuthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		log.Error(ctx, "Error when refresh token", err)
		return nil, err
	}
	email, err := o.getMicrosoftUserInfo(ctx, newToken)
	if err != nil {
		log.Error(ctx, "Error when get user info", err)
		return nil, err
	}
	refreshToken := newToken.RefreshToken
	if newToken.RefreshToken == "" {
		refreshToken = emailProviderEntity.OAuthRefreshToken
	}
	return &response.OAuthInfoResponseDto{
		AccessToken:  newToken.AccessToken,
		RefreshToken: refreshToken,
		Email:        *email,
		ExpiredAt:    newToken.Expiry.Unix(),
	}, nil
}

// --- BUILD OAUTH2 URL ---
func (o *OutlookProviderAdapter) GetOAuthUrl() (*response.OAuthUrlResponseDto, error) {
	var params = url.Values{}
	params.Add("client_id", o.props.ClientID)
	params.Add("redirect_uri", o.props.RedirectURI)
	params.Add("response_type", o.props.ResponseType)
	params.Add("scope", o.props.Scopes)
	params.Add("access_type", o.props.AccessType)
	params.Add("prompt", "consent")
	oauthUrl := o.props.BaseOAuthURL + "?" + params.Encode()
	return &response.OAuthUrlResponseDto{Url: oauthUrl}, nil
}

// --- SEND EMAIL VIA MICROSOFT GRAPH API ---
func (o *OutlookProviderAdapter) SendEmail(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity, emailData *request.EmailDataDto) error {
	token := &oauth2.Token{
		AccessToken:  emailProviderEntity.OAuthToken,
		TokenType:    "Bearer",
		RefreshToken: emailProviderEntity.OAuthRefreshToken,
	}
	msClient := o.msOAuthConfig.Client(ctx, token)

	apiUrl := "https://graph.microsoft.com/v1.0/me/sendMail"
	payload := o.buildGraphSendMailPayload(emailProviderEntity.Email, emailData)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(ctx, "Error marshaling email payload", err)
		return err
	}

	resp, err := msClient.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(ctx, "Error sending email", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 { // 202 Accepted = success
		log.Error(ctx, "Error sending email, status code: %d", resp.StatusCode)
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return common.ErrUnauthorized
		}
		return fmt.Errorf("failed to send email: status code %d", resp.StatusCode)
	}
	return nil
}

// --- GET MICROSOFT USER EMAIL ---
func (o *OutlookProviderAdapter) getMicrosoftUserInfo(ctx context.Context, accessToken *oauth2.Token) (*string, error) {
	msClient := o.msOAuthConfig.Client(ctx, accessToken)
	resp, err := msClient.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		log.Error(ctx, "Error when get user info", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error(ctx, "Error when get user info, status code: %d", resp.StatusCode)
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return nil, common.ErrUnauthorized
		}
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var userInfo struct {
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Error(ctx, "Error when decode user info", err)
		return nil, err
	}
	// Prefer "Mail" but fallback to "UserPrincipalName"
	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}
	return &email, nil
}

// --- BUILD SENDMAIL PAYLOAD ---
func (o *OutlookProviderAdapter) buildGraphSendMailPayload(from string, data *request.EmailDataDto) map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]interface{}{
			"subject": data.Subject,
			"body": map[string]string{
				"contentType": "HTML",
				"content":     data.Body,
			},
			"toRecipients": []map[string]interface{}{
				{"emailAddress": map[string]string{"address": data.To}},
			},
			"from": map[string]interface{}{
				"emailAddress": map[string]string{"address": from},
			},
		},
		"saveToSentItems": "true",
	}
}

// --- CONSTRUCTOR ---
func NewOutlookProviderAdapter(httpClient client.ContextualHttpClient, props *properties.OutlookProviderProperties) strategy.IEmailProviderStrategy {
	return &OutlookProviderAdapter{
		httpClient:    httpClient,
		props:         props,
		msOAuthConfig: properties.NewMicrosoftOAuthConfig(props),
	}
}
