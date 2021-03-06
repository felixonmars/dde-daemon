package appearance

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.Appearance"
	dbusPath = "/com/deepin/daemon/Appearance"
	dbusIFC  = "com.deepin.daemon.Appearance"
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) setPropTheme(value string) {
	if m.Theme == value {
		return
	}

	m.Theme = value
	dbus.NotifyChange(m, "Theme")
}

func (m *Manager) setPropFontSize(size int32) {
	if m.FontSize == size {
		return
	}

	m.FontSize = size
	dbus.NotifyChange(m, "FontSize")
}
