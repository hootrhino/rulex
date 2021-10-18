%{
package parser

import (
    //"fmt"
)
%}

%union {
   field    string
   length   int
}

%token <field>    F
%token <length>   L

%token '[' ']' ',' ':' ' ' '\n'

%%

start: '[' pairs ']' {

};

pairs: pair | pairs pair
;

pair: F ':' L {
   $$=$1
}
;
%%