package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type TOKEN int

const (
	ABSTRACT TOKEN = iota + 1
	ASSERT
	BOOLEAN
	BREAK
	BYTE
	COLON
	CASE
	CATCH
	CHAR
	CLASS
	CONST
	CONTINUE
	DEFAULT
	DO
	DOT
	DOUBLE
	ELSE
	ENUM
	EXPORTS
	EXTENDS
	FINAL
	FINALLY
	FLOAT
	FOR
	GOTO
	IF
	IMPLEMENTS
	IMPORT
	INSTANCEOF
	INT
	INTERFACE
	LONG
	MODULE
	NATIVE
	NEW
	OPEN
	OPENS
	PACKAGE
	PRIVATE
	PROTECTED
	PROVIDES
	PUBLIC
	REQUIRES
	RETURN
	SHORT
	STATIC
	STRICTFP
	SUPER
	SWITCH
	SYNCHRONIZED
	THIS
	THROW
	THROWS
	TO
	TRANSIENT
	TRANSITIVE
	TRY
	USES
	VOID
	VOLATILE
	WHILE
	WITH
	RIGHT_BRACKET
	LEFT_BRACKET
	RIGHT_PARENTHESES
	LEFT_PARENTHESES
	RIGHT_BRACES
	LEFT_BRACES
	COMMA
	SEMICOLON
	ASTERICK
	ASSIGN
	ASSIGN_ADD
	ASSIGN_SUBTRACT
	ASSIGN_MULTIPLY
	ASSIGN_DIVIDE
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE
	MODULO
	INCREMENT
	DECREMENT
	EQUAL
	NOT_EQUAL
	GREATER_THAN
	LESS_THAN
	GREATER_THAN_EQUAL
	LESS_THAN_EQUAL
	EOF
	LOGICAL_AND
	LOGICAL_OR
	NOT
	TRUE
	FALSE
	TERNARY
	BITWISE_OR
	BITWISE_AND
	BITWISE_XOR
	BITWISE_COMPLEMENT
	RIGHT_SHIFT
	UNSIGNED_RIGHT_SHIFT
	LEFT_SHIFT
	UNSIGNED_LEFT_SHIFT
	IDENTIFIER
	NUMBER
	STRING
	SPACE
	DECORATOR
	ARRAY
)

var KEYWORDS = map[string]TOKEN{
	"abstract":     ABSTRACT,
	"assert":       ASSERT,
	"boolean":      BOOLEAN,
	"break":        BREAK,
	"byte":         BYTE,
	"case":         CASE,
	"catch":        CATCH,
	"char":         CHAR,
	"class":        CLASS,
	"const":        CONST,
	"continue":     CONTINUE,
	"default":      DEFAULT,
	"do":           DO,
	"double":       DOUBLE,
	"else":         ELSE,
	"enum":         ENUM,
	"exports":      EXPORTS,
	"extends":      EXTENDS,
	"final":        FINAL,
	"finally":      FINALLY,
	"float":        FLOAT,
	"for":          FOR,
	"goto":         GOTO,
	"if":           IF,
	"implements":   IMPLEMENTS,
	"import":       IMPORT,
	"instanceof":   INSTANCEOF,
	"int":          INT,
	"interface":    INTERFACE,
	"long":         LONG,
	"module":       MODULE,
	"native":       NATIVE,
	"new":          NEW,
	"open":         OPEN,
	"opens":        OPENS,
	"package":      PACKAGE,
	"private":      PRIVATE,
	"protected":    PROTECTED,
	"provides":     PROVIDES,
	"public":       PUBLIC,
	"requires":     REQUIRES,
	"return":       RETURN,
	"short":        SHORT,
	"static":       STATIC,
	"strictfp":     STRICTFP,
	"super":        SUPER,
	"switch":       SWITCH,
	"synchronized": SYNCHRONIZED,
	"this":         THIS,
	"throw":        THROW,
	"throws":       THROWS,
	"to":           TO,
	"transient":    TRANSIENT,
	"transitive":   TRANSITIVE,
	"try":          TRY,
	"uses":         USES,
	"void":         VOID,
	"volatile":     VOLATILE,
	"while":        WHILE,
	"with":         WITH,
	"EOF":          EOF,
}

