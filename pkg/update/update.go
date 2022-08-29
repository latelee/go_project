/*
升级切换记录：

路径：A进程运行，收到升级指令，拷贝为A.old，并切换为到子进程A.old。子进程A.old将A.new拷贝为A，运行A。如果运行失败，则将A.old拷贝为A，再运行。

1、切换升级时，用0退出码（startProcessAndExit），
systemd主要配置：
Type=forking
RestartPreventExitStatus=42 SIGKILL
Restart=always // !! 此值保持
可达到目的。

2、切换升级时，用42退出，
systemd主要配置：
Type=forking
RestartPreventExitStatus=42 SIGKILL
Restart=always // !! 此值保持
无法达到目的。即一切换子进程，原来的进程被systemd启动了。

3、保持0退出码，systemd添加KillMode：
Type=notify
KillMode=process
Restart=always

无法达到目的。  

*/

// 注：update实际只是一些函数的封装，不算是业务模块，此处为演示

package update

import (
    _ "fmt"
    "webdemo/pkg/com"
    "webdemo/pkg/klog"

    _ "io"
    "os"
    _ "os/exec"
    _ "os/signal"
    _ "log"
    "time"
    _ "reflect"
    _ "flag"
    _ "github.com/gin-gonic/gin"
    _ "net/http"
    "syscall"
)

type update struct {
    enable bool
    // 后可加其它字段
}

// TODO：使用常量
// 分别是原文件、升级文件、旧文件（与原文件完全相同，备份用）
var appname = "./app"
var appnameup = "./app.new"
var appnameold = "./app.old"

func StartProcessAndWait(argv []string, fn func(bool)) bool {
    //ppid := syscall.Getppid()
    //klog.Printf("ppid: %v\n", ppid);

    sysattrs := syscall.SysProcAttr{}
    execAttrs := os.ProcAttr{
                Env:   os.Environ(),
                Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
                Sys: &sysattrs}
    proc, err := os.StartProcess(argv[0], argv, &execAttrs)
	if err != nil {
		klog.Printf("start %v failed,err: %v\n", argv, err);
		return false
	}

    go func() {
        stat, err := proc.Wait();
        klog.Printf("Wait ret: %v\n", err);
        pid := proc.Pid
        klog.Printf("Child process %d exit with reason: %v\n", pid, err);
        // 是否已退出
        b := stat.Exited();
        if b {
            klog.Printf("will kill pid %v exitcode: %v.\n", pid, stat.Success());
            proc.Kill() // 似乎要不要都可以，先保留
            // 无论成功或失败，都使用外部函数处理
            fn(stat.Success());
        }
        proc.Release();
    }()

    return true;
}

// 创建子进程后，立即退出
// 参数argv与main函数的argv完全等效，ms为延时用的毫秒数
func startProcessAndExit(argv []string, ms time.Duration, code int) {
    //ppid := syscall.Getppid()
    //klog.Printf("ppid: %v\n", ppid);

    //sysattrs := syscall.SysProcAttr{Setsid: true, Setpgid: true} // 这些标志设置后跑不起来
    sysattrs := syscall.SysProcAttr{}
    execAttrs := os.ProcAttr{
                Env:   os.Environ(),
                Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
                Sys: &sysattrs}

    //klog.Printf("will start: %v\n", argv);
    proc, err := os.StartProcess(argv[0], argv, &execAttrs)
	if err != nil {
		klog.Printf("start %v failed,err: %v\n", argv, err);
		return
	}

    // 注：这里不建议延时太久，如果父进程不及时退出，无法拷贝同名文件
    com.Sleep(ms);
    // 释放资源，自动退出
    klog.Printf("pid %v exit.\n", syscall.Getpid());
    proc.Release()
    os.Exit(code)
}

func checkOKFailed(code bool) {
    // 如果失败，运行旧程序
    if (code == false) {
        klog.Printf("run app failed, will run old app.\n");
        // 将文件改为原始的
        com.CopyFile(appname, appnameold);

        klog.Printf("run %v...\n", appname);
        mycmds := []string{appname};

        startProcessAndExit(mycmds, 10, 0);
    }
}

func EnterUpgradeApp() bool {
    klog.Printf("EnterUpgradeApp: pwd: %v\n", com.GetPWD());
    if !com.IsExist(appnameup) {
        klog.Printf("no upgrade file exist %v.\n", appnameup);
        return false
    }
    // 将本程序文件拷贝一份，后加 .old
    com.CopyFile(appnameold, appname);

    klog.Printf("run %v...\n", appnameold);
    mycmds := []string{appnameold, "-m", "upgrade"};

    startProcessAndExit(mycmds, 10, 0);

    return true;
}

// 进行升级处理
func ProcessUpgrade() {
    klog.Printf("ProcessUpgrade: pwd: %v\n", com.GetPWD());
    //ppid := syscall.Getppid()
    //klog.Printf("child ppid: %v\n", ppid);
    
    com.Sleep(2000);
    
    // 2. 重命名
    klog.Printf("copy %v to %v\n", appnameup, appname);
    com.CopyFile(appname, appnameup);

    klog.Printf("run %v...\n", appname);
    //mycmds := []string{appname, "-m", "normaltest"};
    mycmds := []string{appname};

    // 防止新程序无法运行
    ret := StartProcessAndWait(mycmds, checkOKFailed);
    if !ret {
        klog.Printf("start %v failed, run old one.\n", mycmds);
        checkOKFailed(false);
    }
    // 监控，如果业务程序OK，会一直跑，此处退出，如果业务程序不OK，在前面退出时处理
    // 注：需要根据经验值判断业务程序是否正常，如：业务程序1分钟内退出为异常，1分钟都正常，则正常
    // 在此处循环等待，如果超过范围，则所有程序无法运行
    cnt := 0;
    for {
        klog.Printf(".\n");
        com.Sleep(2000);
        cnt ++;
        if cnt >= 4 {
            klog.Printf("upgrade app exit for timeout.\n");
            os.Exit(0);
        }
    }
}
