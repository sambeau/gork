%{
package main

import (
	// "fmt"
	// "os"

    "github.com/imdario/mergo"
)

var game Game

%}

%union {
  n int
  s string
  t Text
  h Handle
  d Texts
  g Game
  l Location
  i Traits
  x Block
  e Node
  // k Token
}


%left COMMA
%token <k> AS BEGIN BY DATE DOWN EAST ELSE END EXIT GAME IF INTRO IS LB LOCATION 
%token <k> NORTH NOT QUOTE RB SAY SOUTH START THEN TITLE TO UP VERSION WEST

%token <n> NUM 
%token <s> NAME
%token <t> TEXT VNUM 
%token <h> HANDLEBAR

%type <g> gexps gexp
%type <l> lexps lexp
%type <i> traitlist trait
%type <d> text texts
%type <x> expression
%type <e> node cond

%%
game:	/* empty */
		| GAME BEGIN gexps END {
			// fmt.Println($3)
			// fmt.Println()
			// $3.Pprint(os.Stdout,"")
			// fmt.Println()
			game = $3
		}
		;

gexps:	gexp 
		| gexps gexp { 
			mergo.Merge(&$$, $1); 
			mergo.Merge(&$$, $2);

			$$.Title.Merge($2.Title)
			$$.By.Merge($2.By)
			$$.Date.Merge($2.Date)
			$$.Version.Merge($2.Version)

			$$.Description.Merge($2.Description)

		}
		;

gexp:	/* empty */ { $$ = $$ }
		| TITLE texts { $$.Title.Add($2) }
		| BY texts { $$.By.Add($2) }
		| DATE texts { $$.Date.Add($2) }
		| VERSION VNUM {$$.Version.Add($2) }
		| START NAME { $$.Current = $2 }
		| IS traitlist {
			if($$.Traits == nil){
				$$.Traits = Traits{}
			}
			$$.Traits.Merge($2)
		}
		| texts { $$.Description.Add($1) }
		| LOCATION NAME BEGIN lexps END {
			if($$.Locations == nil){
				$$.Locations = map[string]Location{}
			}
			$4.Name = $2
			$$.Locations[$2] = $4
		}
		;

lexps:	lexp
		| lexps lexp {
			mergo.Merge(&$1, $2); 
			mergo.Merge(&$$, $1);
			$$.Description.Merge($2.Description)
		}
		;

lexp:	/* empty */ { $$ = $$ }
		| IS traitlist {
			if($$.Traits == nil){
				$$.Traits = Traits{}
			}
			$$.Traits.Merge($2)
		}
		| texts { $$.Description.Add($1) }
		;

traitlist: trait { $$=$1 }
		 | traitlist COMMA trait {$1.Merge($3);$$=$1}
		 ;

trait   : NOT NAME { $$=Traits{$2:false} }
		| NAME { $$=Traits{$1:true} }
		;

texts   : QUOTE text QUOTE { $$=$2 }
		// | BEGIN expression END {$$.Add($2)}
		| LB expression RB {$$.Add($2)}
		;

text    : /* empty */ {}
		| text TEXT { $$=$1;$$.Add($2) } 
		| text HANDLEBAR { $$=$1;$$.Add($2) } 
		| TEXT { $$.Add($1) }
		| HANDLEBAR { $$.Add($1) }
	    ;

expression: /* empty */ {}
		| expression node {$$=$1;$$.Add($2)}
		| node {$$.Add($1)}
		;

node:  texts {$$.node = $1}
		| IF cond THEN expression {$$.node=IfNode{Cond:$2, Then:$4}}
		| IF cond THEN expression ELSE expression {$$.node=IfNode{Cond:$2, Then:$4, Else: $6}}
		;

cond: 	IS NAME {$$.node = IsCond{Pos{0,0},$2,true}}
		| IS NOT NAME {$$.node = IsCond{Pos{0,0},$3,false}}
		;
%%
