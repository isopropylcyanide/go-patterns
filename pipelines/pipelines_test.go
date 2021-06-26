package main

import (
	"testing"
)

func BenchmarkRudimentaryBatch(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rudimentaryBatchPipeline()
	}
}

func BenchmarkRudimentaryStream(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rudimentaryStreamPipeline()
	}
}

func BenchmarkChannelStream(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		channelStreamPipeline()
	}
}
