package main

import (
	"os"
	"path/filepath"
	"fmt"
	//"github.com/AfLnk/thriftgoparser/proto"
	"runtime/debug"
	"github.com/AfLnk/thriftgoparser/proto"
	//"strings"
	//"bufio"
	//"time"
)

var(

	G_debug = 1
	G_warn = 2
	G_strict = 255
	G_verbose = 1

	/**
	 * Flags to control code generation
	 */
	gen_recurse = false;

	/**
	  * Whether or not negative field keys are accepted.
 	  */
	G_allow_neg_field_keys = true

	/**
 	  * Whether or not 64-bit constants will generate a warning.
 	  */
	G_allow_64bit_consts = 0

	/**
	 * Search path for inclusions
	 * 在哪些路径搜索 -I dir
	 */
	G_incl_searchpath []string

	G_curdir string
	G_curpath string

	//G_parse_mode int32
)

func CheckIsDir(path string) bool{
	return true
}

/**
 * Skips UTF-8 BOM if there is one
 */
func SkipUtf8Bom(f *os.File) bool{
	b := make([]byte, 1)
	f.Read(b)
	if b[0] == 0xEF {
		f.Read(b)
		if b[0] == 0xBB {
			f.Read(b)
			if b[0] == 0xBF {
				return true
			}
		}
	}

	f.Seek(0, 0)
	return false
}

/**
 * Reset program doctext information after processing a file
 */
func reset_program_doctext_info() {
	if (proto.G_program_doctext_candidate != "") {
		proto.G_program_doctext_candidate = ""
	}
	proto.G_program_doctext_lineno = 0
	proto.G_program_doctext_status = proto.INVALID
	fmt.Printf("%s\n", "program doctext set to INVALID");
}

func parseIt(program, parent *proto.TProgram){
	// Get scope file path
	path := program.GetPath()

	// Set current dir global, which is used in the include_file function
	G_curdir = proto.DirName(path)
	G_curdir = path

	// Open the file
	// skip UTF-8 BOM if there is one
	yyin, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open input file: \"%s\"\n", path)
	}

	if SkipUtf8Bom(yyin) {
		fmt.Printf("Skipped UTF-8 BOM at %s\n", path)
	}

	// Create new scope and scan for includes
	fmt.Printf("Scanning %s for includes\n", path);

	//G_parse_mode = INCLUDES;
	proto.G_program = program
	proto.G_scope = program.Scope()

	lexer := NewFileLexer(path)

	ret := yyParse(lexer)
	if 0 != ret{
		fmt.Printf("Parser error during include pass:%d", ret)
	}

	// Recursively parse all the include programs
	for _, v := range program.GetIncludes(){
		parseIt(v, program)
	}

	// reset program doctext status before parsing a new file
	reset_program_doctext_info();

	// Parse the program file
	proto.G_parse_mode = proto.PROGRAM;
	proto.G_program = program
	proto.G_scope = program.Scope()
	if parent != nil{
		proto.G_parent_scope = parent.Scope()
	}else{
		proto.G_parent_scope = nil
	}
	proto.G_parent_prefix = program.GetName() + "."
	G_curpath = path

	// Open the file
	// skip UTF-8 BOM if there is one
	if(SkipUtf8Bom(yyin)) {
		fmt.Printf("Skipped UTF-8 BOM at %s\n", path)
	}

	fmt.Printf("Parsing %s for types\n", path)

	ret = yyParse(lexer)
	if 0 != ret{
		fmt.Printf("Parser error during include pass:%d", ret)
	}
}

func generateIt(program *proto.TProgram) {
	// Oooohh, recursive code generation, hot!!
	if (gen_recurse) {
		for _, v := range program.GetIncludes(){
			// Propagate output path from parent to child programs
			generateIt(v)
		}
	}

	// Generate code!
	fmt.Printf("Program: %s\n", program.GetPath())

	gen := NewGenerator(program)
	gen.Generate_program()
}

func Parse(file_path, psm string) error{
	// Setup time string
	//timeNow := time.Now()
	//timeStr := timeNow.String()

	// 判断文件是否存在
	info, err := os.Stat(file_path)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("file:%s not exist", file_path)
	}

	// Instance of the global parse tree
	program := proto.NewProgram(file_path, info.Name(), psm)

	proto.InitGlobals()

	// Parse it!
	parseIt(program, nil)

	// The current path is not really relevant when we are doing generation.
	// Reset the variable to make warning messages clearer.
	G_curpath = "generation";
	// Reset yylineno for the heck of it.  Use 1 instead of 0 because
	// That is what shows up during argument parsing.

	// Generate it!
	generateIt(program)
	// Finished

	return nil
}

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	InitLoader("root","dev", "10.8.124.136", 3306)

	if err := filepath.Walk(".", func(pth string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ok, err := filepath.Match("*.idl", filepath.Base(pth))
			if err != nil {
				panic(err)
			}

			if ok {
				Parse(pth,"ee.lobster.thriftparser")
				return nil
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}
}