package infrastructure

import (
	"context"
	"test-go/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOManager struct {
	Client     *minio.Client
	BucketName string
	Connected  bool
}

func NewMinIOManager(cfg config.MinIOConfig) (*MinIOManager, error) {
	if cfg.Endpoint == "" {
		return &MinIOManager{Connected: false}, nil
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return &MinIOManager{Connected: false}, err
	}

	// Basic check
	_, err = client.ListBuckets(context.Background())
	if err != nil {
		return &MinIOManager{Connected: false}, err
	}

	return &MinIOManager{
		Client:     client,
		BucketName: cfg.BucketName,
		Connected:  true,
	}, nil
}

func (m *MinIOManager) GetStatus() map[string]interface{} {
	if m == nil || !m.Connected {
		return map[string]interface{}{
			"connected": false,
			"error":     "Not configured or connection failed",
		}
	}

	// Get bucket usage (approximate via listing, simplified for now)
	// In production, you might use Prometheus or MinIO admin API, but for simple stats:
	ctx := context.Background()
	exists, err := m.Client.BucketExists(ctx, m.BucketName)
	if err != nil || !exists {
		return map[string]interface{}{
			"connected":   true,
			"bucket_name": m.BucketName,
			"status":      "Bucket not found",
		}
	}

	// Count objects (up to 1000 for quick check)
	objectCh := m.Client.ListObjects(ctx, m.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	count := 0
	var size int64 = 0
	for obj := range objectCh {
		if obj.Err == nil {
			count++
			size += obj.Size
		}
		if count >= 1000 {
			break // Limit for performance
		}
	}

	return map[string]interface{}{
		"connected":     true,
		"bucket_name":   m.BucketName,
		"object_count":  count,
		"total_size_kb": size / 1024,
		"status":        "Healthy",
		"endpoint":      m.Client.EndpointURL().String(),
	}
}
