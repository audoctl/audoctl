package id

import (
	"github.com/google/uuid"
)

type ID string

func New() ID {
	return ID(uuid.NewString())
}

func NewNil() ID {
	return ID("00000000-0000-0000-0000-000000000000")
}

func Parse(s string) (ID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}

	return ID(uid.String()), nil
}

func (i *ID) IsInvalid() bool {
	if i == nil {
		return true
	}
	uid, err := uuid.Parse(string(*i))
	if err != nil || uid == uuid.Nil {
		return true
	}

	return false
}

func (i *ID) Equal(id *ID) bool {
	if i == nil && id == nil {
		return true
	}
	if i == nil || id == nil {
		return false
	}

	return *i == *id
}

func (i *ID) EqualStr(id string) bool {
	if i == nil && len([]rune(id)) <= 0 {
		return true
	}
	if i == nil || len([]rune(id)) <= 0 {
		return false
	}

	return id == id
}

func (i *ID) String() string {
	return string(*i)
}
