package pipelines

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRudimentaryBatchPipeline(t *testing.T) {
	expectedOutput := []int{6, 10, 14, 18}
	output := RudimentaryBatchPipeline()
	assert.True(t, reflect.DeepEqual(output, expectedOutput))
}

func TestRudimentaryStreamPipeline(t *testing.T) {
	expectedOutput := []int{6, 10, 14, 18}
	output := RudimentaryStreamPipeline()
	assert.True(t, reflect.DeepEqual(output, expectedOutput))
}

func TestChannelStreamPipeline(t *testing.T) {
	done := make(chan interface{})
	// regardless of what stage a pipeline is in, closing done, will close it
	defer close(done)

	//expectedOutput := []int{6, 10, 14, 18}
	actualOutput := make([]int, 0)
	for v := range ChannelStreamPipeline(done) {
		fmt.Println("From channel ", v)
		actualOutput = append(actualOutput, v)
	}
	//assert.True(t, reflect.DeepEqual(actualOutput, expectedOutput))
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
	ch := make(chan interface{})
	for i := 0; i < b.N; i++ {
		ChannelStreamPipeline(ch)
	}
}
