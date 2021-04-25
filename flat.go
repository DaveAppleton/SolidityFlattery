package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var pragma string

type importRec struct {
	FullPath  string
	Code      []string
	Uses      map[string]bool
	Created   bool // has been depended on or processed
	Processed bool // has been processed for includes
	Resolved  bool // all includes have been resolved
	Written   bool // written out
}

var imports map[string]importRec

func loadAndSplitFile(fileName string) (newFiles bool, err error) {
	fileName = strings.Replace(fileName, "\"", "", 2)      // remove quotes
	fileName = strings.Replace(fileName, ";", "", 2)       // and final semicolon
	fileName = strings.TrimSpace(filepath.Clean(fileName)) // resolve ../ etc
	thisPath := filepath.Dir(fileName)
	fmt.Println("++++", thisPath)

	_, err = os.Stat(thisPath)
	if err != nil {
		fmt.Println("err ", thisPath)
		node, err := doWeHaveNodeModule(thisPath)
		if err == nil {
			fmt.Println("using ", node+"/node_modules/"+thisPath)
			thisPath = node + "/node_modules/" + thisPath
			fileName = node + "/node_modules/" + fileName
		}
	}

	shortName := filepath.Base(fileName) // just the file name
	if imports[shortName].Processed {
		fmt.Println(shortName, " already done")
		return
	}
	thisRec := importRec{FullPath: fileName, Created: true, Uses: make(map[string]bool)}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {

		if fileName[:1] == "." {
			fmt.Println(fileName, err)
			return
		}
		fileName = filepath.Clean("./" + fileName)
		data, err = ioutil.ReadFile(fileName)
		if err != nil {
			dir, newerr := os.Getwd()
			if newerr != nil {
				log.Fatal(err)
			}
			fmt.Println(dir)
			fmt.Println("*", "["+fileName+"]", err)
			return
		}

	}
	fmt.Println("Processing ", shortName)
	contents := string(data)
	lines := strings.Split(contents, "\n")
	if len(pragma) == 0 {
		pragma = lines[0]
	}
	noImports := true
	for li, line := range lines {
		if li == 0 {
			continue // skip "pragma solidity ^0.4.x;"
		}
		if starts("import", line) {
			noImports = false
			afterImport := after("import ", line)

			if afterImport[0] == '@' {
				if !imports[afterImport].Created {
					newFiles = true
					imports[afterImport] = importRec{FullPath: afterImport,
						Created: true,
						Uses:    make(map[string]bool)}
					fmt.Println("--> uses new file ", afterImport)
				}
				thisRec.Uses[afterImport] = true
			} else {
				afta := thisPath + "/" + afterImport
				afta = filepath.Clean(afta)
				bafta := filepath.Base(afta)
				if !imports[bafta].Created {
					newFiles = true
					imports[bafta] = importRec{FullPath: afta,
						Created: true,
						Uses:    make(map[string]bool)}
					fmt.Println("--> uses new file ", afta)
				}
				thisRec.Uses[bafta] = true
			}

		}
		if starts("contract", line) || starts("library", line) || starts("interface", line) || starts("abstract contract", line) {
			thisRec.Code = lines[li:]
			fmt.Println("has ", len(lines[li:]), " lines")
			break
		}
	}
	thisRec.Processed = true
	imports[shortName] = thisRec
	thisRec.Resolved = noImports
	if noImports {
		fmt.Println(shortName, " has no dependencies")
	}
	return
}

func findCurrentParent() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Dir(cwd), nil
}

func findDirectoryParent(cwd string) string {
	return filepath.Dir(cwd)
}

func doWeHaveNodeModule(path string) (dir string, err error) {
	start, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		nodeDir := start + "/node_modules/" + path
		f, err := os.Stat(nodeDir)
		if err != nil {
			start = findDirectoryParent(start)
		} else if f.IsDir() {

			return start, nil

		}
		if len(start) < 3 {
			return "", errors.New("cannot find node")
		}
	}
}

var fName string
var oName string

func main() {
	fmt.Println("Solididy File Flattener 1.20 (c) David Appleton 2018 to 2020 and beyond")
	fmt.Println("contact : calistralabs@gmail.com")
	fmt.Println("https://github.com/DaveAppleton/SolidityFlattery")
	fmt.Println("released under Apache 2.0 licence")
	flag.StringVar(&fName, "input", "", "base file to flatten")
	flag.StringVar(&oName, "output", "", "output file")
	flag.Parse()
	if len(fName) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	if len(oName) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	if _, err := os.Stat(oName + ".sol"); err == nil {
		fmt.Println("error : ", oName+".sol already exists")
		fmt.Println("we can't have you accidentally deleting files!!!")
		os.Exit(1)
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   "./" + oName + ".log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	oName += ".sol"
	fmt.Println("ANALYSIS")
	imports = make(map[string]importRec)

	f, _ := os.Create(oName)
	defer f.Close()
	w := bufio.NewWriter(f)
	log.Printf("pre processing %s\n", fName)
	newFiles, err := loadAndSplitFile(fName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// keep processing dependent files
	for {
		repeat := false
		for _, iRec := range imports {
			if iRec.Processed {
				continue
			}
			log.Printf("pre processing %s\n", iRec.FullPath)
			newFiles, err = loadAndSplitFile(iRec.FullPath)
			repeat = repeat || newFiles
		}
		if !repeat {
			break
		}
	}

	absoluteFile, _ := filepath.Abs(fName)

	fmt.Fprintln(w, pragma)

	fmt.Fprintln(w, "// produced by the Solididy File Flattener (c) David Appleton 2018 - 2020 and beyond")
	fmt.Fprintln(w, "// contact : calistralabs@gmail.com")
	fmt.Fprintln(w, "// source  : https://github.com/DaveAppleton/SolidityFlattery")
	fmt.Fprintln(w, "// released under Apache 2.0 licence")
	fmt.Fprintln(w, "// input ", absoluteFile)
	fmt.Fprintln(w, "// flattened : ", time.Now().UTC().Format(time.RFC850))

	fmt.Println("Writing output.")
	for {
		completed := true
		for key, mp := range imports {
			if mp.Written {
				continue
			}
			completed = false
			if mp.Resolved {
				count := 0
				for _, line := range mp.Code {
					fmt.Fprintln(w, line)
					count++
				}
				mp.Written = true
				imports[key] = mp
				log.Printf("Written %s (%d) lines\n", key, count)
				continue
			}
			amResolved := true
			for k2, _ := range mp.Uses {
				if !imports[filepath.Base(k2)].Written {
					amResolved = false
				}
			}
			if amResolved {
				mp.Resolved = true
				imports[key] = mp
				log.Println("Resolved ", key)
				continue
			}
			log.Println(key, "remains unresolved")
		}
		if completed {
			break
		}
	}
	err = w.Flush()

}
