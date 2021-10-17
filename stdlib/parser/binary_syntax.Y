%{
package parser

import (
   "fmt"
)
%}

%union {
   field    string
   length   int
   operator uint8
}

%token <field>    F
%token <length>   N
%token <operator> '[' ':' ']' ','

%start expr

%%

expr: '[' F ':' N ']' {
   fmt.Println("=>",$1, $2)
};
%%