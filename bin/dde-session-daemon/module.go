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

package main

import (
	_ "pkg.deepin.io/dde/daemon/appearance"
	_ "pkg.deepin.io/dde/daemon/audio"
	_ "pkg.deepin.io/dde/daemon/bluetooth"
	_ "pkg.deepin.io/dde/daemon/clipboard"
	//_ "pkg.deepin.io/dde/daemon/dock"
	"gir/gio-2.0"
	_ "pkg.deepin.io/dde/daemon/inputdevices"
	_ "pkg.deepin.io/dde/daemon/keybinding"
	_ "pkg.deepin.io/dde/daemon/launcher"
	_ "pkg.deepin.io/dde/daemon/mounts"
	_ "pkg.deepin.io/dde/daemon/mpris"
	_ "pkg.deepin.io/dde/daemon/network"
	_ "pkg.deepin.io/dde/daemon/power"
	_ "pkg.deepin.io/dde/daemon/screenedge"
	_ "pkg.deepin.io/dde/daemon/screensaver"
	_ "pkg.deepin.io/dde/daemon/sessionwatcher"
	_ "pkg.deepin.io/dde/daemon/soundeffect"
	_ "pkg.deepin.io/dde/daemon/systeminfo"
	_ "pkg.deepin.io/dde/daemon/timedate"
)

var (
	daemonSettings = gio.NewSettings("com.deepin.dde.daemon")
)

// TODO:
// func listenDaemonSettings() {
// 	daemonSettings.Connect("changed", func(s *gio.Settings, name string) {
// 		// gsettings key names must keep consistent with module names
// 		enable := daemonSettings.GetBoolean(name)
// 		loader.Enable(name, enable)
// 		if enable {
// 			loader.Start(name)
// 		} else {
// 			loader.Stop(name)
// 		}
// 	})
// 	daemonSettings.GetBoolean("mounts")
// }
