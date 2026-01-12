// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/kopexa-grc/common/blob"
	"github.com/kopexa-grc/common/blob/s3store"
	"github.com/stretchr/testify/require"
)

func LoadS3TestConfig() s3store.S3Config {
	return s3store.S3Config{
		AccessKeyID:     "a",
		SecretAccessKey: "b",
		Region:          s3store.S3_DEFAULT_REGION,
		Endpoint:        "c",
		ContainerName:   "dev-blob-storage",
		UsePathStyle:    false,
	}
}

func LoadAzureTestConfig() blob.AzureConfig {
	return blob.AzureConfig{
		AccountName: "test-account",
		AccountKey:  "dGVzdC1rZXk=",
		Endpoint:    "https://test.blob.core.windows.net",
	}
}

func TestS3StoreUpload(t *testing.T) {
	t.Skip()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketProvider, err := blob.New(
		&blob.Config{
			S3:    LoadS3TestConfig(),
			Azure: LoadAzureTestConfig(),
		},
	)
	require.NoError(t, err, "failed to create bucket provider")

	publicBlob, err := bucketProvider.Public(blob.ProviderS3)
	require.NoError(t, err, "failed to get public blob")

	// Test data
	testKey := "test-upload.txt"
	testContent := []byte("Hello, S3!")
	testReader := bytes.NewReader(testContent)

	// Upload the file
	err = publicBlob.Upload(ctx, testKey, testReader, &blob.WriterOptions{
		ContentType: "text/plain",
	})
	require.NoError(t, err, "failed to upload file")

	// get signed url
	signedURL, err := publicBlob.SignedURL(ctx, testKey, &blob.SignedURLOptions{
		Expiry: 15 * time.Minute,
		Method: "GET",
	})
	require.NoError(t, err, "failed to get signed URL")
	require.NotEmpty(t, signedURL, "signed URL is empty")

	// download the file using Signed URL
	resp, err := http.Get(signedURL)
	require.NoError(t, err, "failed to download file using signed URL")
	defer resp.Body.Close()

	downloadedContent, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read downloaded content")
	require.Equal(t, testContent, downloadedContent, "downloaded content mismatch")

	// Cleanup
	err = publicBlob.Delete(ctx, testKey)
	require.NoError(t, err, "failed to delete test file")

	t.Logf("✓ Upload test passed for key: %s", testKey)
}

func TestS3UploadWithSignedURL(t *testing.T) {
	t.Skip()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketProvider, err := blob.New(
		&blob.Config{
			S3:    LoadS3TestConfig(),
			Azure: LoadAzureTestConfig(),
		},
	)
	require.NoError(t, err, "failed to create bucket provider")

	spaceBlob, err := bucketProvider.Space("abcdef", blob.ProviderS3)
	require.NoError(t, err, "failed to get space blob")

	// Test data
	testKey := "test-upload-signed-url.txt"
	testContent := []byte("Hello, S3 via Signed URL!")
	testReader := bytes.NewReader(testContent)

	// Get signed URL for upload
	signedURL, err := spaceBlob.SignedURL(ctx, testKey, &blob.SignedURLOptions{
		Expiry:      15 * time.Minute,
		Method:      "PUT",
		ContentType: "text/plain",
	})
	require.NoError(t, err, "failed to get signed URL for upload")
	require.NotEmpty(t, signedURL, "signed URL is empty")

	t.Log(signedURL)

	// Upload the file using Signed URL
	req, err := http.NewRequest(http.MethodPut, signedURL, testReader)
	require.NoError(t, err, "failed to create upload request")
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to upload file using signed URL")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status code from upload")

	// get signed url
	downloadSignedURL, err := spaceBlob.SignedURL(ctx, testKey, &blob.SignedURLOptions{
		Expiry: 15 * time.Minute,
		Method: "GET",
	})
	require.NoError(t, err, "failed to get signed URL")
	require.NotEmpty(t, downloadSignedURL, "signed URL is empty")

	// download the file using Signed URL
	resp, err = http.Get(downloadSignedURL)
	require.NoError(t, err, "failed to download file using signed URL")
	defer resp.Body.Close()

	downloadedContent, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read downloaded content")
	require.Equal(t, testContent, downloadedContent, "downloaded content mismatch")

	// Cleanup
	err = spaceBlob.Delete(ctx, testKey)
	require.NoError(t, err, "failed to delete test file")

	t.Logf("✓ Upload with Signed URL test passed for key: %s", testKey)
}

