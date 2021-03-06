/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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

package grub2

import (
	dbusGrub2ext "dbus/com/deepin/daemon/grub2ext"
)

// dbus api wrapper for grub2ext

func newDbusGrub2Ext() (grub2ext *dbusGrub2ext.Grub2Ext, err error) {
	grub2ext, err = dbusGrub2ext.NewGrub2Ext("com.deepin.daemon.Grub2Ext", "/com/deepin/daemon/Grub2Ext")
	if err != nil {
		logger.Error(err)
	}
	return
}

func grub2extDoWriteThemeMainFile(themeFileContent string) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteThemeMainFile(themeFileContent)
}

func grub2extDoGenerateGrubMenu() {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoGenerateGrubMenu()
}

func grub2extDoGenerateThemeBackground(screenWidth, screenHeight uint16) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoGenerateThemeBackground(screenWidth, screenHeight)
}

func grub2extDoResetThemeBackground() {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoResetThemeBackground()
}

func grub2extDoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoSetThemeBackgroundSourceFile(imageFile, screenWidth, screenHeight)
}

func grub2extDoWriteConfig(fileContent string) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteConfig(fileContent)
}

func grub2extDoWriteGrubSettings(fileContent string) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteGrubSettings(fileContent)
}

func grub2extDoWriteThemeTplFile(jsonContent string) {
	grub2ext, err := newDbusGrub2Ext()
	if err != nil {
		return
	}
	grub2ext.DoWriteThemeTplFile(jsonContent)
}
