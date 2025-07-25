// Copyright 2016 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/open-policy-agent/opa/v1/ast/internal/tokens"
	astJSON "github.com/open-policy-agent/opa/v1/ast/json"
	"github.com/open-policy-agent/opa/v1/util"
)

// DefaultRootDocument is the default root document.
//
// All package directives inside source files are implicitly prefixed with the
// DefaultRootDocument value.
var DefaultRootDocument = VarTerm("data")

// InputRootDocument names the document containing query arguments.
var InputRootDocument = VarTerm("input")

// SchemaRootDocument names the document containing external data schemas.
var SchemaRootDocument = VarTerm("schema")

// FunctionArgRootDocument names the document containing function arguments.
// It's only for internal usage, for referencing function arguments between
// the index and topdown.
var FunctionArgRootDocument = VarTerm("args")

// FutureRootDocument names the document containing new, to-become-default,
// features.
var FutureRootDocument = VarTerm("future")

// RegoRootDocument names the document containing new, to-become-default,
// features in a future versioned release.
var RegoRootDocument = VarTerm("rego")

// RootDocumentNames contains the names of top-level documents that can be
// referred to in modules and queries.
//
// Note, the schema document is not currently implemented in the evaluator so it
// is not registered as a root document name (yet).
var RootDocumentNames = NewSet(
	DefaultRootDocument,
	InputRootDocument,
)

// DefaultRootRef is a reference to the root of the default document.
//
// All refs to data in the policy engine's storage layer are prefixed with this ref.
var DefaultRootRef = Ref{DefaultRootDocument}

// InputRootRef is a reference to the root of the input document.
//
// All refs to query arguments are prefixed with this ref.
var InputRootRef = Ref{InputRootDocument}

// SchemaRootRef is a reference to the root of the schema document.
//
// All refs to schema documents are prefixed with this ref. Note, the schema
// document is not currently implemented in the evaluator so it is not
// registered as a root document ref (yet).
var SchemaRootRef = Ref{SchemaRootDocument}

// RootDocumentRefs contains the prefixes of top-level documents that all
// non-local references start with.
var RootDocumentRefs = NewSet(
	NewTerm(DefaultRootRef),
	NewTerm(InputRootRef),
)

// SystemDocumentKey is the name of the top-level key that identifies the system
// document.
const SystemDocumentKey = String("system")

// ReservedVars is the set of names that refer to implicitly ground vars.
var ReservedVars = NewVarSet(
	DefaultRootDocument.Value.(Var),
	InputRootDocument.Value.(Var),
)

// Wildcard represents the wildcard variable as defined in the language.
var Wildcard = &Term{Value: Var("_")}

// WildcardPrefix is the special character that all wildcard variables are
// prefixed with when the statement they are contained in is parsed.
const WildcardPrefix = "$"

// Keywords contains strings that map to language keywords.
var Keywords = KeywordsForRegoVersion(DefaultRegoVersion)

var KeywordsV0 = [...]string{
	"not",
	"package",
	"import",
	"as",
	"default",
	"else",
	"with",
	"null",
	"true",
	"false",
	"some",
}

var KeywordsV1 = [...]string{
	"not",
	"package",
	"import",
	"as",
	"default",
	"else",
	"with",
	"null",
	"true",
	"false",
	"some",
	"if",
	"contains",
	"in",
	"every",
}

func KeywordsForRegoVersion(v RegoVersion) []string {
	switch v {
	case RegoV0:
		return KeywordsV0[:]
	case RegoV1, RegoV0CompatV1:
		return KeywordsV1[:]
	}
	return nil
}

// IsKeyword returns true if s is a language keyword.
func IsKeyword(s string) bool {
	return IsInKeywords(s, Keywords)
}

func IsInKeywords(s string, keywords []string) bool {
	return slices.Contains(keywords, s)
}

// IsKeywordInRegoVersion returns true if s is a language keyword.
func IsKeywordInRegoVersion(s string, regoVersion RegoVersion) bool {
	switch regoVersion {
	case RegoV0:
		for _, x := range KeywordsV0 {
			if x == s {
				return true
			}
		}
	case RegoV1, RegoV0CompatV1:
		for _, x := range KeywordsV1 {
			if x == s {
				return true
			}
		}
	}

	return false
}

type (
	// Node represents a node in an AST. Nodes may be statements in a policy module
	// or elements of an ad-hoc query, expression, etc.
	Node interface {
		fmt.Stringer
		Loc() *Location
		SetLoc(*Location)
	}

	// Statement represents a single statement in a policy module.
	Statement interface {
		Node
	}
)

type (

	// Module represents a collection of policies (defined by rules)
	// within a namespace (defined by the package) and optional
	// dependencies on external documents (defined by imports).
	Module struct {
		Package     *Package       `json:"package"`
		Imports     []*Import      `json:"imports,omitempty"`
		Annotations []*Annotations `json:"annotations,omitempty"`
		Rules       []*Rule        `json:"rules,omitempty"`
		Comments    []*Comment     `json:"comments,omitempty"`
		stmts       []Statement
		regoVersion RegoVersion
	}

	// Comment contains the raw text from the comment in the definition.
	Comment struct {
		// TODO: these fields have inconsistent JSON keys with other structs in this package.
		Text     []byte
		Location *Location
	}

	// Package represents the namespace of the documents produced
	// by rules inside the module.
	Package struct {
		Path     Ref       `json:"path"`
		Location *Location `json:"location,omitempty"`
	}

	// Import represents a dependency on a document outside of the policy
	// namespace. Imports are optional.
	Import struct {
		Path     *Term     `json:"path"`
		Alias    Var       `json:"alias,omitempty"`
		Location *Location `json:"location,omitempty"`
	}

	// Rule represents a rule as defined in the language. Rules define the
	// content of documents that represent policy decisions.
	Rule struct {
		Default     bool           `json:"default,omitempty"`
		Head        *Head          `json:"head"`
		Body        Body           `json:"body"`
		Else        *Rule          `json:"else,omitempty"`
		Location    *Location      `json:"location,omitempty"`
		Annotations []*Annotations `json:"annotations,omitempty"`

		// Module is a pointer to the module containing this rule. If the rule
		// was NOT created while parsing/constructing a module, this should be
		// left unset. The pointer is not included in any standard operations
		// on the rule (e.g., printing, comparison, visiting, etc.)
		Module *Module `json:"-"`

		generatedBody bool
	}

	// Head represents the head of a rule.
	Head struct {
		Name      Var       `json:"name,omitempty"`
		Reference Ref       `json:"ref,omitempty"`
		Args      Args      `json:"args,omitempty"`
		Key       *Term     `json:"key,omitempty"`
		Value     *Term     `json:"value,omitempty"`
		Assign    bool      `json:"assign,omitempty"`
		Location  *Location `json:"location,omitempty"`

		keywords       []tokens.Token
		generatedValue bool
	}

	// Args represents zero or more arguments to a rule.
	Args []*Term

	// Body represents one or more expressions contained inside a rule or user
	// function.
	Body []*Expr

	// Expr represents a single expression contained inside the body of a rule.
	Expr struct {
		With      []*With   `json:"with,omitempty"`
		Terms     any       `json:"terms"`
		Index     int       `json:"index"`
		Generated bool      `json:"generated,omitempty"`
		Negated   bool      `json:"negated,omitempty"`
		Location  *Location `json:"location,omitempty"`

		generatedFrom *Expr
		generates     []*Expr
	}

	// SomeDecl represents a variable declaration statement. The symbols are variables.
	SomeDecl struct {
		Symbols  []*Term   `json:"symbols"`
		Location *Location `json:"location,omitempty"`
	}

	Every struct {
		Key      *Term     `json:"key"`
		Value    *Term     `json:"value"`
		Domain   *Term     `json:"domain"`
		Body     Body      `json:"body"`
		Location *Location `json:"location,omitempty"`
	}

	// With represents a modifier on an expression.
	With struct {
		Target   *Term     `json:"target"`
		Value    *Term     `json:"value"`
		Location *Location `json:"location,omitempty"`
	}
)

