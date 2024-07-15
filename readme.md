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

First, install Go if you haven't already from [go.dev](go.dev)

Then, run this command:
```sh
go install "github.com/voidwyrm-2/CBuild@latest"
```

## Gor Todo
- [x] Variable assignment
- [x] Expressions
- [x] Labels and `jumpto`
- [ ] If statements
- [ ] Built-in function calls
- [ ] Custom functions
- [ ] Containers(basically just structs)