package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptUseCase_EncryptTrackingID(t *testing.T) {
	trackingID := "46c37190-0921-49be-838c-8d78fb67cf5e"
	props := properties.EncryptProperties{EncryptTrackingIDKey: "dGVzdF9lbmNyeXB0X2tleQ=="}
	encryptUseCase := NewEncryptUseCase(&props)
	ctx := context.Background()
	encryptedTrackingID, err := encryptUseCase.EncryptTrackingID(ctx, trackingID)
	assert.NoError(t, err)
	fmt.Println(encryptedTrackingID)
}
