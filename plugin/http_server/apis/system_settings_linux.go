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
		settingsApi.POST("/eth", server.DefaultApiServer.AddRoute(SetEthNetwork))
		settingsApi.GET("/eth", server.DefaultApiServer.AddRoute(GetEthNetwork))
		settingsApi.GET("/connection", server.DefaultApiServer.AddRoute(GetCurrentNetConnection))
		//
		settingsApi.GET("/time", server.DefaultApiServer.AddRoute(GetSystemTime))
		settingsApi.POST("/time", server.DefaultApiServer.AddRoute(SetSystemTime))
		//
		settingsApi.GET("/wifi", server.DefaultApiServer.AddRoute(GetWifi))
		settingsApi.POST("/wifi", server.DefaultApiServer.AddRoute(SetWifi))
		//
		settingsApi.GET("/volume", server.DefaultApiServer.AddRoute(GetVolume))
		settingsApi.POST("/volume", server.DefaultApiServer.AddRoute(SetVolume))
		//
		settingsApi.POST("/timezone", server.DefaultApiServer.AddRoute(SetSystemTimeZone))
		settingsApi.GET("/timezone", server.DefaultApiServer.AddRoute(GetSystemTimeZone))
		settingsApi.POST("/iproute", server.DefaultApiServer.AddRoute(SetDefaultRoute))
		settingsApi.PUT("/ntp", server.DefaultApiServer.AddRoute(UpdateTimeByNtp))
	}
}
