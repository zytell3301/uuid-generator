package uuid_generator

import "testing"

// You can change  these number for benchmarking generator in different settings
var workerCount = 3
var bufferSize = 100

func BenchmarkGenerator_GenerateV4(b *testing.B) {
	generator, _ := NewGenerator("", bufferSize, workerCount)
	for i := 0; i < b.N; i++ {
		generator.GenerateV4()
	}
}
