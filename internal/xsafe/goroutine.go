package xsafe

import (
	"context"
	"log"

	"github.com/cloudwego/gopkg/concurrency/gopool"
)

type PanicHandler func(value any)

var panicHandler PanicHandler = func(value any) {
	log.Printf("panic recovered: %v", value)
}

func SetPanicHandler(handler PanicHandler) {
	if handler == nil {
		panicHandler = func(value any) {
			log.Printf("panic recovered: %v", value)
		}
		return
	}
	panicHandler = handler
}

func InstallGoPoolPanicHandler() {
	gopool.SetPanicHandler(func(ctx context.Context, value interface{}) {
		panicHandler(value)
	})
}

func DeferRecover(handler PanicHandler) {
	if value := recover(); value != nil {
		if handler != nil {
			handler(value)
			return
		}
		panicHandler(value)
	}
}

func Go(fn func()) {
	go func() {
		defer DeferRecover(nil)
		fn()
	}()
}
