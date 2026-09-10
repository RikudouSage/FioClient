#pragma once

#include <stddef.h>
#include <stdint.h>

typedef enum {
	FioClientSuccess = 0,
	FioClientFailure = 1,
} FioClientResult;

typedef uint64_t FioClientHandle;
typedef FioClientHandle FioClientContextHandle;
typedef FioClientHandle FioClientAccountHandle;

