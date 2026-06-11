package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Record 单条安装记录
type Record struct {
	ToolID        string    `json:"tool_id"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	InstalledAt   time.Time `json:"installed_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	InstallMethod string    `json:"install_method"`
}

// Records 安装记录管理
type Records struct {
	mu       sync.RWMutex
	filePath string
	data     map[string]Record // key = toolID
}

// NewRecords 创建安装记录实例，自动加载已有记录
func NewRecords() (*Records, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(appData, "WinDevReady")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	r := &Records{
		filePath: filepath.Join(dir, "install_records.json"),
		data:     make(map[string]Record),
	}
	// 尝试加载已有记录（文件不存在则忽略）
	_ = r.load()
	return r, nil
}

// load 从磁盘加载记录
func (r *Records) load() error {
	raw, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}
	var list []Record
	if err := json.Unmarshal(raw, &list); err != nil {
		return err
	}
	for _, rec := range list {
		r.data[rec.ToolID] = rec
	}
	return nil
}

// Save 持久化到磁盘
func (r *Records) Save() error {
	r.mu.RLock()
	var list []Record
	for _, rec := range r.data {
		list = append(list, rec)
	}
	r.mu.RUnlock()

	raw, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, raw, 0644)
}

// Get 获取单条记录
func (r *Records) Get(toolID string) (Record, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.data[toolID]
	return rec, ok
}

// GetAll 获取全部记录
func (r *Records) GetAll() []Record {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []Record
	for _, rec := range r.data {
		list = append(list, rec)
	}
	return list
}

// Upsert 新增或更新记录
func (r *Records) Upsert(rec Record) {
	r.mu.Lock()
	if existing, ok := r.data[rec.ToolID]; ok {
		rec.InstalledAt = existing.InstalledAt
		rec.UpdatedAt = time.Now()
	}
	r.data[rec.ToolID] = rec
	r.mu.Unlock()
}

// Remove 删除指定工具的记录
func (r *Records) Remove(toolID string) {
	r.mu.Lock()
	delete(r.data, toolID)
	r.mu.Unlock()
}

// FilePath 返回记录文件路径（供 UI 显示）
func (r *Records) FilePath() string {
	return r.filePath
}