// SetModuleRegoVersion sets the RegoVersion for the Module.
func SetModuleRegoVersion(mod *Module, v RegoVersion) {
	mod.regoVersion = v
}

// Compare returns an integer indicating whether mod is less than, equal to,
// or greater than other.
func (mod *Module) Compare(other *Module) int {
	if mod == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}
	if cmp := mod.Package.Compare(other.Package); cmp != 0 {
		return cmp
	}
	if cmp := importsCompare(mod.Imports, other.Imports); cmp != 0 {
		return cmp
	}
	if cmp := annotationsCompare(mod.Annotations, other.Annotations); cmp != 0 {
		return cmp
	}
	return rulesCompare(mod.Rules, other.Rules)
}

// Copy returns a deep copy of mod.
func (mod *Module) Copy() *Module {
	cpy := *mod
	cpy.Rules = make([]*Rule, len(mod.Rules))

	nodes := make(map[Node]Node, len(mod.Rules)+len(mod.Imports)+1 /* package */)

	for i := range mod.Rules {
		cpy.Rules[i] = mod.Rules[i].Copy()
		cpy.Rules[i].Module = &cpy
		nodes[mod.Rules[i]] = cpy.Rules[i]
	}

	cpy.Imports = make([]*Import, len(mod.Imports))
	for i := range mod.Imports {
		cpy.Imports[i] = mod.Imports[i].Copy()
		nodes[mod.Imports[i]] = cpy.Imports[i]
	}

	cpy.Package = mod.Package.Copy()
	nodes[mod.Package] = cpy.Package

	cpy.Annotations = make([]*Annotations, len(mod.Annotations))
	for i, a := range mod.Annotations {
		cpy.Annotations[i] = a.Copy(nodes[a.node])
	}

	cpy.Comments = make([]*Comment, len(mod.Comments))
	for i := range mod.Comments {
		cpy.Comments[i] = mod.Comments[i].Copy()
	}

	cpy.stmts = make([]Statement, len(mod.stmts))
	for i := range mod.stmts {
		cpy.stmts[i] = nodes[mod.stmts[i]]
	}

	return &cpy
}

// Equal returns true if mod equals other.
func (mod *Module) Equal(other *Module) bool {
	return mod.Compare(other) == 0
}

func (mod *Module) String() string {
	byNode := map[Node][]*Annotations{}
	for _, a := range mod.Annotations {
		byNode[a.node] = append(byNode[a.node], a)
	}

	appendAnnotationStrings := func(buf []string, node Node) []string {
		if as, ok := byNode[node]; ok {
			for i := range as {
				buf = append(buf, "# METADATA")
				buf = append(buf, "# "+as[i].String())
			}
		}
		return buf
	}

	buf := []string{}
	buf = appendAnnotationStrings(buf, mod.Package)
	buf = append(buf, mod.Package.String())

	if len(mod.Imports) > 0 {
		buf = append(buf, "")
		for _, imp := range mod.Imports {
			buf = appendAnnotationStrings(buf, imp)
			buf = append(buf, imp.String())
		}
	}
	if len(mod.Rules) > 0 {
		buf = append(buf, "")
		for _, rule := range mod.Rules {
			buf = appendAnnotationStrings(buf, rule)
			buf = append(buf, rule.stringWithOpts(toStringOpts{regoVersion: mod.regoVersion}))
		}
	}
	return strings.Join(buf, "\n")
}

// RuleSet returns a RuleSet containing named rules in the mod.
func (mod *Module) RuleSet(name Var) RuleSet {
	rs := NewRuleSet()
	for _, rule := range mod.Rules {
		if rule.Head.Name.Equal(name) {
			rs.Add(rule)
		}
	}
	return rs
}

// UnmarshalJSON parses bs and stores the result in mod. The rules in the module
// will have their module pointer set to mod.
func (mod *Module) UnmarshalJSON(bs []byte) error {

	// Declare a new type and use a type conversion to avoid recursively calling
	// Module#UnmarshalJSON.
	type module Module

	if err := util.UnmarshalJSON(bs, (*module)(mod)); err != nil {
		return err
	}

	WalkRules(mod, func(rule *Rule) bool {
		rule.Module = mod
		return false
	})

	return nil
}

func (mod *Module) regoV1Compatible() bool {
	return mod.regoVersion == RegoV1 || mod.regoVersion == RegoV0CompatV1
}

func (mod *Module) RegoVersion() RegoVersion {
	return mod.regoVersion
}

// SetRegoVersion sets the RegoVersion for the module.
// Note: Setting a rego-version that does not match the module's rego-version might have unintended consequences.
func (mod *Module) SetRegoVersion(v RegoVersion) {
	mod.regoVersion = v
}

// NewComment returns a new Comment object.
func NewComment(text []byte) *Comment {
	return &Comment{
		Text: text,
	}
}

