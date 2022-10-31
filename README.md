# Crypto to Fiat converter

GRPC service written in Golang which converts crypto amount to fiat(currency) value

## Getting started

1. Generate proto files:
```
make proto-gen
```
2. Fetch dependencies:
```
go mod tidy
```
3. Run service(you should have go installed):
```
go run ./cmd/main.go
```

## Features

- Config, logger and Dependency injection
- GRPC as transport layer
- In-memory store with partial invalidation, check out Cache section for details
- Clean architecture(I haven't implement dedicated structs for each layer, but the most interchangeable layers has its own structures)
- Unit test example(`intenral/service/cache/memory/memory_test.go`)
- Linter
- Github actions pipeline

## Cache

### Straight forward solution

We can store everything, and we don't need to handle expiration, we just update all data periodically.
But this solution has two constrains:
- Storage size - lets assume we are unable to store everything in cache
- Cache update - when we have a large amount of data in cache, we need to figure out the update mechanism, 
- because we cannot lock everything during update

### Predefined cache plus single requests solution

We load predefined tokens and currencies in memory on app start and keep updating it periodically.
And each user request we also store in cache.

During user request, there is no need to wait for Upsert(Insert/Update) completion, we can send the result to user immediately and finish the Upsert in gorutine

#### Cache invalidation

All cache entries has expiration window and on request we check if data is expired, 
if yes, we request actual data from Price provider and store it in cache. 
This is valid for not predefined tokens and currencies, predefined one will be updated automatically.

##### Bottleneck

Sooner or later we have a risk of being run out of cache storage, because of data that was requested only once and stored in cache forever.
We can implement `Cleaner service` which will remove outdated data from cache.
We can use Queue for that, each cache entry goes also to Queue like this:
```
type QueueEntry struct {
	Token string
	Currency string
	ExirationDate time.Time
}
```
Cleaner service will go over the queue as long as Expiration date is in the past and removes from map data.

#### Predefined tokens and currencies

We can start from predefined list of tokens and currencies, 
but we can also collect requests metrics and cache actual most frequent tokens and currencies


## Pagination

Pagination is implemented only for `GetTokenList`.

## TODO
1. Retry-Backoff logic for Price Provider
2. Metrics collection for convert requests and composing most frequent tokens and currencies based on actual data
3. Cache invalidation for not predefined tokens and currencies
4. Docker
5. Integration tests