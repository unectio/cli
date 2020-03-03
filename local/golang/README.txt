Aux environment to compile and test your functions locally.

To use:
- copy you code.go here
- type make
- run ./function

The function binary is compiled with supported options:

-h        : Print this help
-m        : Method (e.g. POST or GET)
-p        : URL path to pass to function
-b <file> : File from which to read body (- is stdin)
-t        : Content-type to set (works with -b only)
-k        : Key to use
-c        : Set of claims

Unparsed (so called "positional") arguments are expected to be
in foo=bar format and are passed as request arguments
