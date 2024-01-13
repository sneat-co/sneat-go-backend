package dto4userus

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/strongo/validation"
	"slices"
	"strconv"
	"strings"
)

type TelegramAuthData struct {
	ID        int64  `json:"id"`
	AuthDate  int    `json:"auth_date"`
	Username  string `json:"username,omitempty"`
	Hash      string `json:"hash"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
}

func (v TelegramAuthData) Validate() error {
	if v.ID == 0 {
		return validation.NewErrRequestIsMissingRequiredField("id")
	}
	if v.AuthDate == 0 {
		return validation.NewErrRequestIsMissingRequiredField("auth_date")
	}
	if v.Hash == "" {
		return validation.NewErrRequestIsMissingRequiredField("hash")
	}
	if _, err := hex.DecodeString(v.Hash); err != nil {
		return fmt.Errorf("error decoding hash: %w", err)
	}
	return nil
}

func (v TelegramAuthData) String() string {
	return fmt.Sprintf("TelegramAuthData{ID:%d, Username:%s, FirstName:%s, LastName:%s, AuthDate:%d, PhotoURL:%s}",
		v.ID, v.Username, v.FirstName, v.LastName, v.AuthDate, v.PhotoURL,
	)
}

func (v TelegramAuthData) DataCheckString() string {
	values := make([]string, 0, 7)
	values = append(values, "auth_date="+strconv.Itoa(v.AuthDate))
	values = append(values, "id="+strconv.FormatInt(v.ID, 10))

	if v.Username != "" {
		values = append(values, "username="+v.Username)
	}
	if v.FirstName != "" {
		values = append(values, "first_name="+v.FirstName)
	}
	if v.LastName != "" {
		values = append(values, "last_name="+v.LastName)
	}
	if v.PhotoURL != "" {
		values = append(values, "photo_url="+v.PhotoURL)
	}
	slices.Sort(values)
	return strings.Join(values, "\n")
}

// IsHashMatchesData validates SHA256 of v.DataCheckString() matches v.Hash
func (v TelegramAuthData) IsHashMatchesData(secretKey []byte) (bool, error) {
	// Convert the provided hash to bytes
	// Do it as a first step as if hash is invalid we don't need to compute HMAC
	hashBytes, err := hex.DecodeString(v.Hash)
	if err != nil {
		return false, fmt.Errorf("error decoding hash: %w", err)
	}

	// Create an HMAC with SHA256 hash function using the secret key
	h := hmac.New(sha256.New, secretKey)

	dataCheckString := v.DataCheckString()
	if _, err = h.Write([]byte(dataCheckString)); err != nil {
		return false, errors.New("error writing data to HMAC")
	}

	computedHash := h.Sum(nil)

	// Compare the computed hash with the provided hash
	if hmac.Equal(hashBytes, computedHash) {
		return true, nil
	}

	return false, nil
}
