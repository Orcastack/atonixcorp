package tracing

import "log"

func Trace(msg string) {
    log.Println("[TRACE]", msg)
}
