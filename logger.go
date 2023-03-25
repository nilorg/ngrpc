package ngrpc

import (
	"context"
	"log"
)

// Logger logger
type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Debugln(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Infoln(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warnln(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Errorln(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Fatalln(ctx context.Context, args ...interface{})
}

type StdLogger struct {
}

func (StdLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[Debug] "+format, args...)
}
func (StdLogger) Debugln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[Debug]",
	}
	nArgs = append(nArgs, args...)
	log.Println(nArgs...)
}
func (StdLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}
func (StdLogger) Infoln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[INFO]",
	}
	nArgs = append(nArgs, args...)
	log.Println(nArgs...)
}
func (StdLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[Warn] "+format, args...)
}
func (StdLogger) Warnln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[Warn]",
	}
	nArgs = append(nArgs, args...)
	log.Println(nArgs...)
}
func (StdLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[Error] "+format, args...)
}
func (StdLogger) Errorln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[Error]",
	}
	nArgs = append(nArgs, args...)
	log.Println(nArgs...)
}
func (StdLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	log.Fatalf("[Fatal] "+format, args...)
}
func (StdLogger) Fatalln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[Fatal]",
	}
	nArgs = append(nArgs, args...)
	log.Fatalln(nArgs...)
}
