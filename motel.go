package motel_reservation

import (
	"errors"
	"sort"
	"time"
)

type Room struct {
	Beds  uint32
	Floor uint32
}

type Reservation struct {
	Beds               uint32
	Arrival            time.Time
	Departure          time.Time
	Pets               uint32
	HandicapAccessible bool
}

type Motel interface {
	Reserve(request Reservation) (cost float64, err error)
}

type motel struct {
	rooms        []Room
	reservations []Reservation
}

type eventType int

const (
	Arrival eventType = iota
	Departure
)

type event struct {
	date time.Time
	typeOfEvent eventType
}

func (this *motel) Reserve(request Reservation) (cost float64, err error) {
	reservationsInScope := filterReservations(request, this.reservations)

	// Construct events for arrivals & departures through the stay
	events := make([]event, 0)
	for _, reservation := range reservationsInScope {
		events = append(events, event{reservation.Arrival, Arrival})
		events = append(events, event{reservation.Departure, Departure})
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].date.Before(events[j].date)
	})
	events = filterEvents(request, events)

	// Compute occupancy throughout stay. Note that we're also checking validity *before* the reservation time - starting
	// occupancy count must be computed anyway.
	eligableRooms := len(filterRooms(request, this.rooms))
	if eligableRooms == 0 {
		return 0.0, errors.New("No eligible rooms")
	}

	occupancyCount := 0
	for _, event := range events {
		if event.typeOfEvent == Arrival {
			occupancyCount += 1
		} else {
			occupancyCount -= 1
		}

		if occupancyCount >= eligableRooms {
			return 0.0, errors.New("Overbooked!")
		}
	}

	return computeCost(request), nil
}

func filterReservations(request Reservation, reservations []Reservation) (ret []Reservation) {
	for _, reservation := range reservations {
		if requestPredicate(request, reservation) {
			ret = append(ret, reservation)
		}
	}
	return
}

func requestPredicate(request Reservation, reservation Reservation) bool {
	var floor uint32
	if requiresFirstFloor(reservation) {
		floor = 1
	} else {
		floor = 2
	}
	return roomPredicate(Room{reservation.Beds, floor}, request)
}

func roomPredicate(room Room, request Reservation) bool {
	if (request.Pets) > 2 {
		return false
	} else {
		return request.Beds == room.Beds && (!requiresFirstFloor(request) || room.Floor == 1)
	}
}

func filterRooms(request Reservation, rooms []Room) (ret []Room) {
	for _, room := range rooms {
		if roomPredicate(room, request) {
			ret = append(ret, room)
		}
	}
	return
}

func filterEvents(request Reservation, events []event) (ret []event) {
	for _, event := range events {
		if coincident(event.date, request.Departure) {
			ret = append(ret, event)
		}
	}
	return
}

func coincident(one time.Time, two time.Time) bool {
	return one.Before(two) || one == two
}

func requiresFirstFloor(request Reservation) bool {
	return request.Pets > 0 || request.HandicapAccessible
}

func computeCost(request Reservation) float64 {
	var base uint32
	switch request.Beds {
	case 1: base = 50
	case 2: base = 75
	case 3: base = 90
	}
	petCost := request.Pets * 20
	return float64(base + petCost)
}
