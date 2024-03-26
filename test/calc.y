%{
package main
import "fmt"
%}

%union {
num float64
}

%token <num> NUMBER
%left '+' '-'
%left '*' '/'
%right UMINUS

%type <num> exp

%%

input:
| input line
;

line:
exp '
' { fmt.Println($1) }
| error '
' { yyerrok }
;

exp:
NUMBER
| exp '+' exp { $1 + $3 }
| exp '-' exp { $1 - $3 }
| exp '*' exp { $1 * $3 }
| exp '/' exp { $1 / $3 }
| '-' exp %prec UMINUS { -$2 }
;

%%

func main() {
fmt.Println("Enter an expression:")
yyParse(NewLexer(os.Stdin))
}


