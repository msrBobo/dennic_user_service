package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BaseUseCase struct{}

func (u *BaseUseCase) Error(msg string, err error) error {
	if len(strings.TrimSpace(msg)) != 0 {
		return fmt.Errorf("%v: %w", msg, err)
	}
	return err
}

func (u *BaseUseCase) BeforeUpdRequest(updatedAt *string) {

	if *updatedAt != "" {
		*updatedAt = time.Now().UTC().String()
	}
}

func (u *BaseUseCase) BeforeCreateRequest(guid *string, createdAt *string) {
	if guid != nil {
		*guid = uuid.New().String()
	}

	if createdAt != nil {
		*createdAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	}
}

