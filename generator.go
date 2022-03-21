package uuid_generator

import (
	"github.com/google/uuid"
)

type Generator struct {
	space uuid.UUID

	v4Buffer     chan uuid.UUID
	v4StopSignal chan struct{}
	bufferSize   int
	workerCount  int
}

func NewGenerator(space string, bufferSize int, workerCount int) (*Generator, error) {
	generator := Generator{
		bufferSize:   bufferSize,
		workerCount:  workerCount,
		v4Buffer:     make(chan uuid.UUID, bufferSize),
		v4StopSignal: make(chan struct{}),
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
		go g.v4Generator()
	}
}

func (g Generator) v4Generator() {
	for {
		select {
		case <-g.v4StopSignal:
			return
		default:
			uuid, err := uuid.NewRandom()
			switch err == nil {
			case true:
				g.v4Buffer <- uuid
			}
		}
	}
}

// Increases workers count by given number
func (g Generator) IncreaseWorkersBy(count int) {
	for i := 0; i < count; i++ {
		go g.v4Generator()
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
	return &uuid
}