// Loc returns the location of the comment in the definition.
func (c *Comment) Loc() *Location {
	if c == nil {
		return nil
	}
	return c.Location
}

// SetLoc sets the location on c.
func (c *Comment) SetLoc(loc *Location) {
	c.Location = loc
}

func (c *Comment) String() string {
	return "#" + string(c.Text)
}

// Copy returns a deep copy of c.
func (c *Comment) Copy() *Comment {
	cpy := *c
	cpy.Text = make([]byte, len(c.Text))
	copy(cpy.Text, c.Text)
	return &cpy
}

// Equal returns true if this comment equals the other comment.
// Unlike other equality checks on AST nodes, comment equality
// depends on location.
func (c *Comment) Equal(other *Comment) bool {
	return c.Location.Equal(other.Location) && bytes.Equal(c.Text, other.Text)
}

// Compare returns an integer indicating whether pkg is less than, equal to,
// or greater than other.
func (pkg *Package) Compare(other *Package) int {
	return termSliceCompare(pkg.Path, other.Path)
}

// Copy returns a deep copy of pkg.
func (pkg *Package) Copy() *Package {
	cpy := *pkg
	cpy.Path = pkg.Path.Copy()
	return &cpy
}

// Equal returns true if pkg is equal to other.
func (pkg *Package) Equal(other *Package) bool {
	return pkg.Compare(other) == 0
}

// Loc returns the location of the Package in the definition.
func (pkg *Package) Loc() *Location {
	if pkg == nil {
		return nil
	}
	return pkg.Location
}

// SetLoc sets the location on pkg.
func (pkg *Package) SetLoc(loc *Location) {
	pkg.Location = loc
}

func (pkg *Package) String() string {
	if pkg == nil {
		return "<illegal nil package>"
	} else if len(pkg.Path) <= 1 {
		return fmt.Sprintf("package <illegal path %q>", pkg.Path)
	}
	// Omit head as all packages have the DefaultRootDocument prepended at parse time.
	path := make(Ref, len(pkg.Path)-1)
	path[0] = VarTerm(string(pkg.Path[1].Value.(String)))
	copy(path[1:], pkg.Path[2:])
	return fmt.Sprintf("package %v", path)
}

func (pkg *Package) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"path": pkg.Path,
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.Package {
		if pkg.Location != nil {
			data["location"] = pkg.Location
		}
	}

	return json.Marshal(data)
}

// IsValidImportPath returns an error indicating if the import path is invalid.
// If the import path is valid, err is nil.
func IsValidImportPath(v Value) (err error) {
	switch v := v.(type) {
	case Var:
		if !v.Equal(DefaultRootDocument.Value) && !v.Equal(InputRootDocument.Value) {
			return fmt.Errorf("invalid path %v: path must begin with input or data", v)
		}
	case Ref:
		if err := IsValidImportPath(v[0].Value); err != nil {
			return fmt.Errorf("invalid path %v: path must begin with input or data", v)
		}
		for _, e := range v[1:] {
			if _, ok := e.Value.(String); !ok {
				return fmt.Errorf("invalid path %v: path elements must be strings", v)
			}
		}
	default:
		return fmt.Errorf("invalid path %v: path must be ref or var", v)
	}
	return nil
}

// Compare returns an integer indicating whether imp is less than, equal to,
// or greater than other.
func (imp *Import) Compare(other *Import) int {
	if imp == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}
	if cmp := Compare(imp.Path, other.Path); cmp != 0 {
		return cmp
	}

	return VarCompare(imp.Alias, other.Alias)
}

// Copy returns a deep copy of imp.
func (imp *Import) Copy() *Import {
	cpy := *imp
	cpy.Path = imp.Path.Copy()
	return &cpy
}

// Equal returns true if imp is equal to other.
func (imp *Import) Equal(other *Import) bool {
	return imp.Compare(other) == 0
}

// Loc returns the location of the Import in the definition.
func (imp *Import) Loc() *Location {
	if imp == nil {
		return nil
	}
	return imp.Location
}

// SetLoc sets the location on imp.
func (imp *Import) SetLoc(loc *Location) {
	imp.Location = loc
}

// Name returns the variable that is used to refer to the imported virtual
// document. This is the alias if defined otherwise the last element in the
// path.
func (imp *Import) Name() Var {
	if len(imp.Alias) != 0 {
		return imp.Alias
	}
	switch v := imp.Path.Value.(type) {
	case Var:
		return v
	case Ref:
		if len(v) == 1 {
			return v[0].Value.(Var)
		}
		return Var(v[len(v)-1].Value.(String))
	}
	panic("illegal import")
}

func (imp *Import) String() string {
	buf := []string{"import", imp.Path.String()}
	if len(imp.Alias) > 0 {
		buf = append(buf, "as", imp.Alias.String())
	}
	return strings.Join(buf, " ")
}

func (imp *Import) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"path": imp.Path,
	}

	if len(imp.Alias) != 0 {
		data["alias"] = imp.Alias
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.Import {
		if imp.Location != nil {
			data["location"] = imp.Location
		}
	}

	return json.Marshal(data)
}

// Compare returns an integer indicating whether rule is less than, equal to,
// or greater than other.
func (rule *Rule) Compare(other *Rule) int {
	if rule == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}
	if cmp := rule.Head.Compare(other.Head); cmp != 0 {
		return cmp
	}
	if rule.Default != other.Default {
		if !rule.Default {
			return -1
		}
		return 1
	}
	if cmp := rule.Body.Compare(other.Body); cmp != 0 {
		return cmp
	}

	if cmp := annotationsCompare(rule.Annotations, other.Annotations); cmp != 0 {
		return cmp
	}

	return rule.Else.Compare(other.Else)
}

// Copy returns a deep copy of rule.
func (rule *Rule) Copy() *Rule {
	cpy := *rule
	cpy.Head = rule.Head.Copy()
	cpy.Body = rule.Body.Copy()

	if len(cpy.Annotations) > 0 {
		cpy.Annotations = make([]*Annotations, len(rule.Annotations))
		for i, a := range rule.Annotations {
			cpy.Annotations[i] = a.Copy(&cpy)
		}
	}

	if cpy.Else != nil {
		cpy.Else = rule.Else.Copy()
	}
	return &cpy
}

// Equal returns true if rule is equal to other.
func (rule *Rule) Equal(other *Rule) bool {
	return rule.Compare(other) == 0
}

// Loc returns the location of the Rule in the definition.
func (rule *Rule) Loc() *Location {
	if rule == nil {
		return nil
	}
	return rule.Location
}

