package zap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapInstance *Zap

type Zap struct {
	Logger *zap.Logger
	config *Config
}

type Config struct {
	zapConfig  *zap.Config
	zapOptions []zap.Option
}

func New(config *zap.Config, options ...zap.Option) (*Zap, error) {
	if zapInstance != nil {
		return zapInstance, nil
	}
	zapInstance = &Zap{
		config: &Config{
			zapConfig:  config,
			zapOptions: options,
		},
	}
	if err := zapInstance.Init(); err != nil {
		return nil, err
	}
	return zapInstance, nil
}

func Update(config *zap.Config, options ...zap.Option) error {
	z := &Zap{
		config: &Config{
			zapConfig:  config,
			zapOptions: options,
		},
	}
	if err := z.Init(); err != nil {
		return err
	}
	zapInstance = z
	return nil
}

func (z *Zap) Init() error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	z.config.zapConfig.EncoderConfig = encoderConfig
	logger, err := z.config.zapConfig.Build()
	if err != nil {
		return err
	}

	z.config.zapOptions = append(z.config.zapOptions)
	z.Logger = logger.WithOptions(z.config.zapOptions...)

	return nil
}

func (z *Zap) GetLogger(ctx context.Context) (*zap.Logger, error) {
	var fields []zap.Field
	if tag := z.getTagFromContext(ctx); tag != "" {
		fields = append(fields, zap.String("tag", tag))
	}
	if reqID := z.getRequestIDFromContext(ctx); reqID != "" {
		fields = append(fields, zap.String("requestID", reqID))
	}
	if userID := z.getUserIDFromContext(ctx); userID != "" {
		fields = append(fields, zap.String("userID", userID))
	}
	if z.Logger == nil {
		if err := z.Init(); err != nil {
			return nil, err
		}
	}
	newLogger := z.Logger.With(fields...)

	return newLogger, nil
}

func (z *Zap) SetTagInContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, TagKey, tag)
}

func (z *Zap) getTagFromContext(ctx context.Context) string {
	if tag := ctx.Value(TagKey); tag != nil {
		if tagStr, ok := tag.(string); ok {
			return tagStr
		}
	}
	return MainTag
}

func (z *Zap) SetRequestIDInContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func (z *Zap) getRequestIDFromContext(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if requestID, ok := requestID.(string); ok {
			return requestID
		}
	}
	return ""
}

func (z *Zap) SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func (z *Zap) getUserIDFromContext(ctx context.Context) string {
	if userID := ctx.Value(UserIDKey); userID != nil {
		if userID, ok := userID.(string); ok {
			return userID
		}
	}
	return ""
}

func (z *Zap) Sync() error {
	return z.Logger.Sync()
}
