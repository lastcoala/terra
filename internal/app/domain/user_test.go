package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// ---------------------------------------------------------------------------
// ValidateEmail
// ---------------------------------------------------------------------------

func TestUser_ValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid simple", "user@example.com", false},
		{"valid with plus", "user+tag@mail.example.org", false},
		{"valid with dots", "first.last@sub.domain.io", false},
		{"missing @", "userexample.com", true},
		{"missing domain", "user@", true},
		{"missing TLD", "user@domain", true},
		{"empty string", "", true},
		{"double @", "user@@example.com", true},
		{"spaces", "user @example.com", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := User{Email: tc.email}
			err := u.ValidateEmail()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ValidateGender
// ---------------------------------------------------------------------------

func TestUser_ValidateGender(t *testing.T) {
	tests := []struct {
		name    string
		gender  string
		wantErr bool
	}{
		{"male lowercase", "male", false},
		{"female lowercase", "female", false},
		{"male uppercase", "Male", false},
		{"female uppercase", "FEMALE", false},
		{"mixed case", "fEmAlE", false},
		{"empty string", "", true},
		{"other value", "other", true},
		{"non-binary", "non-binary", true},
		{"numeric", "1", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := User{Gender: tc.gender}
			err := u.ValidateGender()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ValidatePassword
// ---------------------------------------------------------------------------

func TestUser_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "Secret1234", false},
		{"all requirements met", "Abcdef1!", false},
		{"too short", "Ab1", true},
		{"exactly 7 chars", "Secret1", true},
		{"missing uppercase", "secret1234", true},
		{"missing lowercase", "SECRET1234", true},
		{"missing digit", "SecretABC", true},
		{"empty string", "", true},
		{"spaces count as chars but missing digit", "Secret  ", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := User{Password: tc.password}
			err := u.ValidatePassword()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// EncryptPassword
// ---------------------------------------------------------------------------

func TestUser_EncryptPassword(t *testing.T) {
	t.Run("hashes the password and stores it", func(t *testing.T) {
		u := User{Password: "secret123"}
		original := u.Password

		err := u.EncryptPassword()
		require.NoError(t, err)

		// Password field must have changed
		assert.NotEqual(t, original, u.Password)

		// The hash must be a valid bcrypt hash of the original password
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(original))
		assert.NoError(t, err, "stored hash should match the original password")
	})

	t.Run("subsequent encrypt produces different hashes", func(t *testing.T) {
		u1 := User{Password: "secret123"}
		u2 := User{Password: "secret123"}

		require.NoError(t, u1.EncryptPassword())
		require.NoError(t, u2.EncryptPassword())

		// bcrypt salts are random – same input yields different ciphertext
		assert.NotEqual(t, u1.Password, u2.Password)
	})
}

// ---------------------------------------------------------------------------
// CheckPassword
// ---------------------------------------------------------------------------

func TestUser_CheckPassword(t *testing.T) {
	t.Run("correct password returns no error", func(t *testing.T) {
		u := User{Password: "myPassword!"}
		require.NoError(t, u.EncryptPassword())

		err := u.CheckPassword("myPassword!")
		assert.NoError(t, err)
	})

	t.Run("wrong password returns error", func(t *testing.T) {
		u := User{Password: "myPassword!"}
		require.NoError(t, u.EncryptPassword())

		err := u.CheckPassword("wrongPassword")
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// Validate
// ---------------------------------------------------------------------------

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    User{Email: "alice@example.com", Gender: "female", Password: "Secret123"},
			wantErr: false,
		},
		{
			name:    "invalid password",
			user:    User{Email: "alice@example.com", Gender: "female", Password: "weak"},
			wantErr: true,
		},
		{
			name:    "invalid email",
			user:    User{Email: "not-an-email", Gender: "male"},
			wantErr: true,
		},
		{
			name:    "invalid gender",
			user:    User{Email: "alice@example.com", Gender: "unknown"},
			wantErr: true,
		},
		{
			name:    "both invalid",
			user:    User{Email: "", Gender: ""},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate(true)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