var OPERATORS = map[string]TOKEN{
	"]":     RIGHT_BRACKET,
	"[":     LEFT_BRACKET,
	")":     RIGHT_PARENTHESES,
	"(":     LEFT_PARENTHESES,
	"}":     RIGHT_BRACES,
	"{":     LEFT_BRACES,
	",":     COMMA,
	":":     COLON,
	".":     DOT,
	";":     SEMICOLON,
	"*":     ASTERICK,
	"=":     ASSIGN,
	"+=":    ASSIGN_ADD,
	"-=":    ASSIGN_SUBTRACT,
	"*=":    ASSIGN_MULTIPLY,
	"/=":    ASSIGN_DIVIDE,
	"+":     ADD,
	"-":     SUBTRACT,
	"/":     DIVIDE,
	"%":     MODULO,
	"++":    INCREMENT,
	"--":    DECREMENT,
	"==":    EQUAL,
	"!=":    NOT_EQUAL,
	">":     GREATER_THAN,
	"<":     LESS_THAN,
	">=":    GREATER_THAN_EQUAL,
	"<=":    LESS_THAN_EQUAL,
	"&&":    LOGICAL_AND,
	"||":    LOGICAL_OR,
	"!":     NOT,
	"true":  TRUE,
	"false": FALSE,
	"?":     TERNARY,
	"|":     BITWISE_OR,
	"&":     BITWISE_AND,
	"^":     BITWISE_XOR,
	"~":     BITWISE_COMPLEMENT,
	">>":    RIGHT_SHIFT,
	">>>":   UNSIGNED_RIGHT_SHIFT,
	"<<":    LEFT_SHIFT,
	"<<<":   UNSIGNED_LEFT_SHIFT,
}

var NON_JS_KEYWORD = []TOKEN{
	PACKAGE,
	PUBLIC,
	ABSTRACT,
	PRIVATE,
	PROTECTED,
	STATIC,
	FINAL,
}

var NATIVE_JAVA_TYPES = []TOKEN{
	VOID,
	DOUBLE,
	FLOAT,
	STRING,
	CHAR,
	SHORT,
	BOOLEAN,
	BYTE,
	LONG,
	INT,
}

type TokenGroup struct {
	token TOKEN
	value string
}

func main() {
	var path = os.Args[1]
	var path_ptr = flag.String("path", "", "path to a Java file")
	var is_watching = flag.Bool("watch", false, "watch the file")
	flag.Parse()

	if *path_ptr != "" {
		path = *path_ptr
	}

	_ = is_watching

	if filepath.Ext(path) != ".java" {
		return
	}

	file, err := readFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	tokens, err := tokenize(file)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	transformed, err := transform(tokens, "", 0, Class{"", 0})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	saveFile(transformed, path[0:len(path)-5]+".js")

}

func readFile(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(f), err
}

