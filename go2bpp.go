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

var config = map[string]int{
	"loops": 5,
	"goto":  1,
}

var tmpUsed = 0
var funcs map[string]func([]ast.Expr) (string, string)
var compareKinds = map[token.Token]empty{
	token.EQL: empty{},
	token.LSS: empty{},
	token.GTR: empty{},
}
var funcReturns map[string]string

func (c *Go2bpp) initFuncs() {
	funcs = make(map[string]func([]ast.Expr) (string, string))
	funcs["concat"] = func(elems []ast.Expr) (string, string) {
		args := make([]interface{}, len(elems))
		pres := ""
		var pr string
		tmplt := "[CONCAT"
		for i, val := range elems {
			args[i], pr = parseExpr(val)
			tmplt += " %s"
			pres += pr + "\n"
		}
		return pres + fmt.Sprintf(tmplt+"]", args...), ""
	}
	funcs["print"] = func(params []ast.Expr) (string, string) {
		return parseExprs(params), ""
	}
	funcs["choose"] = func(params []ast.Expr) (string, string) {
		return fmt.Sprintf("[CHOOSE %s]", parseExprs(params)), ""
	}
	funcs["repeat"] = func(params []ast.Expr) (string, string) {
		gt, pre := parseExpr(params[0])
		gt2, pre2 := parseExpr(params[1])
		return fmt.Sprintf("%s\n$s\n[REPEAT %s %s]", pre, pre2, gt, gt2), ""
	}
	funcs["randint"] = func(params []ast.Expr) (string, string) {
		if len(params) < 2 {
			return "Not enough parameters!", ""
		}
		gt, pre := parseExpr(params[0])
		gt2, pre2 := parseExpr(params[1])
		return fmt.Sprintf("%s%s[RANDINT %s %s]", pre, pre2, gt, gt2), ""
	}
	funcs["randfloat"] = func(params []ast.Expr) (string, string) {
		if len(params) < 2 {
			return "Not enough parameters!", ""
		}
		gt, pre := parseExpr(params[0])
		gt2, pre2 := parseExpr(params[1])
		return fmt.Sprintf("%s%s[RANDOM %s %s]", pre, pre2, gt, gt2), ""
	}
	funcs["floor"] = func(params []ast.Expr) (string, string) {
		gt, pre := parseExpr(params[0])
		if len(pre) > 0 {
			pre += "\n"
		}
		return fmt.Sprintf("%s[FLOOR %s]", pre, gt), ""
	}
	funcs["ceil"] = func(params []ast.Expr) (string, string) {
		gt, pre := parseExpr(params[0])
		if len(pre) > 0 {
			pre += "\n"
		}
		return fmt.Sprintf("%s[CEIL %s]", pre, gt), ""
	}
	funcs["round"] = func(params []ast.Expr) (string, string) {
		gt, pre := parseExpr(params[0])
		if len(pre) > 0 {
			pre += "\n"
		}
		return fmt.Sprintf("%s[ROUND %s]", pre, gt), ""
	}
	funcs["arg"] = func(params []ast.Expr) (string, string) {
		gt, pre := parseExpr(params[0])
		if len(pre) > 0 {
			pre += "\n"
		}
		return fmt.Sprintf("%s[ARGS %s]", pre, gt), ""
	}
	for k := range funcs {
		c.BuiltinFuncs += "<code>" + k + "</code>, "
	}
	c.BuiltinFuncs = c.BuiltinFuncs[:len(c.BuiltinFuncs)-2] + "."
	for k := range config {
		c.BuiltinConfig += "<code>" + k + "</code>, "
	}
	c.BuiltinConfig = c.BuiltinConfig[:len(c.BuiltinConfig)-2] + "."
}

