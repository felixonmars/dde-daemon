/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package accounts

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_m     *Manager
	logger = log.NewLogger("daemon/accounts")
)

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("accounts", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _m != nil {
		return nil
	}

	logger.BeginTracing()
	_m = NewManager()
	err := dbus.InstallOnSystem(_m)
	if err != nil {
		logger.Error("Install manager dbus failed:", err)
		if _m.watcher != nil {
			_m.watcher.EndWatch()
			_m.watcher = nil
		}
		return err
	}

	_m.installUsers()
	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	_m.destroy()
	_m = nil

	return nil
}

func doEmitSuccess(pid uint32, action string) {
	if _m == nil {
		return
	}
	dbus.Emit(_m, "Success", pid, action)
}

func doEmitError(pid uint32, action, reason string) {
	if _m == nil {
		return
	}

	dbus.Emit(_m, "Error", pid, action, reason)
}