// SetLoc sets the location on rule.
func (rule *Rule) SetLoc(loc *Location) {
	rule.Location = loc
}

// Path returns a ref referring to the document produced by this rule. If rule
// is not contained in a module, this function panics.
// Deprecated: Poor handling of ref rules. Use `(*Rule).Ref()` instead.
func (rule *Rule) Path() Ref {
	if rule.Module == nil {
		panic("assertion failed")
	}
	return rule.Module.Package.Path.Extend(rule.Head.Ref().GroundPrefix())
}

// Ref returns a ref referring to the document produced by this rule. If rule
// is not contained in a module, this function panics. The returned ref may
// contain variables in the last position.
func (rule *Rule) Ref() Ref {
	if rule.Module == nil {
		panic("assertion failed")
	}
	return rule.Module.Package.Path.Extend(rule.Head.Ref())
}

func (rule *Rule) String() string {
	regoVersion := DefaultRegoVersion
	if rule.Module != nil {
		regoVersion = rule.Module.RegoVersion()
	}
	return rule.stringWithOpts(toStringOpts{regoVersion: regoVersion})
}

type toStringOpts struct {
	regoVersion RegoVersion
}

func (o toStringOpts) RegoVersion() RegoVersion {
	if o.regoVersion == RegoUndefined {
		return DefaultRegoVersion
	}
	return o.regoVersion
}

func (rule *Rule) stringWithOpts(opts toStringOpts) string {
	buf := []string{}
	if rule.Default {
		buf = append(buf, "default")
	}
	buf = append(buf, rule.Head.stringWithOpts(opts))
	if !rule.Default {
		switch opts.RegoVersion() {
		case RegoV1, RegoV0CompatV1:
			buf = append(buf, "if")
		}
		buf = append(buf, "{", rule.Body.String(), "}")
	}
	if rule.Else != nil {
		buf = append(buf, rule.Else.elseString(opts))
	}
	return strings.Join(buf, " ")
}

func (rule *Rule) isFunction() bool {
	return len(rule.Head.Args) > 0
}

func (rule *Rule) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"head": rule.Head,
		"body": rule.Body,
	}

	if rule.Default {
		data["default"] = true
	}

	if rule.Else != nil {
		data["else"] = rule.Else
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.Rule {
		if rule.Location != nil {
			data["location"] = rule.Location
		}
	}

	if len(rule.Annotations) != 0 {
		data["annotations"] = rule.Annotations
	}

	return json.Marshal(data)
}

func (rule *Rule) elseString(opts toStringOpts) string {
	var buf []string

	buf = append(buf, "else")

	value := rule.Head.Value
	if value != nil {
		buf = append(buf, "=", value.String())
	}

	switch opts.RegoVersion() {
	case RegoV1, RegoV0CompatV1:
		buf = append(buf, "if")
	}

	buf = append(buf, "{", rule.Body.String(), "}")

	if rule.Else != nil {
		buf = append(buf, rule.Else.elseString(opts))
	}

	return strings.Join(buf, " ")
}

// NewHead returns a new Head object. If args are provided, the first will be
// used for the key and the second will be used for the value.
func NewHead(name Var, args ...*Term) *Head {
	head := &Head{
		Name:      name, // backcompat
		Reference: []*Term{NewTerm(name)},
	}
	if len(args) == 0 {
		return head
	}
	head.Key = args[0]
	if len(args) == 1 {
		return head
	}
	head.Value = args[1]
	if head.Key != nil && head.Value != nil {
		head.Reference = head.Reference.Append(args[0])
	}
	return head
}

// VarHead creates a head object, initializes its Name and Location and returns the new head.
// NOTE: The JSON options argument is no longer used, and kept only for backwards compatibility.
func VarHead(name Var, location *Location, _ *astJSON.Options) *Head {
	h := NewHead(name)
	h.Reference[0].Location = location
	return h
}

// RefHead returns a new Head object with the passed Ref. If args are provided,
// the first will be used for the value.
func RefHead(ref Ref, args ...*Term) *Head {
	head := &Head{}
	head.SetRef(ref)
	if len(ref) < 2 {
		head.Name = ref[0].Value.(Var)
	}
	if len(args) >= 1 {
		head.Value = args[0]
	}
	return head
}

// DocKind represents the collection of document types that can be produced by rules.
type DocKind byte

const (
	// CompleteDoc represents a document that is completely defined by the rule.
	CompleteDoc = iota

	// PartialSetDoc represents a set document that is partially defined by the rule.
	PartialSetDoc

	// PartialObjectDoc represents an object document that is partially defined by the rule.
	PartialObjectDoc
) // TODO(sr): Deprecate?

// DocKind returns the type of document produced by this rule.
func (head *Head) DocKind() DocKind {
	if head.Key != nil {
		if head.Value != nil {
			return PartialObjectDoc
		}
		return PartialSetDoc
	} else if head.HasDynamicRef() {
		return PartialObjectDoc
	}
	return CompleteDoc
}

type RuleKind byte

const (
	SingleValue = iota
	MultiValue
)

// RuleKind returns the type of rule this is
func (head *Head) RuleKind() RuleKind {
	// NOTE(sr): This is bit verbose, since the key is irrelevant for single vs
	//           multi value, but as good a spot as to assert the invariant.
	switch {
	case head.Value != nil:
		return SingleValue
	case head.Key != nil:
		return MultiValue
	default:
		panic("unreachable")
	}
}

// Ref returns the Ref of the rule. If it doesn't have one, it's filled in
// via the Head's Name.
func (head *Head) Ref() Ref {
	if len(head.Reference) > 0 {
		return head.Reference
	}
	return Ref{&Term{Value: head.Name}}
}

// SetRef can be used to set a rule head's Reference
func (head *Head) SetRef(r Ref) {
	head.Reference = r
}

// Compare returns an integer indicating whether head is less than, equal to,
// or greater than other.
func (head *Head) Compare(other *Head) int {
	if head == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}
	if head.Assign && !other.Assign {
		return -1
	} else if !head.Assign && other.Assign {
		return 1
	}
	if cmp := Compare(head.Args, other.Args); cmp != 0 {
		return cmp
	}
	if cmp := Compare(head.Reference, other.Reference); cmp != 0 {
		return cmp
	}
	if cmp := VarCompare(head.Name, other.Name); cmp != 0 {
		return cmp
	}
	if cmp := Compare(head.Key, other.Key); cmp != 0 {
		return cmp
	}
	return Compare(head.Value, other.Value)
}

