package storage

import (
	"context"
	"cruise_booking_system/internal/data"
	"io"

	"github.com/minio/minio-go/v7"
)

func UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Ensure bucket exists
	exists, err := data.MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = data.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}

	info, err := data.MinioClient.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	// Return URL (presigned or public)
	// For now just return object name or path
	return info.Key, nil
}
