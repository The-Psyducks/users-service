package auth

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func IsGoogleTokenValid(idToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("error initializing firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting firebase auth client: %w", err)
	}

	_, err = client.VerifyIDToken(ctx, idToken)
	if err != nil {
		if ok := auth.IsIDTokenInvalid(err); ok {
			return false, nil
		}
		return false, fmt.Errorf("error verificando el token de ID: %w", err)
	}

	return true, nil
}
