package main

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"strings"
)

func main() {
	if (len(os.Args) != 2) {
		fmt.Fprintf(os.Stderr, "Usage: generate_ast <output directory>\n")
		os.Exit(64)
	}
	var output_dir string = os.Args[1]
	defineAst(output_dir, "expr", []string{
		"Binary : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal : any value",
		"Unary : Token operator, Expr right",
	})
}

func defineAst(output_dir string, base_name string, types []string) {
	path := output_dir + "/" + base_name + ".go"
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// TODO: handle errors?
	_, err = writer.WriteString("package main\n")
	_, err = writer.WriteString("\n")
	_, err = writer.WriteString("type " + base_name + " struct {\n")
	_, err = writer.WriteString("}\n")
	for _, t := range types {
		class_name := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(writer, base_name, class_name, fields)
	}
	
}

func defineType(writer *bufio.Writer, base_name string, class_name string, field_list string) {
	// TODO: handle errors?
	writer.WriteString("type " + class_name + " struct {\n")
	writer.WriteString("}\n")
}
