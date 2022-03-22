package uuid_generator

import (
	"crypto/rand"
	"github.com/google/uuid"
	"time"
)

type Generator struct {
	space uuid.UUID

	v4Buffer                chan uuid.UUID
	v4StopSignal            chan struct{}
	bufferSize              int
	workerCount             int
	useCustomReader         bool
	readerCheckInterval     int
	readerCheckerStopSignal chan struct{}
}

func NewGenerator(space string, bufferSize int, workerCount int, readerCheckInterval int) (*Generator, error) {
	generator := Generator{
		bufferSize:          bufferSize,
		workerCount:         workerCount,
		v4Buffer:            make(chan uuid.UUID, bufferSize),
		v4StopSignal:        make(chan struct{}),
		readerCheckInterval: readerCheckInterval,
	}
	generator.startV4Workers()
	generator.StartReaderChecker(generator.readerCheckInterval)
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
	for i := 0; i < g.workerCount; i++ {
		go g.v4Generator()
	}
}

// Checker won't get started if a non positive interval supplied.
// Deactivate checker if it is causing any issues like performance issues or etc.
func (g Generator) checkReaderAvailability() {
	switch g.readerCheckInterval <= 0 {
	case true:
		return
	}

	for {
		select {
		case <-g.readerCheckerStopSignal:
			g.readerCheckInterval = 0
			return
		default:
			_, err := rand.Reader.Read(make([]byte, 1))
			switch err != nil {
			case true:
				g.useCustomReader = true
			default:
				g.useCustomReader = false
			}
		}
		time.Sleep(time.Duration(g.readerCheckInterval) * time.Second)
	}
}

// If you notice any issues with checker, you can easily stop it by using this method
func (g Generator) StopReaderChecker() error {
	switch g.readerCheckInterval <= 0 {
	case true:
		return CheckerAlreadyStoppedError
	}
	g.readerCheckerStopSignal <- struct{}{}
	return nil
}

// If you haven't started reader checker on initialization or want to
// start the checker after stopping it, you can use this method instead of
// creating another instance of Generator
func (g Generator) StartReaderChecker(interval int) error {
	switch interval <= 0 {
	case true:
		return InvalidCheckerIntervalSuppliedError
	}
	switch g.readerCheckInterval > 0 {
	case true:
		return CheckerAlreadyStartedError
	}
	g.readerCheckInterval = interval
	go g.checkReaderAvailability()
	return nil
}

// Use this method to change checking interval instead of creating another
// Generator instance
func (g Generator) SetCheckerInterval(interval int) error {
	switch g.readerCheckInterval <= 0 {
	case true:
		return CheckerAlreadyStoppedError
	}
	g.readerCheckInterval = interval
	return nil
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
	g.workerCount = g.workerCount + count
}

// Decreases workers count by given number
func (g Generator) DecreaseWorkersBy(count int) {
	for i := 0; i < count; i++ {
		g.v4StopSignal <- struct{}{}
	}
	g.workerCount = g.workerCount - count
}

// Sets worker count to given number
func (g Generator) ChangeWorkerCount(count int) {
	switch g.workerCount-count < 0 {
	case true:
		g.IncreaseWorkersBy(count - g.workerCount)
		break
	default:
		g.DecreaseWorkersBy(g.workerCount - count)
	}
}

// PAY ATTENTION THAT ANY CHANGES TO BUFFER SIZE WILL CAUSE THE LOSS OF
// GENERATED IDS AND RESTARTING ALL WORKERS.
// THIS WILL CAUSE ALL YOU CODES THAT ARE REQUESTING FOR ID TO BE BLOCKED
// FOR A FEW MILLIS

// Increase generator's buffer size by given number.
func (g Generator) IncreaseBufferSizeBy(size int) {
	g.SetBufferSize(g.bufferSize + size)
}

// Decrease generator's buffer size by given number
func (g Generator) DecreaseBufferSizeBy(size int) {
	g.SetBufferSize(g.bufferSize - size)
}

// Sets buffer size to the given number
func (g Generator) SetBufferSize(size int) {
	g.stopV4Workers()
	g.bufferSize = size
	g.v4Buffer = make(chan uuid.UUID, g.bufferSize)
	g.startV4Workers()
}

func (g Generator) stopV4Workers() {
	for i := 0; i < g.workerCount; i++ {
		g.v4StopSignal <- struct{}{}
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
