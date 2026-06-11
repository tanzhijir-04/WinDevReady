package network

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"WinDevReady/internal/logger"
)

// Status 网络状态
type Status struct {
	NPMReachable  bool   // 能否访问 npm registry
	GitHubReachable bool // 能否访问 GitHub
	ProxyActive   bool   // 系统代理是否开启
	ProxyPort     string // 代理端口
	RegistryURL   string // 当前使用的 registry 地址
}

// Detector 网络检测器
type Detector struct {
	log *logger.Logger
}

// New 创建网络检测器
func New(log *logger.Logger) *Detector {
	return &Detector{log: log}
}

// CheckAll 执行完整网络检测
func (d *Detector) CheckAll() Status {
	status := Status{}

	// 1. 检测 npm registry
	d.log.Info("network", "正在检测 npm registry 连通性...")
	status.NPMReachable = d.pingURL("https://registry.npmjs.org/")
	if status.NPMReachable {
		d.log.Success("network", "npm registry 可达")
		status.RegistryURL = "https://registry.npmjs.org/"
	} else {
		d.log.Warn("network", "npm registry 不可达，将切换到国内镜像")
		// 尝试国内镜像
		if d.pingURL("https://registry.npmmirror.com/") {
			status.RegistryURL = "https://registry.npmmirror.com/"
			d.log.Success("network", "npmmirror 镜像可用")
		} else if d.pingURL("https://mirrors.tuna.tsinghua.edu.cn/npm/") {
			status.RegistryURL = "https://mirrors.tuna.tsinghua.edu.cn/npm/"
			d.log.Success("network", "清华源可用")
		}
	}

	// 2. 检测 GitHub
	d.log.Info("network", "正在检测 GitHub 连通性...")
	status.GitHubReachable = d.pingURL("https://github.com/")
	if status.GitHubReachable {
		d.log.Success("network", "GitHub 可达")
	} else {
		d.log.Warn("network", "GitHub 不可达，部分功能可能受限")
	}

	// 3. 检测系统代理
	d.log.Info("network", "正在检测系统代理...")
	status.ProxyActive, status.ProxyPort = d.detectProxy()
	if status.ProxyActive {
		d.log.Success("network", fmt.Sprintf("检测到系统代理，端口: %s", status.ProxyPort))
	} else {
		d.log.Info("network", "未检测到系统代理")
	}

	return status
}

// ApplyNPMRegistry 设置 npm registry
func (d *Detector) ApplyNPMRegistry(registryURL string) error {
	if registryURL == "" {
		return nil
	}
	cmd := exec.Command("npm", "config", "set", "registry", registryURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("设置 npm registry 失败: %s, 输出: %s", err, string(output))
	}
	d.log.Success("network", fmt.Sprintf("npm registry 已切换为 %s", registryURL))
	return nil
}

// ApplyProxy 为 npm/pip/git 配置代理
func (d *Detector) ApplyProxy(port string) error {
	if port == "" {
		return nil
	}
	proxyURL := fmt.Sprintf("http://127.0.0.1:%s", port)
	var errs []string

	// npm proxy
	if out, err := exec.Command("npm", "config", "set", "proxy", proxyURL).CombinedOutput(); err != nil {
		errs = append(errs, fmt.Sprintf("npm proxy: %s %s", err, out))
	}
	if out, err := exec.Command("npm", "config", "set", "https-proxy", proxyURL).CombinedOutput(); err != nil {
		errs = append(errs, fmt.Sprintf("npm https-proxy: %s %s", err, out))
	}

	// git proxy
	if out, err := exec.Command("git", "config", "--global", "http.proxy", proxyURL).CombinedOutput(); err != nil {
		errs = append(errs, fmt.Sprintf("git http.proxy: %s %s", err, out))
	}
	if out, err := exec.Command("git", "config", "--global", "https.proxy", proxyURL).CombinedOutput(); err != nil {
		errs = append(errs, fmt.Sprintf("git https.proxy: %s %s", err, out))
	}

	// pip proxy（通过环境变量或 pip.conf）
	if runtime.GOOS == "windows" {
		if err := d.setPipProxyWindows(proxyURL); err != nil {
			errs = append(errs, fmt.Sprintf("pip proxy: %s", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("部分代理配置失败:\n%s", strings.Join(errs, "\n"))
	}
	d.log.Success("network", fmt.Sprintf("已为 npm/git/pip 配置代理: %s", proxyURL))
	return nil
}

// setPipProxyWindows 在 Windows 上设置 pip 全局代理
func (d *Detector) setPipProxyWindows(proxyURL string) error {
	pipConfDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	pipDir := pipConfDir + "\\pip"
	if err := os.MkdirAll(pipDir, 0755); err != nil {
		return err
	}
	content := fmt.Sprintf("[global]\nproxy = %s\nhttps-proxy = %s\n", proxyURL, proxyURL)
	return os.WriteFile(pipDir+"\\pip.ini", []byte(content), 0644)
}

// pingURL 检测 URL 是否可达（2秒超时）
func (d *Detector) pingURL(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

// detectProxy 检测 Windows 系统代理
func (d *Detector) detectProxy() (bool, string) {
	// 读取注册表获取系统代理设置
	cmd := exec.Command("reg", "query",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		"/v", "ProxyEnable")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}

	// 检查 ProxyEnable 是否为 1
	if !strings.Contains(string(out), "0x1") {
		return false, ""
	}

	// 读取代理服务器地址
	cmd = exec.Command("reg", "query",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		"/v", "ProxyServer")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}

	// 解析端口号
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "ProxyServer") {
			parts := strings.Split(line, "    ")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if strings.Contains(p, ":") {
					addrParts := strings.Split(p, ":")
					if len(addrParts) > 1 {
						return true, addrParts[len(addrParts)-1]
					}
				}
			}
		}
	}

	// 尝试直接从 net 检测常见代理端口
	for _, port := range []string{"7890", "1080", "8080", "10809"} {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return true, port
		}
	}

	return false, ""
}
