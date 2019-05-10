package database

import (
	"fmt"
	"testing"
)

func BenchmarkSetHash(b *testing.B) {
	p := DbPerson{30, 1, "一个普通人"}
	for i := 0; i < b.N; i++ {
		p.Age = i
		setPerson(fmt.Sprintf("gaoyang_%d", i), &p)
	}
}
func BenchmarkGetHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getPerson("gaoyang")
	}
}
func TestHash(t *testing.T) {
	p1 := DbPerson{30, 1, "一个普通人"}
	setPerson("gaoyang", &p1)

	p2 := getPerson("gaoyang")
	if p2 == nil {
		t.Error("person is nil")
	}
	if p1.Age != p2.Age ||
		p1.Desc != p2.Desc ||
		p1.Sex != p2.Sex {
		t.Errorf("%#v != %#v", p1, p2)
	}
	// if !reflect.DeepEqual(p1, p2) {

	// }
}
