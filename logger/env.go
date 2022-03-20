package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	executionEnvironmentVariable = "ENV"

	k8sEnvironment  = "kubernetes"
	k8sServiceName  = "SERVICE_NAME"
	k8sVersion      = "VERSION"
	k8sNodeName     = "NODE_NAME"
	k8sPodName      = "POD_NAME"
	k8sPodNamespace = "POD_NAMESPACE"
	k8sPodIp        = "POD_IP"
	k8sPodSa        = "POD_SERVICE_ACCOUNT"

	logLevelEnvVariable = "LOG_LEVEL"
)

func NewWithEnv() *zap.Logger {
	logLevel, logLevelErr := getLogLevelFromEnv()
	logger := New(logLevel)
	if logLevelErr != nil {
		logger.Warn(
			"logger initialization with environment",
			zap.Error(logLevelErr),
		)
	}
	if os.Getenv(executionEnvironmentVariable) == k8sEnvironment {
		logger = setupKubernetesEnvironment(logger)
	}
	return logger
}

func setupKubernetesEnvironment(logger *zap.Logger) *zap.Logger {
	return logger.With(
		zap.String(strings.ToLower(k8sPodIp), os.Getenv(k8sPodIp)),
		zap.String(strings.ToLower(k8sNodeName), os.Getenv(k8sNodeName)),
		zap.String(strings.ToLower(k8sPodNamespace), os.Getenv(k8sPodNamespace)),
		zap.String(strings.ToLower(k8sPodName), os.Getenv(k8sPodName)),
		zap.String(strings.ToLower(k8sPodSa), os.Getenv(k8sPodSa)),
		zap.String(strings.ToLower(k8sVersion), os.Getenv(k8sVersion)),
		zap.String(strings.ToLower(k8sServiceName), os.Getenv(k8sServiceName)),
	)
}

func getLogLevelFromEnv() (zapcore.Level, error) {
	logLevelEnv := os.Getenv(logLevelEnvVariable)
	var logLevel zapcore.Level
	var err error
	if err = logLevel.UnmarshalText([]byte(logLevelEnv)); err != nil {
		logLevel = zapcore.DebugLevel
	}
	if logLevelEnv == "" {
		err = errors.New("log level is not set")
	}
	return logLevel, err
}
