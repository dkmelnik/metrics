// Package staticlint provides analyzers for static code checking in Go.
//
// Mechanism description of multichecker:
//
// The analyzers provided in this package are run using multichecker,
// which is provided in golang.org/x/tools/go/analysis/multichecker.
// multichecker takes a list of analyzers as input and runs them all on the files
// passed as arguments. More information about multichecker can be found in its documentation.
//
// Each analyzer in this package represents a static code check for a specific type of problem in Go code.
// All analyzers are aimed at detecting typical errors or anti-patterns in the code.
//
// The exitinmain analyzer detects direct calls to os.Exit in the main function.
//
// Description of analyzers:
//
// - assign: Checks assignments in if statements and for-range loops, highlighting potential code issues.
// - copylock: Checks mutex captures for data copying operations, identifying potential concurrency issues.
// - httpresponse: Checks the usage of HTTP response handling functions, identifying potential errors or anti-patterns.
// - loopclosure: Identifies variable captures in closures within loops, which can lead to unexpected behavior due to Go's scoping rules.
// - nilfunc: Checks function calls with arguments equal to nil, which can lead to runtime panics.
// - printf: Checks formatting arguments in printf-like functions, detecting potential errors related to incorrect string formatting.
// - shift: Checks arguments of left and right shift operators, detecting potential issues with overflow and misuse.
// - structtag: Checks the correct usage of struct tags in the code, identifying potential errors or anti-patterns.
// - tests: Checks the correctness of test code, detecting potential issues with test functions and structures.
// - unreachable: Checks unreachable code, identifying parts of the program that will never be executed.
//
// Place the analyzer in the cmd/staticlint directory of your project.
// Add documentation in the godoc format, describing in detail the multichecker execution mechanism,
// as well as each analyzer and its purpose.
package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"honnef.co/go/tools/staticcheck"
)

var exitInMainAnalyzer = &analysis.Analyzer{
	Name: "exitinmain",
	Doc:  "reports direct os.Exit calls in main functions",
	Run:  runExitInMainAnalyzer,
}

func runExitInMainAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if fn.Name.Name != "main" {
				continue
			}

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selExpr.X.(*ast.Ident)
				if !ok || ident.Name != "os" || selExpr.Sel.Name != "Exit" {
					return true
				}
				pass.Reportf(callExpr.Pos(), "direct use of os.Exit in main function")
				return true
			})
		}
	}

	return nil, nil
}

func main() {
	checks := []*analysis.Analyzer{
		assign.Analyzer,
		copylock.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		exitInMainAnalyzer,
	}

	for _, value := range staticcheck.Analyzers {
		checks = append(checks, value.Analyzer)
	}

	multichecker.Main(checks...)
}
