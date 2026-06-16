// Package dynamoclient provides a thread-safe singleton DynamoDB client.
package dynamoclient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ErrNotInitialized is returned by GetClient when Initialize has not been called.
var ErrNotInitialized = errors.New("dynamoclient: client not initialized; call Initialize first")

// Config holds the options for initializing the DynamoDB client.
type Config struct {
	// Region is the AWS region (e.g. "us-east-1").
	Region string

	// Production controls whether to connect to real AWS (true)
	// or to a local DynamoDB endpoint (false).
	Production bool

	// LocalEndpoint is the endpoint used when Production is false.
	// Defaults to "http://localhost:8000" if empty.
	LocalEndpoint string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig(prod bool) Config {
	return Config{
		Region:        os.Getenv("AWS_DEFAULT_REGION"),
		Production:    prod,
		LocalEndpoint: os.Getenv("DYNAMO_ENDPOINT"),
	}
}

// manager holds the singleton state.
type manager struct {
	mu     sync.RWMutex
	client *dynamodb.Client
	cfg    Config
}

var instance = &manager{}

// Initialize configures and creates the singleton DynamoDB client.
// It must be called once before any call to GetClient.
// Subsequent calls are no-ops unless Reset() was called first.
func Initialize(ctx context.Context, cfg Config) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.client != nil {
		return nil
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.LocalEndpoint == "" {
		cfg.LocalEndpoint = "http://localhost:8000"
	}

	var awsCfg aws.Config
	var err error

	if cfg.Production {
		// Em produção: usa credenciais reais da env/IAM normalmente
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
		)
	} else {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
			),
		)
	}

	if err != nil {
		return fmt.Errorf("dynamoclient: failed to load AWS config: %w", err)
	}

	if cfg.Production {
		instance.client = dynamodb.NewFromConfig(awsCfg)
	} else {
		instance.client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.LocalEndpoint)
		})
	}

	instance.cfg = cfg
	return nil
}

// GetClient returns the singleton *dynamodb.Client.
// Returns ErrNotInitialized if Initialize has not been called yet.
func GetClient() (*dynamodb.Client, error) {
	instance.mu.RLock()
	defer instance.mu.RUnlock()

	if instance.client == nil {
		return nil, ErrNotInitialized
	}
	return instance.client, nil
}

// IsInitialized reports whether the singleton has been initialized.
func IsInitialized() bool {
	instance.mu.RLock()
	defer instance.mu.RUnlock()
	return instance.client != nil
}

// CurrentConfig returns the Config that was used to initialize the client.
// Returns the zero value if the client has not been initialized.
func CurrentConfig() Config {
	instance.mu.RLock()
	defer instance.mu.RUnlock()
	return instance.cfg
}

// Reset destroys the singleton so that Initialize can be called again.
// Intended for testing; avoid in production code.
func Reset() {
	instance.mu.Lock()
	defer instance.mu.Unlock()
	instance.client = nil
	instance.cfg = Config{}
}
