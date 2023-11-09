package apis

import "github.com/hootrhino/rulex/plugin/http_server/server"

// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

func LoadSystemSettingsAPI() {
	//
	// 系统设置
	//

	ifacesApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		ifacesApi.GET(("/ifaces"), server.AddRoute(GetNetInterfaces))
		ifacesApi.GET(("/uarts"), server.AddRoute(GetUartList))
	}
	settingsApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		// volume
		settingsApi.GET("/volume", server.AddRoute(GetVolume))
		settingsApi.POST("/volume", server.AddRoute(SetVolume))
	}
	// ethnet
	ethApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		ethApi.POST("/eth", server.AddRoute(SetEthNetwork))
		ethApi.GET("/eth", server.AddRoute(GetEthNetwork))
		ethApi.GET("/connection", server.AddRoute(GetCurrentNetConnection))
	}
	// wifi
	wifiApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		wifiApi.GET("/wifi", server.AddRoute(GetWifi))
		wifiApi.POST("/wifi", server.AddRoute(SetWifi))
		wifiApi.GET("/wifi/scan", server.AddRoute(ScanWIFIWithNmcli))
	}
	// time
	timesApi := server.RouteGroup(server.ContextUrl("/settings"))
	{
		// time
		timesApi.GET("/time", server.AddRoute(GetSystemTime))
		timesApi.POST("/time", server.AddRoute(SetSystemTime))
		timesApi.PUT("/ntp", server.AddRoute(UpdateTimeByNtp))
	}
	// 4g module
	settings4GApi := server.RouteGroup(server.ContextUrl("/4g"))
	{
		settings4GApi.GET("/info", server.AddRoute(Get4GBaseInfo))
		settings4GApi.POST("/restart", server.AddRoute(RhinoPiRestart4G))
		settings4GApi.GET("/apn", server.AddRoute(GetAPN))
		settings4GApi.POST("/apn", server.AddRoute(SetAPN))
	}
	// 软路由相关
	settingsSoftRouterApi := server.RouteGroup(server.ContextUrl("/softRouter"))
	{
		settingsSoftRouterApi.GET("/dhcp", server.AddRoute(GetDHCP))
		settingsSoftRouterApi.POST("/dhcp", server.AddRoute(SetDHCP))
		settingsSoftRouterApi.GET("/dhcp/clients", server.AddRoute(GetDhcpClients))
		// 默认 Ip route
		settingsSoftRouterApi.POST("/iproute", server.AddRoute(SetNewDefaultIpRoute))
		settingsSoftRouterApi.GET("/iproute", server.AddRoute(GetOldDefaultIpRoute))

	}
	// 固件
	settingsFirmware := server.RouteGroup(server.ContextUrl("/firmware"))
	{
		settingsFirmware.POST("/reboot", server.AddRoute(Reboot))
		settingsFirmware.POST("/recoverNew", server.AddRoute(RecoverNew))
		settingsFirmware.POST("/restartRulex", server.AddRoute(ReStartRulex))
		settingsFirmware.POST("/upload", server.AddRoute(UploadFirmWare))
		settingsFirmware.POST("/upgrade", server.AddRoute(UpgradeFirmWare))
		settingsFirmware.GET("/upgradeLog", server.AddRoute(GetUpGradeLog))
		settingsFirmware.GET("/vendorKey", server.AddRoute(GetVendorKey))
	}

}
