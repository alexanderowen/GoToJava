package human

type Human struct {
	age  int
	name string
}

func (h Human) GetAge() int {
	return h.age
}

func (h Human) GetName() string {
	return h.name
}
