package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/g"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

// InitLogger 初始化日志
func InitLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var writer []io.Writer
	if g.Conf.GetBool("Log.Console.Enable") {
		level, err := zerolog.ParseLevel(g.Conf.GetString("Log.Console.Level"))
		if err != nil {
			panic("日志级别错误:" + err.Error())
		}
		// 配置控制台日志
		if strings.ToLower(g.Conf.GetString("Env")) == "development" {
			// 开发模式
			output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
			output.FormatFieldName = func(i interface{}) string {
				return fmt.Sprintf("\n%s: ", i)
			}
			output.FormatFieldValue = func(i interface{}) string {
				val := fmt.Sprintf("%s", i)
				if isStrictJSONObjectOrArray(val) {
					var un any
					_ = json.Unmarshal([]byte(val), &un)
					indent, _ := json.MarshalIndent(un, "", "	")
					val = "\n" + string(indent)
				}
				return val
			}

			writer = append(writer, &zerolog.FilteredLevelWriter{
				Writer: zerolog.LevelWriterAdapter{
					Writer: output,
				},
				Level: level,
			})
		} else {
			// 非开发模式
			writer = append(writer, &zerolog.FilteredLevelWriter{
				Writer: zerolog.LevelWriterAdapter{
					Writer: os.Stdout,
				},
				Level: level,
			})
		}
	}
	if g.Conf.GetBool("Log.File.Enable") {
		level, err := zerolog.ParseLevel(g.Conf.GetString("Log.File.Level"))
		if err != nil {
			panic("日志级别错误:" + err.Error())
		}
		file := &lumberjack.Logger{
			Filename:   g.Conf.GetString("Log.File.Filename"),
			MaxSize:    g.Conf.GetInt("Log.File.MaxSize"),
			MaxAge:     g.Conf.GetInt("Log.File.MaxAge"),
			MaxBackups: g.Conf.GetInt("Log.File.MaxBackups"),
			Compress:   g.Conf.GetBool("Log.File.Compress"),
		}
		writer = append(writer, &zerolog.FilteredLevelWriter{
			Writer: zerolog.LevelWriterAdapter{
				Writer: file,
			},
			Level: level,
		})
	}
	multi := zerolog.MultiLevelWriter(writer...)
	log.Logger = log.Hook(new(NotifyHook))
	log.Logger = log.Output(multi)
	log.Logger = log.With().Caller().Stack().Logger()
}

// 严格校验是否是json字符串
func isStrictJSONObjectOrArray(s string) bool {
	// 先验证是合法JSON
	if !json.Valid([]byte(s)) {
		return false
	}

	// 再检查第一个非空白字符
	trimmed := bytes.TrimLeft([]byte(s), " \t\n\r")
	if len(trimmed) == 0 {
		return false
	}

	// 仅当以 { 或 [ 开头时才算对象/数组
	return trimmed[0] == '{' || trimmed[0] == '['
}

type NotifyHook struct{}

func (h NotifyHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// notify
}
