# Querier

Source code for a take-home task assignment, written by H. Yazdani
from 27 January 2025, through 29th. It is sufficiently documented with
comments and explanations. A handful of test cases are present as well.

The project description is to write an Order processing worker. It
accepts a stream of incoming requests, add each Order to the queue, mock
the processing operation, and set a status for each Order.

*Order* consists of `OrderID`, `Priority`, `Status`, and
`ProcessingTime`.

*Priority* has two possible values, `Normal` and `High`. Orders with
higher value must be processed first, even if they're added later in
time.

*Status* can be in three states. `Pending`, `Processed`, or `Failed`.
The number of Orders with each status should be printed every
**two seconds**.

*ProcessingTime* determines the duration that Order processing should
take. If an Order takes longer than **5 seconds**, it should
automatically get cancelled with the *Failed* status. Successfully
processed Orders will have the *Processed* status. Incoming Orders get
the *Pending* state as default.

- The project was written purely in Golang.
- Flow of incoming requests is simulated by the `streamer` package.
- Dockerfile, and running instructions are provided.

## Run

In the root of the project, build the image from the Dockerfile by
running:

```shell
docker build -t querier .
```

When you have the image ready, run the container with:

```shell
docker run -v ./assets/cfg.yaml:/app/assets/cfg.yaml querier
```

This command binds the config file to the container. Feel free to change
the configurations and experiment with different values in different
scenarios.

Send an interrupt signal (usually with `ctrl + C`) to stop the
application.

## Approach

### Priority Queue

With an ongoing stream of incoming requests, using a queue for holding
the pending Orders is required. Since each Order has a priority assigned
to it, entries with higher priority must take precedence over the ones
with a lower priority.

One approach is to maintain two slices as two priority levels. As long
as the higher level queue is empty, retrieve data from the lower level
queue. While this approach definitely works, there are two major
downsides to it:

#### Rigidity

By defining different data structures for each priority level, we're
losing flexibility to support more priority levels in the future. Each
new level mandates the developers to add a new slice, modify the popping
strategy, and possibly other considerations.

#### Performance

Under a heavy load, queue's size will constantly fluctuate. In turn,
with slices as the underlying type, each size change will call for a new
allocation and its consequent data copying. It goes without saying, this
does not bode well for a performance critical application.

### Heap Tree

Fortunately for us, Go has a killer feature built-in. `container/heap`
package in the standard library provided a production ready framework
to work with.

A heap is a tree with the property that each node is the minimum-valued
node in its subtree. By implementing `heap.Interface`, Push adds items
while Pop removes the highest-priority item from the queue.

There are two caveats remaining, race conditions and type safety. The
queue needs to be protected from concurrent read and writes. Also,
accepting and retrieving values of type `any` is error-prone, and needs
to be safeguarded.

The `priorityQqueue` type is a low level implementation of a heap tree.
While the `Queue` type is an abstraction over priority queue, and
provides concurrency and type safety.

### Worker Pool

Worker pool is implemented to control the number of Orders being
processed at any given time. The number of active workers are
dynamically controlled, changing with the (possible) fluctuation of
incoming Order requests.

Each new Order gets its own goroutine to be processed in. As the number
of running goroutines increases, `max-jobs` variable (defined in the
configurations) will stop this number to grow indefinitely, and limit
the number of active workers.

By setting zero or a negative value to the `max-jobs`, number of active
goroutines **will not** be limited, and will roughly match the number of
simultaneous incoming Orders.

### Concurrency

By leveraging concurrency, we drastically improve performance. This gain
comes with a risk though, with bugs such as race conditions, deadlocks,
starvation, memory leakage, among other possible issues.

Concurrent read and writes are protected by mutex (mutual exclusion)
locks. By acquiring the lock, that goroutine is the only one with access
to the data at that moment, giving us the data integrity guarantee.

This exclusive access can be a double-edged sword, with mutexes
(forgetfully) not being unlocked, or long-running operations. In this
project, all lock operations are paired with a deferred unlock call, and
operations are kept short and concise as possible.

## Structure

Here is a tree of the project's structure.

```text
.
├── assets
│     └── cfg.yaml
├── cmd
│     └── querier
│         └── main.go
├── config
│     └── config.go
├── pkg
│     ├── model
│     │     ├── order.go
│     │     ├── priority.go
│     │     └── status.go
│     ├── order
│     │     ├── order.go
│     │     ├── order_test.go
│     │     ├── priority_queue.go
│     │     ├── priority_queue_test.go
│     │     ├── serde.go
│     │     ├── serde_test.go
│     │     ├── stats
│     │     │     └── stats.go
│     │     └── worker
│     │         ├── pool.go
│     │         ├── server.go
│     │         └── worker.go
│     └── streamer
│         └── streamer.go
├── Dockerfile
├── go.mod
├── go.sum
└── README.md

11 directories, 21 files
```

### `cmd/querier`

Hosts the main package, tasked with starting the application by getting
the configurations and calling other packages. Also, it listens for the
interrupt signal to handle graceful shutdown.

### `config`

Defines configuration structs, and parses the yaml config file.

### `pkg/model`

Describes the domain we're dealing with, including Order, Priority, and
Status. `model` does not involve itself with the application logic.

### `pkg/order`

Home to `priorityQueue` as a low level primitive, `Queue` as an
abstraction over the Order queue, and functions for efficiently encoding
and decoding the Order data.

Benchmark results comparing the custom encode and decoding methods,
versus the JSON serialize and deserializing functions:

```text
goos: linux
goarch: amd64
pkg: github.com/hossein1376/querier/pkg/order
cpu: AMD Ryzen 7 7800X3D 8-Core Processor           
BenchmarkEncode_Custom
BenchmarkEncode_Custom-16   1000000000      0.2104 ns/op    0 B/op      0 allocs/op
BenchmarkEncode_JSON
BenchmarkEncode_JSON-16     8524358         140.4 ns/op     96 B/op     1 allocs/op
BenchmarkDecode_Custom
BenchmarkDecode_Custom-16   1000000000      0.2042 ns/op    0 B/op      0 allocs/op
BenchmarkDecode_JSON
BenchmarkDecode_JSON-16     1626056         735.9 ns/op     240 B/op    5 allocs/op
```

### `pkg/order/stats`

A concurrent safe package for keeping track of pending, processed, and
failed Orders.

### `pkg/order/worker`

Being the deepest nested package in the project; worker pool, HTTP
server, and mocked Order processing function reside here. For each
incoming request, a new Order is enqueued and a subsequent signal is
transmitted. Then, a new goroutine is created to process the first item
in the queue.

### `pkg/streamer`

Mock incoming Order requests with random data.
