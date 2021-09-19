package queuing

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	g "patterns/generators"
	"testing"
)

// Zen: The true utility of queueing is to decouple stages so that the runtime
// of one stage has no impact on the runtime of another. This can lead to a cascading
// effect which can be good/bad depending on the system. Queuing trades utilization
// for lag. It does increase performance but only in applicable situations
// The only applicable systems are
// 	  - A1) If batching requests in a stage saves time
//    - A2) If delays in a stage produce a feedback loop into the system
// Queueing should ideally be implemented
//    - S1) At the entrance to your pipeline
//    - S2) In stages where batching will lead to higher efficiency

func BenchmarkUnbufferedWrite(b *testing.B) {
	performSimpleWrite(b, tmpFileOrFatal())
}

// Bufio queues writes internally into a buffer and implements queuing
// This satisfies the condition A1 as batching leads to performance
func BenchmarkBufferedWrite(b *testing.B) {
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performSimpleWrite(b, bufferedFile)
}

func performSimpleWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for bt := range g.Take(done, g.Repeat(done, byte(0)), b.N) {
		_, _ = writer.Write([]byte{bt.(byte)})
	}

}

func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return file
}
