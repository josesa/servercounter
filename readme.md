## Moving Window Request Counter Server

This project implements a small HTTP server that counts requests on a moving window.

### Configuration

The application has several parameters that can be configured according to requirements. The configurations are loaded via Environment Variables.

```bash
SC_STORAGEPATH # Path for the file where the application can load and save the counter values. (Default: data.txt)
SC_FLUSHINTERVALSECONDS # Number of seconds between counter map cleanup. (Default: 60 * 1.1)
SC_WINDOWSIZESECONDS # Number of seconds for the moving window counter. (Default: 60)
SC_ADDRESS # Address for the HTTP server. (Default: :8080)
```

### To run

```bash
go run cmd/main.go
```

Server is accessible on
http://localhost:8080/request

When hitting the page, the hit counter is increased and the current count is displayed.

### Extensions

The application can, of course, be extended and adapted to different requirements.

- Possibility to trigger periodic save of counters instead of just on termination.
- Better map flush logic to ensure that a compatible ratio between flush and window size.
- Evaluate the possibility of using a sync.Map instead. It has the tradeoff of not using typed map but might have better performance in some situations.
- If multiple servers are to be used, it is more suitable to use an external counter. For example by using Redis.
- Transform the service into a middleware so it could capture hits on all requests. Current count could be added to the Response Headers.
