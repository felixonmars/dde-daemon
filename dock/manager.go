package dock

import "pkg.deepin.io/lib/dbus"
import pkgbus "dbus/org/freedesktop/dbus"
import "time"
import "sync"

var busdaemon *pkgbus.DBusDaemon

const (
// FieldTitle   = "title"
// FieldIcon    = "icon"
// FieldMenu    = "menu"
// FieldAppXids = "app-xids"
//
// FieldStatus   = "app-status"
// ActiveStatus  = "active"
// NormalStatus  = "normal"
// InvalidStatus = "invalid"
)

// EntryProxyerManager为驻留程序以及打开程序的dbus接口管理器。
// 所有已驻留程序以及打开的程序都会生成对应的dbus接口。
type EntryProxyerManager struct {
	// Entries为目前所有打开程序与驻留程序列表。（此属性可能会被废弃掉，调用一个初始化的方法，然后在需要的情况下触发Added信号）。
	Entries       []*EntryProxyer
	entrireLocker sync.Mutex

	// Added在程序需要在前端显示时被触发。
	Added func(dbus.ObjectPath)
	// Removed会在程序不再需要在dock前端显示时触发。
	Removed func(string)
	// TrayInited在trayicon相关内容初始化完成后触发。
	TrayInited func()
}

func (m *EntryProxyerManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/EntryManager",
		Interface:  "dde.dock.EntryManager",
	}
}

func NewEntryProxyerManager() *EntryProxyerManager {
	m := &EntryProxyerManager{}
	return m
}

func (m *EntryProxyerManager) destroy() {
	for _, entry := range m.Entries {
		dbus.UnInstallObject(entry)
	}
	dbus.UnInstallObject(m)
	if busdaemon != nil {
		pkgbus.DestroyDBusDaemon(busdaemon)
		busdaemon = nil
	}
}

func (m *EntryProxyerManager) watchEntries() {
	var err error
	busdaemon, err = pkgbus.NewDBusDaemon("org.freedesktop.DBus", "/org/freedesktop/DBus")
	if err != nil {
		panic(err)
	}

	// register existed entries
	names, err := busdaemon.ListNames()
	if err != nil {
		panic(err)
	}
	for _, n := range names {
		m.registerEntry(n)
	}

	// monitor name lost, name acquire
	busdaemon.ConnectNameOwnerChanged(func(name, oldOwner, newOwner string) {
		// if a new dbus session was installed, the name and newOwner
		// will be not empty, if a dbus session was uninstalled, the
		// name and oldOwner will be not empty
		if len(newOwner) != 0 {
			if name == "com.deepin.dde.TrayManager" {
				dbus.Emit(m, "TrayInited")
			}
			go func() {
				// FIXME: how long time should to wait for
				time.Sleep(60 * time.Millisecond)
				m.entrireLocker.Lock()
				defer m.entrireLocker.Unlock()
				m.registerEntry(name)
			}()
		} else {
			m.entrireLocker.Lock()
			defer m.entrireLocker.Unlock()
			m.unregisterEntry(name)
		}
	})
}

func (m *EntryProxyerManager) registerEntry(name string) {
	if !isEntryNameValid(name) {
		return
	}
	logger.Debugf("register entry: %s", name)
	entryId, ok := getEntryId(name)
	if !ok {
		return
	}
	logger.Debugf("register entry id: %s", entryId)
	entry, err := NewEntryProxyer(entryId)
	if err != nil {
		logger.Warningf("register entry failed: %v", err)
		return
	}
	err = dbus.InstallOnSession(entry)
	if err != nil {
		logger.Warningf("register entry failed: %v", err)
		return
	}
	m.Entries = append(m.Entries, entry)

	dbus.Emit(m, "Added", dbus.ObjectPath(entry.GetDBusInfo().ObjectPath))

	logger.Infof("register entry success: %s", name)
}

func (m *EntryProxyerManager) unregisterEntry(name string) {
	if !isEntryNameValid(name) {
		return
	}
	logger.Debugf("unregister entry: %s", name)
	entryId, ok := getEntryId(name)
	if !ok {
		return
	}
	logger.Debugf("unregister entry id: %s", entryId)

	// find the index
	index := -1
	var entry *EntryProxyer = nil
	for i, e := range m.Entries {
		if e.entryId == entryId {
			index = i
			entry = e
			break
		}
	}

	if index < 0 {
		logger.Warningf("slice out of bounds, entry len: %d, index: %d", len(m.Entries), index)
		return
	}
	logger.Debugf("entry len: %d, index: %d", len(m.Entries), index)

	if entry != nil {
		dbus.UnInstallObject(entry)
	}

	// remove the entry from slice
	copy(m.Entries[index:], m.Entries[index+1:])
	m.Entries[len(m.Entries)-1] = nil
	m.Entries = m.Entries[:len(m.Entries)-1]

	dbus.Emit(m, "Removed", entry.Id)

	logger.Infof("unregister entry success: %s", name)
}