func scan(source []rune) (TokenGroup, error) {

	if source[0] == '/' && (source[1] == '/' || source[1] == '*') {

		var comment string = ""

		if source[1] == '/' {
			for len(source) > 0 && source[0] != '\n' {
				comment += string(source[0])
				source = source[1:]
			}
		}
		if source[1] == '*' {
			for len(source) > 0 && (source[0] != '*' && source[1] != '/') {
				comment += string(source[0])
				source = source[1:]
			}
			if comment[len(comment)-1] != '/' && comment[len(comment)-2] != '*' {
				return TokenGroup{}, errors.New("Unrecognized character found in source " + string(source[0]))
			}
		}
		return TokenGroup{SPACE, comment}, nil
	}

	if source[0] == '"' {
		var string_literal = []rune{'"'}

		source = source[1:]
		for len(source) > 0 && source[0] != '"' {
			string_literal = append(string_literal, source[0])
			source = source[1:]
		}
		string_literal = append(string_literal, '"')

		if string_literal[len(string_literal)-1] != '"' {
			return TokenGroup{}, errors.New("Unrecognized string_literal found in source " + string(source[0]))
		}

		return TokenGroup{STRING, string(string_literal)}, nil
	}

	if source[0] == '\'' {
		var char_literal = []rune{}

		for i := 0; i < 4; i++ {
			if i == 3 && char_literal[1] != '\\' {
				continue
			}
			if len(source) != 0 {
				char_literal = append(char_literal, source[0])
				source = source[1:]
			}
		}

		if char_literal[1] == '\\' && char_literal[3] == '\'' && len(char_literal) == 4 {
			return TokenGroup{STRING, string(char_literal)}, nil
		}

		if char_literal[2] != '\'' || len(char_literal) != 3 {
			return TokenGroup{}, errors.New("Unrecognized char_literal found in source " + string(char_literal[2]))
		}

		return TokenGroup{STRING, string(char_literal)}, nil
	}

	if source[0] == '<' {
		if source[1] == '=' {
			return TokenGroup{LESS_THAN_EQUAL, string(source[0])}, nil
		}
		if source[1] == '<' {
			if source[2] == '<' {
				return TokenGroup{UNSIGNED_LEFT_SHIFT, string(source[0])}, nil
			}
			return TokenGroup{LEFT_SHIFT, string(source[0])}, nil
		}
		return TokenGroup{LESS_THAN, string(source[0])}, nil
	}

	if source[0] == '>' {
		if source[1] == '=' {
			return TokenGroup{GREATER_THAN_EQUAL, string(source[0])}, nil
		}
		if source[1] == '>' {
			if source[2] == '>' {
				return TokenGroup{UNSIGNED_RIGHT_SHIFT, string(source[0])}, nil
			}
			return TokenGroup{RIGHT_SHIFT, string(source[0])}, nil
		}
		return TokenGroup{GREATER_THAN, string(source[0])}, nil
	}

	if source[0] == '=' {
		if source[1] == '=' {
			return TokenGroup{EQUAL, string(source[0])}, nil
		}
		// TODO: ADD ARRAY INIT

		padding_length := 0
		for unicode.IsSpace(source[padding_length+1]) {
			padding_length += 1
		}
		if source[padding_length+1] == '{' {
			arr := "=" + strings.Repeat(" ", padding_length) + "["

			source = source[padding_length+2:]
			for source[0] != '}' {
				arr += string(source[0])
				source = source[1:]
			}
			return TokenGroup{ARRAY, arr + "]"}, nil

		}
		return TokenGroup{ASSIGN, string(source[0])}, nil
	}

	if source[0] == '!' {
		if source[1] == '=' {
			return TokenGroup{NOT_EQUAL, string(source[0])}, nil
		}
		return TokenGroup{NOT, string(source[0])}, nil
	}

	if source[0] == '&' && source[1] == '&' {
		return TokenGroup{LOGICAL_AND, string(source[0])}, nil
	}

	if source[0] == '|' && source[1] == '|' {
		return TokenGroup{LOGICAL_OR, string(source[0])}, nil
	}

	if source[0] == '+' {
		if source[1] == '=' {
			return TokenGroup{ASSIGN_ADD, string(source[0])}, nil
		}
		if source[1] == '+' {
			return TokenGroup{INCREMENT, string(source[0])}, nil
		}
	}

	if source[0] == '-' {
		if source[1] == '=' {
			return TokenGroup{ASSIGN_SUBTRACT, string(source[0])}, nil
		}
		if source[1] == '-' {
			return TokenGroup{DECREMENT, string(source[0])}, nil
		}
	}

	if source[0] == '*' && source[1] == '=' {
		return TokenGroup{ASSIGN_MULTIPLY, string(source[0])}, nil
	}

	if source[0] == '/' && source[1] == '=' {
		return TokenGroup{ASSIGN_DIVIDE, string(source[0])}, nil
	}

	identifier, ok := OPERATORS[string(source[0])]

	if ok {
		return TokenGroup{identifier, string(source[0])}, nil
	}

	if unicode.IsDigit(source[0]) {
		var number string = ""
		for len(source) > 0 && unicode.IsDigit(source[0]) {
			number += string(source[0])
			source = source[1:]
		}
		return TokenGroup{NUMBER, number}, nil
	}
	if unicode.IsLetter(source[0]) || source[0] == '@' {
		var word = []rune{}
		if source[0] == '@' {
			word = append(word, '@')
			source = source[1:]
		}

		for len(source) > 0 && (unicode.IsLetter(source[0]) || source[0] == '<' || source[0] == '>' ||
			source[0] == '-' || source[0] == '_' || source[0] == '[' || source[0] == ']') {
			word = append(word, source[0])
			source = source[1:]

		}
		identifier, ok := KEYWORDS[string(word)]
		if ok {
			if identifier == IMPORT {
				for len(source) > 0 && source[0] != ';' {
					word = append(word, source[0])
					source = source[1:]
				}
			}

			if identifier == FOR || identifier == WHILE || identifier == IF {

				for unicode.IsSpace(source[0]) {
					source = source[1:]
				}

				skipped_first_type := false // HACK

				for len(source) > 0 && source[0] != '{' {
					if identifier == FOR && !skipped_first_type && source[0] != '(' {
						if unicode.IsSpace(source[0]) {
							skipped_first_type = true
							word = append(word, []rune("let ")...)
						}
						source = source[1:]
						continue
					}
					word = append(word, source[0])
					source = source[1:]
				}

				padding_length := len(word) - 1
				for unicode.IsSpace(word[padding_length]) {
					padding_length -= 1
				}

				if word[padding_length] != ')' {
					return TokenGroup{identifier, string(word)}, errors.New("missing \")\"")
				}
			}

			return TokenGroup{identifier, string(word)}, nil
		}
		if word[0] == '@' {
			return TokenGroup{DECORATOR, string(word)}, nil
		}
		return TokenGroup{IDENTIFIER, string(word)}, nil

	}
	if unicode.IsSpace(source[0]) {
		return TokenGroup{SPACE, string(source[0])}, nil
	}

	return TokenGroup{}, errors.New("Unrecognized character found in source " + string(source[0]))

}

