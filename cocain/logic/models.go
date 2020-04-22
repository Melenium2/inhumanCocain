package gates

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Opaque struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Jwt       string    `json:"jwt" db:"jwt"`
	Opaque    string    `json:"opaque" db:"opaque"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewOpaque(userId int, jwt string) *Opaque {
	return &Opaque{
		UserId: userId,
		Jwt: jwt,
		Opaque: generateOpaqueToken(),
		CreatedAt: time.Now().UTC(),
	}
}

func generateOpaqueToken() string {
	salt := time.Now().UTC().String()
	hash := sha256.New().Sum([]byte(salt))
	return hex.EncodeToString(hash)
}

