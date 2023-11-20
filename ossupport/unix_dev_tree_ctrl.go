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

package ossupport

import (
	"os"
	"regexp"
)

const devFolder = "/dev"
const regexFilter = "(ttyS|ttyHS|ttyUSB|ttyACM|ttyAMA|rfcomm|ttyO|ttymxc)[0-9]{1,3}"

func GetPortsListUnix() ([]string, error) {
	files, err := os.ReadDir(devFolder)
	if err != nil {
		return nil, err
	}
	ports := make([]string, 0, len(files))
	regex, err := regexp.Compile(regexFilter)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		// Skip folders
		if f.IsDir() {
			continue
		}
		if !regex.MatchString(f.Name()) {
			continue
		}
		portName := devFolder + "/" + f.Name()
		ports = append(ports, portName)
	}

	return ports, nil
}
