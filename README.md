# TJ - Text Jumble

Search tool etc, etc.

# EXAMPLE

 - Some code to search into
```php
function foo($bar, $baz) {
    $str = "A string"; // Some comment
    return $bar + $baz;
}
```

 - A "parser" definition
```
name:        /[a-zA-Z][A-Za-z0-9]/
body:        "{" ... "}"
class:       "class"
function:    "function" name
method:      class function
variable:    "$" name
call:        name "(" ... ")"
expression:  ... ";"
expression:  ... ","
return:      "return" expression
comment:     "//" ... "\n"
comment:     "/*" ... "*/"
string:      '"' ... '"' :quote "\"
```

 - Perform some searches
```bash
$ tj -method:ci endswith -body -return '$bar'
3:      return $bar + $baz;
$ tj -comment:ci "some"
2:    $str = "A string"; // Some comment
$ tj -comment:regex "[Ss]ome"
2:    $str = "A string"; // Some comment
```

# Implementation

 - Lexer: contains the definitions; map of name : list of definitions
   - Each definition is a list of lexTokens 
   - A lexToken can be:
     - A fixed string
     - A regular expression
     - A reference to another lexToken
     - An ellypsis
   - Rules are flattened out and can be saved and restored in binary format
 - Resolver: gathers which lexer rules are actually needed
   - Example #1: tj -method:ci endswith -body -return '$bar'
   - 'method' will need: 'class' followed by 'function'
   - Inside a 'body' (inside because there is the ellypsis)
   - Find a 'return'
   - Finally return this pattern: $class any $function $name any $body:begin any $return any $body:end
   - Example #2: tj -expression -variable 'foo' -variable 'baz'
   - Return pattern: '$' $name any '$' $name any ';' | '$' $name any '$' $name any ','
 - Matcher: matches text against rules from the resolver
   - Match each token as-is or fail to match the rule
   - If the token is 'any', match next token without failing;
     Fail if any other token after this one is encountered

# TODO

 - How to handle Python blocks?


