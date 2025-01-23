package model

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
	"time"
)



type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}



type Employee struct {
	Id             int
	GivenName      string
	Surname        string
	Address        string
	City           string
	State          string
	PostalCode     string
	roleId         int
	Role           *Role
	HireDate       time.Time
	PayRate        float32
	EmployeeNumber int
}

func GetEmployee(employeeNumber int) (*Employee, error) {
	result := Employee{}
	row := db.QueryRow(
		"SELECT id, given_name, surname, address, city, state, postal_code, role_id, hire_date, pay_rate, employee_number "+
			"FROM employee "+
			"WHERE employee_number = $1", employeeNumber)

	err := row.Scan(&result.Id, &result.GivenName, &result.Surname, &result.Address, &result.City, &result.State, &result.PostalCode,
		&result.roleId, &result.HireDate, &result.PayRate, &result.EmployeeNumber)

	result.Role, _ = GetRole(result.roleId)

	return &result, err
}

func GetEmployeeWithPassword(employeeNumber int, password string) (*Employee, error) {
	result := Employee{}

	hashedPassword := sha512.Sum512([]byte(password))
	b64Pass := base64.StdEncoding.EncodeToString(hashedPassword[:])

	row := db.QueryRow(
		"SELECT id, given_name, surname, address, city, state, postal_code, role_id, hire_date, pay_rate, employee_number "+
			"FROM employee "+
			"WHERE employee_number = $1 "+
			"  AND password = $2", employeeNumber, b64Pass)

	err := row.Scan(&result.Id, &result.GivenName, &result.Surname, &result.Address, &result.City, &result.State, &result.PostalCode,
		&result.roleId, &result.HireDate, &result.PayRate, &result.EmployeeNumber)

	return &result, err
}

func CreateEmployee(emp *Employee) (*Employee, error) {
	row := db.QueryRow(
		"INSERT INTO employee "+
			"(id, given_name, surname, address, city, state, postal_code, "+
			"role_id, hire_date, pay_rate, employee_number) "+
			"VALUES "+
			"(nextval('employee_id_seq'), $1, $2, $3, $4, $5, $6, "+
			"$7, $8, $9, nextval('employee_number_seq')) "+
			"RETURNING id",
		emp.GivenName, emp.Surname, emp.Address, emp.City, emp.State, emp.PostalCode,
		emp.Role.Id, emp.HireDate, emp.PayRate)

	err := row.Scan(&emp.Id)
	log.Print(emp.Id)
	row = db.QueryRow(
		"SELECT employee_number "+
			"FROM employee "+
			"WHERE id = $1", emp.Id)
	row.Scan(&emp.EmployeeNumber)

	return emp, err
}

func GetRole(roleId int) (*Role, error) {
	result := Role{}

	row := db.QueryRow(
		"SELECT id, name "+
			"FROM role "+
			"WHERE id = $1", roleId)

	err := row.Scan(&result.Id, &result.Name)

	return &result, err
}

func GetRoles() ([]*Role, error) {
	result := []*Role{}

	rows, err := db.Query("SELECT id, name " +
		"FROM role")

	if err != nil {
		log.Print(err)
	} else {
		for rows.Next() {
			role := Role{}
			rows.Scan(&role.Id, role.Name)
			result = append(result, &role)
		}
	}

	return result, err
}
