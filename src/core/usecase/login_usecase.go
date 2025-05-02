package usecase

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golibs-starter/golib/log"
	"strconv"
	"time"
)

const (
	RsaPrivateKey    = "RSA PRIVATE KEY"
	RsaPublicKey     = "PUBLIC KEY"
	AccessTokenType  = "AccessToken"
	RefreshTokenType = "RefreshToken"
)

type ILoginUsecase interface {
	Login(ctx context.Context, email, password string) (*response.LoginResponseDto, error)
	GenerateTokenFromRefreshToken(ctx context.Context, refreshToken string) (*response.LoginResponseDto, error)
}
type LoginUsecase struct {
	getUserUseCase      IGetUserUseCase
	hashPasswordUseCase IHashPasswordUseCase
	authProps           *properties.AuthProperties
}

func (l *LoginUsecase) GenerateTokenFromRefreshToken(ctx context.Context, refreshToken string) (*response.LoginResponseDto, error) {
	claims := &dto.ClaimTokenDto{}
	// Parse the refresh token with public key
	publicKey := l.authProps.PublicKey
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != RsaPublicKey {
		log.Error(ctx, "[LoginUsecase] Invalid public key")
		return nil, errors.New("invalid public key")
	}
	publicKeyParsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error parsing public key: %v", err)
		return nil, err
	}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Error(ctx, "[LoginUsecase] Unexpected signing method: %v", token.Header["alg"])
			return nil, common.ErrInvalidToken
		}
		return publicKeyParsed, nil
	})
	if err != nil || !token.Valid {
		log.Error(ctx, "[LoginUsecase] Invalid or expired refresh token: %v", err)
		return nil, common.ErrInvalidToken
	}

	user, err := l.getUserUseCase.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error fetching user: %v", err)
		return nil, err
	}

	privateKey, err := l.generatePrivateKey(ctx)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error generating private key: %v", err)
		return nil, err
	}

	newAccessToken, err := l.generateToken(ctx, user, AccessTokenType, privateKey)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error generating access token: %v", err)
		return nil, err
	}

	return &response.LoginResponseDto{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		UserInfo: response.UserInfoDto{
			FullName: user.FullName,
			Email:    user.Email,
		},
	}, nil
}

func (l *LoginUsecase) Login(ctx context.Context, email, password string) (*response.LoginResponseDto, error) {
	//check if user exist
	user, err := l.getUserUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || l.hashPasswordUseCase.CompareHashAndPassword(user.Password, password) != nil {
		return nil, common.ErrEmailOrPasswordInvalid
	}
	//generate private key
	privateKey, err := l.generatePrivateKey(ctx)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error generating private key: %v", err)
		return nil, err
	}
	//generate access token
	accessToken, err := l.generateToken(ctx, user, AccessTokenType, privateKey)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error generating access token: %v", err)
		return nil, err
	}
	//generate refresh token
	refreshToken, err := l.generateToken(ctx, user, RefreshTokenType, privateKey)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error generating refresh token: %v", err)
		return nil, err
	}
	var rsp response.LoginResponseDto
	rsp.AccessToken = accessToken
	rsp.RefreshToken = refreshToken
	rsp.UserInfo = response.UserInfoDto{
		FullName: user.FullName,
		Email:    user.Email,
	}
	return &rsp, nil

}
func (l *LoginUsecase) generateToken(ctx context.Context, user *entity.UserEntity, tokenType string, privateKey *rsa.PrivateKey) (string, error) {
	var expiredTime int64
	if tokenType == AccessTokenType {
		expiredTime = l.authProps.ExpiredAccessToken
	} else {
		expiredTime = l.authProps.ExpiredRefreshToken
	}
	claimTokenDto := dto.ClaimTokenDto{
		UserId: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredTime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claimTokenDto)

	// Sign the token with the private key
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
}

func (l *LoginUsecase) generatePrivateKey(ctx context.Context) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(l.authProps.PrivateKey))
	if block == nil || block.Type != RsaPrivateKey {
		log.Error(ctx, "[LoginUsecase] Invalid private key")
		return nil, errors.New("invalid private key")
	}
	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error(ctx, "[LoginUsecase] Error parsing private key: %v", err)
		return nil, err
	}
	return privateKey, nil
}
func NewLoginUsecase(getUserUseCase IGetUserUseCase, hashPasswordUseCase IHashPasswordUseCase, authProps *properties.AuthProperties) ILoginUsecase {
	return &LoginUsecase{
		getUserUseCase:      getUserUseCase,
		hashPasswordUseCase: hashPasswordUseCase,
		authProps:           authProps,
	}
}
