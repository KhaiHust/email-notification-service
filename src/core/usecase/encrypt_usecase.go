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
	EncryptVersionTemplate(ctx context.Context, version string) (string, error)
	DecryptVersionTemplate(ctx context.Context, version string) (string, error)
	EncryptAES(ctx context.Context, plainText string) (string, error)
	DecryptAES(ctx context.Context, cipherText string) (string, error)
	EncryptTrackingID(ctx context.Context, trackingID string) (string, error)
	DecryptTrackingID(ctx context.Context, trackingID string) (string, error)
	EncryptDataEmailRequest(ctx context.Context, data string) (string, error)
	DecryptDataEmailRequest(ctx context.Context, data string) (string, error)
}
type EncryptUseCase struct {
	props *properties.EncryptProperties
}

func (e EncryptUseCase) EncryptDataEmailRequest(ctx context.Context, data string) (string, error) {
	block, err := aes.NewCipher([]byte(e.props.EncryptDataEmailRequestKey))
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error creating AES cipher: %v", err)
		return "", err
	}
	return e.encrypt(ctx, data, block)
}

func (e EncryptUseCase) DecryptDataEmailRequest(ctx context.Context, data string) (string, error) {
	dataBytes, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error decoding base64: %v", err)
		return "", err
	}
	return e.decrypt(ctx, err, dataBytes)
}

func (e EncryptUseCase) EncryptTrackingID(ctx context.Context, trackingID string) (string, error) {
	block, err := aes.NewCipher([]byte(e.props.EncryptTrackingIDKey))
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error creating AES cipher: %v", err)
		return "", err
	}
	return e.encrypt(ctx, trackingID, block)
}

func (e EncryptUseCase) DecryptTrackingID(ctx context.Context, trackingID string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(trackingID)
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error decoding base64: %v", err)
		return "", err
	}
	return e.decrypt(ctx, err, data)
}

func (e EncryptUseCase) DecryptVersionTemplate(ctx context.Context, version string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(version)
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error decoding base64: %v", err)
		return "", err
	}
	return e.decrypt(ctx, err, data)
}

func (e EncryptUseCase) EncryptVersionTemplate(ctx context.Context, version string) (string, error) {
	block, err := aes.NewCipher([]byte(e.props.EncryptVersionTemplate))
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error creating AES cipher: %v", err)
		return "", err
	}
	return e.encrypt(ctx, version, block)
}

func (e EncryptUseCase) EncryptAES(ctx context.Context, plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(e.props.EncryptKey))
	if err != nil {
		log.Error(ctx, "[EncryptUseCase] Error creating AES cipher: %v", err)
		return "", err
	}
	return e.encrypt(ctx, plainText, block)
}

func (e EncryptUseCase) encrypt(ctx context.Context, plainText string, block cipher.Block) (string, error) {
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
	return e.decrypt(ctx, err, data)
}

func (e EncryptUseCase) decrypt(ctx context.Context, err error, data []byte) (string, error) {
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
