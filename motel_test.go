package motel_reservation

import (
	"testing"
	"time"
)

var rooms = []Room {
	Room{1, 1},
	Room{2, 1},
	Room{2, 1},
	Room{2, 1},
	Room{2, 2},
	Room{3, 2},
}

func parse(date string) time.Time {
	t, _ := time.Parse("2006-Jan-02", date)
	return t
}

func TestValidReservation(t *testing.T) {
	reservations := []Reservation{
		Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-05"), 0, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-10"), 0, false},
	}
	request := Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false}
	m := motel{rooms, reservations}
	cost, err := m.Reserve(request)
	if err != nil {
		t.Error("Failed")
	}
	if (cost != 75.0) {
		t.Errorf("Incorrect cost: %f", cost)
	}
}

func TestValidReservation2(t *testing.T) {
	reservations := []Reservation{
		Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
		Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-15"), 0, false},
		Reservation{3, parse("2020-Jan-01"), parse("2020-Jan-15"), 0, false},
	}
	request := Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false}
	m := motel{rooms, reservations}
	cost, err := m.Reserve(request)
	if err != nil {
		t.Error("Failed")
	}
	if (cost != 75.0) {
		t.Errorf("Incorrect cost: %f", cost)
	}
}

func TestRejectOverbooking(t *testing.T) {
	reservations := []Reservation{
		Reservation{3, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
		Reservation{3, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
	}
	request := Reservation{3, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false}
	m := motel{rooms, reservations}
	_, err := m.Reserve(request)
	if err == nil {
		t.Error("Failed")
	}
}

func TestRejectOverbooking2(t *testing.T) {
	reservations := []Reservation{
		Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false},
	}
	request := Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, false}
	m := motel{rooms, reservations}
	_, err := m.Reserve(request)
	if err == nil {
		t.Error("Failed")
	}
}

func TestRejectHandicapOverbooking(t *testing.T) {
	reservations := []Reservation{
		Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, true},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, true},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, true},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, true},
	}
	request := Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 0, true}
	m := motel{rooms, reservations}
	_, err := m.Reserve(request)
	if err == nil {
		t.Error("Failed")
	}
}

func TestRejectPetOverbooking(t *testing.T) {
	reservations := []Reservation{
		Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 2, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 2, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 2, false},
		Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 2, false},
	}
	request := Reservation{2, parse("2020-Jan-01"), parse("2020-Jan-03"), 2, false}
	m := motel{rooms, reservations}
	_, err := m.Reserve(request)
	if err == nil {
		t.Error("Failed")
	}
}

func TestMaxPets(t *testing.T) {
	var reservations []Reservation
	request := Reservation{1, parse("2020-Jan-01"), parse("2020-Jan-03"), 3, false}
	m := motel{rooms, reservations}
	_, err := m.Reserve(request)
	if err == nil {
		t.Error("Failed")
	}
}
