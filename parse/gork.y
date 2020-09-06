%{
package main

import (
	"fmt"
	"os"

    "github.com/imdario/mergo"
)

%}

%union {
  n  int
  s  string
  a  []string
  d  Descriptions
  g  Game
  l  Location
  t  Traits
  x  Block
}


%left COMMA
%token AS BEGIN BY DATE DOWN EAST END EXIT GAME INTRO IS LOCATION 
%token NORTH NOT SOUTH START TITLE TO UP VERSION WEST

%token <n> NUM 
%token <s> NAME STRING VNUM
%token <a> STRINGS

%type <g> gexps gexp
%type <l> lexps lexp
%type <t> traitlist trait
%type <d> xstrings
%type <x> exe

%%
game:	/* empty */
		| GAME BEGIN gexps END {
//			fmt.Println($3)
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
			$$.Description.Merge($2.Description)
			$$.By.Merge($2.By)
			$$.Date.Merge($2.Date)
			$$.Version.Merge($2.Version)
		}
		;

gexp:	/* empty */ { $$ = $$ }
		| xstrings { $$.Description.Merge($1) }
		| TITLE STRINGS { $$.Title.AddStringArray($2) }
		| BY STRINGS { $$.By.AddStringArray($2)}
		| DATE STRINGS { $$.Date.AddStringArray($2) }
		| VERSION VNUM { $$.Version.AddString($2)}
		| START NAME { $$.Current = $2 }
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
		| STRINGS { $$.Description.AddStringArray($1) }
		;

traitlist: trait { $$=$1 }
		| traitlist COMMA trait {$1.Merge($3);$$=$1}
		;

trait:	NOT NAME { $$=Traits{$2:false} }
		| NAME { $$=Traits{$1:true} }
		;

xstrings: STRINGS { $$.AddStringArray($1) }
		| BEGIN exe END { $$=$2.Eval() }
		;

exe: 	/* empty */ { $$ = $$ }
		| STRINGS { $$.AddStringArray($1);fmt.Printf("[$1:%v]]",$1);fmt.Printf("[$$:%v]]",$$) }
		;

%%
