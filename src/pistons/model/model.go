package model

import "log"


type Model struct {
	Id int 
	Name string 
}



func SearchForModels(makeId int, term string) ([]Model, error) {
	result := []Model{}


	rows, err := db.Query("SELECT m.id, m.name FROM model m WHERE m.make_id = $1 AND lower(m.name) LIKE lower($2 || '%%') ORDER BY name",makeId, term)
	if err != nil {
		log.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			model := Model{}
			rows.Scan(&model.Id, &model.Name)
			result = append(result, model)
		}
	}

	return result, err 
}



func GetModel(modelId int) (*Model, error) {
	result := Model{}

	row := db.QueryRow("SELECT id, name FROM model m WHERE id = $1", modelId)
	err := row.Scan(&result.Id, &result.Name)

	if err != nil {
		log.Println(err)
	}

	return &result, err 
}


func FindYearsForModel(modelId int) ([]Year, error) {
	result := []Year{}
	rows, err := db.Query(
		"SELECT y.id, y.year "+
		"FROM year y "+
		"JOIN model_year my "+
		"  ON my.year_id = y.id "+
		"JOIN model m "+
		"  ON m.id = my.model_id "+
		"WHERE m.id = $1 "+
		"ORDER BY y.year", modelId)


	if err != nil {
		log.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			year := Year{}
			err := rows.Scan(&year.Id, &year.Value)
			if err != nil {
				log.Println(err)
			} else {
				result = append(result, year)
			}
		}
	}

	return result, err 
}