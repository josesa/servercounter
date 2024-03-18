## Moving Window Request Counter Server

This project implements a small HTTP server that counts requests on a moving window.

### To run

```
go run cmd/main.go
```

### Extensions

The application can, of course, be extended and adapted to different requirements.

- Possibility to trigger periodic save of counters instead of just on termination.
- Better map flush logic to ensure that a compatible ratio between flush and window size.
- Evaluate the possibility of using a sync.Map instead. It has the tradeoff of not using typed map but might have better performance in some situations.
- If multiple servers are to be used, an external counter, implemented using redis, can be used.
