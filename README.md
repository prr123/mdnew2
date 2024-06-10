# github markdown parser

why? goldmark lacks documentation 
gomarkdown does not have extendible inline parsers

parser does two parses

1. blocks
2. each block is parsed for inline chars (inline parsers)

- creates an ast tree
- uses block parsers for different blocks
- all parsers are extendable
- block parsers are stored in a go map triggered by the first letter of a line
- inline parsers are stored in the inline map