// Copy returns a deep copy of head.
func (head *Head) Copy() *Head {
	cpy := *head
	cpy.Reference = head.Reference.Copy()
	cpy.Args = head.Args.Copy()
	cpy.Key = head.Key.Copy()
	cpy.Value = head.Value.Copy()
	cpy.keywords = nil
	return &cpy
}

// Equal returns true if this head equals other.
func (head *Head) Equal(other *Head) bool {
	return head.Compare(other) == 0
}

func (head *Head) String() string {
	return head.stringWithOpts(toStringOpts{})
}

func (head *Head) stringWithOpts(opts toStringOpts) string {
	buf := strings.Builder{}
	buf.WriteString(head.Ref().String())
	containsAdded := false

	switch {
	case len(head.Args) != 0:
		buf.WriteString(head.Args.String())
	case len(head.Reference) == 1 && head.Key != nil:
		switch opts.RegoVersion() {
		case RegoV0:
			buf.WriteRune('[')
			buf.WriteString(head.Key.String())
			buf.WriteRune(']')
		default:
			containsAdded = true
			buf.WriteString(" contains ")
			buf.WriteString(head.Key.String())
		}
	}
	if head.Value != nil {
		if head.Assign {
			buf.WriteString(" := ")
		} else {
			buf.WriteString(" = ")
		}
		buf.WriteString(head.Value.String())
	} else if !containsAdded && head.Name == "" && head.Key != nil {
		buf.WriteString(" contains ")
		buf.WriteString(head.Key.String())
	}
	return buf.String()
}

func (head *Head) MarshalJSON() ([]byte, error) {
	var loc *Location
	if astJSON.GetOptions().MarshalOptions.IncludeLocation.Head && head.Location != nil {
		loc = head.Location
	}

	// NOTE(sr): we do this to override the rendering of `head.Reference`.
	// It's still what'll be used via the default means of encoding/json
	// for unmarshaling a json object into a Head struct!
	type h Head
	return json.Marshal(struct {
		h
		Ref      Ref       `json:"ref"`
		Location *Location `json:"location,omitempty"`
	}{
		h:        h(*head),
		Ref:      head.Ref(),
		Location: loc,
	})
}

// Vars returns a set of vars found in the head.
func (head *Head) Vars() VarSet {
	vis := NewVarVisitor()
	// TODO: improve test coverage for this.
	if head.Args != nil {
		vis.WalkArgs(head.Args)
	}
	if head.Key != nil {
		vis.Walk(head.Key)
	}
	if head.Value != nil {
		vis.Walk(head.Value)
	}
	if len(head.Reference) > 0 {
		vis.WalkRef(head.Reference[1:])
	}
	return vis.vars
}

// Loc returns the Location of head.
func (head *Head) Loc() *Location {
	if head == nil {
		return nil
	}
	return head.Location
}

// SetLoc sets the location on head.
func (head *Head) SetLoc(loc *Location) {
	head.Location = loc
}

func (head *Head) HasDynamicRef() bool {
	pos := head.Reference.Dynamic()
	return pos > 0 && (pos < len(head.Reference))
}

// Copy returns a deep copy of a.
func (a Args) Copy() Args {
	cpy := Args{}
	for _, t := range a {
		cpy = append(cpy, t.Copy())
	}
	return cpy
}

func (a Args) String() string {
	buf := make([]string, 0, len(a))
	for _, t := range a {
		buf = append(buf, t.String())
	}
	return "(" + strings.Join(buf, ", ") + ")"
}

// Loc returns the Location of a.
func (a Args) Loc() *Location {
	if len(a) == 0 {
		return nil
	}
	return a[0].Location
}

// SetLoc sets the location on a.
func (a Args) SetLoc(loc *Location) {
	if len(a) != 0 {
		a[0].SetLocation(loc)
	}
}

// Vars returns a set of vars that appear in a.
func (a Args) Vars() VarSet {
	vis := NewVarVisitor()
	vis.WalkArgs(a)
	return vis.vars
}

// NewBody returns a new Body containing the given expressions. The indices of
// the immediate expressions will be reset.
func NewBody(exprs ...*Expr) Body {
	for i, expr := range exprs {
		expr.Index = i
	}
	return Body(exprs)
}

// MarshalJSON returns JSON encoded bytes representing body.
func (body Body) MarshalJSON() ([]byte, error) {
	// Serialize empty Body to empty array. This handles both the empty case and the
	// nil case (whereas by default the result would be null if body was nil.)
	if len(body) == 0 {
		return []byte(`[]`), nil
	}
	ret, err := json.Marshal([]*Expr(body))
	return ret, err
}

// Append adds the expr to the body and updates the expr's index accordingly.
func (body *Body) Append(expr *Expr) {
	n := len(*body)
	expr.Index = n
	*body = append(*body, expr)
}

// Set sets the expr in the body at the specified position and updates the
// expr's index accordingly.
func (body Body) Set(expr *Expr, pos int) {
	body[pos] = expr
	expr.Index = pos
}

// Compare returns an integer indicating whether body is less than, equal to,
// or greater than other.
//
// If body is a subset of other, it is considered less than (and vice versa).
func (body Body) Compare(other Body) int {
	minLen := min(len(other), len(body))
	for i := range minLen {
		if cmp := body[i].Compare(other[i]); cmp != 0 {
			return cmp
		}
	}
	if len(body) < len(other) {
		return -1
	}
	if len(other) < len(body) {
		return 1
	}
	return 0
}

// Copy returns a deep copy of body.
func (body Body) Copy() Body {
	cpy := make(Body, len(body))
	for i := range body {
		cpy[i] = body[i].Copy()
	}
	return cpy
}

// Contains returns true if this body contains the given expression.
func (body Body) Contains(x *Expr) bool {
	return slices.ContainsFunc(body, x.Equal)
}

// Equal returns true if this Body is equal to the other Body.
func (body Body) Equal(other Body) bool {
	return body.Compare(other) == 0
}

// Hash returns the hash code for the Body.
func (body Body) Hash() int {
	s := 0
	for _, e := range body {
		s += e.Hash()
	}
	return s
}

// IsGround returns true if all of the expressions in the Body are ground.
func (body Body) IsGround() bool {
	for _, e := range body {
		if !e.IsGround() {
			return false
		}
	}
	return true
}

// Loc returns the location of the Body in the definition.
func (body Body) Loc() *Location {
	if len(body) == 0 {
		return nil
	}
	return body[0].Location
}

