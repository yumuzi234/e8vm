## About constants

Currently the consts are only const integers.
true and false should be consts, but are not.
string is not const, which is okay in our current settings.

we need a const evaulation system
consts now have:
- untyped numbers
- booleans
- typed numbers
- possibly floats in the future

current all consts are just untyped numbers, and whenever a compare occurs,
they are changed into bool vars.
this is not right.

here is the question, where should the value be saved in consts?
as a const ref? or as a type?
