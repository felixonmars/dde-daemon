package inputdevices

import (
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/dbus"
	"fmt"
)

var (
	_manager *Manager
	logger = log.NewLogger("daemon/inputdevices")
)

type Daemon struct {
	*ModuleBase
}

func init() {
	Register(NewInputdevicesDaemon(logger))
}
func NewInputdevicesDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("inputdevices", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	logger.BeginTracing()
	_manager = NewManager()

	err := installSessionBus(_manager)
	if err != nil {
		return err
	}

	err = installSessionBus(_manager.kbd)
	if err != nil {
		return err
	}

	err = installSessionBus(_manager.wacom)
	if err != nil {
		return err
	}

	err = installSessionBus(_manager.tpad)
	if err != nil {
		return err
	}
	err = installSessionBus(_manager.mouse)
	if err != nil {
		return err
	}

	go startDeviceListener()
	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	logger.EndTracing()
	endDeviceListener()
	return nil
}

func installSessionBus(obj dbus.DBusObject) error {
	if obj == nil {
		logger.Error("Invalid dbus object: empty")
		return fmt.Errorf("Invalid dbus object")
	}

	err := dbus.InstallOnSession(obj)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return err
	}
	return nil
}
