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

func (g Generator) GenerateV5(name string) string {
	return uuid.NewSHA1(g.space, []byte(name)).String()
}
