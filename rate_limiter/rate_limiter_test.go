package rate_limiter

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/goleak"
)

func TestRateLimitedAPIDriver(t *testing.T) {
	requestCount := 8
	cases := []struct {
		name                 string
		replenishRate        time.Duration
		burstSize            int
		requestWaitPeriodMs  time.Duration // how long requests may wait, -1 for infinite
		expectedCountSuccess int
		expectedCountTimeout int
	}{
		{
			name: `Replenish every 2 seconds with a 5 burst. First 5 requests should
				 be instantaneous. Rest 3 should be rate limited. No failures expected`,
			replenishRate:        time.Second,
			burstSize:            5,  // burst < max work which means some requests would wait
			requestWaitPeriodMs:  -1, // requests wait indefinitely, rate limited but successful
			expectedCountSuccess: 8,
			expectedCountTimeout: 0,
		},
		{
			name: `Replenish every second with a 500 burst. This simulates a no rate limiter
					behaviour as number of requests are lower. No failures expected either`,
			replenishRate:        time.Second,
			burstSize:            1000, // burst > max work, which means no waiting
			requestWaitPeriodMs:  -1,   // requests wait indefinitely
			expectedCountSuccess: 8,
			expectedCountTimeout: 0,
		},
		{
			name: `Replenish every seconds with a 5 burst. First 5 requests should
				 be instantaneous. Rest 3 should be rate limited and fail as they cannot wait.`,
			replenishRate:        time.Second,
			burstSize:            5,                      // burst < max work which means some requests would wait
			requestWaitPeriodMs:  time.Millisecond * 800, // requests wait a max of 100ms which is less than burst
			expectedCountSuccess: 5,
			expectedCountTimeout: 3, // 3 requests should fail as they weren't willing to wait
		},
	}

	for _, tc := range cases {
		conn := Open(tc.replenishRate, tc.burstSize)
		var wg sync.WaitGroup
		wg.Add(requestCount) // total request count is fixed
		log.Printf(tc.name)
		nSuccess, nTimeout := 0, 0

		for i := 0; i < requestCount; i++ {
			var cancelFunc context.CancelFunc
			go func() {
				defer wg.Done()
				var ctx context.Context // boiler plate to create context
				ctx = context.Background()
				if tc.requestWaitPeriodMs == -1 {
				} else {
					ctx, cancelFunc = context.WithTimeout(context.Background(), tc.requestWaitPeriodMs)
					defer cancelFunc()
				}

				if err := conn.ReadFile(ctx); err != nil {
					assert.Errorf(t, err, "would exceed context deadline")
					nTimeout += 1
				} else {
					log.Printf("Readfile")
					nSuccess += 1
				}
			}()
		}
		wg.Wait()
		assert.Equal(t, tc.expectedCountSuccess, nSuccess)
		assert.Equal(t, tc.expectedCountTimeout, nTimeout)
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