func TestS3TypedUpload(t *testing.T) {
	t.Skip()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketProvider, err := blob.New(
		&blob.Config{
			S3:    LoadS3TestConfig(),
			Azure: LoadAzureTestConfig(),
		},
	)
	require.NoError(t, err, "failed to create bucket provider")

	spaceBlob, err := bucketProvider.Space("abcdef", blob.ProviderS3)
	require.NoError(t, err, "failed to get space blob")

	// Test data
	testKey := "test-typed-upload.txt"
	testContent := []byte("Hello, Typed Upload!")
	testReader := bytes.NewReader(testContent)

	// Create typed writer
	writer, err := spaceBlob.NewWriter(ctx, testKey, &blob.WriterOptions{
		ContentType: "image/gif",
	})
	require.NoError(t, err, "failed to create typed writer")

	// Write data
	_, err = io.Copy(writer, testReader)
	require.NoError(t, err, "failed to write data to typed writer")

	// Close writer
	err = writer.Close()
	require.NoError(t, err, "failed to close typed writer")

	// Read back the uploaded content
	reader, err := spaceBlob.NewRangeReader(ctx, testKey, 0, -1, &blob.ReaderOptions{})
	require.NoError(t, err, "failed to create reader for uploaded file")
	require.Equal(t, "image/gif", reader.ContentType())

	// Cleanup
	err = spaceBlob.Delete(ctx, testKey)
	require.NoError(t, err, "failed to delete test file")

	t.Logf("✓ Typed Upload test passed for key: %s", testKey)
}

func TestS3Copy(t *testing.T) {
	t.Skip()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketProvider, err := blob.New(
		&blob.Config{
			S3:    LoadS3TestConfig(),
			Azure: LoadAzureTestConfig(),
		},
	)
	require.NoError(t, err, "failed to create bucket provider")

	spaceBlob, err := bucketProvider.Space("abcdef", blob.ProviderS3)
	require.NoError(t, err, "failed to get space blob")

	// Test data
	testKey := "test-to-copy-upload.txt"
	testContent := []byte("Hello, S3 Copy!")
	testReader := bytes.NewReader(testContent)

	// Upload the file
	err = spaceBlob.Upload(ctx, testKey, testReader, &blob.WriterOptions{
		ContentType: "text/plain",
	})
	require.NoError(t, err, "failed to upload file")

	// Copy the file
	copiedKey := "test-copy.txt"
	err = spaceBlob.Copy(ctx, testKey, copiedKey, &blob.CopyOptions{})
	require.NoError(t, err, "failed to copy file")

	signedURL, err := spaceBlob.SignedURL(ctx, copiedKey, &blob.SignedURLOptions{
		Expiry: 15 * time.Minute,
		Method: "GET",
	})
	require.NoError(t, err, "failed to get signed URL")
	require.NotEmpty(t, signedURL, "signed URL is empty")

	// download the copied file using Signed URL
	resp, err := http.Get(signedURL)
	require.NoError(t, err, "failed to download file using signed URL")
	defer resp.Body.Close()

	downloadedContent, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read downloaded content")
	require.Equal(t, testContent, downloadedContent, "downloaded content mismatch")

	// Cleanup
	err = spaceBlob.Delete(ctx, testKey)
	require.NoError(t, err, "failed to delete test file")

	err = spaceBlob.Delete(ctx, copiedKey)
	require.NoError(t, err, "failed to delete copied test file")
}

func TestS3RangeDownload(t *testing.T) {
	t.Skip()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bucketProvider, err := blob.New(
		&blob.Config{
			S3:    LoadS3TestConfig(),
			Azure: LoadAzureTestConfig(),
		},
	)
	require.NoError(t, err, "failed to create bucket provider")

	spaceBlob, err := bucketProvider.Space("abcdef", blob.ProviderS3)
	require.NoError(t, err, "failed to get space blob")

	// Test data
	testKey := "test-range-download.txt"
	testContent := []byte("Hello, S3 Range Download!")
	testReader := bytes.NewReader(testContent)

	// Upload the file
	err = spaceBlob.Upload(ctx, testKey, testReader, &blob.WriterOptions{
		ContentType: "text/plain",
	})
	require.NoError(t, err, "failed to upload file")

	// Range download
	offset := int64(7)
	length := int64(12) // "S3 Range Do"
	rangeReader, err := spaceBlob.NewRangeReader(ctx, testKey, offset, length, &blob.ReaderOptions{})
	require.NoError(t, err, "failed to create range reader")

	var rc io.ReadCloser
	rangeReader.As(&rc)

	downloadedContent, err := io.ReadAll(rc)
	require.NoError(t, err, "failed to read range downloaded content")
	expectedContent := testContent[offset : offset+length]
	require.Equal(t, expectedContent, downloadedContent, "range downloaded content mismatch")

	// Cleanup
	err = spaceBlob.Delete(ctx, testKey)
	require.NoError(t, err, "failed to delete test file")

	t.Logf("✓ Range Download test passed for key: %s", testKey)
}
