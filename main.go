package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}

}

func run() {

	// call fork-exec to create new process namespace
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// isolate hostname - clone new unix time sharing system
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	must(cmd.Run())
}

func child() {

	fmt.Printf("Running %v,  Process id:  %d\n ", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	must(syscall.Sethostname([]byte("Container")))
	// must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	// must(os.MkdirAll("rootfs/oldrootfs", 0700))
	/*
		int pivot_root(const char *new_root, const char *put_old);
		pivot_root() moves the root file system of the current process to the directory put_old and makes new_root the new root file system of the current process.
	*/
	// must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	must(syscall.Chroot("rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
	syscall.Unmount("/proc", 0)
}

func gracefulShutdown(term chan os.Signal) {
	// term := make(chan os.Signal)
	// signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	// go gracefulShutdown(term)
	termsig := <-term
	fmt.Printf("\nRecieved signal %v\n", termsig)
	os.Exit(0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
