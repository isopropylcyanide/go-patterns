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
| Or Channel | Multiplexes multiple channels not known in advance into one channel | Concurrency in Go |
| Semaphore Worker Pool | Restrict number of worker in the pool using semaphore | Ultimate Go Programming |
| Pipelines | Using channel to create pipelined stages | Concurrency in Go |


