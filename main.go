package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsEmp(s []Employee, e Employee) (bool, Employee) {
	for _, a := range s {
		if a.EmpID == e.EmpID && a.Date == e.Date {
			// fmt.Print("Same ")
			// fmt.Println(a)
			return true, a
		}
	}
	return false, e
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

//Employee Model
type Employee struct {
	Name  string    `json:"Name"`
	EmpID int       `json:"EmpID"`
	Date  time.Time `json:"Date"`
	Time  time.Time `json:"Time"`
}

//User Model
type User struct {
	Name  string `json:"Name"`
	EmpID int    `json:"EmpID"`
}

func main() {
	//Missing 104 and 405 in currentUsersEmp
	officerIDS := []int{101, 102, 103, 104, 105, 106, 107, 108, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 301, 302, 303, 304, 305, 306, 307, 401, 402, 403, 404, 405, 406}
	cleanerEmpIDS := []int{407, 408}
	currentUsersEmp := []User{}

	//Two main
	// leaveEmpList := []Employee{}
	arriveEmpList := []Employee{}

	allOfficersRecors := []Employee{}
	records := readCsvFile("./data.csv")
	users := readCsvFile("./user.csv")

	csvFile, errCSV := os.Create("attendance.csv")

	if errCSV != nil {
		log.Fatal("Creating CSV error :", errCSV)
	}

	_, errWriteHead := fmt.Fprintln(csvFile, "Name,", "Date,", "Time")
	if errWriteHead != nil {
		log.Fatalln("Head write to attendance.csv Error:", errWriteHead)
	}
	defer csvFile.Close()

	for _, r := range records {
		if r[5] != "" {
			empID, errID := strconv.Atoi(r[5])
			if errID != nil {
				log.Fatalln("Error convert EmpID :", errID)
			}

			// Check if officer or not
			isIDContains := contains(officerIDS, empID)
			if isIDContains {
				//Officer
				date, errDate := time.Parse("2006-01-02", r[0])
				time, errTime := time.Parse("15:04:05", r[1])
				if errDate != nil {
					log.Fatalln("Time parse errDate:", errDate)
				}
				if errTime != nil {
					log.Fatalln("Time parse errTime:", errTime)
				}
				currentEmp := Employee{Name: r[4], EmpID: empID, Date: date, Time: time}
				allOfficersRecors = append(allOfficersRecors, currentEmp)

			}
		}
	}

	//Making a copy of original list
	originalCopyAllOfficersRecors := make([]Employee, len(allOfficersRecors))
	copy(originalCopyAllOfficersRecors, allOfficersRecors)

	//Sort
	sort.SliceStable(allOfficersRecors, func(i, j int) bool {
		return allOfficersRecors[i].EmpID < allOfficersRecors[j].EmpID
	})

	for _, ofr := range originalCopyAllOfficersRecors {
		//Find the arrival time
		isEmpContainsArrival, _ := containsEmp(arriveEmpList, ofr)

		if isEmpContainsArrival == false {
			arriveEmpList = append(arriveEmpList, ofr)
		}
	}

	// for i, e := range empList {
	// 	fmt.Println(i, " Name:", e.Name, " Date:", e.Date.String()[0:10], " Day:", e.Date.Weekday(), " Time:", e.Time.Format("2006-01-02 3:4:5 pm")[11:])
	// 	_, errWrite := fmt.Fprintln(csvFile, e.Name, e.Date.String()[0:10], e.Date.Weekday(), e.Time.Format("2006-01-02 3:4:5 pm")[11:])
	// 	if errWrite != nil {
	// 		log.Fatalln("Save to attendance.csv Error:", errWrite)
	// 	}
	// }

	//Current user get for finding the absent user
	for _, u := range users {
		if u[2] != "" {
			empID, errID := strconv.Atoi(u[2])
			if errID != nil {
				log.Fatalln("Error convert EmpID :", errID)
			}

			isIDContainsInClear := contains(cleanerEmpIDS, empID)

			if isIDContainsInClear == false {
				currentUsersEmp = append(currentUsersEmp, User{Name: u[1], EmpID: empID})
			}

		}
	}
}
