/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
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

package main

const tplFrontEndConnProp = `// This file is automatically generated, please don't edit manually.
import QtQuick 2.1
import Deepin.Widgets 1.0
import "../edit"

BaseEditPage {
    id: editPage
    activeExpandIndex: {
        switch (connectionSession.type) {
            case nmConnectionTypeWired: return 23;
            case nmConnectionTypeWireless: return 5;
            case nmConnectionTypeWirelessAdhoc: return 5;
            case nmConnectionTypeWirelessHotspot: return 5;
            case nmConnectionTypePppoe: return 7;
            case nmConnectionTypeMobileGsm: return 2;
            case nmConnectionTypeMobileCdma: return 2;
            case nmConnectionTypeVpnL2tp: return 9;
            case nmConnectionTypeVpnPptp: return 9;
            case nmConnectionTypeVpnVpnc: return 9;
            case nmConnectionTypeVpnOpenvpn: return 9;
            case nmConnectionTypeVpnOpenconnect: return 9;
        }
    }

    EditLineConnectionId {
        id: lineConnectionId
        connectionSession: editPage.connectionSession
        availableSections: editPage.availableSections
        availableKeys: editPage.availableKeys
        errors: editPage.errors
        section: "connection"
        key: "id"
        text: dsTr("Name")
        alwaysUpdate: true
    }
    EditLineSwitchButton {
        id: lineConnectionAutoconnect
        connectionSession: editPage.connectionSession
        availableSections: editPage.availableSections
        availableKeys: editPage.availableKeys
        errors: editPage.errors
        section: "connection"
        key: "autoconnect"
        text: dsTr("Automatically connect")
    }
    EditLineSwitchButton {
        id: lineConnectionVkVpnAutoconnect
        connectionSession: editPage.connectionSession
        availableSections: editPage.availableSections
        availableKeys: editPage.availableKeys
        errors: editPage.errors
        section: "connection"
        key: "vk-vpn-autoconnect"
        text: dsTr("Automatically connect")
    }
    {{range $i, $vsection := .}}{{if .Ignore}}{{else}}{{$id := $vsection.Name | ToVsClassName | printf "section%s"}}
    EditSection{{$vsection.Name | ToVsClassName}} {
        myIndex: {{$i}}
        id: {{$id}}
        activeExpandIndex: editPage.activeExpandIndex
        connectionSession: editPage.connectionSession
        availableSections: editPage.availableSections
        availableKeys: editPage.availableKeys
        errors: editPage.errors
    }
    EditSectionSeparator {relatedSection: {{$id}}}
    {{end}}{{end}}
}
`

const tplFrontEndSection = `// This file is automatically generated, please don't edit manually.
import QtQuick 2.1
import Deepin.Widgets 1.0
import "../edit"

BaseEditSection { {{$sectionId := .Name | ToVsClassName | printf "section%s"}}
    id: {{$sectionId}}
    virtualSection: "{{.Value}}"

    header.sourceComponent: EditDownArrowHeader{
        text: dsTr("{{.DisplayName}}")
    }

    content.sourceComponent: Column { {{range $i, $key := GetAllKeysInVsection .Name}}{{if IsKeyUsedByFrontEnd $key}}{{$sectionValue := $key | ToKeyRelatedSectionValue}}{{$value := $key | ToKeyValue}}
        {{$widget := ToFrontEndWidget $key}}{{$widget}} {
            id: line{{$sectionValue | ToClassName}}{{$value | ToClassName}}
            connectionSession: {{$sectionId}}.connectionSession
            availableSections: {{$sectionId}}.availableSections
            availableKeys: {{$sectionId}}.availableKeys
            errors: {{$sectionId}}.errors
            section: "{{$sectionValue}}"
            key: "{{$value}}"
            text: dsTr("{{$key | ToKeyDisplayName}}"){{range $propKey, $propValue := GetKeyWidgetProps $key}}
            {{$propKey}}: {{$propValue}}{{end}}
        }{{end}}{{end}}
    }
}
`
