package usecase

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/golibs-starter/golib/log"
)

type IEncryptUseCase interface {
	EncryptAES(ctx context.Context, plainText string) (string, error)
	DecryptAES(ctx context.Context, cipherText string) (string, error)
}
type EncryptUseCase struct {
	props *properties.EncryptProperties
}

func (e EncryptUseCase) EncryptAES(ctx context.Context, plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(e.props.EncryptKey))
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error creating AES cipher: %v", err)
		return "", err
	}
	plainTextBytes := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		log.Error(ctx, "[EncryptUseCase] Error generating IV: %v", err)
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (e EncryptUseCase) DecryptAES(ctx context.Context, cipherText string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error decoding base64: %v", err)
		return "", err
	}
	block, err := aes.NewCipher([]byte(e.props.EncryptKey))
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error creating AES cipher: %v", err)
		return "", err
	}
	if len(data) < aes.BlockSize {
		log.Error(ctx, "[EncryptUseCase] Ciphertext too short")
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return string(data), nil
}

func NewEncryptUseCase(props *properties.EncryptProperties) IEncryptUseCase {
	return &EncryptUseCase{
		props: props,
	}
}
