package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"time"
)

type empty struct{}

var tmpUsed = 0
var funcs map[string]func([]ast.Expr) string
var compareKinds = map[token.Token]empty{
	token.EQL: empty{},
	token.LSS: empty{},
	token.GTR: empty{},
}

func (c *Go2bpp) initFuncs() {
	funcs = make(map[string]func([]ast.Expr) string)
	funcs["concat"] = func(elems []ast.Expr) string {
		args := make([]interface{}, len(elems))
		tmplt := "[CONCAT"
		for i, val := range elems {
			args[i] = parseExpr(val)
			tmplt += " %s"
		}
		return fmt.Sprintf(tmplt+"]", args...)
	}
	funcs["print"] = func(params []ast.Expr) string {
		return parseExprs(params)
	}
	funcs["choose"] = func(params []ast.Expr) string {
		return fmt.Sprintf("[CHOOSE %s]", parseExprs(params))
	}
	funcs["repeat"] = func(params []ast.Expr) string {
		return fmt.Sprintf("[REPEAT %s %s]", parseExpr(params[0]), parseExpr(params[1]))
	}
	funcs["randint"] = func(params []ast.Expr) string {
		if len(params) < 2 {
			return "Not enough parameters!"
		}
		return fmt.Sprintf("[RANDINT %s %s]", parseExpr(params[0]), parseExpr(params[1]))
	}
	funcs["randfloat"] = func(params []ast.Expr) string {
		if len(params) < 2 {
			return "Not enough parameters!"
		}
		return fmt.Sprintf("[RANDOM %s %s]", parseExpr(params[0]), parseExpr(params[1]))
	}
	funcs["floor"] = func(params []ast.Expr) string {
		return fmt.Sprintf("[FLOOR %s]", parseExprs(params))
	}
	funcs["ceil"] = func(params []ast.Expr) string {
		return fmt.Sprintf("[CEIL %s]", parseExprs(params))
	}
	funcs["round"] = func(params []ast.Expr) string {
		return fmt.Sprintf("[ROUND %s]", parseExprs(params))
	}
	funcs["arg"] = func(params []ast.Expr) string {
		return fmt.Sprintf("[ARGS %s]", parseExpr(params[0]))
	}
	for k := range funcs {
		c.BuiltinFuncs += "<code>" + k + "</code>, "
	}
	c.BuiltinFuncs = c.BuiltinFuncs[:len(c.BuiltinFuncs)-2] + "."
}

func (c *Go2bpp) parse() string {
	tmpUsed = 0
	if !c.HasLoaded {
		return ""
	}
	code := c.Editor.Call("getValue").String()
	go func() {
		time.Sleep(time.Second / 2)
		c.Editor.Call("destroy")
		c.SetupEditor()
		c.Editor.Call("setValue", code)
	}()
	start := time.Now()
	out := ""
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		return err.Error()
	}
	decls := f.Decls
	for _, decl := range decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if decl.(*ast.FuncDecl).Name.Obj.Name == "main" { // B++ currently doesn't support functions
				out += parseStmt(decl.(*ast.FuncDecl).Body)
			}
		}
	}
	c.TimeTaken = fmt.Sprintf("%v", time.Since(start))
	return out
}

func parseStmt(stmt ast.Stmt) string {
	switch stmt.(type) {
	case *ast.AssignStmt:
		stm := stmt.(*ast.AssignStmt)
		return fmt.Sprintf("[DEFINE %s %s]", stm.Lhs[0].(*ast.Ident).Name, parseExprs(stm.Rhs))
	case *ast.ExprStmt:
		return parseExpr(stmt.(*ast.ExprStmt).X)
	case *ast.IfStmt:
		stm := stmt.(*ast.IfStmt)
		body := strings.Replace(parseStmt(stm.Body), "\n", "", -1)
		el := `""`
		if stm.Else != nil {
			el = strings.Replace(parseStmt(stm.Else), "\n", "", -1)
		}
		return fmt.Sprintf("[IF %s %s %s]", parseExpr(stm.Cond), body, el)
	case *ast.BlockStmt:
		out := ""
		for _, val := range stmt.(*ast.BlockStmt).List {
			out += parseStmt(val) + "\n"
		}
		return out
	}
	return fmt.Sprintf("Unable to parse statement of type %s!", reflect.TypeOf(stmt).Elem().Name())
}

func parseExprs(exprs []ast.Expr) string {
	out := ""
	for _, val := range exprs {
		out += parseExpr(val)
	}
	return out
}

func parseExpr(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.BasicLit:
		return expr.(*ast.BasicLit).Value
	case *ast.BinaryExpr:
		stm := expr.(*ast.BinaryExpr)
		_, exists := compareKinds[stm.Op]
		if exists {
			if stm.Op == token.EQL {
				stm.Op = token.ASSIGN
			}
			return fmt.Sprintf("[COMPARE %s %s %s]", parseExpr(stm.X), stm.Op, parseExpr(stm.Y))
		}
		return fmt.Sprintf("[MATH %s %s %s]", parseExpr(stm.X), stm.Op, parseExpr(stm.Y))
	case *ast.Ident:
		return fmt.Sprintf("[VAR %s]", expr.(*ast.Ident).Name)
	case *ast.ParenExpr:
		return parseExpr(expr.(*ast.ParenExpr).X)
	case *ast.CompositeLit:
		elems := expr.(*ast.CompositeLit).Elts
		args := make([]interface{}, len(elems))
		tmplt := "[ARRAY"
		for i, val := range elems {
			args[i] = parseExpr(val)
			tmplt += " %s"
		}
		return fmt.Sprintf(tmplt+"]", args...)
	case *ast.IndexExpr:
		exp := expr.(*ast.IndexExpr)
		return fmt.Sprintf("[INDEX %s %s]", parseExpr(exp.X), exp.Index.(*ast.BasicLit).Value)
	case *ast.CallExpr:
		call := expr.(*ast.CallExpr)
		fun, exists := funcs[call.Fun.(*ast.Ident).Name]
		if !exists {
			return fmt.Sprintf("No such function %s!", call.Fun.(*ast.Ident).Name)
		}
		return fun(call.Args)
	}
	return fmt.Sprintf("Unable to parse expression of type %s!", reflect.TypeOf(expr).Elem().Name())
}
