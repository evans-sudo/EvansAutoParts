package model

import "log"



type StoreInventory struct {
	Location string `json:"location"`
	Quantity int    `json:"quantity"`
}

type Part struct {
	Id                 int
	PartNumber         string
	Price              float32
	Supplier           string
	Quality            string
	ImageName          string
	StoreInventory     *StoreInventory
	WarehouseInventory int
	TypeName           string
}

type Type struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	subcategoryId int
}

type Subcategory struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	categoryId int
	Types      []*Type `json:"types"`
}

type Category struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Subcategories []*Subcategory `json:"subcategories"`
}



func GetPartCategories() ([]*Category, error) {
	result := []*Category{}
	rows, err := db.Query(
		"SELECT c.id, c.name " +
			"FROM part_category c") 


	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			cat := &Category{}
			cat.Subcategories = []*Subcategory{}
			rows.Scan(&cat.Id, &cat.Name)
			result = append(result, cat)
		}
		rows.Close()

		rows, err := db.Query(
			"SELECT s.id, s.name, s.category_id " +
				"FROM part_subcategory s")


		subcategories := []*Subcategory{}
		if err != nil {
			log.Println(err)
		} else {
			for rows.Next() {
				subcat := &Subcategory{}
				subcat.Types = []*Type{}
				rows.Scan(&subcat.Id, &subcat.Name, &subcat.categoryId)
				subcategories = append(subcategories, subcat)
			}
			rows.Close()

			rows, err := db.Query(
				"SELECT t.id, t.name, t.subcategory_id " +
					"FROM part_type t")


			if err != nil {
				log.Println(err)
			} else {
				types := []*Type{}
				for rows.Next() {
					for rows.Next() {
						partType := &Type{}
						rows.Scan(&partType.Id, &partType.Name, &partType.subcategoryId)
						types = append(types, partType)
					}
				}
				rows.Close()

				for typeIdx := range types {
					for SubcategoryIdx := range subcategories {
						if types[typeIdx].subcategoryId == subcategories[SubcategoryIdx].Id {
							subcategories[SubcategoryIdx].Types = append(subcategories[SubcategoryIdx].Types, types[typeIdx])
						}
					}
				}
				for subcategoryIdx := range subcategories {
					for catIdx := range result {
						if subcategories[subcategoryIdx].categoryId == result[catIdx].Id {
							result[catIdx].Subcategories = append(result[catIdx].Subcategories, subcategories[subcategoryIdx])
						}
					}
				}
			}
		}

	}
	 return result, err
}


func SearchForParts(modelId, yearId, engineId, typeId int) ([]*Part, error) {
	result := []*Part{}

	rows, err := db.Query(
		"SELECT p.id, p.part_number, p.price, p.supplier, "+
		"  p.quality, p.image_name, t.name, coalesce(si.location, ''), coalesce(si.quantity, 0) "+
		"FROM part p "+
		"LEFT JOIN store_inventory si "+
		"  ON si.part_id = p.id "+
		"JOIN model_year_engine_part myep "+
		"  ON myep.part_id = p.id "+
		"JOIN part_type t "+
		"  ON t.id = p.type_id "+
		"WHERE myep.model_id = $1 "+
		"  AND myep.year_id = $2 "+
		"  AND myep.engine_id = $3 "+
		"  AND p.type_id = $4", modelId, yearId, engineId, typeId)
	if err != nil {
		log.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			p := &Part{StoreInventory: &StoreInventory{}}
			rows.Scan(&p.Id, &p.PartNumber, &p.Price, &p.Supplier,
			&p.Quality, &p.ImageName, &p.TypeName, &p.StoreInventory.Location, &p.StoreInventory.Quantity)
			result = append(result, p)
		}
	}

	return result, err
}



func GetPart(partId int) (*Part, error) {
	result := Part{StoreInventory: &StoreInventory{}}

	row := db.QueryRow(
		"SELECT p.id, p.part_number, p.price, p.supplier, "+
		"  p.quality, p.image_name, t.name, coalesce(si.location, ''), coalesce(si.quantity, 0) "+
		"FROM part p "+
		"LEFT JOIN store_inventory si "+
		"  ON si.part_id = p.id "+
		"JOIN part_type t "+
		"  ON t.id = p.type_id "+
		"WHERE p.id = $1 ", partId)

	err := row.Scan(&result.Id, &result.PartNumber, &result.Price, &result.Supplier,
		&result.Quality, &result.ImageName, &result.TypeName, &result.StoreInventory.Location, &result.StoreInventory.Quantity)

	return &result, err 
}
