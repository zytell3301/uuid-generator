package uuid_generator

import (
	"github.com/google/uuid"
)

type Generator struct {
	space uuid.UUID

	v4Buffer         chan uuid.UUID
	v4StopSignal     chan struct{}
	v4GenerateSignal chan struct{}
	bufferSize       int
	workerCount      int
}

func NewGenerator(space string, bufferSize int, workerCount int) (*Generator, error) {
	generator := Generator{
		bufferSize:       bufferSize,
		workerCount:      workerCount,
		v4Buffer:         make(chan uuid.UUID, bufferSize),
		v4GenerateSignal: make(chan struct{}),
		v4StopSignal:     make(chan struct{}),
	}
	generator.startV4Workers()
	switch space == "" {
	case true:
		return &generator, nil
	}

	uuidSpace, err := uuid.Parse(space)
	switch err != nil {
	case true:
		return nil, err
	}
	generator.space = uuidSpace
	return &generator, nil
}

func (g Generator) startV4Workers() {
	for i := 0; i < g.bufferSize; i++ {
		g.v4Generator()
	}
}

func (g Generator) refillBuffer() {
	for i := 0; i < (g.bufferSize - len(g.v4Buffer)); i++ {
		g.v4GenerateSignal <- struct{}{}
	}
}

func (g Generator) v4Generator() {
	for {
		select {
		case <-g.v4GenerateSignal:
			uuid, err := uuid.NewRandom()
			switch err == nil {
			case true:
				g.v4Buffer <- uuid
			}
		case <-g.v4StopSignal:
			return
		}
	}
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
func (g Generator) GenerateV4() *uuid.UUID {
	uuid := <-g.v4Buffer
	g.refillBuffer()
	return &uuid
}
