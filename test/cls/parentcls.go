package main

import (
	"fmt"
	"os"
)

type Person struct {
	name string
	age  int
}

type Police struct {
	Person
	serverage int
}

func (p *Person) Call() {
	fmt.Fprintf(os.Stdout, "call person [%s] age[%d]\n", p.name, p.age)
	return
}

func (p *Police) Call() {
	p.Person.Call()
	fmt.Fprintf(os.Stdout, "call police serverage [%d]\n", p.serverage)
}

func NewPerson(name string, age int) *Person {
	p := &Person{name: name, age: age}
	return p
}

func NewPolice(name string, age int, serverage int) *Police {
	p := &Police{Person: *NewPerson(name, age), serverage: serverage}
	return p
}

func main() {
	p := NewPolice("jack", 33, 5)
	p.Call()
}
