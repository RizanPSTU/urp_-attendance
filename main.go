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

func containsInAttendance(sal []Attendance, u User) bool {
	for _, sa := range sal {
		if u.EmpID == sa.EmpID {
			return true
		}
	}
	return false
}

func findTheUserWithA(ul []User, a Attendance) int {
	for i, u := range ul {
		if u.EmpID == a.EmpID {
			return i
		}
	}
	return -1
}

func findTheUserWithU(ul []User, u User) int {
	for i, u := range ul {
		if u.EmpID == u.EmpID {
			return i
		}
	}
	return -1
}

func containsEmpAndFixingArrival(s []Employee, e Employee) (bool, Employee) {
	for i, a := range s {
		if a.EmpID == e.EmpID && a.Date == e.Date {
			if a.Time.After(e.Time) {
				arriveEmpList[i] = e
			}
			return true, a
		}
	}
	return false, e
}
func containsEmpAndFixingLeave(s []Employee, e Employee) (bool, Employee) {
	for i, a := range s {
		if a.EmpID == e.EmpID && a.Date == e.Date {
			if a.Time.Before(e.Time) {
				leaveEmpList[i] = e
			}
			return true, a
		}
	}
	return false, e
}

func giveStatus(at time.Time, lt time.Time) string {

	if at == lt {
		return "Invalid Arrival or Leave Time"
	}
	timeA, _ := time.Parse("15:04:05", "09:00:00")
	timeL, _ := time.Parse("15:04:05", "17:00:00")

	if at.After(timeA) && lt.Before(timeL) {
		return "Late Arrived & Early leaved"
	}

	if at.After(timeA) == false && lt.Before(timeL) {
		return "Early leaved"
	}

	if at.After(timeA) && lt.Before(timeL) == false {
		return "Late Arrived"
	}

	return "Present Office Full Time"
}

func writeToAttendanceAbsent(f *os.File, e Attendance) {
	position := findTheUserWithA(currentUsersEmp, e)
	if position != -1 {
		e.Name = currentUsersEmp[position].Name
	}
	_, errWriteA := fmt.Fprintln(f, index, ",", e.Name, ",", e.EmpID, ",", dateFix(e.Date), ",", e.Day, ",", timeFix(e.ArrivalTime), ",", timeFix(e.LeaveTime), ",", giveStatus(e.ArrivalTime, e.LeaveTime))
	if errWriteA != nil {
		log.Fatalln("Save to attendance_with_absent.csv Error:", errWriteA)
	}
	index++
}

