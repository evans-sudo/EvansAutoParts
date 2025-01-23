package model

type Year struct {
	Id int 
	Value int 
}


func GetYear(yearId int) (*Year, error) {
	result := Year{}

	row := db.QueryRow("SELECT y.id, y.year FROM year y WHERE y.id = $1", yearId)

	err := row.Scan(&result.Id, &result.Value)

	return &result, err 
}