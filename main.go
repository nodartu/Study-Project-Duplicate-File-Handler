package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	UsersInput      []string
	FileFormat      string
	SortingOption   string
	DupOrNot        string
	SameFormatPaths = make(map[float64][]string)
	AllKeys         []float64
	counter         = 1
	deletedSum      float64
)

func main() {

	input()
	if input() {
		fmt.Println("Enter file format:")
		fmt.Scanln(&FileFormat)
		WalkFiles()
		DeletSingleFiles(SameFormatPaths)
		KeysToSlice(SameFormatPaths)
		SortAllPaths(SameFormatPaths)
		FormatOption()
		CheckOrNot()
		DeleteOrNot()
	}

}

func input() bool {
	UsersInput = os.Args
	if len(UsersInput) == 1 {
		fmt.Println("Directory is not specified")
		return false
	}
	return true
}

func FormatOption() {
	fmt.Printf("\nSize sorting options:\n" +
		"1. Descending\n" +
		"2. Ascending\n\n" +
		"Enter a sorting option:\n")

	for {
		fmt.Scan(&SortingOption)
		switch SortingOption {
		case "1":
			DescendingSort()
			break
		case "2":
			AscendingSort()
			break
		case "/exit":
			os.Exit(0)
		default:
			fmt.Printf("\nWrong option\n")
			continue
		}
		break
	}
}

func KeysToSlice(m map[float64][]string) {
	for key, _ := range m {
		AllKeys = append(AllKeys, key)
	}

}

func SortAllPaths(m map[float64][]string) {
	for _, val := range m {
		sort.Slice(val, func(i, j int) bool {
			return val[i] < val[j]
		})
	}
}
func DescendingSort() {
	sort.Sort(sort.Reverse(sort.Float64Slice(AllKeys)))
	PrintResult(SameFormatPaths)
}

func AscendingSort() {
	sort.Float64s(AllKeys)
	PrintResult(SameFormatPaths)
}

func PrintResult(m map[float64][]string) {
	for i := 0; i < len(AllKeys); i++ {
		key := AllKeys[i]
		val := m[key]
		fmt.Println()
		fmt.Println(key, "bytes")
		for _, str := range val {
			fmt.Println(str)
		}
	}
}

func WalkFiles() {
	err := filepath.Walk(UsersInput[1], func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() && strings.HasSuffix(path, FileFormat) {
			key := float64(info.Size())
			SameFormatPaths[key] = append(SameFormatPaths[key], path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func DeletSingleFiles(m map[float64][]string) {
	for key, val := range m {
		if len(val) < 2 {
			delete(m, key)
		}
	}
}

func CheckOrNot() {
	for {
		fmt.Printf("\nCheck for duplicates?\n")
		fmt.Scan(&DupOrNot)
		switch DupOrNot {
		case "yes":
			CheckDuplicates(SameFormatPaths)
			PrintFinalMap(SameFormatPaths)
			break
		case "no":
			break
		default:
			fmt.Printf("\nWrong option\n")
		}
		break
	}
}

func CheckDuplicates(m map[float64][]string) {
	for key, val := range m {
		NewStrings := AreDuplicates(val)
		val = nil
		m[key] = NewStrings
	}
}

func AreDuplicates(str []string) []string {
	NewMap := make(map[string][]string, len(str))
	var NewString []string
	for _, val := range str {
		Hash := HashTheFileByPath(val)
		NewMap[Hash] = append(NewMap[Hash], val)
	}
	for i, solo := range NewMap {
		if len(solo) < 2 {
			delete(NewMap, i)
		}
	}
	for key, value := range NewMap {
		NewString = append(NewString, key)
		for _, path := range value {
			NewString = append(NewString, path)
		}
	}
	return NewString
}

func PrintFinalMap(m map[float64][]string) {
	for i := 0; i < len(AllKeys); i++ {
		key := AllKeys[i]
		val := m[key]
		fmt.Println()
		fmt.Println(key, "bytes")
		for j, str := range val {
			if strings.HasPrefix(str, "Hash") {
				fmt.Println(str)
			} else {
				val[j] = strconv.Itoa(counter) + ". " + str
				fmt.Println(val[j])
				counter++
			}
		}
	}
}

func HashTheFileByPath(s string) string {
	file, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		log.Fatal(err)
	}
	return "Hash: " + hex.EncodeToString(md5Hash.Sum(nil))
}

func DeleteOrNot() {
	for {
		fmt.Printf("\nDelete files?\n")
		fmt.Scan(&DupOrNot)
		switch DupOrNot {
		case "yes":
			SelectToDelete()
			break
		case "no":
			os.Exit(0)
		default:
			fmt.Printf("\nWrong option\n")
			continue
		}
		break
	}
}

func SelectToDelete() {
	var PathNumbers []string
	for {
		correct := true
		fmt.Printf("\nEnter file numbers to delete:\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		PathNumbers = strings.Fields(scanner.Text())
		if len(PathNumbers) == 0 {
			correct = false
		}

		for _, val := range PathNumbers {
			NumVal, _ := strconv.Atoi(val)
			switch {
			case isNumber(val) && NumVal <= counter && correct == true:
				continue
			default:
				correct = false
			}
		}

		if correct == false {
			fmt.Printf("\nWrong format\n")
		} else {
			DeletingFiles(PathNumbers, SameFormatPaths)
			break
		}
	}
}

func DeletingFiles(s []string, m map[float64][]string) {
	deletedSum = 0
	for _, NumberOfPathToDelete := range s {
		fmt.Println(NumberOfPathToDelete)
		for i, _ := range AllKeys {
			SizeOfFiles := AllKeys[i]
			KeyOfMap := m[SizeOfFiles]
			for _, LineInMapValue := range KeyOfMap {
				if strings.HasPrefix(LineInMapValue, NumberOfPathToDelete) {
					LineInMapValue = LineInMapValue[3:]
					DeleteTheFileByPath(LineInMapValue)
					deletedSum += SizeOfFiles
				}
			}
		}
	}
	fmt.Println("Total freed up space:", deletedSum, "bytes")
}

func DeleteTheFileByPath(s string) {
	err := os.RemoveAll(s)
	if err != nil {
		log.Fatal(err)
	}
}

func isNumber(n string) bool {
	for _, num := range n {
		if !unicode.IsNumber(num) {
			return false
		}
	}
	return true
}
