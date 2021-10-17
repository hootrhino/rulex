%{
package parser

import (
    "fmt"
)
%}

%token '[' ':' ']' ','
%token LENGTH [0-9]+
%token FIELD [a-Z]+

%start main

%%

main: '[' KL ']' {

} | '[' KL ',' KL ']' {

}
KL : FIELD ':' LENGTH {
  
}
;