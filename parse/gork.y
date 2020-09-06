%{
package main

import (
	"fmt"
	"os"

    "github.com/imdario/mergo"
)

%}

%union {
  n int
  s string
  t Text
  d Descriptions
  g Game
  l Location
  i Traits
  // x Block
  k Token
}


%left COMMA
%token <k> AS BEGIN BY DATE DOWN EAST END EXIT GAME INTRO IS LOCATION 
%token <k> NORTH NOT QUOTE SOUTH START TITLE TO UP VERSION WEST

%token <n> NUM 
%token <s> NAME
%token <t> TEXT VNUM

%type <g> gexps gexp
%type <l> lexps lexp
%type <i> traitlist trait
%type <d> text texts
// %type <x> exe

%%
game:	/* empty */
		| GAME BEGIN gexps END {
			fmt.Println($3)
			fmt.Println()
			$3.Pprint(os.Stdout,"")
			fmt.Println()
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

trait:	NOT NAME { $$=Traits{$2:false} }
		| NAME { $$=Traits{$1:true} }
		;

texts: QUOTE text QUOTE { $$=$2 }
		;

text: TEXT { $$.Add($1) }
     | text TEXT { $$=$1;$$.Add($2) } 
     ;

// exe: 	/* empty */ { $$ = $$ }
// 		| text { $$.Add($1);fmt.Printf("[$1:%v]]",$1);fmt.Printf("[$$:%v]]",$$) }
// 		;

%%
