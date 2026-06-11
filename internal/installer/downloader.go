package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Downloader 文件下载器
type Downloader struct {
	tempDir string
}

// NewDownloader 创建下载器，初始化临时目录
func NewDownloader() (*Downloader, error) {
	dir, err := os.MkdirTemp("", "windevready-dl-*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}
	return &Downloader{tempDir: dir}, nil
}

// DownloadWithFallback 按优先级依次尝试多个下载地址，返回本地文件路径
func (d *Downloader) DownloadWithFallback(urls []string, filename string, progressFn func(downloaded, total int64)) (string, error) {
	if len(urls) == 0 {
		return "", fmt.Errorf("未提供下载地址")
	}

	var lastErr error
	for _, url := range urls {
		path, err := d.downloadSingle(url, filename, progressFn)
		if err == nil {
			return path, nil
		}
		lastErr = err
	}
	return "", fmt.Errorf("所有下载地址均失败，最后错误: %w", lastErr)
}

// downloadSingle 从单个 URL 下载文件
func (d *Downloader) downloadSingle(url, filename string, progressFn func(downloaded, total int64)) (string, error) {
	// 如果未指定文件名，从 URL 提取
	if filename == "" {
		filename = filepath.Base(url)
	}

	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	// 创建临时文件
	savePath := filepath.Join(d.tempDir, filename)
	out, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	total := resp.ContentLength
	if total <= 0 {
		total = 0
	}

	// 带进度的拷贝
	buf := make([]byte, 32*1024) // 32KB 缓冲区
	var downloaded int64
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			written, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return "", fmt.Errorf("写入文件失败: %w", writeErr)
			}
			downloaded += int64(written)
			if progressFn != nil {
				progressFn(downloaded, total)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return "", fmt.Errorf("读取数据失败: %w", readErr)
		}
	}

	return savePath, nil
}

// Cleanup 清理临时下载目录
func (d *Downloader) Cleanup() {
	if d.tempDir != "" {
		os.RemoveAll(d.tempDir)
	}
}
