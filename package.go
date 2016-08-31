/*
Package jq queries (deeply nested) json documents.

The functions come in pairs:  there is a wordier
XxxError function that will report errors, and a
convenience Xxx function that will not.

For both functions, a missing object in the json
is NOT an error.  Only objects of wrong types are.
Say if you want to get at a boolean, but you get
a collection instead.
*/
package jq
