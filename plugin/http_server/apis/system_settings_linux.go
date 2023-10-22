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
	settingsApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/settings"))
	{
		// volume
		settingsApi.GET("/volume", server.DefaultApiServer.AddRoute(GetVolume))
		settingsApi.POST("/volume", server.DefaultApiServer.AddRoute(SetVolume))
	}
	// ethnet
	ethApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/settings"))
	{
		ethApi.POST("/eth", server.DefaultApiServer.AddRoute(SetEthNetwork))
		ethApi.GET("/eth", server.DefaultApiServer.AddRoute(GetEthNetwork))
		ethApi.GET("/connection", server.DefaultApiServer.AddRoute(GetCurrentNetConnection))
	}
	// wifi
	wifiApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/settings"))
	{
		wifiApi.GET("/wifi", server.DefaultApiServer.AddRoute(GetWifi))
		wifiApi.POST("/wifi", server.DefaultApiServer.AddRoute(SetWifi))
	}
	// time
	timesApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/settings"))
	{
		// time
		timesApi.GET("/time", server.DefaultApiServer.AddRoute(GetSystemTime))
		timesApi.POST("/time", server.DefaultApiServer.AddRoute(SetSystemTime))
		// timezone
		timesApi.POST("/timezone", server.DefaultApiServer.AddRoute(SetSystemTimeZone))
		timesApi.GET("/timezone", server.DefaultApiServer.AddRoute(GetSystemTimeZone))
		timesApi.PUT("/ntp", server.DefaultApiServer.AddRoute(UpdateTimeByNtp))
	}
	// 4g module
	settings4GApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/4g"))
	{
		settings4GApi.GET("/info", server.DefaultApiServer.AddRoute(Get4GBaseInfo))
		settings4GApi.GET("/csq", server.DefaultApiServer.AddRoute(Get4GCSQ))
		settings4GApi.GET("/cops", server.DefaultApiServer.AddRoute(Get4GCOPS))
		settings4GApi.GET("/iccid", server.DefaultApiServer.AddRoute(Get4GICCID))
	}
	// 软路由相关
	settingsSoftRouterApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/softRouter"))
	{
		settingsSoftRouterApi.GET("/dhcpClients", server.DefaultApiServer.AddRoute(GetDhcpClients))
		settingsSoftRouterApi.POST("/iproute", server.DefaultApiServer.AddRoute(SetNewDefaultIpRoute))
		settingsSoftRouterApi.GET("/iproute", server.DefaultApiServer.AddRoute(GetOldDefaultIpRoute))

	}
	// 固件
	settingsFirmware := server.DefaultApiServer.GetGroup(server.ContextUrl("/firmware"))
	{
		settingsFirmware.POST("/reboot", server.DefaultApiServer.AddRoute(Reboot))
		settingsFirmware.POST("/restartRulex", server.DefaultApiServer.AddRoute(ReStartRulex))
		settingsFirmware.POST("/upload", server.DefaultApiServer.AddRoute(UploadFirmWare))
	}

}
