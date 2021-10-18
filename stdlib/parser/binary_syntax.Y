%{
      package parser
      import(
         "fmt"
      )

      type Pair struct {
         k string
         l int
      }
      func (p Pair) String() string{
         return fmt.Sprintf("@Pair[%s,%v]", p.k, p.l)
      }

%}

%union {
   k string
   l int
   pair   Pair
}

%type <pair>      pair


%token  <k>        k
%token  <l>        l
%token  lb         '<'
%token  rb         '>'
%token  as         ':'

%start pair

%%

pair:  lb k as l rb
    {
     $$ = Pair{k: $2, l: $4}
    }
    ;

%%