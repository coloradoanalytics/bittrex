package bittrex

type BookRow struct {
	Quantity float64
	Price    float64 `json:"Rate"`
}

type BookRowUpdate struct {
	Quantity float64
	Price    float64 `json:"Rate"`
	Type     int     //0, 1, 2 means ADD, REMOVE, UPDATE
}

type BookRowsAscending []BookRow

func (r BookRowsAscending) Len() int {
	return len(r)
}

func (r BookRowsAscending) Less(i, j int) bool {
	return r[i].Price < r[j].Price
}

func (r BookRowsAscending) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type BookRowsDescending []BookRow

func (r BookRowsDescending) Len() int {
	return len(r)
}

func (r BookRowsDescending) Less(i, j int) bool {
	return r[i].Price > r[j].Price
}

func (r BookRowsDescending) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func separateRowUpdates(rows []BookRowUpdate) ([]BookRowUpdate, []BookRowUpdate, []BookRowUpdate) {
	add := []BookRowUpdate{}
	remove := []BookRowUpdate{}
	update := []BookRowUpdate{}
	for _, r := range rows {
		switch r.Type {
		case 0:
			add = append(add, r)
		case 1:
			remove = append(remove, r)
		case 2:
			update = append(update, r)
		}
	}
	return add, remove, update
}
