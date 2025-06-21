// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore

const (
	InfoBlobSuffix string = ".info"
)

const (
	maxRetryDelay = 5000
	retryDelay    = 100
	maxRetries    = 5
)

const (
	defaultUploadBlockSize = 8 * 1024 * 1024 // configure the upload buffer size
	defaultUploadBuffers   = 5               // configure the number of rotating buffers that are used when uploading (for degree of parallelism)
)
