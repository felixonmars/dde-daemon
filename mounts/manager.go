/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
	"time"
)

const (
	mediaHandlerSchema = "org.gnome.desktop.media-handling"

	refrashDuration = time.Minute * 10
)

type Manager struct {
	DiskList DiskInfos
	Error    func(string, string) // (id, message)

	monitor *gio.VolumeMonitor
	setting *gio.Settings
	quit    chan struct{}
}

func newManager() *Manager {
	var m = new(Manager)
	m.monitor = gio.VolumeMonitorGet()
	m.DiskList = m.listDisk()
	m.init()
	m.setting, _ = utils.CheckAndNewGSettings(mediaHandlerSchema)
	m.quit = make(chan struct{})
	return m
}

func (m *Manager) init() {
	for _, info := range m.DiskList {
		if info.object.Type != diskObjVolume || info.Type != DiskTypeRemovable {
			continue
		}

		logger.Debug("[init] will mount volume:", info.Name, info.Id)
		err := m.Mount(info.Id)
		if err != nil {
			logger.Warningf("Mount '%s - %s' failed: %v",
				info.Name, info.Id, err)
		}
	}
	m.refrashDiskList()
}

func (m *Manager) destroy() {
	if m.quit != nil {
		close(m.quit)
		m.quit = nil
	}

	if m.setting != nil {
		m.setting.Unref()
		m.setting = nil
	}

	m.DiskList.destroy()
	if m.monitor != nil {
		m.monitor.Unref()
		m.monitor = nil
	}
}

func (m *Manager) emitError(id, msg string) {
	dbus.Emit(m, "Error", id, msg)
}

func (m *Manager) refrashDiskList() {
	m.DiskList.destroy()
	m.setPropDiskList(m.listDisk())
}

func (m *Manager) updateDiskInfo() {
	for {
		select {
		case <-time.After(refrashDuration):
			m.refrashDiskList()
		case <-m.quit:
			return
		}
	}
}
