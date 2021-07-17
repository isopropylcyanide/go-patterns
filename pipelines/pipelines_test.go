package pipelines

import (
	"testing"
)

func TestRudimentaryBatchPipeline(t *testing.T) {
	RudimentaryBatchPipeline()
}

func TestRudimentaryStreamPipeline(t *testing.T) {
	RudimentaryStreamPipeline()
}

func TestChannelStreamPipeline(t *testing.T) {
	ChannelStreamPipeline()
}

func BenchmarkRudimentaryBatch(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RudimentaryBatchPipeline()
	}
}

func BenchmarkRudimentaryStream(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RudimentaryStreamPipeline()
	}
}

func BenchmarkChannelStream(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChannelStreamPipeline()
	}
}
