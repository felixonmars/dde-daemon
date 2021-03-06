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

package network

import (
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/iso"
	"pkg.deepin.io/lib/mobileprovider"
)

const mobileProviderValueCustom = "<custom>"

func newMobileConnectionForDevice(id, uuid string, devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, err error) {
	logger.Infof("new mobile connection, id=%s, uuid=%s, devPath=%s", id, uuid, devPath)

	// guess default plan for mobile device
	countryCode, _ := iso.GetLocaleCountryCode()
	serviceType := getMobileDeviceServicType(devPath)
	plan, err := getDefaultPlanForMobileDevice(countryCode, serviceType)
	if err != nil {
		return
	}

	data := newMobileConnectionData("mobile", uuid, serviceType)
	addSettingSection(data, sectionCache)
	logicSetSettingVkMobileCountry(data, countryCode)
	logicSetSettingVkMobileProvider(data, plan.ProviderName)
	logicSetSettingVkMobilePlan(data, mobileprovider.MarshalPlan(plan))
	refileSectionCache(data)

	if active {
		cpath, _, err = nmAddAndActivateConnection(data, devPath)
	} else {
		cpath, err = nmAddConnection(data)
	}
	return
}
func getDefaultPlanForMobileDevice(countryCode, serviceType string) (plan mobileprovider.Plan, err error) {
	if serviceType == connectionMobileGsm {
		plan, err = mobileprovider.GetDefaultGSMPlanForCountry(countryCode)
	} else {
		plan, err = mobileprovider.GetDefaultCDMAPlanForCountry(countryCode)
	}
	if err != nil {
		logger.Error(err)
	}
	return
}
func getMobileDeviceServicType(devPath dbus.ObjectPath) (serviceType string) {
	capabilities := nmGetDeviceModemCapabilities(devPath)
	if (capabilities & NM_DEVICE_MODEM_CAPABILITY_LTE) == capabilities {
		// all LTE modems treated as GSM/UMTS
		serviceType = connectionMobileGsm
	} else if (capabilities & NM_DEVICE_MODEM_CAPABILITY_GSM_UMTS) == capabilities {
		serviceType = connectionMobileGsm
	} else if (capabilities & NM_DEVICE_MODEM_CAPABILITY_CDMA_EVDO) == capabilities {
		serviceType = connectionMobileCdma
	} else {
		logger.Errorf("Unknown modem capabilities (0x%x)", capabilities)
	}
	return
}

func newMobileConnectionData(id, uuid, serviceType string) (data connectionData) {
	data = make(connectionData)

	addSettingSection(data, sectionConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionAutoconnect(data, true)

	logicSetSettingVkMobileServiceType(data, serviceType)

	addSettingSection(data, sectionPpp)
	logicSetSettingVkPppEnableLcpEcho(data, true)

	addSettingSection(data, sectionSerial)
	setSettingSerialBaud(data, 115200)

	initSettingSectionIpv4(data)

	return
}

func initSettingSectionGsm(data connectionData) {
	setSettingConnectionType(data, NM_SETTING_GSM_SETTING_NAME)
	addSettingSection(data, sectionGsm)
	setSettingGsmNumber(data, "*99#")
	setSettingGsmPasswordFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	setSettingGsmPinFlags(data, NM_SETTING_SECRET_FLAG_NONE)
}

// Get available keys
func getSettingGsmAvailableKeys(data connectionData) (keys []string) {
	if getSettingVkMobileProvider(data) == mobileProviderValueCustom {
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_NUMBER)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_USERNAME)
		if isSettingRequireSecret(getSettingGsmPasswordFlags(data)) {
			keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_PASSWORD)
		}
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_APN)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_NETWORK_ID)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_HOME_ONLY)
		keys = appendAvailableKeys(data, keys, sectionGsm, NM_SETTING_GSM_PIN)
	}
	return
}

// Get available values
func getSettingGsmAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_GSM_PASSWORD_FLAGS:
		values = availableValuesSettingSecretFlags
	case NM_SETTING_GSM_APN:
	}
	return
}

// Check whether the values are correct
func checkSettingGsmValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	ensureSettingGsmApnNoEmpty(data, errs)
	ensureSettingGsmNumberNoEmpty(data, errs)
	return
}

func syncMoibleConnectionId(data connectionData) {
	// sync connection name
	if !isSettingSectionExists(data, sectionCache) {
		return
	}
	providerName := getSettingVkMobileProvider(data)
	if providerName == mobileProviderValueCustom {
		switch getSettingVkMobileServiceType(data) {
		case connectionMobileGsm:
			setSettingConnectionId(data, Tr("Custom")+" GSM")
		case connectionMobileCdma:
			setSettingConnectionId(data, Tr("Custom")+" CDMA")
		}
	} else {
		if plan, err := mobileprovider.UnmarshalPlan(getSettingVkMobilePlan(data)); err == nil {
			if plan.IsGSM {
				if len(plan.Name) > 0 {
					setSettingConnectionId(data, providerName+" "+plan.Name)
				} else {
					setSettingConnectionId(data, providerName+" "+Tr("Default"))
				}
			} else {
				setSettingConnectionId(data, providerName)
			}
		}
	}
}
