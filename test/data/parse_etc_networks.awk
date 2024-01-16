# // Copyright (C) 2024 wwhai
# //
# // This program is free software: you can redistribute it and/or modify
# // it under the terms of the GNU Affero General Public License as
# // published by the Free Software Foundation, either version 3 of the
# // License, or (at your option) any later version.
# //
# // This program is distributed in the hope that it will be useful,
# // but WITHOUT ANY WARRANTY; without even the implied warranty of
# // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# // GNU Affero General Public License for more details.
# //
# // You should have received a copy of the GNU Affero General Public License
# // along with this program.  If not, see <http://www.gnu.org/licenses/>.

#!/usr/bin/awk -f

# Start with array creation
BEGIN {
    printf "[";
}

# New lease: start object
/^lease/ {
    # If that ip is unknown, create a new JSON object
    if (!known[$2]) {
        # if this object is not the first, print a comma
        if (!notFirst) {
            notFirst=1;
        } else {
            printf ",";
        }

        # create a new JSON object with the first key being the IP (column 2)
        printf "{\"ip\":\"%s\"", $2; known[$2]=1;

        # print subsequent lines, see below
        p=1;
    }
}

# If printing is enabled print line as a JSON key/value pair
p && /^  / {
    # Print key (first word)
    printf ",\"%s\":", $1;

    # Clean up the rest of the line: trim whitespace and trailing ;, remove " and escape \
    $1="";
    gsub(/\\/, "\\\\", $0);
    gsub(/"/, "", $0);
    gsub(/^[\t ]*/, "", $0);
    gsub(/;$/, "", $0);
    printf "\"%s\"", $0;
}

# End of lease: close JSON object and disable printing
/^\}$/ {
    if (p) {
        printf "}"
    }
    p=0
}

# Close the JSON array
END {
    print "]";
}