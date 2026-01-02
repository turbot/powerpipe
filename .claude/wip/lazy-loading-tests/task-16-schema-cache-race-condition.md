# Task 16: Fix Schema Cache Race Condition in pipe-fittings

## Objective

Address the race condition in pipe-fittings' `parse.getResourceSchema()` function which uses a global map for schema caching without synchronization.

## Context

- Discovered during Task 10 (Concurrent Access Tests)
- The race occurs in `parse/schema.go` in the `getResourceSchema()` function
- A global map is used to cache HCL schemas but concurrent access is not synchronized
- This affects concurrent `GetResource()` calls in lazy loading mode
- Currently documented in `internal/workspace/concurrent_test.go`

## Race Detection Output

```
WARNING: DATA RACE
Write at 0x00c00056f0b0 by goroutine 83:
  runtime.mapaccess2_faststr()
      internal/runtime/maps/runtime_faststr_swiss.go:162 +0x29c
  github.com/turbot/pipe-fittings/v2/parse.getResourceSchema()
      /Users/nathan/src/pipe-fittings/parse/schema.go:44 +0x380
  github.com/turbot/pipe-fittings/v2/parse.DecodeHclBody()
      /Users/nathan/src/pipe-fittings/parse/decode_body.go:28 +0x15c
  github.com/turbot/powerpipe/internal/resourceloader.(*Loader).decodeResourceBlock()
      /Users/nathan/src/powerpipe/internal/resourceloader/parser.go:140 +0x180

Previous read at 0x00c00056f0b0 by goroutine 60:
  runtime.mapaccess1_faststr()
      internal/runtime/maps/runtime_faststr_swiss.go:103 +0x28c
  github.com/turbot/pipe-fittings/v2/parse.getResourceSchema()
      /Users/nathan/src/pipe-fittings/parse/schema.go:41 +0x32c
```

## Dependencies

- Requires changes to pipe-fittings repository (external dependency)
- Independent of Task 15 (viper race condition)

## Acceptance Criteria

- [ ] Add synchronization to schema cache in pipe-fittings
- [ ] `TestConcurrent_GetResource` passes with `-race` flag
- [ ] No performance regression for single-threaded schema lookups
- [ ] No breaking changes to pipe-fittings API

## Proposed Solution

### Option 1: Use sync.Map for Schema Cache
```go
// In parse/schema.go
var schemaCache sync.Map // map[string]*hcl.BodySchema

func getResourceSchema(resourceType string) *hcl.BodySchema {
    if schema, ok := schemaCache.Load(resourceType); ok {
        return schema.(*hcl.BodySchema)
    }

    // Build schema...
    newSchema := buildSchema(resourceType)
    schemaCache.Store(resourceType, newSchema)
    return newSchema
}
```
**Pros**: Thread-safe, optimized for read-heavy workloads
**Cons**: Slightly more memory overhead

### Option 2: Add RWMutex
```go
var (
    schemaCache   = make(map[string]*hcl.BodySchema)
    schemaCacheMu sync.RWMutex
)

func getResourceSchema(resourceType string) *hcl.BodySchema {
    schemaCacheMu.RLock()
    if schema, ok := schemaCache[resourceType]; ok {
        schemaCacheMu.RUnlock()
        return schema
    }
    schemaCacheMu.RUnlock()

    schemaCacheMu.Lock()
    defer schemaCacheMu.Unlock()
    // Double-check after acquiring write lock
    if schema, ok := schemaCache[resourceType]; ok {
        return schema
    }
    // Build and cache schema...
}
```
**Pros**: Standard Go pattern, explicit control
**Cons**: Slightly more complex double-check pattern

### Option 3: sync.Once Per Schema Type
```go
var schemaOnce = make(map[string]*sync.Once)
var schemaCache = make(map[string]*hcl.BodySchema)
var schemaMu sync.Mutex

func getResourceSchema(resourceType string) *hcl.BodySchema {
    schemaMu.Lock()
    once, ok := schemaOnce[resourceType]
    if !ok {
        once = &sync.Once{}
        schemaOnce[resourceType] = once
    }
    schemaMu.Unlock()

    once.Do(func() {
        schema := buildSchema(resourceType)
        schemaMu.Lock()
        schemaCache[resourceType] = schema
        schemaMu.Unlock()
    })

    schemaMu.Lock()
    defer schemaMu.Unlock()
    return schemaCache[resourceType]
}
```
**Pros**: Guarantees single initialization per type
**Cons**: More complex, likely overkill

## Recommended Approach

**Option 1 (sync.Map)** is recommended because:
1. Schema cache is read-heavy (schemas are cached once, read many times)
2. sync.Map is optimized for this pattern
3. Simple to implement with minimal code changes
4. No performance regression for the common case

## Files to Modify in pipe-fittings

- `parse/schema.go` - Contains `getResourceSchema()` and the global schema cache

## Implementation Steps

1. **Fork/Branch**: Create branch in pipe-fittings
2. **Locate**: Find schema cache declaration in `parse/schema.go`
3. **Refactor**: Change to `sync.Map`
4. **Test**: Add concurrent test for schema access
5. **Benchmark**: Verify no performance regression
6. **PR**: Submit pull request to pipe-fittings
7. **Update**: Once merged, update powerpipe's pipe-fittings dependency
8. **Verify**: Re-run `TestConcurrent_GetResource` with `-race` flag

## Priority

**Medium** - This race condition affects concurrent lazy loading in production scenarios where multiple goroutines may load resources simultaneously.

## Related

- Task 10: Concurrent Access Tests (discovered this issue)
- Task 15: Viper Race Condition Fix (separate race condition)

## Notes

- This is an upstream issue in pipe-fittings, not powerpipe
- The race can occur in production during concurrent dashboard server requests
- Fix will benefit all pipe-fittings consumers, not just powerpipe
