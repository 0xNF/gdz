package main

import "github.com/0xNF/gdz/internal/fs"

func main() {
	c := fs.NewConfig()
	fs.Get(&c)
}
