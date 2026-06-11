package logger

import (
	"fmt"
	"sync"
	"time"
)

// Level 日志级别
type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
	LevelSuccess
)

// Entry 单条日志
type Entry struct {
	Time    time.Time
	Level   Level
	ToolID  string // 关联的工具 ID（可选）
	Message string
}

// Logger 流式日志管理器，支持多订阅者
type Logger struct {
	mu          sync.RWMutex
	subscribers []chan Entry
}

// New 创建日志实例
func New() *Logger {
	return &Logger{}
}

// Subscribe 订阅日志流，返回一个 channel
func (l *Logger) Subscribe() chan Entry {
	ch := make(chan Entry, 64)
	l.mu.Lock()
	l.subscribers = append(l.subscribers, ch)
	l.mu.Unlock()
	return ch
}

// Unsubscribe 取消订阅
func (l *Logger) Unsubscribe(ch chan Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, sub := range l.subscribers {
		if sub == ch {
			l.subscribers = append(l.subscribers[:i], l.subscribers[i+1:]...)
			break
		}
	}
	close(ch)
}

// emit 发送日志到所有订阅者
func (l *Logger) emit(level Level, toolID, msg string) {
	entry := Entry{
		Time:    time.Now(),
		Level:   level,
		ToolID:  toolID,
		Message: msg,
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, ch := range l.subscribers {
		// 非阻塞写入，避免慢消费者拖垮日志
		select {
		case ch <- entry:
		default:
		}
	}
}

// Info 普通信息日志
func (l *Logger) Info(toolID, msg string) {
	l.emit(LevelInfo, toolID, msg)
}

// Warn 警告日志
func (l *Logger) Warn(toolID, msg string) {
	l.emit(LevelWarn, toolID, msg)
}

// Error 错误日志
func (l *Logger) Error(toolID, msg string) {
	l.emit(LevelError, toolID, msg)
}

// Success 成功日志
func (l *Logger) Success(toolID, msg string) {
	l.emit(LevelSuccess, toolID, msg)
}

// FormatEntry 将日志条目格式化为可显示的字符串
func FormatEntry(e Entry) string {
	prefix := ""
	switch e.Level {
	case LevelInfo:
		prefix = "[信息]"
	case LevelWarn:
		prefix = "[警告]"
	case LevelError:
		prefix = "[错误]"
	case LevelSuccess:
		prefix = "[成功]"
	}
	if e.ToolID != "" {
		return fmt.Sprintf("%s [%s] %s %s", e.Time.Format("15:04:05"), e.ToolID, prefix, e.Message)
	}
	return fmt.Sprintf("%s %s %s", e.Time.Format("15:04:05"), prefix, e.Message)
}
