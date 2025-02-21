package model

import "log"

type Make struct {
	Id   int
	Name string
}

func SearchForMakes(term string) ([]Make, error) {
	result := []Make{}
	rows, err := db.Query("SELECT m.id, m.name FROM make m WHERE lower(m.name) LIKE lower($1 || '%%') ORDER BY m.name", term)

	if err != nil {
		log.Print("Error searching database: " + err.Error())
	} else {
		defer rows.Close()
		for rows.Next() {
			make := Make{}
			rows.Scan(&make.Id, &make.Name)
			result = append(result, make)
		}
	}

	return result, err 
}


func GetMake(id int) (*Make, error) {
	result := Make{}
	row := db.QueryRow("SELECT m.id, m.name FROM make m WHERE m.id = $1", id)

	err := row.Scan(&result.Id, result.Name)

	if err != nil {
		log.Println(err)
	}

	return &result, err
}
