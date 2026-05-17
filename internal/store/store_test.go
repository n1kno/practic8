package store

import (
	"testing"
)

func TestAddAndGetAll(t *testing.T) {
	// очистим records перед тестом (не экспортируется, поэтому через Add не сбрасываем)
	// можно протестировать последовательное добавление
	initialLen := len(GetAll())
	Add(Record{FirstName: "A", LastName: "B", Profession: "C"})
	records := GetAll()
	if len(records) != initialLen+1 {
		t.Errorf("expected %d records, got %d", initialLen+1, len(records))
	}
	last := records[len(records)-1]
	if last.FirstName != "A" || last.Profession != "C" {
		t.Errorf("unexpected record: %+v", last)
	}
}
