package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/adapter/http/client/dto"
	"github.com/KhaiHust/email-notification-service/adapter/http/strategy"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/web/client"
	"golang.org/x/oauth2"
	"net/url"
	"strings"
)

type GmailProviderAdapter struct {
	httpClient        client.ContextualHttpClient
	props             *properties.GmailProviderProperties
	googleOAuthConfig *oauth2.Config
}

func (g GmailProviderAdapter) SendEmail(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity, emailData *request.EmailDataDto) error {
	message := g.buildMessage(emailProviderEntity.Email, emailData)
	token := &oauth2.Token{
		AccessToken:  emailProviderEntity.OAuthToken,
		TokenType:    "Bearer",
		RefreshToken: emailProviderEntity.OAuthRefreshToken,
	}

	ggClient := g.googleOAuthConfig.Client(ctx, token)
	rawMessage := base64.URLEncoding.EncodeToString([]byte(message))

	url := "https://gmail.googleapis.com/gmail/v1/users/me/messages/send"
	payload := map[string]string{"raw": rawMessage}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(ctx, "Error marshaling email payload", err)
		return err
	}

	resp, err := ggClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(ctx, "Error sending email", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(ctx, "Error sending email, status code: %d", resp.StatusCode)
		return fmt.Errorf("failed to send email: status code %d", resp.StatusCode)
	}
	return nil
}
func (g GmailProviderAdapter) GetOAuthInfo(ctx context.Context, code string) (*response.OAuthInfoResponseDto, error) {
	token, err := g.googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		log.Error(ctx, "Error when exchange code to access token", err)
		return nil, err
	}
	email, err := g.getGmailUserInfo(ctx, token)
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

func (g GmailProviderAdapter) GetOAuthUrl() (*response.OAuthUrlResponseDto, error) {
	var params = url.Values{}
	params.Add("client_id", g.props.ClientID)
	params.Add("redirect_uri", g.props.RedirectURI)
	params.Add("response_type", g.props.ResponseType)
	params.Add("scope", g.props.Scopes)
	params.Add("access_type", g.props.AccessType)
	oauthUrl := g.props.BaseOAuthURL + "?" + params.Encode()
	return &response.OAuthUrlResponseDto{Url: oauthUrl}, nil
}

func NewGmailProviderAdapter(httpClient client.ContextualHttpClient, props *properties.GmailProviderProperties) strategy.IEmailProviderStrategy {
	return &GmailProviderAdapter{
		httpClient:        httpClient,
		props:             props,
		googleOAuthConfig: properties.NewGoogleOAuthConfig(props),
	}
}
func (g *GmailProviderAdapter) getGmailUserInfo(ctx context.Context, accessToken *oauth2.Token) (*string, error) {
	ggClient := g.googleOAuthConfig.Client(ctx, accessToken)
	resp, err := ggClient.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		log.Error(ctx, "Error when get user info", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Error(ctx, "Error when get user info, status code: %d", resp.StatusCode)
		return nil, err
	}
	var userInfo dto.GoogleGetInfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Error(ctx, "Error when decode user info", err)
		return nil, err
	}
	return &userInfo.Email, nil
}
func (g *GmailProviderAdapter) buildMessage(from string, data *request.EmailDataDto) string {
	return "From: " + from + "\r\n" +
		"To: " + strings.Join(data.Tos, ", ") + "\r\n" +
		"Subject: " + data.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		data.Body
}
