package jail

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"syscall"

	"github.com/dotcloud/docker/daemon/execdriver"
)

const DriverName = "jail"

func init() {
	execdriver.RegisterInitFunc(DriverName, func(args *execdriver.InitArgs) error {
		runtime.LockOSThread()

		path, err := exec.LookPath(args.Args[0])
		if err != nil {
			log.Printf("Unable to locate %v", args.Args[0])
			os.Exit(127)
		}
		if err := syscall.Exec(path, args.Args, os.Environ()); err != nil {
			return fmt.Errorf("dockerinit unable to execute %s - %s", path, err)
		}
		panic("Unreachable")
	})
}

type driver struct {
	root     string
	initPath string
}

func NewDriver(root, initPath string) (*driver, error) {
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}

	return &driver{
		root:     root,
		initPath: initPath,
	}, nil
}

func (d *driver) Name() string {
	return DriverName
}

func copyFile(src string, dest string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, content, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) Run(c *execdriver.Command, pipes *execdriver.Pipes, startCallback execdriver.StartCallback) (int, error) {
	if err := execdriver.SetTerminal(c, pipes); err != nil {
		return -1, err
	}

	root := c.Rootfs

	init := path.Join(root, ".dockerinit")
	if err := copyFile(os.Args[0], init); err != nil {
		return -1, err
	}

	devDir := path.Join(root, "dev")
	if err := os.MkdirAll(devDir, 0755); err != nil {
		return -1, err
	}

	params := []string{
		"/usr/sbin/jail",
		"-c",
		"name=" + c.ID,
		"path=" + root,
		"command=" + c.InitPath,
		"-driver",
		DriverName,
	}

	if c.User != "" {
		params = append(params, "-u", c.User)
	}

	if c.Privileged {
		params = append(params, "-privileged")
	}

	if c.WorkingDir != "" {
		params = append(params, "-w", c.WorkingDir)
	}

	params = append(params, "--", c.Entrypoint)
	params = append(params, c.Arguments...)

	c.Path = "/usr/sbin/jail"
	c.Args = params

	if err := c.Run(); err != nil {
		return -1, err
	}

	return getExitCode(c), nil
}

func getExitCode(c *execdriver.Command) int {
	if c.ProcessState == nil {
		return -1
	}
	return c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
}

func (d *driver) Kill(c *execdriver.Command, sig int) error {
	return nil
}

func (d *driver) Pause(c *execdriver.Command) error {
	return nil
}

func (d *driver) Unpause(c *execdriver.Command) error {
	return nil
}

func (d *driver) Terminate(c *execdriver.Command) error {
	return nil
}

func (d *driver) GetPidsForContainer(id string) ([]int, error) {
	return nil, nil
}

type info struct {
	ID     string
	driver *driver
}

func (d *driver) Info(id string) execdriver.Info {
	return &info{ID: id, driver: d}
}

func (info *info) IsRunning() bool {
	if err := exec.Command("jls", "-j", info.ID).Run(); err != nil {
		return true
	}

	return false
}
