package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"lintang/go_producer/biz/domain"
	"lintang/go_producer/config"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/natefinch/lumberjack"

	hertzzap "github.com/hertz-contrib/logger/zap"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

// pake hertzlogger gak kayak pake uber/zap logger beneran
func InitZapLogger(cfg *config.Config) *hertzzap.Logger {
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.EncodeDuration = zapcore.SecondsDurationEncoder
	productionCfg.EncodeCaller = zapcore.ShortCallerEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// log encooder (json for prod, console for dev)
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	// loglevel
	logDevLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	logLevelProd := zap.NewAtomicLevelAt(zap.InfoLevel)

	//write sycer
	writeSyncerStdout, writeSyncerFile := getLogWriter(cfg.MaxBackups, cfg.MaxAge)

	prodCfg := hertzzap.CoreConfig{
		Enc: fileEncoder,
		Ws:  writeSyncerFile,
		Lvl: logLevelProd,
	}

	devCfg := hertzzap.CoreConfig{
		Enc: consoleEncoder,
		Ws:  writeSyncerStdout,
		Lvl: logDevLevel,
	}
	logsCores := []hertzzap.CoreConfig{
		prodCfg,
		devCfg,
	}
	coreConsole := zapcore.NewCore(consoleEncoder, writeSyncerStdout, logDevLevel)
	coreFile := zapcore.NewCore(fileEncoder, writeSyncerFile, logLevelProd)
	core := zapcore.NewTee(
		coreConsole,
		coreFile,
	)
	lg = zap.New(core)
	zap.ReplaceGlobals(lg)

	prodAndDevLogger := hertzzap.NewLogger(hertzzap.WithZapOptions(zap.WithFatalHook(zapcore.WriteThenPanic)),
		hertzzap.WithCores(logsCores...))

	return prodAndDevLogger
}

func getLogWriter(maxBackup, maxAge int) (writeSyncerStdout zapcore.WriteSyncer, writeSyncerFile zapcore.WriteSyncer) {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename: "./logs/app.log",

		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	})
	stdout := zapcore.AddSync(os.Stdout)

	return stdout, file
}

type ValidateError struct {
	ErrType, FailField, Msg string
}

// Error implements error interface.
func (e *ValidateError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.ErrType + ": expr_path=" + e.FailField + ", cause=invalid"
}

type BindError struct {
	ErrType, FailField, Msg string
}

// Error implements error interface.
func (e *BindError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.ErrType + ": expr_path=" + e.FailField + ", cause=invalid"
}

func CreateCustomValidationError() *binding.ValidateConfig {
	validateConfig := &binding.ValidateConfig{}
	validateConfig.SetValidatorErrorFactory(func(failField, msg string) error {
		err := ValidateError{
			ErrType:   "validateErr",
			FailField: "[validateFailField]: " + failField,
			Msg:       msg,
		}

		return &err
	})
	return validateConfig
}

func AccessLog() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Request.URI().Path()[:])
		query := string(ctx.Request.URI().QueryString()[:])
		ctx.Next(c)
		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", ctx.Response.StatusCode()),
			zap.String("method", string(ctx.Request.Header.Method())),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", string(ctx.Host())),
			zap.String("errors", ctx.Errors.String()),
			zap.Duration("cost", cost),
		)
	}
}

type JobReq struct {
	Name           string            `json:"name"`
	DisplayName    string            `json:"displayname"`
	Schedule       string            `json:"schedule"`
	Timezone       string            `json:"timezone"`
	Owner          string            `json:"owner"`
	OwnerEmail     string            `json:"owner_email"`
	Disabled       bool              `json:"disabled"`
	Concurrency    string            `json:"concurrency"`
	Executor       string            `json:"executor"`
	ExecutorConfig map[string]string `json:"executor_config"`
}

func InstallCURLInDkron() error {
	at := time.Now().Add(2 * time.Second)

	payload, err := json.Marshal(JobReq{
		Name:        "insatll curl",
		DisplayName: "insatll curl",
		Schedule:    fmt.Sprintf("@at " + at.Format(time.RFC3339)),
		Timezone:    "Asia/Jakarta",
		Owner:       "lintang birda saputra",
		OwnerEmail:  "lintangbirdasaputra23@gmail.com",
		Disabled:    false,
		Concurrency: "allow",
		Executor:    "shell",
		ExecutorConfig: map[string]string{
			"command": `sh /curl/curl.sh'`,
		},
	})
	if err != nil {
		zap.L().Error("Marshal JSON", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	req, err := http.NewRequest("POST", "http://dkron:8080/v1/jobs", bytes.NewBuffer(payload))

	if err != nil {
		zap.L().Error("NewRequest ", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("client.Do(req) ", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	defer resp.Body.Close()
	zap.L().Info("Successfully installl curl in dkron")
	return nil
}
