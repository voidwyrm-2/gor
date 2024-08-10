# Gor
Gor(from "Go R") is a dynamically-typed, R-inspired language

Here's a "hello world" example:
```
? comments are done with '?'

? assign our message to the variable 'hello'
hello <- "Hello, Catdog!"

? print 'hello' to the screen
puts(hello)
```

## Installation
**Requirements**
* Go >= 1.22.1

First, install Go if you haven't already from [go.dev](go.dev) or using homebrew(`brew install go`)

Then, run this command:
```sh
go install "github.com/voidwyrm-2/gor@latest"
```

## Changelog for 0.5(aka, the "WOW I CAN WRITE GO BETTER THAN A MONKEY, ISN'T THAT INCREDIBLE?" update)
- Expressions are actually usable(no paranthese though(it scarwy))
- Removed a bunch of bloat from the main.go file
- Streamlined the parser a bit
- IF STATEMENTS WOOOOOOOOOOOO

## Gor Todo
- [x] Variable assignment
- [x] Expressions
- [x] Labels and `jumpto`
- [x] Make my expression parsing not suck
- [x] If statements
- [ ] Built-in function calls
- [ ] Custom functions and calls
- [ ] ~~Structs but written badly~~ containers

### Hey if you know of any optimazations I can do in the code base, make an issue(please I beg of you my code is so non-preformant)