func (c *Go2bpp) parse() string {
	tmpUsed = 0
	funcReturns = make(map[string]string)
	if !c.HasLoaded {
		return ""
	}
	c.initFuncs()
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
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return err.Error()
	}
	for _, val := range f.Comments {
		txt := val.Text()
		if strings.HasPrefix(txt, "config") {
			txt = strings.Replace(txt, "config:", "", 1)
			var title string
			var content int
			_, err = fmt.Sscanf(txt, "config %s %d", &title, &content)
			if err != nil {
				return err.Error()
			}
			config[title] = content
		}
	}
	if config["goto"] == 1 {
		out += "[GOTO main]\n"
	}
	decls := f.Decls
	for _, decl := range decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if decl.(*ast.FuncDecl).Name.Obj.Name == "main" {
				if config["goto"] == 1 {
					out += "[SECTION main]\n"
				}
				out += parseStmt(decl.(*ast.FuncDecl).Body, decl.(*ast.FuncDecl).Name.Obj.Name)
			} else if config["goto"] == 1 {
				fn := decl.(*ast.FuncDecl)
				out += fmt.Sprintf("[SECTION %s]\n", fn.Name.Obj.Name)
				out += parseStmt(fn.Body, decl.(*ast.FuncDecl).Name.Obj.Name)
				funcs[decl.(*ast.FuncDecl).Name.Obj.Name] = func(params []ast.Expr) (string, string) {
					txt := fmt.Sprintf("[DEFINE returnpos tmp%d]\n[GOTO %s]\n[SECTION tmp%d]", tmpUsed, fn.Name.Obj.Name, tmpUsed)
					tmpUsed++
					if len(params) != len(fn.Type.Params.List) {
						return "not enough args", ""
					}
					prmSet := ""
					for i := range params {
						prmName := fn.Type.Params.List[i].Type.(*ast.Ident).Name
						gt, pre := parseExpr(params[i])
						prmSet += fmt.Sprintf("%s\n[DEFINE %s %s]\n", pre, prmName, gt)
					}
					return fmt.Sprintf("[VAR %s]", funcReturns[fn.Name.Obj.Name]), prmSet + txt
				}
				out += "[GOTO [VAR returnpos]]\n"
			}
		}
	}
	c.TimeTaken = fmt.Sprintf("%v", time.Since(start))
	if config["goto"] == 1 {
		out += "[SECTION end]\n"
	}
	for strings.Count(out, "\n\n") > 0 {
		out = strings.Replace(out, "\n\n", "\n", -1)
	}
	return out
}

func parseStmt(stmt ast.Stmt, funcName string) string {
	switch stmt.(type) {
	case *ast.AssignStmt:
		stm := stmt.(*ast.AssignStmt)
		gt, pre := parseExpr(stm.Rhs[0])
		return fmt.Sprintf("%s\n[DEFINE %s %s]", pre, stm.Lhs[0].(*ast.Ident).Name, gt)
	case *ast.ExprStmt:
		gt, pr := parseExpr(stmt.(*ast.ExprStmt).X)
		return pr + "\n" + gt
	case *ast.IfStmt:
		stm := stmt.(*ast.IfStmt)
		if config["goto"] == 1 {
			tmpIfEnd := tmpUsed
			tmpUsed++
			gt, pre := parseExpr(stm.Cond)
			body := parseStmt(stm.Body, funcName)
			body = fmt.Sprintf("%s\n[GOTO tmp%d]", body, tmpIfEnd)
			tmpIf := tmpUsed
			tmpUsed++
			tmpEl := 0
			el := `""`
			if stm.Else != nil {
				tmpEl = tmpUsed
				tmpUsed++
				el = fmt.Sprintf("[GOTO tmp%d]", tmpEl)
			}
			txt := fmt.Sprintf("%s\n[IF %s [GOTO tmp%d] %s]\n[SECTION tmp%d]\n%s", pre, gt, tmpIf, el, tmpIf, body)
			if stm.Else != nil {
				txt += fmt.Sprintf("\n[SECTION tmp%d]\n%s\n[GOTO tmp%d]", tmpEl, parseStmt(stm.Else, funcName), tmpIfEnd)
			}
			txt += fmt.Sprintf("\n[SECTION tmp%d]", tmpIfEnd)
			return txt
		} else {
			body := strings.Replace(parseStmt(stm.Body, funcName), "\n", "", -1)
			el := `""`
			if stm.Else != nil {
				el = strings.Replace(parseStmt(stm.Else, funcName), "\n", "", -1)
			}
			gt, pre := parseExpr(stm.Cond)
			return fmt.Sprintf("%s\n[IF %s %s %s]", pre, gt, body, el)
		}
	case *ast.BlockStmt:
		out := ""
		for _, val := range stmt.(*ast.BlockStmt).List {
			out += parseStmt(val, funcName) + "\n"
		}
		return out
	case *ast.ForStmt:
		stm := stmt.(*ast.ForStmt)
		if config["goto"] == 1 {
			out := ""
			out += parseStmt(stm.Init, funcName) + "\n"
			tmpBody := tmpUsed
			tmpUsed++
			out += fmt.Sprintf("[SECTION tmp%d]\n", tmpBody)
			out += parseStmt(stm.Body, funcName)
			out += parseStmt(stm.Post, funcName) + "\n"
			gt, pre := parseExpr(stm.Cond)
			out += fmt.Sprintf("%s\n[IF %s %s %s]\n", pre, gt, fmt.Sprintf("[GOTO tmp%d]", tmpBody), `""`)
			return out
		} else {
			out := ""
			out += parseStmt(stm.Init, funcName) + "\n"
			for i := 0; i < config["loops"]; i++ {
				out += parseStmt(&ast.IfStmt{
					Body: stm.Body,
					Cond: stm.Cond,
				}, funcName) + "\n"
				out += parseStmt(stm.Post, funcName) + "\n"
			}
			return out
		}
	case *ast.IncDecStmt:
		stm := stmt.(*ast.IncDecStmt)
		vr := stm.X.(*ast.Ident).Name
		if stm.Tok == token.INC {
			return fmt.Sprintf("[DEFINE %s [MATH [VAR %s] + 1]]", vr, vr)
		}
		return fmt.Sprintf("[DEFINE %s [MATH [VAR %s] - 1]]", vr, vr)
	case *ast.ReturnStmt:
		if funcName == "main" {
			return "[GOTO end]"
		}
		stm := stmt.(*ast.ReturnStmt)
		gt, pre := parseExpr(stm.Results[0])
		funcReturns[funcName] = fmt.Sprintf("ret%d", tmpUsed)
		tmpUsed++
		return fmt.Sprintf("%s\n[DEFINE ret%d %s]", pre, tmpUsed-1, gt)
	}
	return fmt.Sprintf("Unable to parse statement of type %s!", reflect.TypeOf(stmt).Elem().Name())
}

