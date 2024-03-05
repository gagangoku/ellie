grammar Ellie;

prog: lhs '=' expression EOF;

lhs: TERM_TABLE | TERM_TABLE_COL | TERM_DB | TERM_DB_COL | TERM_CONST | TERM_STAT;

expression: expression logicalFn expression
 | expression mathFn expression
 | expression setFn expression
 | function
 | bool
 | TEXT
 | NUMBER
 | TERM_TABLE_COL | TERM_TABLE | TERM_DB_COL | TERM_DB | TERM_CONST | TERM_STAT
 | ID
 | '(' expression ')'
 ;

function: ID '(' arguments? ')';

logicalFn: 'AND' | 'OR';
mathFn: '+' | '-' | '/' | '*' | '^' | '<' | '<=' | '>' | '>=' | '==' | '!=';
setFn: 'CONTAINS' | 'IN' ;

arguments: expression ( ',' expression )*;

TERM_TABLE     : ( 'tbl.' [_a-zA-Z][_a-zA-Z0-9]* '.' [_a-zA-Z][_a-zA-Z0-9]* );
TERM_TABLE_COL : ( 'tbl.' [_a-zA-Z][_a-zA-Z0-9]* );
TERM_DB        : ( 'db.' [_a-zA-Z][_a-zA-Z0-9]* );
TERM_DB_COL    : ( 'db.' [_a-zA-Z][_a-zA-Z0-9]* '.' [_a-zA-Z][_a-zA-Z0-9]* );
TERM_CONST     : ( 'const.' [_a-zA-Z][_a-zA-Z0-9]* );
TERM_STAT      : ( 'stat.' [_a-zA-Z][_a-zA-Z0-9]* );

bool           : TRUE | FALSE;
TRUE           : 'true' | 'TRUE';
FALSE          : 'false' | 'FALSE';
NUMBER         : '-'? ( [0-9]* '.' )? [0-9]+;
ID             : [a-zA-Z_] [a-zA-Z0-9_]*;
TEXT           : '\'' ~[\r\n']* '\'' | '"' ~[\r\n"]* '"';
SPACE          : [ \t\r\n]+ -> skip;
