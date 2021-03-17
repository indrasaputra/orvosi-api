package tool

import (
	"context"

	"github.com/indrasaputra/orvosi-api/entity"
	"google.golang.org/api/idtoken"
)

// IDTokenDecoder responsibles for decoding google's ID token.
type IDTokenDecoder struct {
	audience string
}

// NewIDTokenDecoder creates an instance of IDTokenDecoder.
func NewIDTokenDecoder(audience string) *IDTokenDecoder {
	return &IDTokenDecoder{
		audience: audience,
	}
}

// Decode decodes google token.
func (id *IDTokenDecoder) Decode(googleToken string) (*entity.User, *entity.Error) {
	payload, err := idtoken.Validate(context.Background(), googleToken, id.audience)
	if err != nil {
		return nil, entity.WrapError(entity.ErrInvalidGoogleToken, err.Error())
	}

	user := payloadToUser(payload)
	return user, nil
}

func payloadToUser(payload *idtoken.Payload) *entity.User {
	id, _ := payload.Claims["sub"].(string)
	name, _ := payload.Claims["name"].(string)
	email, _ := payload.Claims["email"].(string)

	return &entity.User{
		GoogleID: id,
		Name:     name,
		Email:    email,
	}
}