func writeAbsent(f *os.File, u User, e Attendance) {
	_, errWriteA := fmt.Fprintln(f, index, ",", u.Name, ",", u.EmpID, ",", dateFix(e.Date), ",", e.Day, ",", "Null", ",", "Null", ",", "Absent")
	if errWriteA != nil {
		log.Fatalln("Save to attendance_with_absent.csv Error:", errWriteA)
	}
	index++
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

func dateFix(d time.Time) string {
	return d.String()[0:10]
}

func timeFix(d time.Time) string {
	return d.Format("2006-01-02 3:4:5 pm")[11:]
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

//Attendance Model
type Attendance struct {
	Name        string    `json:"Name"`
	EmpID       int       `json:"EmpID"`
	Date        time.Time `json:"Date"`
	Day         string    `json:"Day"`
	ArrivalTime time.Time `json:"ArrivalTime"`
	LeaveTime   time.Time `json:"LeaveTime"`
	Status      string    `json:"Status"`
}

var (
	//main
	leaveEmpList                = []Employee{}
	arriveEmpList               = []Employee{}
	attendanceList              = []Attendance{}
	singleDateAttendantList     = []Attendance{}
	currentUsersEmp             = []User{}
	index                   int = 1
)

func main() {
	//Missing 104 and 405 in currentUsersEmp
	officerIDS := []int{101, 102, 103, 104, 105, 106, 107, 108, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 301, 302, 303, 304, 305, 306, 307, 401, 402, 403, 404, 405, 406}
	cleanerEmpIDS := []int{407, 408}

	allOfficersRecors := []Employee{}
	records := readCsvFile("./data.csv")
	users := readCsvFile("./user.csv")

	//Crating model array from raw data
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
	originalCopyAllOfficersRecords := make([]Employee, len(allOfficersRecors))
	copy(originalCopyAllOfficersRecords, allOfficersRecors)

	//Sort by date (Early is first)
	sort.Slice(originalCopyAllOfficersRecords, func(i, j int) bool {
		return originalCopyAllOfficersRecords[i].Date.After(originalCopyAllOfficersRecords[j].Date)
	})

	//Sort by date (Early is first) View
	// for _, t := range originalCopyAllOfficersRecords {
	// 	fmt.Println(t)
	// }

	//Sort by EmpID
	sort.SliceStable(allOfficersRecors, func(i, j int) bool {
		return allOfficersRecors[i].EmpID < allOfficersRecors[j].EmpID
	})

	//Finding arrival
	for _, ofr := range originalCopyAllOfficersRecords {
		//Find the arrival time
		isEmpContainsArrival, _ := containsEmpAndFixingArrival(arriveEmpList, ofr)

		if isEmpContainsArrival == false {
			arriveEmpList = append(arriveEmpList, ofr)
		}
	}

	//Finding leave
	for _, ofr := range originalCopyAllOfficersRecords {
		//Find the leave time
		isEmpContainsLeave, _ := containsEmpAndFixingLeave(leaveEmpList, ofr)

		if isEmpContainsLeave == false {
			leaveEmpList = append(leaveEmpList, ofr)
		}
	}

	//File for attendance
	csvFile, errCSV := os.Create("attendance.csv")
	if errCSV != nil {
		log.Fatal("Creating CSV error :", errCSV)
	}
	_, errWriteHead := fmt.Fprintln(csvFile, "Index,", "Name,", "Date,", "Day,", "Arrival Time,", "Leave Time,", "Status")
	if errWriteHead != nil {
		log.Fatalln("Head write to attendance.csv Error:", errWriteHead)
	}

	if len(arriveEmpList) == len(leaveEmpList) {
		fmt.Println("Same len :)")
		for i, e := range arriveEmpList {
			attendanceList = append(attendanceList, Attendance{Name: e.Name, EmpID: e.EmpID, Date: e.Date, Day: e.Date.Weekday().String(), ArrivalTime: e.Time, LeaveTime: leaveEmpList[i].Time, Status: giveStatus(e.Time, leaveEmpList[i].Time)})
			fmt.Println(i+1, " Name:", e.Name, " Date:", dateFix(e.Date), " Day:", e.Date.Weekday(), " Arrival Time:", timeFix(e.Time), " Leave Time:", timeFix(leaveEmpList[i].Time), " Status:", giveStatus(e.Time, leaveEmpList[i].Time))
			index := i + 1
			_, errWrite := fmt.Fprintln(csvFile, index, ",", e.Name, ",", dateFix(e.Date), ",", e.Date.Weekday(), ",", timeFix(e.Time), ",", timeFix(leaveEmpList[i].Time), ",", giveStatus(e.Time, leaveEmpList[i].Time))
			if errWrite != nil {
				log.Fatalln("Save to attendance.csv Error:", errWrite)
			}
		}
	}
	csvFile.Close()

	//File for attendance_with_absent
	csvFileA, errCSVA := os.Create("attendance_with_absent.csv")
	if errCSVA != nil {
		log.Fatal("Creating CSV error :", errCSVA)
	}
	_, errWriteHeadA := fmt.Fprintln(csvFile, "Index,", "Name,", "Date,", "Day,", "Arrival Time,", "Leave Time,", "Status")
	if errWriteHead != nil {
		log.Fatalln("Head write to attendance_with_absent.csv Error:", errWriteHeadA)
	}

	defer csvFileA.Close()

	// var userInput int
	// fmt.Println("Generate absant: Press 1")
	// fmt.Scanln(&userInput)
	// fmt.Println("You entered :", userInput)

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

	take := attendanceList[0].Date
	atlIndex := 0
	atEnd := false
	for !atEnd {
		ate := attendanceList[atlIndex]
		if take == ate.Date {
			singleDateAttendantList = append(singleDateAttendantList, ate)
			atlIndex++
		} else {
			take = ate.Date
			for _, dat := range singleDateAttendantList {
				writeToAttendanceAbsent(csvFileA, dat)
			}

			for _, u := range currentUsersEmp {
				isPresent := containsInAttendance(singleDateAttendantList, u)
				if isPresent == false {
					writeAbsent(csvFileA, u, singleDateAttendantList[0])
				}
			}
			singleDateAttendantList = []Attendance{}
		}

		if atlIndex == len(arriveEmpList) {
			atEnd = true
			for _, dat := range singleDateAttendantList {
				writeToAttendanceAbsent(csvFileA, dat)
			}

			for _, u := range currentUsersEmp {
				isPresent := containsInAttendance(singleDateAttendantList, u)
				if isPresent == false {
					writeAbsent(csvFileA, u, singleDateAttendantList[0])
				}
			}
			singleDateAttendantList = []Attendance{}
		}
	}

}
