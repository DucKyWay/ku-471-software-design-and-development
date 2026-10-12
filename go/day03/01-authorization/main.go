package main

import "fmt"

type Role string

const (
	Customer Role = "customer"
	Staff    Role = "staff"
	Manager  Role = "manager"
)

type User struct {
	ID   string
	Role Role
}

func (u User) IsStaff() bool      { return u.Role == Staff }
func (u User) CanCancelAny() bool { return u.Role == Staff || u.Role == Manager }

type Status string

const (
	Pending   Status = "pending"
	Confirmed Status = "confirmed"
	Cooking   Status = "cooking"
)

type Order struct {
	UserID string
	Status Status
}

func (o Order) OwnedBy(id string) bool { return o.UserID == id }
func (o Order) CanCancel() bool        { return o.Status == Pending || o.Status == Confirmed }

func (o Order) CanBeCancelledBy(actor User) bool {
	if !o.CanCancel() {
		return false
	}
	return actor.CanCancelAny() || o.OwnedBy(actor.ID)
}

func (o Order) CanChangeStatusBy(actor User) bool { return actor.IsStaff() }

func main() {
	alice := User{ID: "u1", Role: Customer}
	bob := User{ID: "u2", Role: Customer}
	manager := User{ID: "u3", Role: Manager}
	order := Order{UserID: alice.ID, Status: Pending}

	fmt.Printf("owner can cancel: %t\n", order.CanBeCancelledBy(alice))
	fmt.Printf("other customer can cancel: %t\n", order.CanBeCancelledBy(bob))
	fmt.Printf("manager can cancel another user's order: %t\n", order.CanBeCancelledBy(manager))
	fmt.Printf("manager can change status: %t\n", order.CanChangeStatusBy(manager))
	fmt.Printf("manager can cancel after cooking starts: %t\n", (Order{UserID: alice.ID, Status: Cooking}).CanBeCancelledBy(manager))
}
