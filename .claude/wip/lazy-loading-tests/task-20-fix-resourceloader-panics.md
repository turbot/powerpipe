# Task 20: Fix Resource Loader Test Panics

## Status: Complete ✅

## Priority: HIGH

## Problem

Resource loader error handling tests were panicking with nil pointer dereference:

```
--- FAIL: TestErrorHandling_MissingDirectory (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference
    /Users/nathan/src/powerpipe/internal/resourceloader/error_test.go:33 +0x78
```

## Solution

The error_test.go file was completely rewritten with proper error handling tests that:

1. **Use correct API** - `createTestLoader()` helper properly initializes all required components
2. **Handle nil returns** - All tests check for errors before accessing returned values
3. **Test panic prevention** - Dedicated `TestError_NoPanic*` tests verify no panics occur
4. **Cover all error scenarios** - 20 comprehensive error handling tests

## Tests Verified (All Passing)

- `TestError_LoadNonexistentResource` - Non-existent resource returns error
- `TestError_LoadInvalidResourceName` - Invalid names return errors
- `TestError_LoadWrongType` - Type mismatches handled correctly
- `TestError_LoadBenchmarkWrongType` - Benchmark type assertion safe
- `TestError_FileNotFound` - Missing files return errors
- `TestError_PreloadWithSomeFailures` - Preload fails on any error
- `TestError_PreloadWithErrorCallback` - Error callbacks work
- `TestError_PreloadContextCancellation` - Context cancellation handled
- `TestError_CircularDependencyDetection` - Cycles detected
- `TestError_CircularDependencyInGetDependencyOrder` - Order fails on cycles
- `TestError_MissingDependency` - Missing deps handled gracefully
- `TestError_ResolverEmptyIndex` - Empty index doesn't panic
- `TestError_ResolverConcurrentAccess` - Thread-safe operations
- `TestError_NoPanicOnNilCache` - Cache always initialized
- `TestError_NoPanicOnBadIndex` - Bad index entries don't panic
- `TestError_NoPanicOnConcurrentLoad` - Concurrent loads safe
- `TestError_CacheRecoveryAfterClear` - Cache clear recoverable
- `TestError_CacheInvalidateNonexistent` - Invalid cache keys safe
- `TestError_MessageContainsResourceName` - Errors include context
- `TestError_DependencyOrderContainsCycleInfo` - Cycle errors informative

## Verification

```bash
go test ./internal/resourceloader/... -run "TestError_" -v
# PASS - all 20 tests pass
```

## Completed Tasks

1. [x] Review error_test.go and identify the nil pointer issue
2. [x] Find the correct API for resource loading error scenarios
3. [x] Fix tests to handle nil returns or use correct API
4. [x] Verify all error handling tests pass
