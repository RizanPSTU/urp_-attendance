package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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

func main() {
	empIDS := []int{101, 102, 103, 104, 105, 106, 107, 108, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 301, 302, 303, 304, 305, 306, 307, 401, 402, 403, 404, 405, 406}
	cleanerEmpID := []int{407, 408}
	fmt.Println(empIDS)
	fmt.Println(cleanerEmpID)

	records := readCsvFile("./data.csv")
	for i, r := range records {
		if r[5] != "" {
			// fmt.Println(i+1, r)
			empID, errID := strconv.Atoi(r[5])
			if errID != nil {
				log.Fatalln("ERROR convert emp ID :", errID)
			}
			isIDContains := contains(empIDS, empID)
			if isIDContains {
				test, errTime := time.Parse("2006-01-02", r[0])
				test2, err := time.Parse("15:04:05", r[1])
				if errTime != nil {
					log.Fatalln("Time parse :", errTime)
				}
				if err != nil {
					log.Fatalln("Time parse :", err)
				}
				fmt.Println("###################  START  ######################")
				fmt.Println("Parsed time :", test)
				fmt.Println("Parsed time 24 :", test2)
				fmt.Println("Name :", r[4])
				fmt.Println("EmpID :", r[5])
				fmt.Println("###################  END  ######################")
			} else {
				fmt.Println(i, " = Not in the empID :", empID)
			}
		}

	}
}