// SetLoc sets the location on body.
func (body Body) SetLoc(loc *Location) {
	if len(body) != 0 {
		body[0].SetLocation(loc)
	}
}

func (body Body) String() string {
	buf := make([]string, 0, len(body))
	for _, v := range body {
		buf = append(buf, v.String())
	}
	return strings.Join(buf, "; ")
}

// Vars returns a VarSet containing variables in body. The params can be set to
// control which vars are included.
func (body Body) Vars(params VarVisitorParams) VarSet {
	vis := NewVarVisitor().WithParams(params)
	vis.WalkBody(body)
	return vis.Vars()
}

// NewExpr returns a new Expr object.
func NewExpr(terms any) *Expr {
	switch terms.(type) {
	case *SomeDecl, *Every, *Term, []*Term: // ok
	default:
		panic("unreachable")
	}
	return &Expr{
		Negated: false,
		Terms:   terms,
		Index:   0,
		With:    nil,
	}
}

// Complement returns a copy of this expression with the negation flag flipped.
func (expr *Expr) Complement() *Expr {
	cpy := *expr
	cpy.Negated = !cpy.Negated
	return &cpy
}

// ComplementNoWith returns a copy of this expression with the negation flag flipped
// and the with modifier removed. This is the same as calling .Complement().NoWith()
// but without making an intermediate copy.
func (expr *Expr) ComplementNoWith() *Expr {
	cpy := *expr
	cpy.Negated = !cpy.Negated
	cpy.With = nil
	return &cpy
}

// Equal returns true if this Expr equals the other Expr.
func (expr *Expr) Equal(other *Expr) bool {
	return expr.Compare(other) == 0
}

// Compare returns an integer indicating whether expr is less than, equal to,
// or greater than other.
//
// Expressions are compared as follows:
//
// 1. Declarations are always less than other expressions.
// 2. Preceding expression (by Index) is always less than the other expression.
// 3. Non-negated expressions are always less than negated expressions.
// 4. Single term expressions are always less than built-in expressions.
//
// Otherwise, the expression terms are compared normally. If both expressions
// have the same terms, the modifiers are compared.
func (expr *Expr) Compare(other *Expr) int {

	if expr == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}

	o1 := expr.sortOrder()
	o2 := other.sortOrder()
	if o1 < o2 {
		return -1
	} else if o2 < o1 {
		return 1
	}

	switch {
	case expr.Index < other.Index:
		return -1
	case expr.Index > other.Index:
		return 1
	}

	switch {
	case expr.Negated && !other.Negated:
		return 1
	case !expr.Negated && other.Negated:
		return -1
	}

	switch t := expr.Terms.(type) {
	case *Term:
		if cmp := Compare(t.Value, other.Terms.(*Term).Value); cmp != 0 {
			return cmp
		}
	case []*Term:
		if cmp := termSliceCompare(t, other.Terms.([]*Term)); cmp != 0 {
			return cmp
		}
	case *SomeDecl:
		if cmp := Compare(t, other.Terms.(*SomeDecl)); cmp != 0 {
			return cmp
		}
	case *Every:
		if cmp := Compare(t, other.Terms.(*Every)); cmp != 0 {
			return cmp
		}
	}

	return withSliceCompare(expr.With, other.With)
}

func (expr *Expr) sortOrder() int {
	switch expr.Terms.(type) {
	case *SomeDecl:
		return 0
	case *Term:
		return 1
	case []*Term:
		return 2
	case *Every:
		return 3
	}
	return -1
}

// CopyWithoutTerms returns a deep copy of expr without its Terms
func (expr *Expr) CopyWithoutTerms() *Expr {
	cpy := *expr

	if expr.With != nil {
		cpy.With = make([]*With, len(expr.With))
		for i := range expr.With {
			cpy.With[i] = expr.With[i].Copy()
		}
	}

	return &cpy
}

// Copy returns a deep copy of expr.
func (expr *Expr) Copy() *Expr {

	cpy := expr.CopyWithoutTerms()

	switch ts := expr.Terms.(type) {
	case *SomeDecl:
		cpy.Terms = ts.Copy()
	case []*Term:
		cpy.Terms = termSliceCopy(ts)
	case *Term:
		cpy.Terms = ts.Copy()
	case *Every:
		cpy.Terms = ts.Copy()
	}

	return cpy
}

// Hash returns the hash code of the Expr.
func (expr *Expr) Hash() int {
	s := expr.Index
	switch ts := expr.Terms.(type) {
	case *SomeDecl:
		s += ts.Hash()
	case []*Term:
		for _, t := range ts {
			s += t.Value.Hash()
		}
	case *Term:
		s += ts.Value.Hash()
	}
	if expr.Negated {
		s++
	}
	for _, w := range expr.With {
		s += w.Hash()
	}
	return s
}

// IncludeWith returns a copy of expr with the with modifier appended.
func (expr *Expr) IncludeWith(target *Term, value *Term) *Expr {
	cpy := *expr
	cpy.With = append(cpy.With, &With{Target: target, Value: value})
	return &cpy
}

// NoWith returns a copy of expr where the with modifier has been removed.
func (expr *Expr) NoWith() *Expr {
	cpy := *expr
	cpy.With = nil
	return &cpy
}

// IsEquality returns true if this is an equality expression.
func (expr *Expr) IsEquality() bool {
	return isGlobalBuiltin(expr, Var(Equality.Name))
}

// IsAssignment returns true if this an assignment expression.
func (expr *Expr) IsAssignment() bool {
	return isGlobalBuiltin(expr, Var(Assign.Name))
}

// IsCall returns true if this expression calls a function.
func (expr *Expr) IsCall() bool {
	_, ok := expr.Terms.([]*Term)
	return ok
}

// IsEvery returns true if this expression is an 'every' expression.
func (expr *Expr) IsEvery() bool {
	_, ok := expr.Terms.(*Every)
	return ok
}

// IsSome returns true if this expression is a 'some' expression.
func (expr *Expr) IsSome() bool {
	_, ok := expr.Terms.(*SomeDecl)
	return ok
}

// Operator returns the name of the function or built-in this expression refers
// to. If this expression is not a function call, returns nil.
func (expr *Expr) Operator() Ref {
	op := expr.OperatorTerm()
	if op == nil {
		return nil
	}
	return op.Value.(Ref)
}

