# go-patterns
Various patterns and dabble in Golang.

## Sources
- [Concurrency in Go](https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/)
- [Hands on Go Programming](https://books.google.co.in/books/about/Hands_on_Go_Programming.html?id=Q3whEAAAQBAJ&redir_esc=y)
- [Ultimate Go Programming](https://learning.oreilly.com/videos/ultimate-go-programming)

| Pattern | Description | Source | 
| --- | --- | --- | 
| Confinement | Safe concurrent access without channels or mutex | Concurrency in Go | 
| For Select | Common idiom when dealing with selection on multiple channels | Concurrency in Go |
| Error Handling | Separate the concerns of error handling from a producer goroutine | Concurrency in Go |
| Goroutine leaks | You create a goroutine, you better ensure to stop it | Concurrency in Go |
| Channel Patterns | Patterns to multiplex multiple channels | Concurrency in Go |
| Semaphore Worker Pool | Restrict number of worker in the pool using semaphore | Ultimate Go Programming |
| Pipelines | Using channel to create pipelined stages | Concurrency in Go |
| Generators | Using channels to create memory efficient generators for pipelined stages | Concurrency in Go |
| Fan In/Out | Fanning pipeline stages in/out for performance and efficiency | Concurrency in Go |
| Error Propagation | Using wrapped high level errors to propagate errors | Concurrency in Go |
| Heartbeats | A way to signal health in concurrent life for waiting parties | Concurrency in Go |
| Replicated requests | A fault tolerant but expensive way to service requests faster | Concurrency in Go |
| Rate limiter | Constrain access to resource for a finite period for resiliency | Concurrency in Go |
| Healing Goroutines | Mechanism to restart or supervise long running goroutines | Concurrency in Go |


