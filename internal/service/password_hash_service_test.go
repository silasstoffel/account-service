package service_test

import (
	"testing"

	"github.com/silasstoffel/account-service/internal/service"
)

func TestCreateHash(t *testing.T) {
	t.Run("Should return a hash", func(t *testing.T) {
		hash, err := service.CreateHash("password")
		if err != nil {
			t.Errorf("CreateHash() error = %v", err)
			return
		}
		if len(hash) == 0 {
			t.Error("CreateHash() hash is empty")
		}
	})
}

func TestCompareHash(t *testing.T) {
	t.Parallel()
	t.Run("Should compare hash", func(t *testing.T) {
		t.Parallel()
		hash, err := service.CreateHash("password")
		if err != nil {
			t.Errorf("CreateHash() error = %v", err)
			return
		}

		err = service.CompareHash("password", hash)
		if err != nil {
			t.Errorf("CompareHash() error = %v", err)
		}
	})

	t.Run("Should return an error", func(t *testing.T) {
		t.Parallel()
		hash, _ := service.CreateHash("password-ok")
		err := service.CompareHash("password-no-ok", hash)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
