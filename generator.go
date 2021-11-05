package Uuid

import (
	"github.com/google/uuid"
)

type Generator struct {
	space uuid.UUID
}

func NewGenerator(space string) (*Generator, error) {
	uuidSpace, err := uuid.Parse(space)

	switch err != nil {
	case true:
		return nil, err
	}

	return &Generator{space: uuidSpace}, nil
}

// This method treats as a hashing function. But be careful that
// same space and same name will result in same uuid. So the
// guarantee of the uniqueness of the generated UUIDs is application's
// responsibility
func (g Generator) GenerateV5(name string) *uuid.UUID {
	uuid := uuid.NewSHA1(g.space, []byte(name))
	return &uuid
}

// This method returns a fully random UUID (UUID v4).
// Returned UUID will be nil if an error occurred while
// generating the uuid
func (Generator) GenerateV4() (*uuid.UUID, error) {
	uuid, err := uuid.NewRandom()
	switch err != nil {
	case true:
		return nil, err
	}

	return &uuid, nil
}
