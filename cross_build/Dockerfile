# Copyright (C) 2024 wwhai
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
FROM golang:1.21

RUN apt update && \
    apt upgrade && \
    apt install -y sudo make jq


WORKDIR  "/go/src"
RUN git clone https://github.com/hybridgroup/gocv.git
WORKDIR "/go/src/gocv"
# RUN make install
ADD ../ "/go/src/rulex"
WORKDIR "/go/src/rulex"
# RUN ["go","build"]
RUN "bash"