// OperatorTerm returns the name of the function or built-in this expression
// refers to. If this expression is not a function call, returns nil.
func (expr *Expr) OperatorTerm() *Term {
	terms, ok := expr.Terms.([]*Term)
	if !ok || len(terms) == 0 {
		return nil
	}
	return terms[0]
}

// Operand returns the term at the zero-based pos. If the expr does not include
// at least pos+1 terms, this function returns nil.
func (expr *Expr) Operand(pos int) *Term {
	terms, ok := expr.Terms.([]*Term)
	if !ok {
		return nil
	}
	idx := pos + 1
	if idx < len(terms) {
		return terms[idx]
	}
	return nil
}

// Operands returns the built-in function operands.
func (expr *Expr) Operands() []*Term {
	terms, ok := expr.Terms.([]*Term)
	if !ok {
		return nil
	}
	return terms[1:]
}

// IsGround returns true if all of the expression terms are ground.
func (expr *Expr) IsGround() bool {
	switch ts := expr.Terms.(type) {
	case []*Term:
		for _, t := range ts[1:] {
			if !t.IsGround() {
				return false
			}
		}
	case *Term:
		return ts.IsGround()
	}
	return true
}

// SetOperator sets the expr's operator and returns the expr itself. If expr is
// not a call expr, this function will panic.
func (expr *Expr) SetOperator(term *Term) *Expr {
	expr.Terms.([]*Term)[0] = term
	return expr
}

// SetLocation sets the expr's location and returns the expr itself.
func (expr *Expr) SetLocation(loc *Location) *Expr {
	expr.Location = loc
	return expr
}

// Loc returns the Location of expr.
func (expr *Expr) Loc() *Location {
	if expr == nil {
		return nil
	}
	return expr.Location
}

// SetLoc sets the location on expr.
func (expr *Expr) SetLoc(loc *Location) {
	expr.SetLocation(loc)
}

func (expr *Expr) String() string {
	buf := make([]string, 0, 2+len(expr.With))
	if expr.Negated {
		buf = append(buf, "not")
	}
	switch t := expr.Terms.(type) {
	case []*Term:
		if expr.IsEquality() && validEqAssignArgCount(expr) {
			buf = append(buf, fmt.Sprintf("%v %v %v", t[1], Equality.Infix, t[2]))
		} else {
			buf = append(buf, Call(t).String())
		}
	case fmt.Stringer:
		buf = append(buf, t.String())
	}

	for i := range expr.With {
		buf = append(buf, expr.With[i].String())
	}

	return strings.Join(buf, " ")
}

func (expr *Expr) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"terms": expr.Terms,
		"index": expr.Index,
	}

	if len(expr.With) > 0 {
		data["with"] = expr.With
	}

	if expr.Generated {
		data["generated"] = true
	}

	if expr.Negated {
		data["negated"] = true
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.Expr {
		if expr.Location != nil {
			data["location"] = expr.Location
		}
	}

	return json.Marshal(data)
}

// UnmarshalJSON parses the byte array and stores the result in expr.
func (expr *Expr) UnmarshalJSON(bs []byte) error {
	v := map[string]any{}
	if err := util.UnmarshalJSON(bs, &v); err != nil {
		return err
	}
	return unmarshalExpr(expr, v)
}

// Vars returns a VarSet containing variables in expr. The params can be set to
// control which vars are included.
func (expr *Expr) Vars(params VarVisitorParams) VarSet {
	vis := NewVarVisitor().WithParams(params)
	vis.Walk(expr)
	return vis.Vars()
}

// NewBuiltinExpr creates a new Expr object with the supplied terms.
// The builtin operator must be the first term.
func NewBuiltinExpr(terms ...*Term) *Expr {
	return &Expr{Terms: terms}
}

func (expr *Expr) CogeneratedExprs() []*Expr {
	visited := map[*Expr]struct{}{}
	visitCogeneratedExprs(expr, func(e *Expr) bool {
		if expr.Equal(e) {
			return true
		}
		if _, ok := visited[e]; ok {
			return true
		}
		visited[e] = struct{}{}
		return false
	})

	result := make([]*Expr, 0, len(visited))
	for e := range visited {
		result = append(result, e)
	}
	return result
}

func (expr *Expr) BaseCogeneratedExpr() *Expr {
	if expr.generatedFrom == nil {
		return expr
	}
	return expr.generatedFrom.BaseCogeneratedExpr()
}

func visitCogeneratedExprs(expr *Expr, f func(*Expr) bool) {
	if parent := expr.generatedFrom; parent != nil {
		if stop := f(parent); !stop {
			visitCogeneratedExprs(parent, f)
		}
	}
	for _, child := range expr.generates {
		if stop := f(child); !stop {
			visitCogeneratedExprs(child, f)
		}
	}
}

func (d *SomeDecl) String() string {
	if call, ok := d.Symbols[0].Value.(Call); ok {
		if len(call) == 4 {
			return "some " + call[1].String() + ", " + call[2].String() + " in " + call[3].String()
		}
		return "some " + call[1].String() + " in " + call[2].String()
	}
	buf := make([]string, len(d.Symbols))
	for i := range buf {
		buf[i] = d.Symbols[i].String()
	}
	return "some " + strings.Join(buf, ", ")
}

// SetLoc sets the Location on d.
func (d *SomeDecl) SetLoc(loc *Location) {
	d.Location = loc
}

// Loc returns the Location of d.
func (d *SomeDecl) Loc() *Location {
	return d.Location
}

// Copy returns a deep copy of d.
func (d *SomeDecl) Copy() *SomeDecl {
	cpy := *d
	cpy.Symbols = termSliceCopy(d.Symbols)
	return &cpy
}

// Compare returns an integer indicating whether d is less than, equal to, or
// greater than other.
func (d *SomeDecl) Compare(other *SomeDecl) int {
	return termSliceCompare(d.Symbols, other.Symbols)
}

// Hash returns a hash code of d.
func (d *SomeDecl) Hash() int {
	return termSliceHash(d.Symbols)
}

func (d *SomeDecl) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"symbols": d.Symbols,
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.SomeDecl {
		if d.Location != nil {
			data["location"] = d.Location
		}
	}

	return json.Marshal(data)
}

func (q *Every) String() string {
	if q.Key != nil {
		return fmt.Sprintf("every %s, %s in %s { %s }",
			q.Key,
			q.Value,
			q.Domain,
			q.Body)
	}
	return fmt.Sprintf("every %s in %s { %s }",
		q.Value,
		q.Domain,
		q.Body)
}

func (q *Every) Loc() *Location {
	return q.Location
}