func tokenize(file string) ([]TokenGroup, error) {
	// implement a parser for the Java language and return an AST

	source := []rune(file)

	tokens := []TokenGroup{}

	for len(source) > 0 {

		if unicode.IsSpace(source[0]) {
			source = source[1:]
			continue
		}

		token, err := scan(source)
		if err != nil {
			return tokens, err
		}

		// fmt.Println(token)

		if token.token != SPACE {
			tokens = append(tokens, token)
		}

		source = source[len(token.value):]
	}

	return tokens, nil
}

func find(tokens []TOKEN, token TOKEN) (bool, int) {
	for i, t := range tokens {
		if t == token {
			return true, i
		}
	}
	return false, -1
}

func get_indent(indent int) string {
	return strings.Repeat("    ", indent)
}

type Class = struct {
	name   string
	indent int
}

func transform(tokens []TokenGroup, js_source string, indent int, class_data Class) (string, error) {

	if len(tokens) == 0 || tokens[len(tokens)-1].token == EOF {
		return js_source, nil
	}

	token := tokens[0]

	// SKIPS TOKENS
	is_not_js_keyword, _ := find(NON_JS_KEYWORD, token.token)
	if is_not_js_keyword {
		if token.token == PACKAGE {
			for tokens[0].token != SEMICOLON {
				tokens = tokens[1:]
			}
		}

		return transform(tokens[1:], js_source, indent, class_data)
	}

	if token.token == SEMICOLON {
		return transform(tokens[1:], js_source+";\n"+get_indent(indent), indent, class_data)
	}

	is_java_type, _ := find(NATIVE_JAVA_TYPES, token.token)
	if is_java_type || token.token == IDENTIFIER {

		if tokens[1].token == LEFT_PARENTHESES {
			return transform(tokens[1:], js_source+tokens[0].value, indent, class_data)
		} else if tokens[1].token == ASSIGN {
			return transform(tokens[2:], js_source+tokens[0].value+" = ", indent, class_data)
		}

		name := tokens[1]
		type_token := tokens[2]

		if name.token != IDENTIFIER {
			return transform(tokens[1:], js_source+tokens[0].value, indent, class_data)
		}

		pre_key := ""

		switch type_token.token {
		case LEFT_PARENTHESES:
			if class_data.name == "" {
				pre_key = "function "
				pre_key += name.value + "("
			} else if class_data.name == name.value {
				pre_key = "constructor ("
			} else {
				pre_key = name.value + "("
			}
			tokens = tokens[3:]
			isArg := 0
			js_source += pre_key
			for tokens[0].token != RIGHT_PARENTHESES {
				if tokens[0].token == IDENTIFIER {
					isArg += 1
				}
				if isArg%2 == 0 {
					js_source += tokens[0].value
				}
				tokens = tokens[1:]
			}
			return transform(tokens, js_source, indent, class_data)
		case ASSIGN:
			if class_data.name == "" {
				js_source += "let "
			}
			js_source += name.value
		case SEMICOLON:
			if class_data.name == "" {
				js_source += "let "
			}
			js_source += name.value
		case ARRAY:
			fmt.Println(class_data)
			if class_data.name == "" {
				js_source += "let "
			}
			js_source += name.value + " "
		default:
			js_source += name.value
		}

		// if class_data.name != ""

		// switch type_token.token {
		// case LEFT_PARENTHESES:
		// 	if class_data.name == "" {
		// 		js_source += "function "
		// 	}
		// 	js_source += name.value + "("
		// 	if class_data.name == name.value {
		// 		js_source = "constructor ("
		// 	}
		// 	tokens = tokens[3:]
		// 	isArg := 0
		// 	for tokens[0].token != RIGHT_PARENTHESES {
		// 		if tokens[0].token == IDENTIFIER {
		// 			isArg += 1
		// 		}
		// 		if isArg%2 == 0 {
		// 			js_source += tokens[0].value
		// 		}
		// 		tokens = tokens[1:]
		// 	}
		// 	return transform(tokens, js_source, indent, class_data)
		// case ASSIGN:
		// 	js_source += "let " + name.value
		// case SEMICOLON:
		// 	js_source += "let " + name.value
		// case ARRAY:
		// 	js_source += "let " + name.value + " "
		// default:
		// 	js_source += name.value
		// }

		return transform(tokens[2:], js_source, indent, class_data)
	}

	if token.token == IMPORT {
		js_source += "import \"" + token.value[7:] + "\";\n"

		return transform(tokens[2:], js_source, indent, class_data)
	}

	if token.token == ASSERT {
		return transform(tokens[1:], js_source+"console.assert(", indent, class_data)
	}

	if token.token == BREAK {
		return transform(tokens[1:], js_source+"break;\n"+get_indent(indent-1), indent-1, class_data)
	}

	if token.token == CASE {
		return transform(tokens[1:], js_source+"case \n"+get_indent(indent+1), indent+1, class_data)
	}

	if token.token == CATCH {
		return transform(tokens[1:], js_source+"catch", indent+1, class_data)
	}

	if token.token == CLASS {
		return transform(tokens[2:], js_source+"class "+tokens[1].value, indent, Class{tokens[1].value, indent})
	}

	if token.token == LEFT_BRACES {
		return transform(tokens[1:], js_source+" { \n"+get_indent(indent+1), indent+1, class_data)
	}

	if token.token == RIGHT_BRACES {

		if class_data.indent == indent-1 {
			class_data.name = ""
			class_data.indent = 0
		}
		return transform(tokens[1:], js_source[0:len(js_source)-4]+"}\n"+get_indent(indent-1), indent-1, class_data)
	}

	if token.token == RIGHT_PARENTHESES {
		return transform(tokens[1:], js_source+")", indent, class_data)
	}
	if token.token == LEFT_PARENTHESES {
		return transform(tokens[1:], js_source+"(", indent, class_data)
	}

	if token.token == NEW {
		if tokens[1].token != IDENTIFIER {
			return js_source, errors.New("missing IDENTIFIER")
		}
		return transform(tokens[2:], js_source+"new "+tokens[1].value, indent, class_data)
	}

	if token.token == ASSIGN {
		return transform(tokens[1:], js_source+" = ", indent, class_data)
	}

	if token.token == RETURN {
		return transform(tokens[1:], js_source+token.value+" ", indent, class_data)
	}
	if token.token == DECORATOR {
		return transform(tokens[1:], js_source+token.value+"\n"+get_indent(indent), indent, class_data)
	}

	return transform(tokens[1:], js_source+token.value, indent, class_data)

}

func saveFile(content string, path string) {

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
