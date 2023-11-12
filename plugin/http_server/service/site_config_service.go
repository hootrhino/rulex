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

package service

import (
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

func GetSiteConfig() (model.MSiteConfig, error) {
	m := model.MSiteConfig{}
	if err := interdb.DB().Where("uuid=0").First(&m).Error; err != nil {
		return model.MSiteConfig{}, err
	} else {
		return m, nil
	}
}

// 创建 SiteConfig
func InitSiteConfig(SiteConfig model.MSiteConfig) error {
	SiteConfig.UUID = "0" // 默认就一个配置
	return interdb.DB().FirstOrCreate(&SiteConfig).Error
}

// 更新 SiteConfig
func UpdateSiteConfig(SiteConfig model.MSiteConfig) error {
	SiteConfig.UUID = "0" // 默认就一个配置
	return interdb.DB().Where("uuid=0").Model(SiteConfig).Updates(SiteConfig).Error
}
