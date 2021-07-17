package fan_in_fan_out

// Zen: Pipelines are elegant composable stages but they can be slow, very slow.
// In some situations we fan out to process the input from the stage above in parallel.
// This improves runtime of the stage overall and it is said to be fanned out.
// Requirements: The stage shouldn't rely on state/values that it has calculated before.
// Requirements: It takes a long time to run to warrant a fan-out.

// SlowPrimeNumberFinder demos a stage in which we try to find first 10 primes of a stream
// of random integers. Bound to be slow as this stage is processing it sequentially.
func SlowPrimeNumberFinder() {

}

func main() {
}
