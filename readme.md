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

## Changelog for 0.4.1
- The automatic version checker is now a bit better

## Gor Todo
- [x] Variable assignment
- [x] Expressions
- [x] Labels and `jumpto`
- [ ] Make my expression parsing not suck
- [ ] If statements
- [ ] Built-in function calls
- [ ] Custom functions and calls
- [ ] ~~Structs but written badly~~ containers

### Hey if you know of any optimazations I can do in the code base, make an issue(please I beg of you my code is so non-preformant)