func (q *Every) SetLoc(l *Location) {
	q.Location = l
}

// Copy returns a deep copy of d.
func (q *Every) Copy() *Every {
	cpy := *q
	cpy.Key = q.Key.Copy()
	cpy.Value = q.Value.Copy()
	cpy.Domain = q.Domain.Copy()
	cpy.Body = q.Body.Copy()
	return &cpy
}

func (q *Every) Compare(other *Every) int {
	for _, terms := range [][2]*Term{
		{q.Key, other.Key},
		{q.Value, other.Value},
		{q.Domain, other.Domain},
	} {
		if d := Compare(terms[0], terms[1]); d != 0 {
			return d
		}
	}
	return q.Body.Compare(other.Body)
}

// KeyValueVars returns the key and val arguments of an `every`
// expression, if they are non-nil and not wildcards.
func (q *Every) KeyValueVars() VarSet {
	vis := NewVarVisitor()
	if q.Key != nil {
		vis.Walk(q.Key)
	}
	vis.Walk(q.Value)
	return vis.vars
}

func (q *Every) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"key":    q.Key,
		"value":  q.Value,
		"domain": q.Domain,
		"body":   q.Body,
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.Every {
		if q.Location != nil {
			data["location"] = q.Location
		}
	}

	return json.Marshal(data)
}

func (w *With) String() string {
	return "with " + w.Target.String() + " as " + w.Value.String()
}

// Equal returns true if this With is equals the other With.
func (w *With) Equal(other *With) bool {
	return Compare(w, other) == 0
}

// Compare returns an integer indicating whether w is less than, equal to, or
// greater than other.
func (w *With) Compare(other *With) int {
	if w == nil {
		if other == nil {
			return 0
		}
		return -1
	} else if other == nil {
		return 1
	}
	if cmp := Compare(w.Target, other.Target); cmp != 0 {
		return cmp
	}
	return Compare(w.Value, other.Value)
}

// Copy returns a deep copy of w.
func (w *With) Copy() *With {
	cpy := *w
	cpy.Value = w.Value.Copy()
	cpy.Target = w.Target.Copy()
	return &cpy
}

// Hash returns the hash code of the With.
func (w With) Hash() int {
	return w.Target.Hash() + w.Value.Hash()
}

// SetLocation sets the location on w.
func (w *With) SetLocation(loc *Location) *With {
	w.Location = loc
	return w
}

// Loc returns the Location of w.
func (w *With) Loc() *Location {
	if w == nil {
		return nil
	}
	return w.Location
}

// SetLoc sets the location on w.
func (w *With) SetLoc(loc *Location) {
	w.Location = loc
}

func (w *With) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		"target": w.Target,
		"value":  w.Value,
	}

	if astJSON.GetOptions().MarshalOptions.IncludeLocation.With {
		if w.Location != nil {
			data["location"] = w.Location
		}
	}

	return json.Marshal(data)
}

// Copy returns a deep copy of the AST node x. If x is not an AST node, x is returned unmodified.
func Copy(x any) any {
	switch x := x.(type) {
	case *Module:
		return x.Copy()
	case *Package:
		return x.Copy()
	case *Import:
		return x.Copy()
	case *Rule:
		return x.Copy()
	case *Head:
		return x.Copy()
	case Args:
		return x.Copy()
	case Body:
		return x.Copy()
	case *Expr:
		return x.Copy()
	case *With:
		return x.Copy()
	case *SomeDecl:
		return x.Copy()
	case *Every:
		return x.Copy()
	case *Term:
		return x.Copy()
	case *ArrayComprehension:
		return x.Copy()
	case *SetComprehension:
		return x.Copy()
	case *ObjectComprehension:
		return x.Copy()
	case Set:
		return x.Copy()
	case *object:
		return x.Copy()
	case *Array:
		return x.Copy()
	case Ref:
		return x.Copy()
	case Call:
		return x.Copy()
	case *Comment:
		return x.Copy()
	}
	return x
}

// RuleSet represents a collection of rules that produce a virtual document.
type RuleSet []*Rule

// NewRuleSet returns a new RuleSet containing the given rules.
func NewRuleSet(rules ...*Rule) RuleSet {
	rs := make(RuleSet, 0, len(rules))
	for _, rule := range rules {
		rs.Add(rule)
	}
	return rs
}

// Add inserts the rule into rs.
func (rs *RuleSet) Add(rule *Rule) {
	for _, exist := range *rs {
		if exist.Equal(rule) {
			return
		}
	}
	*rs = append(*rs, rule)
}

// Contains returns true if rs contains rule.
func (rs RuleSet) Contains(rule *Rule) bool {
	for i := range rs {
		if rs[i].Equal(rule) {
			return true
		}
	}
	return false
}

// Diff returns a new RuleSet containing rules in rs that are not in other.
func (rs RuleSet) Diff(other RuleSet) RuleSet {
	result := NewRuleSet()
	for i := range rs {
		if !other.Contains(rs[i]) {
			result.Add(rs[i])
		}
	}
	return result
}

// Equal returns true if rs equals other.
func (rs RuleSet) Equal(other RuleSet) bool {
	return len(rs.Diff(other)) == 0 && len(other.Diff(rs)) == 0
}

// Merge returns a ruleset containing the union of rules from rs an other.
func (rs RuleSet) Merge(other RuleSet) RuleSet {
	result := NewRuleSet()
	for i := range rs {
		result.Add(rs[i])
	}
	for i := range other {
		result.Add(other[i])
	}
	return result
}

func (rs RuleSet) String() string {
	buf := make([]string, 0, len(rs))
	for _, rule := range rs {
		buf = append(buf, rule.String())
	}
	return "{" + strings.Join(buf, ", ") + "}"
}

// Returns true if the equality or assignment expression referred to by expr
// has a valid number of arguments.
func validEqAssignArgCount(expr *Expr) bool {
	return len(expr.Operands()) == 2
}

// this function checks if the expr refers to a non-namespaced (global) built-in
// function like eq, gt, plus, etc.
func isGlobalBuiltin(expr *Expr, name Var) bool {
	terms, ok := expr.Terms.([]*Term)
	if !ok {
		return false
	}

	// NOTE(tsandall): do not use Term#Equal or Value#Compare to avoid
	// allocation here.
	ref, ok := terms[0].Value.(Ref)
	if !ok || len(ref) != 1 {
		return false
	}
	if head, ok := ref[0].Value.(Var); ok {
		return head.Equal(name)
	}
	return false
}
