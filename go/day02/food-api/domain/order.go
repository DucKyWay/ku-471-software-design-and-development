package domain

type Status string

const (
	Pending    Status = "pending"
	Confirmed  Status = "confirmed"
	Packing    Status = "packing"
	Cooking    Status = "cooking"
	Delivering Status = "delivering"
	Completed  Status = "completed"
	Cancelled  Status = "cancelled"
)

var flow = []Status{Pending, Confirmed, Packing, Cooking, Delivering, Completed}

type OrderLine struct {
	Item MenuItem `json:"item"`
	Qty  int      `json:"qty"`
}

func (line OrderLine) Total() float64 { return line.Item.Price * float64(line.Qty) }

type Order struct {
	ID     string      `json:"id"`
	Status Status      `json:"status"`
	Lines  []OrderLine `json:"lines"`
}

func (order Order) Total() float64 {
	total := 0.0
	for _, line := range order.Lines {
		total += line.Total()
	}
	return total
}

func (order Order) CanCancel() bool {
	return order.Status == Pending || order.Status == Confirmed
}

func (order Order) CanAdvanceTo(next Status) bool {
	for i, status := range flow {
		if status == order.Status {
			return i+1 < len(flow) && flow[i+1] == next
		}
	}
	return false
}

func (status Status) Valid() bool {
	switch status {
	case Pending, Confirmed, Packing, Cooking, Delivering, Completed, Cancelled:
		return true
	default:
		return false
	}
}
