package domains

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrEmailAlreadyExists = errors.New("email address already exists")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidHash = errors.New("the encoded hash is not in the correct format")

type AccountRepository interface {
	FindByEmail(ctx context.Context, email string) (*Account, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Account, error)
	Create(ctx context.Context, account *Account) error
}

type Account struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
}

type argonParams struct {
	memory      uint32 // Память в КБ (64 MB)
	iterations  uint32 // Количество проходов
	parallelism uint8  // Количество потоков CPU
	saltLength  uint32 // Длина соли в байтах
	keyLength   uint32 // Длина итогового хэша в байтах
}

var p = argonParams{
	memory:      64 * 1024, // 64 MB
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

func NewAccount(email string, password string) (*Account, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         RoleUser,
	}, nil
}

func (self *Account) VerifyPassword(password string) error {
	parts := strings.Split(self.PasswordHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return ErrInvalidHash
	}

	params := &argonParams{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidHash
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidHash
	}
	params.keyLength = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) != 1 {
		return ErrInvalidPassword
	}

	return nil
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}
