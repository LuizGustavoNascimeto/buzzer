package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type CloudWatchLogger struct {
	client        *cloudwatchlogs.Client
	logGroupName  string
	logStreamName string
	sequenceToken *string
	mu            sync.Mutex // protege sequenceToken em chamadas concorrentes
}

type LogEntry struct {
	Timestamp  time.Time     `json:"timestamp"`
	Level      string        `json:"level"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	ClientIP   string        `json:"client_ip"`
	Message    string        `json:"message,omitempty"`
	Error      string        `json:"error,omitempty"`
}

var (
	instance *CloudWatchLogger
	once     sync.Once
	initErr  error
)

type LoggerConfig struct {
	Region        string
	LogGroupName  string
	LogStreamName string
}

// GetInstance retorna a instância única do logger.
// Na primeira chamada, inicializa com o config fornecido.
// Nas chamadas seguintes, ignora o config e retorna a instância existente.
func GetInstance(cfg *LoggerConfig) (*CloudWatchLogger, error) {
	once.Do(func() {
		instance, initErr = newCloudWatchLogger(cfg)
	})
	return instance, initErr
}

// MustGetInstance é igual ao GetInstance, mas entra em pânico se houver erro.
// Ideal para ser usada no main, onde uma falha é irrecuperável.
func MustGetInstance(cfg *LoggerConfig) *CloudWatchLogger {
	logger, err := GetInstance(cfg)
	if err != nil {
		panic(fmt.Sprintf("falha ao inicializar CloudWatchLogger: %v", err))
	}
	return logger
}

// --- Construtor interno (não exportado) ---

func newCloudWatchLogger(cfg *LoggerConfig) (*CloudWatchLogger, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar config AWS: %w", err)
	}

	logger := &CloudWatchLogger{
		client:        cloudwatchlogs.NewFromConfig(awsCfg),
		logGroupName:  cfg.LogGroupName,
		logStreamName: cfg.LogStreamName,
	}

	if err := logger.ensureLogGroupAndStream(); err != nil {
		return nil, err
	}

	return logger, nil
}

// --- Métodos da instância ---

func (l *CloudWatchLogger) ensureLogGroupAndStream() error {
	ctx := context.TODO()

	l.client.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &l.logGroupName,
	})

	l.client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &l.logGroupName,
		LogStreamName: &l.logStreamName,
	})

	return nil
}

func (l *CloudWatchLogger) SendLog(entry LogEntry) error {
	message, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &l.logGroupName,
		LogStreamName: &l.logStreamName,
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(string(message)),
				Timestamp: aws.Int64(time.Now().UnixMilli()),
			},
		},
		SequenceToken: l.sequenceToken,
	}

	output, err := l.client.PutLogEvents(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("erro ao enviar log: %w", err)
	}

	l.sequenceToken = output.NextSequenceToken
	return nil
}
