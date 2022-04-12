package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/visualfc/atk/tk"
)

type Window struct {
	*tk.Window
}

type Lines struct {
	LineNum int64
	Value   int64
}

func NewWindow() *Window {
	mw := &Window{tk.RootWindow()}
	filetype := tk.FileType{
		Info: "Project's program",
		Ext:  ".lsf",
	}

	vbox := tk.NewVPackLayout(mw)
	lay1 := tk.NewHPackLayout(mw)
	lbl1 := tk.NewLabel(mw, `This program fix G2 and G3 problem.`)
	lbl2 := tk.NewLabel(mw, `WARNING: It tested only in my production,
	 don't blame us if it eats your kittens.`)
	chooseFile := tk.NewEntry(mw)
	chooseFile.SetWidth(30)
	errlog := tk.NewText(mw)
	btnCheck := tk.NewButton(mw, "Check")
	btnFix := tk.NewButton(mw, "Fix")
	var file string

	//Button Check event

	btnCheck.OnCommand(func() {
		var errors string = "Bad lines in file:\n"
		file, _ = tk.GetOpenFile(mw, "Choose", []tk.FileType{filetype}, "", "")

		if len(file) > 0 {
			chooseFile.SetText(file)
			if CheckMax(file) != nil {
				errors = errors + "Error: Max @VAR,@BOX,@TXT: ID > 64785\n"
			}

			if len(Check(file)) > 0 {

				for _, line := range Check(file) {
					errors = errors + "Error in line: " + strconv.FormatInt(line.LineNum, 10) + "\n"
					btnFix.SetState(tk.StateActive)
				}
			} else {
				errors = errors + "No bad lines.\n"
			}
		}
		errlog.SetText(errors)
	})
	btnFix.SetState(tk.StateDisable)

	//Button Fix event
	btnFix.OnCommand(func() {
		for _, fix := range Check(file) {
			err := InsertStringToFile(file, "#", int(fix.LineNum-1))
			if err != nil {
				fmt.Println(file)
				panic(err)
			}
			errlog.SetText("OK!")
		}
		btnFix.SetState(tk.StateDisable)
	})
	lay1.AddWidgets(chooseFile, btnCheck, btnFix)
	lay1.SetPaddingN(2, 2)
	vbox.AddWidgets(lbl1, lbl2, lay1, errlog)
	mw.ResizeN(400, 500)
	return mw
}

func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func contains(s []int64, str int64) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	tk.MainLoop(func() {
		mw := NewWindow()
		mw.SetTitle("IsaGRAF 3.x repair tools")
		mw.Center(nil)
		mw.ShowNormal()
	})
}

//Check bad communications and return it.
func Check(filepath string) []Lines {
	var ResultArr []Lines
	var LinesArr []Lines
	var LinesComm []int64

	rl, _ := regexp.Compile(`^@(ARC|NOT):*`)
	rc, _ := regexp.Compile(`([0-9]+)==`)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lineNum int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		if rl.MatchString(scanner.Text()) {
			lineSplitRg := RegSplit(scanner.Text(), ",")[0]
			lineSplit := strings.Split(lineSplitRg, ":")[1]
			line, err := strconv.ParseInt(lineSplit, 0, 64)
			if err != nil {
				panic(err)
			}
			LinesArr = append(LinesArr, Lines{lineNum, line})
		}
		if rc.MatchString(scanner.Text()) {

			lineSplitComm := RegSplit(scanner.Text(), "==")[0]
			comm, err := strconv.ParseInt(strings.Replace(lineSplitComm, " ", "", -1), 0, 64)
			if err != nil {
				panic(err)
			}
			LinesComm = append(LinesComm, comm)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, line := range LinesArr {

		if !contains(LinesComm, line.Value) {
			ResultArr = append(ResultArr, line)
		}
	}
	return ResultArr

}

//Check maximum  limit from VAR|BOX|TXT Ids
func CheckMax(filepath string) error {
	rl, _ := regexp.Compile(`^@(VAR|BOX|TXT):*`)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lineNum int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		if rl.MatchString(scanner.Text()) {
			lineSplitRg := RegSplit(scanner.Text(), ",")[0]
			lineSplit := strings.Split(lineSplitRg, ":")[1]
			line, err := strconv.ParseInt(lineSplit, 0, 64)
			if err != nil {
				panic(err)
			}
			if line > 64785 {
				return errors.New("MAX ID used")
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return nil

}
