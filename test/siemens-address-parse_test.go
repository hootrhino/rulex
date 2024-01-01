// Copyright (C) 2024 wwhai
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

package test

import (
	"testing"

	"github.com/hootrhino/rulex/utils"
)

// go test -timeout 30s -run ^Test_parse_siemens_address github.com/hootrhino/rulex/test -v -count=1
func Test_parse_siemens_address(t *testing.T) {
	I := "I0.1"
	Q := "Q0.1"
	DBD := "DB1.DBD12"
	DBX := "DB2.DBX34"

	{
		V, err := utils.ParseADDR_I(I)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(V)
	}
	{
		V, err := utils.ParseADDR_Q(Q)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(V)
	}
	{
		V, err := utils.ParseDB(DBD)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(V)
	}
	{
		V, err := utils.ParseDB(DBX)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(V)
	}
}

// go test -timeout 30s -run ^Test_parse_siemens_address_error github.com/hootrhino/rulex/test -v -count=1
func Test_parse_siemens_address_error(t *testing.T) {
	I := "I0"
	Q := "Q0"
	DBD := "DB1"
	DBX := "DB2"

	{
		V, err := utils.ParseADDR_I(I)
		if err != nil {
			t.Log(err)
		}
		t.Log(V)
	}
	{
		V, err := utils.ParseADDR_Q(Q)
		if err != nil {
			t.Log(err)
		}
		t.Log(V)
	}
	{
		V, err := utils.ParseDB(DBD)
		if err != nil {
			t.Log(err)
		}
		t.Log(V)
	}
	{
		V, err := utils.ParseDB(DBX)
		if err != nil {
			t.Log(err)
		}
		t.Log(V)
	}
}
