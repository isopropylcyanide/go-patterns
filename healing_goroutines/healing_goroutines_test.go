package healing_goroutines

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	goleak.VerifyTestMain(m)
}

func TestIrresponsibleWardWithoutSteward(t *testing.T) {
	done := make(chan interface{})
	mainDuration := 2 * time.Second

	time.AfterFunc(mainDuration, func() {
		log.Printf("main: I can't wait anylonger than %s. Halting\n", mainDuration)
		close(done)
	})
	// Here, we show the simplest case where no one monitors this goroutine.
	// The main should ideally give up after a while which should signal the ward to let go
	<-doIrresponsibleWork(done, mainDuration)
	assert.Equal(t, "1", haltCounter.String())
}

func TestIrresponsibleWardWithSteward(t *testing.T) {
	done := make(chan interface{})
	mainDuration := 6 * time.Second
	monitorDuration := 2 * time.Second
	haltCounter.Store(0)

	// here we add a monitoring function to run our goroutine
	workWithSteward := newSteward(monitorDuration, doIrresponsibleWork)
	time.AfterFunc(mainDuration, func() {
		log.Printf("main: I can't wait anylonger than %s. Halting\n", mainDuration)
		close(done)
	})

	for range workWithSteward(done, monitorDuration) {
	}
	assert.Equal(t, "3", haltCounter.String())
}