func parseExprs(exprs []ast.Expr) string {
	pres := ""
	out := ""
	for _, val := range exprs {
		gt, pr := parseExpr(val)
		out += gt
		pres += pr + "\n"
	}
	return pres + out
}

func parseExpr(expr ast.Expr) (string, string) {
	switch expr.(type) {
	case *ast.BasicLit:
		return expr.(*ast.BasicLit).Value, ""
	case *ast.BinaryExpr:
		stm := expr.(*ast.BinaryExpr)
		get, pre := parseExpr(stm.X)
		get2, pre2 := parseExpr(stm.Y)
		_, exists := compareKinds[stm.Op]
		if exists {
			if stm.Op == token.EQL {
				stm.Op = token.ASSIGN
			}
			return fmt.Sprintf("%s%s[COMPARE %s %s %s]", pre, pre2, get, stm.Op, get2), ""
		}
		return fmt.Sprintf("%s%s[MATH %s %s %s]", pre, pre2, get, stm.Op, get2), ""
	case *ast.Ident:
		return fmt.Sprintf("[VAR %s]", expr.(*ast.Ident).Name), ""
	case *ast.ParenExpr:
		return parseExpr(expr.(*ast.ParenExpr).X)
	case *ast.CompositeLit:
		elems := expr.(*ast.CompositeLit).Elts
		args := make([]interface{}, len(elems))
		pre := ""
		var pr string
		tmplt := "[ARRAY"
		for i, val := range elems {
			args[i], pr = parseExpr(val)
			pre += pr + "\n"
			tmplt += " %s"
		}
		return pr + "\n" + fmt.Sprintf(tmplt+"]", args...), ""
	case *ast.IndexExpr:
		exp := expr.(*ast.IndexExpr)
		get, pr := parseExpr(exp.X)
		get2, pr2 := parseExpr(exp.Index)
		return fmt.Sprintf("%s%s[INDEX %s %s]", get, get2, pr, pr2), ""
	case *ast.CallExpr:
		call := expr.(*ast.CallExpr)
		fun, exists := funcs[call.Fun.(*ast.Ident).Name]
		if !exists {
			return fmt.Sprintf("No such function %s!", call.Fun.(*ast.Ident).Name), ""
		}
		return fun(call.Args)
	}
	return fmt.Sprintf("Unable to parse expression of type %s!", reflect.TypeOf(expr).Elem().Name()), ""
}
