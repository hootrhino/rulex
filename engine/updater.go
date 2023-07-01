package engine

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	"github.com/hootrhino/rulex/glogger"
)

/**
GET /rulex?version=v1.5xx&arch=amd64&os=linux
Response:
headers=[
	'Rulex-Md5': 'xxx'
	'Rulex-Version':'v1.6'
]
body=exec_file

*/

// checkForUpdates rulex热更新，应当在单独的线程执行
func (e *RuleEngine) checkForUpdates(sig chan<- string) {
	// 设置轮询检查更新的定时器
	var t = time.NewTicker(10 * time.Minute)
	// 定时检查新版本
	for range t.C {
		glogger.GLogger.Info("checking for updates")
		data, err := requsetNewVerion(e.Config.UpdateServer)
		if err != nil {
			glogger.GLogger.Errorf("request new version failed, %s", err.Error())
			continue
		}
		if err = writeToFile(newRulex, data); err != nil {
			glogger.GLogger.Errorf("writeToFile: %s", err.Error())
			continue
		}

		sig <- newRulex
	}
}

const newRulex = "./rulex_new"

func requsetNewVerion(version string) ([]byte, error) {
	// 构建请求
	req, err := http.NewRequest(http.MethodGet, version, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = fmt.Sprintf("arch=%s&os=%s&version=%s",
		runtime.GOARCH, runtime.GOOS, version)
	cli := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response: %s", resp.Status)
	}
	newVersion := resp.Header.Get("Rulex-Version")
	md5Str := resp.Header.Get("Rulex-Md5")
	if len(newVersion) == 0 || len(md5Str) == 0 {
		if resp.Body != nil {
			resp.Body.Close()
		}
		return nil, fmt.Errorf("response: version='%s', md5='%s'", newVersion, md5Str)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	md5Data := md5.Sum(data)
	md5Str1 := hex.EncodeToString(md5Data[:])
	if md5Str != md5Str1 {
		return nil, fmt.Errorf("response: inconsistent MD5, '%s'!='%s'", md5Str, md5Str1)
	}

	glogger.GLogger.Infof("received new version '%s'", version)
	return data, nil
}

func writeToFile(fileName string, data []byte) error {
	info, err := os.Stat(fileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if info != nil && !info.IsDir() {
		os.Remove(info.Name())
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(data); err != nil {
		return err
	}

	return file.Chmod(0777)
}

// update 执行更新流程
func update(e *RuleEngine, name string) bool {
	fmt.Println("rulex updating")

	e.Version()
	// 1. 关闭当前服务
	e.Stop()

	// 2. 开启新服务
	var proc = Process{
		Name: name,
		Env:  true,
		Argv: os.Args,
	}
	pid, err := proc.Start()
	// 3. 更新失败，重新启动服务
	if err != nil {
		fmt.Printf("update failed, error is '%s'\n", err.Error())
		e.Start()
		return false
	}

	fmt.Printf("update succeed, new process is '%d'\n", pid)
	return true
}

// Wait 阻塞当前进程，等待关闭信号
func (e *RuleEngine) Wait(signals ...os.Signal) {
	sigs := make(chan os.Signal, 10)
	signal.Notify(sigs, signals...)

	var updateSig chan string
	fmt.Println("wait", e.Config.UpdateServer, e.Config.LogPath)
	if len(e.Config.UpdateServer) > 0 {
		updateSig = make(chan string, 1)
		go e.checkForUpdates(updateSig)
	}

	//	读取信号
	for {
		select {
		case <-sigs: //关闭信号，直接退出
			e.Stop()
			return
		case name := <-updateSig: // 检查更新
			// 更新成功，直接退出
			if update(e, name) {
				return
			}
		}
	}
}

// / Process 进程参数
type Process struct {
	Name       string         // 要加载的可执行文件路径
	Listener   []net.Listener // 要继承的链接
	Env        bool           // 是否继承环境变量?
	ExcludeEnv []string       // 要忽略的环境变量key
	Argv       []string       // 传递给新进程的命令行参数
}

// Filer 获取 os.File
type Filer interface {
	File() (f *os.File, err error)
}

// environ 继承环境变量
func (p Process) environ() []string {
	if !p.Env {
		return nil
	}

	var envs = os.Environ()
	var excludeEnv = append(p.ExcludeEnv, rulexInheritedFd)
	var envList = make([]string, 0, len(envs)+2)
	for _, s := range envs {
		var exclude bool
		for _, prefix := range excludeEnv {
			if len(s) > len(prefix) && s[len(prefix)] == '=' &&
				s[0:len(prefix)] == prefix {
				exclude = true
				break
			}
		}
		if exclude {
			continue
		}
		envList = append(envList, s)
	}

	return envList
}

func (p Process) files() []*os.File {
	var files = make([]*os.File, 0, len(p.Listener)+minInheritedFdCount)
	files = append(files, os.Stdin, os.Stdout, os.Stderr)

	// 从lisnter 转为 os.File
	for _, ln := range p.Listener {
		filer, ok := ln.(Filer)
		if !ok {
			continue
		}

		// 获取到文件描述符
		file, err := filer.File()
		if err != nil {
			continue
		}
		files = append(files, file)
	}

	return files
}

const (
	rulexInherited      = "RULEX_INHERITED"
	rulexInheritedFd    = "RULEX_INHERITED_FD"
	minInheritedFdCount = 3
)

// Start 启动新的进程
func (p Process) Start() (int, error) {
	// 查找可执行文件是否存在
	file, err := exec.LookPath(p.Name)
	if err != nil {
		return -1, err
	}

	// 构建文件描述符列表
	var files = p.files()
	// File()函数会调用系统调用copy一份fd，这里要进行关闭
	// 但是 stdin，stdout，stderr是直接引用的，不能关闭
	defer func() {
		for i := minInheritedFdCount; i < len(files); i++ {
			fd := files[i]
			_ = fd.Close()
		}
	}()

	// 获取到，要继承的环境变量
	env := p.environ()

	// 将继承的，文件描述符的数量写入到env
	if len(files) > minInheritedFdCount {
		env = append(env, rulexInheritedFd+"="+strconv.Itoa(len(files)))
	}

	// 获取到当前的执行路径
	wd, _ := os.Getwd()
	// 构建参数
	attr := os.ProcAttr{
		Dir:   wd,
		Env:   env,
		Files: files,
	}

	// 启动新进程
	proc, err := os.StartProcess(file, p.Argv, &attr)
	if err != nil {
		return -1, fmt.Errorf("start_process: %w", err)
	}

	return proc.Pid, nil
}
