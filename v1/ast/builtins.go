// Copyright 2016 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package ast

import (
	"strings"

	"github.com/open-policy-agent/opa/v1/types"
)

// Builtins is the registry of built-in functions supported by OPA.
// Call RegisterBuiltin to add a new built-in.
var Builtins []*Builtin

// RegisterBuiltin adds a new built-in function to the registry.
// NOTE: The underlying map storing built-ins is **not** thread-safe,
// and it's recommended to call this only during initialization, and never
// later. Registering built-ins after that point is unsupported and will
// likely lead to concurrent map read/write panics.
func RegisterBuiltin(b *Builtin) {
	Builtins = append(Builtins, b)
	BuiltinMap[b.Name] = b
	if len(b.Infix) > 0 {
		BuiltinMap[b.Infix] = b

		InternStringTerm(b.Infix)
	}

	InternStringTerm(b.Name)
	if strings.Contains(b.Name, ".") {
		InternStringTerm(strings.Split(b.Name, ".")...)
	}
}

// DefaultBuiltins is the registry of built-in functions supported in OPA
// by default. When adding a new built-in function to OPA, update this
// list.
var DefaultBuiltins = [...]*Builtin{
	// Unification/equality ("=")
	Equality,

	// Assignment (":=")
	Assign,

	// Membership, infix "in": `x in xs`
	Member,
	MemberWithKey,

	// Comparisons
	GreaterThan,
	GreaterThanEq,
	LessThan,
	LessThanEq,
	NotEqual,
	Equal,

	// Arithmetic
	Plus,
	Minus,
	Multiply,
	Divide,
	Ceil,
	Floor,
	Round,
	Abs,
	Rem,

	// Bitwise Arithmetic
	BitsOr,
	BitsAnd,
	BitsNegate,
	BitsXOr,
	BitsShiftLeft,
	BitsShiftRight,

	// Binary
	And,
	Or,

	// Aggregates
	Count,
	Sum,
	Product,
	Max,
	Min,
	Any,
	All,

	// Arrays
	ArrayConcat,
	ArraySlice,
	ArrayReverse,

	// Conversions
	ToNumber,

	// Casts (DEPRECATED)
	CastObject,
	CastNull,
	CastBoolean,
	CastString,
	CastSet,
	CastArray,

	// Regular Expressions
	RegexIsValid,
	RegexMatch,
	RegexMatchDeprecated,
	RegexSplit,
	GlobsMatch,
	RegexTemplateMatch,
	RegexFind,
	RegexFindAllStringSubmatch,
	RegexReplace,

	// Sets
	SetDiff,
	Intersection,
	Union,

	// Strings
	AnyPrefixMatch,
	AnySuffixMatch,
	Concat,
	FormatInt,
	IndexOf,
	IndexOfN,
	Substring,
	Lower,
	Upper,
	Contains,
	StringCount,
	StartsWith,
	EndsWith,
	Split,
	Replace,
	ReplaceN,
	Trim,
	TrimLeft,
	TrimPrefix,
	TrimRight,
	TrimSuffix,
	TrimSpace,
	Sprintf,
	StringReverse,
	RenderTemplate,

	// Numbers
	NumbersRange,
	NumbersRangeStep,
	RandIntn,

	// Encoding
	JSONMarshal,
	JSONMarshalWithOptions,
	JSONUnmarshal,
	JSONIsValid,
	Base64Encode,
	Base64Decode,
	Base64IsValid,
	Base64UrlEncode,
	Base64UrlEncodeNoPad,
	Base64UrlDecode,
	URLQueryDecode,
	URLQueryEncode,
	URLQueryEncodeObject,
	URLQueryDecodeObject,
	YAMLMarshal,
	YAMLUnmarshal,
	YAMLIsValid,
	HexEncode,
	HexDecode,

	// Object Manipulation
	ObjectUnion,
	ObjectUnionN,
	ObjectRemove,
	ObjectFilter,
	ObjectGet,
	ObjectKeys,
	ObjectSubset,

	// JSON Object Manipulation
	JSONFilter,
	JSONRemove,
	JSONPatch,

	// Tokens
	JWTDecode,
	JWTVerifyRS256,
	JWTVerifyRS384,
	JWTVerifyRS512,
	JWTVerifyPS256,
	JWTVerifyPS384,
	JWTVerifyPS512,
	JWTVerifyES256,
	JWTVerifyES384,
	JWTVerifyES512,
	JWTVerifyHS256,
	JWTVerifyHS384,
	JWTVerifyHS512,
	JWTDecodeVerify,
	JWTEncodeSignRaw,
	JWTEncodeSign,

	// Time
	NowNanos,
	ParseNanos,
	ParseRFC3339Nanos,
	ParseDurationNanos,
	Format,
	Date,
	Clock,
	Weekday,
	AddDate,
	Diff,

	// Crypto
	CryptoX509ParseCertificates,
	CryptoX509ParseAndVerifyCertificates,
	CryptoX509ParseAndVerifyCertificatesWithOptions,
	CryptoMd5,
	CryptoSha1,
	CryptoSha256,
	CryptoX509ParseCertificateRequest,
	CryptoX509ParseRSAPrivateKey,
	CryptoX509ParseKeyPair,
	CryptoParsePrivateKeys,
	CryptoHmacMd5,
	CryptoHmacSha1,
	CryptoHmacSha256,
	CryptoHmacSha512,
	CryptoHmacEqual,

	// Graphs
	WalkBuiltin,
	ReachableBuiltin,
	ReachablePathsBuiltin,

	// Sort
	Sort,

	// Types
	IsNumber,
	IsString,
	IsBoolean,
	IsArray,
	IsSet,
	IsObject,
	IsNull,
	TypeNameBuiltin,

	// HTTP
	HTTPSend,

	// GraphQL
	GraphQLParse,
	GraphQLParseAndVerify,
	GraphQLParseQuery,
	GraphQLParseSchema,
	GraphQLIsValid,
	GraphQLSchemaIsValid,

	// JSON Schema
	JSONSchemaVerify,
	JSONMatchSchema,

	// Cloud Provider Helpers
	ProvidersAWSSignReqObj,

	// Rego
	RegoParseModule,
	RegoMetadataChain,
	RegoMetadataRule,

	// OPA
	OPARuntime,

	// Tracing
	Trace,

	// Networking
	NetCIDROverlap,
	NetCIDRIntersects,
	NetCIDRContains,
	NetCIDRContainsMatches,
	NetCIDRExpand,
	NetCIDRMerge,
	NetLookupIPAddr,
	NetCIDRIsValid,

	// Glob
	GlobMatch,
	GlobQuoteMeta,

	// Units
	UnitsParse,
	UnitsParseBytes,

	// UUIDs
	UUIDRFC4122,
	UUIDParse,

	// SemVers
	SemVerIsValid,
	SemVerCompare,

	// Printing
	Print,
	InternalPrint,

	// Testing
	InternalTestCase,
}

// BuiltinMap provides a convenient mapping of built-in names to
// built-in definitions.
var BuiltinMap map[string]*Builtin

// Deprecated: Builtins can now be directly annotated with the
// Nondeterministic property, and when set to true, will be ignored
// for partial evaluation.
var IgnoreDuringPartialEval = []*Builtin{
	RandIntn,
	UUIDRFC4122,
	JWTDecodeVerify,
	JWTEncodeSignRaw,
	JWTEncodeSign,
	NowNanos,
	HTTPSend,
	OPARuntime,
	NetLookupIPAddr,
}

/**
 * Unification
 */

// Equality represents the "=" operator.
var Equality = &Builtin{
	Name:  "eq",
	Infix: "=",
	Decl: types.NewFunction(
		types.Args(types.A, types.A),
		types.B,
	),
	canSkipBctx: true,
}

/**
 * Assignment
 */

// Assign represents the assignment (":=") operator.
var Assign = &Builtin{
	Name:  "assign",
	Infix: ":=",
	Decl: types.NewFunction(
		types.Args(types.A, types.A),
		types.B,
	),
	canSkipBctx: true,
}

// Member represents the `in` (infix) operator.
var Member = &Builtin{
	Name:  "internal.member_2",
	Infix: "in",
	Decl: types.NewFunction(
		types.Args(
			types.A,
			types.A,
		),
		types.B,
	),
	canSkipBctx: true,
}

// MemberWithKey represents the `in` (infix) operator when used
// with two terms on the lhs, i.e., `k, v in obj`.
var MemberWithKey = &Builtin{
	Name:  "internal.member_3",
	Infix: "in",
	Decl: types.NewFunction(
		types.Args(
			types.A,
			types.A,
			types.A,
		),
		types.B,
	),
	canSkipBctx: true,
}

/**
 * Comparisons
 */
var comparison = category("comparison")

var GreaterThan = &Builtin{
	Name:       "gt",
	Infix:      ">",
	Categories: comparison,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A),
			types.Named("y", types.A),
		),
		types.Named("result", types.B).Description("true if `x` is greater than `y`; false otherwise"),
	),
	canSkipBctx: true,
}

var GreaterThanEq = &Builtin{
	Name:       "gte",
	Infix:      ">=",
	Categories: comparison,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A),
			types.Named("y", types.A),
		),
		types.Named("result", types.B).Description("true if `x` is greater or equal to `y`; false otherwise"),
	),
	canSkipBctx: true,
}

// LessThan represents the "<" comparison operator.
var LessThan = &Builtin{
	Name:       "lt",
	Infix:      "<",
	Categories: comparison,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A),
			types.Named("y", types.A),
		),
		types.Named("result", types.B).Description("true if `x` is less than `y`; false otherwise"),
	),
	canSkipBctx: true,
}

var LessThanEq = &Builtin{
	Name:       "lte",
	Infix:      "<=",
	Categories: comparison,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A),
			types.Named("y", types.A),
		),
		types.Named("result", types.B).Description("true if `x` is less than or equal to `y`; false otherwise"),
	),
	canSkipBctx: true,
}

var NotEqual = &Builtin{
	Name:       "neq",
	Infix:      "!=",
	Categories: comparison,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A),
			types.Named("y", types.A),
		),
		types.Named("result", types.B).Description("true if `x` is not equal to `y`; false otherwise"),
	),
	canSkipBctx: true,
}

// Equal represents the "==" comparison operator.
var Equal = &Builtin{
	Name:       "equal",
	Infix:      "==",
	Categories: comparison,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A),
			types.Named("y", types.A),
		),
		types.Named("result", types.B).Description("true if `x` is equal to `y`; false otherwise"),
	),
	canSkipBctx: true,
}

/**
 * Arithmetic
 */
var number = category("numbers")

var Plus = &Builtin{
	Name:        "plus",
	Infix:       "+",
	Description: "Plus adds two numbers together.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N),
			types.Named("y", types.N),
		),
		types.Named("z", types.N).Description("the sum of `x` and `y`"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Minus = &Builtin{
	Name:        "minus",
	Infix:       "-",
	Description: "Minus subtracts the second number from the first number or computes the difference between two sets.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewAny(types.N, types.SetOfAny)),
			types.Named("y", types.NewAny(types.N, types.SetOfAny)),
		),
		types.Named("z", types.NewAny(types.N, types.SetOfAny)).Description("the difference of `x` and `y`"),
	),
	Categories:  category("sets", "numbers"),
	canSkipBctx: true,
}

var Multiply = &Builtin{
	Name:        "mul",
	Infix:       "*",
	Description: "Multiplies two numbers.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N),
			types.Named("y", types.N),
		),
		types.Named("z", types.N).Description("the product of `x` and `y`"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Divide = &Builtin{
	Name:        "div",
	Infix:       "/",
	Description: "Divides the first number by the second number.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the dividend"),
			types.Named("y", types.N).Description("the divisor"),
		),
		types.Named("z", types.N).Description("the result of `x` divided by `y`"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Round = &Builtin{
	Name:        "round",
	Description: "Rounds the number to the nearest integer.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the number to round"),
		),
		types.Named("y", types.N).Description("the result of rounding `x`"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Ceil = &Builtin{
	Name:        "ceil",
	Description: "Rounds the number _up_ to the nearest integer.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the number to round"),
		),
		types.Named("y", types.N).Description("the result of rounding `x` _up_"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Floor = &Builtin{
	Name:        "floor",
	Description: "Rounds the number _down_ to the nearest integer.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the number to round"),
		),
		types.Named("y", types.N).Description("the result of rounding `x` _down_"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Abs = &Builtin{
	Name:        "abs",
	Description: "Returns the number without its sign.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the number to take the absolute value of"),
		),
		types.Named("y", types.N).Description("the absolute value of `x`"),
	),
	Categories:  number,
	canSkipBctx: true,
}

var Rem = &Builtin{
	Name:        "rem",
	Infix:       "%",
	Description: "Returns the remainder for of `x` divided by `y`, for `y != 0`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N),
			types.Named("y", types.N),
		),
		types.Named("z", types.N).Description("the remainder"),
	),
	Categories:  number,
	canSkipBctx: true,
}

/**
 * Bitwise
 */

var BitsOr = &Builtin{
	Name:        "bits.or",
	Description: "Returns the bitwise \"OR\" of two integers.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the first integer"),
			types.Named("y", types.N).Description("the second integer"),
		),
		types.Named("z", types.N).Description("the bitwise OR of `x` and `y`"),
	),
	canSkipBctx: true,
}

var BitsAnd = &Builtin{
	Name:        "bits.and",
	Description: "Returns the bitwise \"AND\" of two integers.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the first integer"),
			types.Named("y", types.N).Description("the second integer"),
		),
		types.Named("z", types.N).Description("the bitwise AND of `x` and `y`"),
	),
	canSkipBctx: true,
}

var BitsNegate = &Builtin{
	Name:        "bits.negate",
	Description: "Returns the bitwise negation (flip) of an integer.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the integer to negate"),
		),
		types.Named("z", types.N).Description("the bitwise negation of `x`"),
	),
	canSkipBctx: true,
}

var BitsXOr = &Builtin{
	Name:        "bits.xor",
	Description: "Returns the bitwise \"XOR\" (exclusive-or) of two integers.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the first integer"),
			types.Named("y", types.N).Description("the second integer"),
		),
		types.Named("z", types.N).Description("the bitwise XOR of `x` and `y`"),
	),
	canSkipBctx: true,
}

var BitsShiftLeft = &Builtin{
	Name:        "bits.lsh",
	Description: "Returns a new integer with its bits shifted `s` bits to the left.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the integer to shift"),
			types.Named("s", types.N).Description("the number of bits to shift"),
		),
		types.Named("z", types.N).Description("the result of shifting `x` `s` bits to the left"),
	),
	canSkipBctx: true,
}

var BitsShiftRight = &Builtin{
	Name:        "bits.rsh",
	Description: "Returns a new integer with its bits shifted `s` bits to the right.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.N).Description("the integer to shift"),
			types.Named("s", types.N).Description("the number of bits to shift"),
		),
		types.Named("z", types.N).Description("the result of shifting `x` `s` bits to the right"),
	),
	canSkipBctx: true,
}

/**
 * Sets
 */

var sets = category("sets")

var And = &Builtin{
	Name:        "and",
	Infix:       "&",
	Description: "Returns the intersection of two sets.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.SetOfAny).Description("the first set"),
			types.Named("y", types.SetOfAny).Description("the second set"),
		),
		types.Named("z", types.SetOfAny).Description("the intersection of `x` and `y`"),
	),
	Categories:  sets,
	canSkipBctx: true,
}

// Or performs a union operation on sets.
var Or = &Builtin{
	Name:        "or",
	Infix:       "|",
	Description: "Returns the union of two sets.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.SetOfAny),
			types.Named("y", types.SetOfAny),
		),
		types.Named("z", types.SetOfAny).Description("the union of `x` and `y`"),
	),
	Categories:  sets,
	canSkipBctx: true,
}

var Intersection = &Builtin{
	Name:        "intersection",
	Description: "Returns the intersection of the given input sets.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("xs", types.NewSet(types.SetOfAny)).Description("set of sets to intersect"),
		),
		types.Named("y", types.SetOfAny).Description("the intersection of all `xs` sets"),
	),
	Categories:  sets,
	canSkipBctx: true,
}

var Union = &Builtin{
	Name:        "union",
	Description: "Returns the union of the given input sets.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("xs", types.NewSet(types.SetOfAny)).Description("set of sets to merge"),
		),
		types.Named("y", types.SetOfAny).Description("the union of all `xs` sets"),
	),
	Categories:  sets,
	canSkipBctx: true,
}

/**
 * Aggregates
 */

var aggregates = category("aggregates")

var Count = &Builtin{
	Name:        "count",
	Description: " Count takes a collection or string and returns the number of elements (or characters) in it.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("collection", types.NewAny(
				types.SetOfAny,
				types.NewArray(nil, types.A),
				types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
				types.S,
			)).Description("the set/array/object/string to be counted"),
		),
		types.Named("n", types.N).Description("the count of elements, key/val pairs, or characters, respectively."),
	),
	Categories:  aggregates,
	canSkipBctx: true,
}

var Sum = &Builtin{
	Name:        "sum",
	Description: "Sums elements of an array or set of numbers.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("collection", types.NewAny(
				types.SetOfNum,
				types.NewArray(nil, types.N),
			)).Description("the set or array of numbers to sum"),
		),
		types.Named("n", types.N).Description("the sum of all elements"),
	),
	Categories:  aggregates,
	canSkipBctx: true,
}

var Product = &Builtin{
	Name:        "product",
	Description: "Multiplies elements of an array or set of numbers",
	Decl: types.NewFunction(
		types.Args(
			types.Named("collection", types.NewAny(
				types.SetOfNum,
				types.NewArray(nil, types.N),
			)).Description("the set or array of numbers to multiply"),
		),
		types.Named("n", types.N).Description("the product of all elements"),
	),
	Categories:  aggregates,
	canSkipBctx: true,
}

var Max = &Builtin{
	Name:        "max",
	Description: "Returns the maximum value in a collection.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("collection", types.NewAny(
				types.SetOfAny,
				types.NewArray(nil, types.A),
			)).Description("the set or array to be searched"),
		),
		types.Named("n", types.A).Description("the maximum of all elements"),
	),
	Categories:  aggregates,
	canSkipBctx: true,
}

var Min = &Builtin{
	Name:        "min",
	Description: "Returns the minimum value in a collection.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("collection", types.NewAny(
				types.SetOfAny,
				types.NewArray(nil, types.A),
			)).Description("the set or array to be searched"),
		),
		types.Named("n", types.A).Description("the minimum of all elements"),
	),
	Categories:  aggregates,
	canSkipBctx: true,
}

/**
 * Sorting
 */

var Sort = &Builtin{
	Name:        "sort",
	Description: "Returns a sorted array.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("collection", types.NewAny(
				types.NewArray(nil, types.A),
				types.SetOfAny,
			)).Description("the array or set to be sorted"),
		),
		types.Named("n", types.NewArray(nil, types.A)).Description("the sorted array"),
	),
	Categories:  aggregates,
	canSkipBctx: true,
}

/**
 * Arrays
 */

var ArrayConcat = &Builtin{
	Name:        "array.concat",
	Description: "Concatenates two arrays.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewArray(nil, types.A)).Description("the first array"),
			types.Named("y", types.NewArray(nil, types.A)).Description("the second array"),
		),
		types.Named("z", types.NewArray(nil, types.A)).Description("the concatenation of `x` and `y`"),
	),
	canSkipBctx: true,
}

var ArraySlice = &Builtin{
	Name:        "array.slice",
	Description: "Returns a slice of a given array. If `start` is greater or equal than `stop`, `slice` is `[]`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("arr", types.NewArray(nil, types.A)).Description("the array to be sliced"),
			types.Named("start", types.N).Description("the start index of the returned slice; if less than zero, it's clamped to 0"),
			types.Named("stop", types.N).Description("the stop index of the returned slice; if larger than `count(arr)`, it's clamped to `count(arr)`"),
		),
		types.Named("slice", types.NewArray(nil, types.A)).Description("the subslice of `array`, from `start` to `end`, including `arr[start]`, but excluding `arr[end]`"),
	),
	canSkipBctx: true,
} // NOTE(sr): this function really needs examples

var ArrayReverse = &Builtin{
	Name:        "array.reverse",
	Description: "Returns the reverse of a given array.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("arr", types.NewArray(nil, types.A)).Description("the array to be reversed"),
		),
		types.Named("rev", types.NewArray(nil, types.A)).Description("an array containing the elements of `arr` in reverse order"),
	),
	canSkipBctx: true,
}

/**
 * Conversions
 */
var conversions = category("conversions")

var ToNumber = &Builtin{
	Name:        "to_number",
	Description: "Converts a string, bool, or number value to a number: Strings are converted to numbers using `strconv.Atoi`, Boolean `false` is converted to 0 and `true` is converted to 1.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewAny(
				types.N,
				types.S,
				types.B,
				types.Nl,
			)).Description("value to convert"),
		),
		types.Named("num", types.N).Description("the numeric representation of `x`"),
	),
	Categories:  conversions,
	canSkipBctx: true,
}

/**
 * Regular Expressions
 */

var RegexMatch = &Builtin{
	Name:        "regex.match",
	Description: "Matches a string against a regular expression.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("regular expression"),
			types.Named("value", types.S).Description("value to match against `pattern`"),
		),
		types.Named("result", types.B).Description("true if `value` matches `pattern`"),
	),
}

var RegexIsValid = &Builtin{
	Name:        "regex.is_valid",
	Description: "Checks if a string is a valid regular expression: the detailed syntax for patterns is defined by https://github.com/google/re2/wiki/Syntax.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("regular expression"),
		),
		types.Named("result", types.B).Description("true if `pattern` is a valid regular expression"),
	),
	canSkipBctx: true,
}

var RegexFindAllStringSubmatch = &Builtin{
	Name:        "regex.find_all_string_submatch_n",
	Description: "Returns all successive matches of the expression.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("regular expression"),
			types.Named("value", types.S).Description("string to match"),
			types.Named("number", types.N).Description("number of matches to return; `-1` means all matches"),
		),
		types.Named("output", types.NewArray(nil, types.NewArray(nil, types.S))).Description("array of all matches"),
	),
	canSkipBctx: false,
}

var RegexTemplateMatch = &Builtin{
	Name:        "regex.template_match",
	Description: "Matches a string against a pattern, where there pattern may be glob-like",
	Decl: types.NewFunction(
		types.Args(
			types.Named("template", types.S).Description("template expression containing `0..n` regular expressions"),
			types.Named("value", types.S).Description("string to match"),
			types.Named("delimiter_start", types.S).Description("start delimiter of the regular expression in `template`"),
			types.Named("delimiter_end", types.S).Description("end delimiter of the regular expression in `template`"),
		),
		types.Named("result", types.B).Description("true if `value` matches the `template`"),
	),
	canSkipBctx: true,
} // TODO(sr): example:`regex.template_match("urn:foo:{.*}", "urn:foo:bar:baz", "{", "}")`` returns ``true``.

var RegexSplit = &Builtin{
	Name:        "regex.split",
	Description: "Splits the input string by the occurrences of the given pattern.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("regular expression"),
			types.Named("value", types.S).Description("string to match"),
		),
		types.Named("output", types.NewArray(nil, types.S)).Description("the parts obtained by splitting `value`"),
	),
	canSkipBctx: false,
}

// RegexFind takes two strings and a number, the pattern, the value and number of match values to
// return, -1 means all match values.
var RegexFind = &Builtin{
	Name:        "regex.find_n",
	Description: "Returns the specified number of matches when matching the input against the pattern.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("regular expression"),
			types.Named("value", types.S).Description("string to match"),
			types.Named("number", types.N).Description("number of matches to return, if `-1`, returns all matches"),
		),
		types.Named("output", types.NewArray(nil, types.S)).Description("collected matches"),
	),
	canSkipBctx: false,
}

// GlobsMatch takes two strings regexp-style strings and evaluates to true if their
// intersection matches a non-empty set of non-empty strings.
// Examples:
//   - "a.a." and ".b.b" -> true.
//   - "[a-z]*" and [0-9]+" -> not true.
var GlobsMatch = &Builtin{
	Name: "regex.globs_match",
	Description: `Checks if the intersection of two glob-style regular expressions matches a non-empty set of non-empty strings.
The set of regex symbols is limited for this builtin: only ` + "`.`, `*`, `+`, `[`, `-`, `]` and `\\` are treated as special symbols.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("glob1", types.S).Description("first glob-style regular expression"),
			types.Named("glob2", types.S).Description("second glob-style regular expression"),
		),
		types.Named("result", types.B).Description("true if the intersection of `glob1` and `glob2` matches a non-empty set of non-empty strings"),
	),
	canSkipBctx: true,
}

/**
 * Strings
 */
var stringsCat = category("strings")

var AnyPrefixMatch = &Builtin{
	Name:        "strings.any_prefix_match",
	Description: "Returns true if any of the search strings begins with any of the base strings.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("search", types.NewAny(
				types.S,
				types.SetOfStr,
				types.NewArray(nil, types.S),
			)).Description("search string(s)"),
			types.Named("base", types.NewAny(
				types.S,
				types.SetOfStr,
				types.NewArray(nil, types.S),
			)).Description("base string(s)"),
		),
		types.Named("result", types.B).Description("result of the prefix check"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var AnySuffixMatch = &Builtin{
	Name:        "strings.any_suffix_match",
	Description: "Returns true if any of the search strings ends with any of the base strings.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("search", types.NewAny(
				types.S,
				types.SetOfStr,
				types.NewArray(nil, types.S),
			)).Description("search string(s)"),
			types.Named("base", types.NewAny(
				types.S,
				types.SetOfStr,
				types.NewArray(nil, types.S),
			)).Description("base string(s)"),
		),
		types.Named("result", types.B).Description("result of the suffix check"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Concat = &Builtin{
	Name:        "concat",
	Description: "Joins a set or array of strings with a delimiter.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("delimiter", types.S).Description("string to use as a delimiter"),
			types.Named("collection", types.NewAny(
				types.SetOfStr,
				types.NewArray(nil, types.S),
			)).Description("strings to join"),
		),
		types.Named("output", types.S).Description("the joined string"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var FormatInt = &Builtin{
	Name:        "format_int",
	Description: "Returns the string representation of the number in the given base after rounding it down to an integer value.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("number", types.N).Description("number to format"),
			types.Named("base", types.N).Description("base of number representation to use"),
		),
		types.Named("output", types.S).Description("formatted number"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var IndexOf = &Builtin{
	Name:        "indexof",
	Description: "Returns the index of a substring contained inside a string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("haystack", types.S).Description("string to search in"),
			types.Named("needle", types.S).Description("substring to look for"),
		),
		types.Named("output", types.N).Description("index of first occurrence, `-1` if not found"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var IndexOfN = &Builtin{
	Name:        "indexof_n",
	Description: "Returns a list of all the indexes of a substring contained inside a string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("haystack", types.S).Description("string to search in"),
			types.Named("needle", types.S).Description("substring to look for"),
		),
		types.Named("output", types.NewArray(nil, types.N)).Description("all indices at which `needle` occurs in `haystack`, may be empty"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Substring = &Builtin{
	Name:        "substring",
	Description: "Returns the  portion of a string for a given `offset` and a `length`.  If `length < 0`, `output` is the remainder of the string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to extract substring from"),
			types.Named("offset", types.N).Description("offset, must be positive"),
			types.Named("length", types.N).Description("length of the substring starting from `offset`"),
		),
		types.Named("output", types.S).Description("substring of `value` from `offset`, of length `length`"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Contains = &Builtin{
	Name:        "contains",
	Description: "Returns `true` if the search string is included in the base string",
	Decl: types.NewFunction(
		types.Args(
			types.Named("haystack", types.S).Description("string to search in"),
			types.Named("needle", types.S).Description("substring to look for"),
		),
		types.Named("result", types.B).Description("result of the containment check"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var StringCount = &Builtin{
	Name:        "strings.count",
	Description: "Returns the number of non-overlapping instances of a substring in a string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("search", types.S).Description("string to search in"),
			types.Named("substring", types.S).Description("substring to look for"),
		),
		types.Named("output", types.N).Description("count of occurrences, `0` if not found"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var StartsWith = &Builtin{
	Name:        "startswith",
	Description: "Returns true if the search string begins with the base string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("search", types.S).Description("search string"),
			types.Named("base", types.S).Description("base string"),
		),
		types.Named("result", types.B).Description("result of the prefix check"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var EndsWith = &Builtin{
	Name:        "endswith",
	Description: "Returns true if the search string ends with the base string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("search", types.S).Description("search string"),
			types.Named("base", types.S).Description("base string"),
		),
		types.Named("result", types.B).Description("result of the suffix check"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Lower = &Builtin{
	Name:        "lower",
	Description: "Returns the input string but with all characters in lower-case.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string that is converted to lower-case"),
		),
		types.Named("y", types.S).Description("lower-case of x"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Upper = &Builtin{
	Name:        "upper",
	Description: "Returns the input string but with all characters in upper-case.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string that is converted to upper-case"),
		),
		types.Named("y", types.S).Description("upper-case of x"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Split = &Builtin{
	Name:        "split",
	Description: "Split returns an array containing elements of the input string split on a delimiter.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string that is split"),
			types.Named("delimiter", types.S).Description("delimiter used for splitting"),
		),
		types.Named("ys", types.NewArray(nil, types.S)).Description("split parts"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Replace = &Builtin{
	Name:        "replace",
	Description: "Replace replaces all instances of a sub-string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string being processed"),
			types.Named("old", types.S).Description("substring to replace"),
			types.Named("new", types.S).Description("string to replace `old` with"),
		),
		types.Named("y", types.S).Description("string with replaced substrings"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var ReplaceN = &Builtin{
	Name: "strings.replace_n",
	Description: `Replaces a string from a list of old, new string pairs.
Replacements are performed in the order they appear in the target string, without overlapping matches.
The old string comparisons are done in argument order.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("patterns", types.NewObject(
				nil,
				types.NewDynamicProperty(
					types.S,
					types.S)),
			).Description("replacement pairs"),
			types.Named("value", types.S).Description("string to replace substring matches in"),
		),
		types.Named("output", types.S).Description("string with replaced substrings"),
	),
	canSkipBctx: true,
}

var RegexReplace = &Builtin{
	Name:        "regex.replace",
	Description: `Find and replaces the text using the regular expression pattern.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("s", types.S).Description("string being processed"),
			types.Named("pattern", types.S).Description("regex pattern to be applied"),
			types.Named("value", types.S).Description("regex value"),
		),
		types.Named("output", types.S).Description("string with replaced substrings"),
	),
	canSkipBctx: false,
}

var Trim = &Builtin{
	Name:        "trim",
	Description: "Returns `value` with all leading or trailing instances of the `cutset` characters removed.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to trim"),
			types.Named("cutset", types.S).Description("string of characters that are cut off"),
		),
		types.Named("output", types.S).Description("string trimmed of `cutset` characters"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var TrimLeft = &Builtin{
	Name:        "trim_left",
	Description: "Returns `value` with all leading instances of the `cutset` characters removed.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to trim"),
			types.Named("cutset", types.S).Description("string of characters that are cut off on the left"),
		),
		types.Named("output", types.S).Description("string left-trimmed of `cutset` characters"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var TrimPrefix = &Builtin{
	Name:        "trim_prefix",
	Description: "Returns `value` without the prefix. If `value` doesn't start with `prefix`, it is returned unchanged.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to trim"),
			types.Named("prefix", types.S).Description("prefix to cut off"),
		),
		types.Named("output", types.S).Description("string with `prefix` cut off"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var TrimRight = &Builtin{
	Name:        "trim_right",
	Description: "Returns `value` with all trailing instances of the `cutset` characters removed.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to trim"),
			types.Named("cutset", types.S).Description("string of characters that are cut off on the right"),
		),
		types.Named("output", types.S).Description("string right-trimmed of `cutset` characters"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var TrimSuffix = &Builtin{
	Name:        "trim_suffix",
	Description: "Returns `value` without the suffix. If `value` doesn't end with `suffix`, it is returned unchanged.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to trim"),
			types.Named("suffix", types.S).Description("suffix to cut off"),
		),
		types.Named("output", types.S).Description("string with `suffix` cut off"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var TrimSpace = &Builtin{
	Name:        "trim_space",
	Description: "Return the given string with all leading and trailing white space removed.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("string to trim"),
		),
		types.Named("output", types.S).Description("string leading and trailing white space cut off"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var Sprintf = &Builtin{
	Name:        "sprintf",
	Description: "Returns the given string, formatted.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("format", types.S).Description("string with formatting verbs"),
			types.Named("values", types.NewArray(nil, types.A)).Description("arguments to format into formatting verbs"),
		),
		types.Named("output", types.S).Description("`format` formatted by the values in `values`"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var StringReverse = &Builtin{
	Name:        "strings.reverse",
	Description: "Reverses a given string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to reverse"),
		),
		types.Named("y", types.S).Description("reversed string"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

var RenderTemplate = &Builtin{
	Name: "strings.render_template",
	Description: `Renders a templated string with given template variables injected. For a given templated string and key/value mapping, values will be injected into the template where they are referenced by key.
	For examples of templating syntax, see https://pkg.go.dev/text/template`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("a templated string"),
			types.Named("vars", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("a mapping of template variable keys to values"),
		),
		types.Named("result", types.S).Description("rendered template with template variables injected"),
	),
	Categories:  stringsCat,
	canSkipBctx: true,
}

/**
 * Numbers
 */

// RandIntn returns a random number 0 - n
// Marked non-deterministic because it relies on RNG internally.
var RandIntn = &Builtin{
	Name:        "rand.intn",
	Description: "Returns a random integer between `0` and `n` (`n` exclusive). If `n` is `0`, then `y` is always `0`. For any given argument pair (`str`, `n`), the output will be consistent throughout a query evaluation.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("str", types.S).Description("seed string for the random number"),
			types.Named("n", types.N).Description("upper bound of the random number (exclusive)"),
		),
		types.Named("y", types.N).Description("random integer in the range `[0, abs(n))`"),
	),
	Categories:       number,
	Nondeterministic: true,
	canSkipBctx:      false,
}

var NumbersRange = &Builtin{
	Name:        "numbers.range",
	Description: "Returns an array of numbers in the given (inclusive) range. If `a==b`, then `range == [a]`; if `a > b`, then `range` is in descending order.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("a", types.N).Description("the start of the range"),
			types.Named("b", types.N).Description("the end of the range (inclusive)"),
		),
		types.Named("range", types.NewArray(nil, types.N)).Description("the range between `a` and `b`"),
	),
	canSkipBctx: false, // needed for context timeout check
}

var NumbersRangeStep = &Builtin{
	Name: "numbers.range_step",
	Description: `Returns an array of numbers in the given (inclusive) range incremented by a positive step.
	If "a==b", then "range == [a]"; if "a > b", then "range" is in descending order.
	If the provided "step" is less then 1, an error will be thrown.
	If "b" is not in the range of the provided "step", "b" won't be included in the result.
	`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("a", types.N).Description("the start of the range"),
			types.Named("b", types.N).Description("the end of the range (inclusive)"),
			types.Named("step", types.N).Description("the step between numbers in the range"),
		),
		types.Named("range", types.NewArray(nil, types.N)).Description("the range between `a` and `b` in `step` increments"),
	),
	canSkipBctx: false, // needed for context timeout check
}

/**
 * Units
 */

var UnitsParse = &Builtin{
	Name: "units.parse",
	Description: `Converts strings like "10G", "5K", "4M", "1500m", and the like into a number.
This number can be a non-integer, such as 1.5, 0.22, etc. Scientific notation is supported,
allowing values such as "1e-3K" (1) or "2.5e6M" (2.5 million M).

Supports standard metric decimal and binary SI units (e.g., K, Ki, M, Mi, G, Gi, etc.) where
m, K, M, G, T, P, and E are treated as decimal units and Ki, Mi, Gi, Ti, Pi, and Ei are treated as
binary units.

Note that 'm' and 'M' are case-sensitive to allow distinguishing between "milli" and "mega" units
respectively. Other units are case-insensitive.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("the unit to parse"),
		),
		types.Named("y", types.N).Description("the parsed number"),
	),
	canSkipBctx: true,
}

var UnitsParseBytes = &Builtin{
	Name: "units.parse_bytes",
	Description: `Converts strings like "10GB", "5K", "4mb", or "1e6KB" into an integer number of bytes.

Supports standard byte units (e.g., KB, KiB, etc.) where KB, MB, GB, and TB are treated as decimal
units, and KiB, MiB, GiB, and TiB are treated as binary units. Scientific notation is supported,
enabling values like "1.5e3MB" (1500MB) or "2e6GiB" (2 million GiB).

The bytes symbol (b/B) in the unit is optional; omitting it will yield the same result (e.g., "Mi"
and "MiB" are equivalent).`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("the byte unit to parse"),
		),
		types.Named("y", types.N).Description("the parsed number"),
	),
	canSkipBctx: true,
}

//
/**
 * Type
 */

// UUIDRFC4122 returns a version 4 UUID string.
// Marked non-deterministic because it relies on RNG internally.
var UUIDRFC4122 = &Builtin{
	Name:        "uuid.rfc4122",
	Description: "Returns a new UUIDv4.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("k", types.S).Description("seed string"),
		),
		types.Named("output", types.S).Description("a version 4 UUID; for any given `k`, the output will be consistent throughout a query evaluation"),
	),
	Nondeterministic: true,
	canSkipBctx:      false,
}

var UUIDParse = &Builtin{
	Name:        "uuid.parse",
	Description: "Parses the string value as an UUID and returns an object with the well-defined fields of the UUID if valid.",
	Categories:  nil,
	Decl: types.NewFunction(
		types.Args(
			types.Named("uuid", types.S).Description("UUID string to parse"),
		),
		types.Named("result", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("Properties of UUID if valid (version, variant, etc). Undefined otherwise."),
	),
	Relation:    false,
	canSkipBctx: true,
}

/**
 * JSON
 */

var objectCat = category("object")

var JSONFilter = &Builtin{
	Name: "json.filter",
	Description: "Filters the object. " +
		"For example: `json.filter({\"a\": {\"b\": \"x\", \"c\": \"y\"}}, [\"a/b\"])` will result in `{\"a\": {\"b\": \"x\"}}`). " +
		"Paths are not filtered in-order and are deduplicated before being evaluated.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			)).Description("object to filter"),
			types.Named("paths", types.NewAny(
				types.NewArray(
					nil,
					types.NewAny(
						types.S,
						types.NewArray(
							nil,
							types.A,
						),
					),
				),
				types.NewSet(
					types.NewAny(
						types.S,
						types.NewArray(
							nil,
							types.A,
						),
					),
				),
			)).Description("JSON string paths"),
		),
		types.Named("filtered", types.A).Description("remaining data from `object` with only keys specified in `paths`"),
	),
	Categories:  objectCat,
	canSkipBctx: true,
}

var JSONRemove = &Builtin{
	Name: "json.remove",
	Description: "Removes paths from an object. " +
		"For example: `json.remove({\"a\": {\"b\": \"x\", \"c\": \"y\"}}, [\"a/b\"])` will result in `{\"a\": {\"c\": \"y\"}}`. " +
		"Paths are not removed in-order and are deduplicated before being evaluated.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			)).Description("object to remove paths from"),
			types.Named("paths", types.NewAny(
				types.NewArray(
					nil,
					types.NewAny(
						types.S,
						types.NewArray(
							nil,
							types.A,
						),
					),
				),
				types.NewSet(
					types.NewAny(
						types.S,
						types.NewArray(
							nil,
							types.A,
						),
					),
				),
			)).Description("JSON string paths"),
		),
		types.Named("output", types.A).Description("result of removing all keys specified in `paths`"),
	),
	Categories:  objectCat,
	canSkipBctx: true,
}

var JSONPatch = &Builtin{
	Name: "json.patch",
	Description: "Patches an object according to RFC6902. " +
		"For example: `json.patch({\"a\": {\"foo\": 1}}, [{\"op\": \"add\", \"path\": \"/a/bar\", \"value\": 2}])` results in `{\"a\": {\"foo\": 1, \"bar\": 2}`. " +
		"The patches are applied atomically: if any of them fails, the result will be undefined. " +
		"Additionally works on sets, where a value contained in the set is considered to be its path.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.A).Description("the object to patch"), // TODO(sr): types.A?
			types.Named("patches", types.NewArray(
				nil,
				types.NewObject(
					[]*types.StaticProperty{
						{Key: "op", Value: types.S},
						{Key: "path", Value: types.A},
					},
					types.NewDynamicProperty(types.A, types.A),
				),
			)).Description("the JSON patches to apply"),
		),
		types.Named("output", types.A).Description("result obtained after consecutively applying all patch operations in `patches`"),
	),
	Categories:  objectCat,
	canSkipBctx: true,
}

var ObjectSubset = &Builtin{
	Name: "object.subset",
	Description: "Determines if an object `sub` is a subset of another object `super`." +
		"Object `sub` is a subset of object `super` if and only if every key in `sub` is also in `super`, " +
		"**and** for all keys which `sub` and `super` share, they have the same value. " +
		"This function works with objects, sets, arrays and a set of array and set." +
		"If both arguments are objects, then the operation is recursive, e.g. " +
		"`{\"c\": {\"x\": {10, 15, 20}}` is a subset of `{\"a\": \"b\", \"c\": {\"x\": {10, 15, 20, 25}, \"y\": \"z\"}`. " +
		"If both arguments are sets, then this function checks if every element of `sub` is a member of `super`, " +
		"but does not attempt to recurse. If both arguments are arrays, " +
		"then this function checks if `sub` appears contiguously in order within `super`, " +
		"and also does not attempt to recurse. If `super` is array and `sub` is set, " +
		"then this function checks if `super` contains every element of `sub` with no consideration of ordering, " +
		"and also does not attempt to recurse.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("super", types.NewAny(types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			), types.SetOfAny,
				types.NewArray(nil, types.A),
			)).Description("object to test if sub is a subset of"),
			types.Named("sub", types.NewAny(types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			), types.SetOfAny,
				types.NewArray(nil, types.A),
			)).Description("object to test if super is a superset of"),
		),
		types.Named("result", types.A).Description("`true` if `sub` is a subset of `super`"),
	),
	canSkipBctx: true,
}

var ObjectUnion = &Builtin{
	Name: "object.union",
	Description: "Creates a new object of the asymmetric union of two objects. " +
		"For example: `object.union({\"a\": 1, \"b\": 2, \"c\": {\"d\": 3}}, {\"a\": 7, \"c\": {\"d\": 4, \"e\": 5}})` will result in `{\"a\": 7, \"b\": 2, \"c\": {\"d\": 4, \"e\": 5}}`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("a", types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			)).Description("left-hand object"),
			types.Named("b", types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			)).Description("right-hand object"),
		),
		types.Named("output", types.A).Description("a new object which is the result of an asymmetric recursive union of two objects where conflicts are resolved by choosing the key from the right-hand object `b`"),
	), // TODO(sr): types.A?  ^^^^^^^ (also below)
	canSkipBctx: true,
}

var ObjectUnionN = &Builtin{
	Name: "object.union_n",
	Description: "Creates a new object that is the asymmetric union of all objects merged from left to right. " +
		"For example: `object.union_n([{\"a\": 1}, {\"b\": 2}, {\"a\": 3}])` will result in `{\"b\": 2, \"a\": 3}`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("objects", types.NewArray(
				nil,
				types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			)).Description("list of objects to merge"),
		),
		types.Named("output", types.A).Description("asymmetric recursive union of all objects in `objects`, merged from left to right, where conflicts are resolved by choosing the key from the right-hand object"),
	),
	canSkipBctx: true,
}

var ObjectRemove = &Builtin{
	Name:        "object.remove",
	Description: "Removes specified keys from an object.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			)).Description("object to remove keys from"),
			types.Named("keys", types.NewAny(
				types.NewArray(nil, types.A),
				types.SetOfAny,
				types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			)).Description("keys to remove from x"),
		),
		types.Named("output", types.A).Description("result of removing the specified `keys` from `object`"),
	),
	canSkipBctx: true,
}

var ObjectFilter = &Builtin{
	Name: "object.filter",
	Description: "Filters the object by keeping only specified keys. " +
		"For example: `object.filter({\"a\": {\"b\": \"x\", \"c\": \"y\"}, \"d\": \"z\"}, [\"a\"])` will result in `{\"a\": {\"b\": \"x\", \"c\": \"y\"}}`).",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(
				nil,
				types.NewDynamicProperty(types.A, types.A),
			)).Description("object to filter keys"),
			types.Named("keys", types.NewAny(
				types.NewArray(nil, types.A),
				types.SetOfAny,
				types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			)).Description("keys to keep in `object`"),
		),
		types.Named("filtered", types.A).Description("remaining data from `object` with only keys specified in `keys`"),
	),
	canSkipBctx: true,
}

var ObjectGet = &Builtin{
	Name: "object.get",
	Description: "Returns value of an object's key if present, otherwise a default. " +
		"If the supplied `key` is an `array`, then `object.get` will search through a nested object or array using each key in turn. " +
		"For example: `object.get({\"a\": [{ \"b\": true }]}, [\"a\", 0, \"b\"], false)` results in `true`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))).Description("object to get `key` from"),
			types.Named("key", types.A).Description("key to lookup in `object`"),
			types.Named("default", types.A).Description("default to use when lookup fails"),
		),
		types.Named("value", types.A).Description("`object[key]` if present, otherwise `default`"),
	),
	canSkipBctx: true,
}

var ObjectKeys = &Builtin{
	Name: "object.keys",
	Description: "Returns a set of an object's keys. " +
		"For example: `object.keys({\"a\": 1, \"b\": true, \"c\": \"d\")` results in `{\"a\", \"b\", \"c\"}`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))).Description("object to get keys from"),
		),
		types.Named("value", types.SetOfAny).Description("set of `object`'s keys"),
	),
	canSkipBctx: true,
}

/*
 *  Encoding
 */
var encoding = category("encoding")

var JSONMarshal = &Builtin{
	Name:        "json.marshal",
	Description: "Serializes the input term to JSON.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("the term to serialize"),
		),
		types.Named("y", types.S).Description("the JSON string representation of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var JSONMarshalWithOptions = &Builtin{
	Name: "json.marshal_with_options",
	Description: "Serializes the input term JSON, with additional formatting options via the `opts` parameter. " +
		"`opts` accepts keys `pretty` (enable multi-line/formatted JSON), `prefix` (string to prefix lines with, default empty string) and `indent` (string to indent with, default `\\t`).",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("the term to serialize"),
			types.Named("opts", types.NewObject(
				[]*types.StaticProperty{
					types.NewStaticProperty("pretty", types.B),
					types.NewStaticProperty("indent", types.S),
					types.NewStaticProperty("prefix", types.S),
				},
				types.NewDynamicProperty(types.S, types.A),
			)).Description("encoding options"),
		),
		types.Named("y", types.S).Description("the JSON string representation of `x`, with configured prefix/indent string(s) as appropriate"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var JSONUnmarshal = &Builtin{
	Name:        "json.unmarshal",
	Description: "Deserializes the input string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("a JSON string"),
		),
		types.Named("y", types.A).Description("the term deserialized from `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var JSONIsValid = &Builtin{
	Name:        "json.is_valid",
	Description: "Verifies the input string is a valid JSON document.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("a JSON string"),
		),
		types.Named("result", types.B).Description("`true` if `x` is valid JSON, `false` otherwise"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var Base64Encode = &Builtin{
	Name:        "base64.encode",
	Description: "Serializes the input string into base64 encoding.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to encode"),
		),
		types.Named("y", types.S).Description("base64 serialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var Base64Decode = &Builtin{
	Name:        "base64.decode",
	Description: "Deserializes the base64 encoded input string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to decode"),
		),
		types.Named("y", types.S).Description("base64 deserialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var Base64IsValid = &Builtin{
	Name:        "base64.is_valid",
	Description: "Verifies the input string is base64 encoded.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to check"),
		),
		types.Named("result", types.B).Description("`true` if `x` is valid base64 encoded value, `false` otherwise"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var Base64UrlEncode = &Builtin{
	Name:        "base64url.encode",
	Description: "Serializes the input string into base64url encoding.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to encode"),
		),
		types.Named("y", types.S).Description("base64url serialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var Base64UrlEncodeNoPad = &Builtin{
	Name:        "base64url.encode_no_pad",
	Description: "Serializes the input string into base64url encoding without padding.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to encode"),
		),
		types.Named("y", types.S).Description("base64url serialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var Base64UrlDecode = &Builtin{
	Name:        "base64url.decode",
	Description: "Deserializes the base64url encoded input string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to decode"),
		),
		types.Named("y", types.S).Description("base64url deserialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var URLQueryDecode = &Builtin{
	Name:        "urlquery.decode",
	Description: "Decodes a URL-encoded input string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("the URL-encoded string"),
		),
		types.Named("y", types.S).Description("URL-encoding deserialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var URLQueryEncode = &Builtin{
	Name:        "urlquery.encode",
	Description: "Encodes the input string into a URL-encoded string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("the string to encode"),
		),
		types.Named("y", types.S).Description("URL-encoding serialization of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var URLQueryEncodeObject = &Builtin{
	Name:        "urlquery.encode_object",
	Description: "Encodes the given object into a URL encoded query string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("object", types.NewObject(
				nil,
				types.NewDynamicProperty(
					types.S,
					types.NewAny(
						types.S,
						types.NewArray(nil, types.S),
						types.SetOfStr,
					),
				),
			),
			).Description("the object to encode"),
		),
		types.Named("y", types.S).Description("the URL-encoded serialization of `object`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var URLQueryDecodeObject = &Builtin{
	Name:        "urlquery.decode_object",
	Description: "Decodes the given URL query string into an object.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("the query string"),
		),
		types.Named("object", types.NewObject(nil, types.NewDynamicProperty(
			types.S,
			types.NewArray(nil, types.S)))).Description("the resulting object"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var YAMLMarshal = &Builtin{
	Name:        "yaml.marshal",
	Description: "Serializes the input term to YAML.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("the term to serialize"),
		),
		types.Named("y", types.S).Description("the YAML string representation of `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var YAMLUnmarshal = &Builtin{
	Name:        "yaml.unmarshal",
	Description: "Deserializes the input string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("a YAML string"),
		),
		types.Named("y", types.A).Description("the term deserialized from `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

// YAMLIsValid verifies the input string is a valid YAML document.
var YAMLIsValid = &Builtin{
	Name:        "yaml.is_valid",
	Description: "Verifies the input string is a valid YAML document.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("a YAML string"),
		),
		types.Named("result", types.B).Description("`true` if `x` is valid YAML, `false` otherwise"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var HexEncode = &Builtin{
	Name:        "hex.encode",
	Description: "Serializes the input string using hex-encoding.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("string to encode"),
		),
		types.Named("y", types.S).Description("serialization of `x` using hex-encoding"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

var HexDecode = &Builtin{
	Name:        "hex.decode",
	Description: "Deserializes the hex-encoded input string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("a hex-encoded string"),
		),
		types.Named("y", types.S).Description("deserialized from `x`"),
	),
	Categories:  encoding,
	canSkipBctx: true,
}

/**
 * Tokens
 */
var tokensCat = category("tokens")

var JWTDecode = &Builtin{
	Name:        "io.jwt.decode",
	Description: "Decodes a JSON Web Token and outputs it as an object.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token to decode"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			types.S,
		}, nil)).Description("`[header, payload, sig]`, where `header` and `payload` are objects; `sig` is the hexadecimal representation of the signature on the token."),
	),
	Categories:  tokensCat,
	canSkipBctx: true,
}

var JWTVerifyRS256 = &Builtin{
	Name:        "io.jwt.verify_rs256",
	Description: "Verifies if a RS256 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyRS384 = &Builtin{
	Name:        "io.jwt.verify_rs384",
	Description: "Verifies if a RS384 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyRS512 = &Builtin{
	Name:        "io.jwt.verify_rs512",
	Description: "Verifies if a RS512 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyPS256 = &Builtin{
	Name:        "io.jwt.verify_ps256",
	Description: "Verifies if a PS256 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyPS384 = &Builtin{
	Name:        "io.jwt.verify_ps384",
	Description: "Verifies if a PS384 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyPS512 = &Builtin{
	Name:        "io.jwt.verify_ps512",
	Description: "Verifies if a PS512 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyES256 = &Builtin{
	Name:        "io.jwt.verify_es256",
	Description: "Verifies if a ES256 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyES384 = &Builtin{
	Name:        "io.jwt.verify_es384",
	Description: "Verifies if a ES384 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyES512 = &Builtin{
	Name:        "io.jwt.verify_es512",
	Description: "Verifies if a ES512 JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("certificate", types.S).Description("PEM encoded certificate, PEM encoded public key, or the JWK key (set) used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyHS256 = &Builtin{
	Name:        "io.jwt.verify_hs256",
	Description: "Verifies if a HS256 (secret) JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("secret", types.S).Description("plain text secret used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyHS384 = &Builtin{
	Name:        "io.jwt.verify_hs384",
	Description: "Verifies if a HS384 (secret) JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("secret", types.S).Description("plain text secret used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

var JWTVerifyHS512 = &Builtin{
	Name:        "io.jwt.verify_hs512",
	Description: "Verifies if a HS512 (secret) JWT signature is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified"),
			types.Named("secret", types.S).Description("plain text secret used to verify the signature"),
		),
		types.Named("result", types.B).Description("`true` if the signature is valid, `false` otherwise"),
	),
	Categories:  tokensCat,
	canSkipBctx: false,
}

// Marked non-deterministic because it relies on time internally.
var JWTDecodeVerify = &Builtin{
	Name: "io.jwt.decode_verify",
	Description: `Verifies a JWT signature under parameterized constraints and decodes the claims if it is valid.
Supports the following algorithms: HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384 and PS512.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("jwt", types.S).Description("JWT token whose signature is to be verified and whose claims are to be checked"),
			types.Named("constraints", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("claim verification constraints"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.B,
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
		}, nil)).Description("`[valid, header, payload]`:  if the input token is verified and meets the requirements of `constraints` then `valid` is `true`; `header` and `payload` are objects containing the JOSE header and the JWT claim set; otherwise, `valid` is `false`, `header` and `payload` are `{}`"),
	),
	Categories:       tokensCat,
	Nondeterministic: true,
	canSkipBctx:      false,
}

var tokenSign = category("tokensign")

// Marked non-deterministic because it relies on RNG internally.
var JWTEncodeSignRaw = &Builtin{
	Name:        "io.jwt.encode_sign_raw",
	Description: "Encodes and optionally signs a JSON Web Token.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("headers", types.S).Description("JWS Protected Header"),
			types.Named("payload", types.S).Description("JWS Payload"),
			types.Named("key", types.S).Description("JSON Web Key (RFC7517)"),
		),
		types.Named("output", types.S).Description("signed JWT"),
	),
	Categories:       tokenSign,
	Nondeterministic: true,
	canSkipBctx:      false,
}

// Marked non-deterministic because it relies on RNG internally.
var JWTEncodeSign = &Builtin{
	Name:        "io.jwt.encode_sign",
	Description: "Encodes and optionally signs a JSON Web Token. Inputs are taken as objects, not encoded strings (see `io.jwt.encode_sign_raw`).",
	Decl: types.NewFunction(
		types.Args(
			types.Named("headers", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("JWS Protected Header"),
			types.Named("payload", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("JWS Payload"),
			types.Named("key", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("JSON Web Key (RFC7517)"),
		),
		types.Named("output", types.S).Description("signed JWT"),
	),
	Categories:       tokenSign,
	Nondeterministic: true,
	canSkipBctx:      false,
}

/**
 * Time
 */

// Marked non-deterministic because it relies on time directly.
var NowNanos = &Builtin{
	Name:        "time.now_ns",
	Description: "Returns the current time since epoch in nanoseconds.",
	Decl: types.NewFunction(
		nil,
		types.Named("now", types.N).Description("nanoseconds since epoch"),
	),
	Nondeterministic: true,
	canSkipBctx:      false,
}

var ParseNanos = &Builtin{
	Name:        "time.parse_ns",
	Description: "Returns the time in nanoseconds parsed from the string in the given format. `undefined` if the result would be outside the valid time range that can fit within an `int64`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("layout", types.S).Description("format used for parsing, see the [Go `time` package documentation](https://golang.org/pkg/time/#Parse) for more details"),
			types.Named("value", types.S).Description("input to parse according to `layout`"),
		),
		types.Named("ns", types.N).Description("`value` in nanoseconds since epoch"),
	),
	canSkipBctx: true,
}

var ParseRFC3339Nanos = &Builtin{
	Name:        "time.parse_rfc3339_ns",
	Description: "Returns the time in nanoseconds parsed from the string in RFC3339 format. `undefined` if the result would be outside the valid time range that can fit within an `int64`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("value", types.S).Description("input string to parse in RFC3339 format"),
		),
		types.Named("ns", types.N).Description("`value` in nanoseconds since epoch"),
	),
	canSkipBctx: true,
}

var ParseDurationNanos = &Builtin{
	Name:        "time.parse_duration_ns",
	Description: "Returns the duration in nanoseconds represented by a string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("duration", types.S).Description("a duration like \"3m\"; see the [Go `time` package documentation](https://golang.org/pkg/time/#ParseDuration) for more details"),
		),
		types.Named("ns", types.N).Description("the `duration` in nanoseconds"),
	),
	canSkipBctx: true,
}

var Format = &Builtin{
	Name:        "time.format",
	Description: "Returns the formatted timestamp for the nanoseconds since epoch.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewAny(
				types.N,
				types.NewArray([]types.Type{types.N, types.S}, nil),
				types.NewArray([]types.Type{types.N, types.S, types.S}, nil),
			)).Description("a number representing the nanoseconds since the epoch (UTC); or a two-element array of the nanoseconds, and a timezone string; or a three-element array of ns, timezone string and a layout string or golang defined formatting constant (see golang supported time formats)"),
		),
		types.Named("formatted timestamp", types.S).Description("the formatted timestamp represented for the nanoseconds since the epoch in the supplied timezone (or UTC)"),
	),
	canSkipBctx: true,
}

var Date = &Builtin{
	Name:        "time.date",
	Description: "Returns the `[year, month, day]` for the nanoseconds since epoch.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewAny(
				types.N,
				types.NewArray([]types.Type{types.N, types.S}, nil),
			)).Description("a number representing the nanoseconds since the epoch (UTC); or a two-element array of the nanoseconds, and a timezone string"),
		),
		types.Named("date", types.NewArray([]types.Type{types.N, types.N, types.N}, nil)).Description("an array of `year`, `month` (1-12), and `day` (1-31)"),
	),
	canSkipBctx: true,
}

var Clock = &Builtin{
	Name:        "time.clock",
	Description: "Returns the `[hour, minute, second]` of the day for the nanoseconds since epoch.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewAny(
				types.N,
				types.NewArray([]types.Type{types.N, types.S}, nil),
			)).Description("a number representing the nanoseconds since the epoch (UTC); or a two-element array of the nanoseconds, and a timezone string"),
		),
		types.Named("output", types.NewArray([]types.Type{types.N, types.N, types.N}, nil)).
			Description("the `hour`, `minute` (0-59), and `second` (0-59) representing the time of day for the nanoseconds since epoch in the supplied timezone (or UTC)"),
	),
	canSkipBctx: true,
}

var Weekday = &Builtin{
	Name:        "time.weekday",
	Description: "Returns the day of the week (Monday, Tuesday, ...) for the nanoseconds since epoch.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.NewAny(
				types.N,
				types.NewArray([]types.Type{types.N, types.S}, nil),
			)).Description("a number representing the nanoseconds since the epoch (UTC); or a two-element array of the nanoseconds, and a timezone string"),
		),
		types.Named("day", types.S).Description("the weekday represented by `ns` nanoseconds since the epoch in the supplied timezone (or UTC)"),
	),
	canSkipBctx: true,
}

var AddDate = &Builtin{
	Name:        "time.add_date",
	Description: "Returns the nanoseconds since epoch after adding years, months and days to nanoseconds. Month & day values outside their usual ranges after the operation and will be normalized - for example, October 32 would become November 1. `undefined` if the result would be outside the valid time range that can fit within an `int64`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("ns", types.N).Description("nanoseconds since the epoch"),
			types.Named("years", types.N).Description("number of years to add"),
			types.Named("months", types.N).Description("number of months to add"),
			types.Named("days", types.N).Description("number of days to add"),
		),
		types.Named("output", types.N).Description("nanoseconds since the epoch representing the input time, with years, months and days added"),
	),
	canSkipBctx: true,
}

var Diff = &Builtin{
	Name:        "time.diff",
	Description: "Returns the difference between two unix timestamps in nanoseconds (with optional timezone strings).",
	Decl: types.NewFunction(
		types.Args(
			types.Named("ns1", types.NewAny(
				types.N,
				types.NewArray([]types.Type{types.N, types.S}, nil),
			)).Description("nanoseconds since the epoch; or a two-element array of the nanoseconds, and a timezone string"),
			types.Named("ns2", types.NewAny(
				types.N,
				types.NewArray([]types.Type{types.N, types.S}, nil),
			)).Description("nanoseconds since the epoch; or a two-element array of the nanoseconds, and a timezone string"),
		),
		types.Named("output", types.NewArray([]types.Type{types.N, types.N, types.N, types.N, types.N, types.N}, nil)).Description("difference between `ns1` and `ns2` (in their supplied timezones, if supplied, or UTC) as array of numbers: `[years, months, days, hours, minutes, seconds]`"),
	),
	canSkipBctx: true,
}

/**
 * Crypto.
 */

var CryptoX509ParseCertificates = &Builtin{
	Name: "crypto.x509.parse_certificates",
	Description: `Returns zero or more certificates from the given encoded string containing
DER certificate data.

If the input is empty, the function will return null. The input string should be a list of one or more
concatenated PEM blocks. The whole input of concatenated PEM blocks can optionally be Base64 encoded.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("certs", types.S).Description("base64 encoded DER or PEM data containing one or more certificates or a PEM string of one or more certificates"),
		),
		types.Named("output", types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.S, types.A)))).Description("parsed X.509 certificates represented as objects"),
	),
	canSkipBctx: true,
}

var CryptoX509ParseAndVerifyCertificates = &Builtin{
	Name: "crypto.x509.parse_and_verify_certificates",
	Description: `Returns one or more certificates from the given string containing PEM
or base64 encoded DER certificates after verifying the supplied certificates form a complete
certificate chain back to a trusted root.

The first certificate is treated as the root and the last is treated as the leaf,
with all others being treated as intermediates.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("certs", types.S).Description("base64 encoded DER or PEM data containing two or more certificates where the first is a root CA, the last is a leaf certificate, and all others are intermediate CAs"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.B,
			types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))),
		}, nil)).Description("array of `[valid, certs]`: if the input certificate chain could be verified then `valid` is `true` and `certs` is an array of X.509 certificates represented as objects; if the input certificate chain could not be verified then `valid` is `false` and `certs` is `[]`"),
	),
	canSkipBctx: true,
}

var CryptoX509ParseAndVerifyCertificatesWithOptions = &Builtin{
	Name: "crypto.x509.parse_and_verify_certificates_with_options",
	Description: `Returns one or more certificates from the given string containing PEM
or base64 encoded DER certificates after verifying the supplied certificates form a complete
certificate chain back to a trusted root. A config option passed as the second argument can
be used to configure the validation options used.

The first certificate is treated as the root and the last is treated as the leaf,
with all others being treated as intermediates.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("certs", types.S).Description("base64 encoded DER or PEM data containing two or more certificates where the first is a root CA, the last is a leaf certificate, and all others are intermediate CAs"),
			types.Named("options", types.NewObject(
				nil,
				types.NewDynamicProperty(types.S, types.A),
			)).Description("object containing extra configs to verify the validity of certificates. `options` object supports four fields which maps to same fields in [x509.VerifyOptions struct](https://pkg.go.dev/crypto/x509#VerifyOptions). `DNSName`, `CurrentTime`: Nanoseconds since the Unix Epoch as a number, `MaxConstraintComparisons` and `KeyUsages`. `KeyUsages` is list and can have possible values as in: `\"KeyUsageAny\"`, `\"KeyUsageServerAuth\"`, `\"KeyUsageClientAuth\"`, `\"KeyUsageCodeSigning\"`, `\"KeyUsageEmailProtection\"`, `\"KeyUsageIPSECEndSystem\"`, `\"KeyUsageIPSECTunnel\"`, `\"KeyUsageIPSECUser\"`, `\"KeyUsageTimeStamping\"`, `\"KeyUsageOCSPSigning\"`, `\"KeyUsageMicrosoftServerGatedCrypto\"`, `\"KeyUsageNetscapeServerGatedCrypto\"`, `\"KeyUsageMicrosoftCommercialCodeSigning\"`, `\"KeyUsageMicrosoftKernelCodeSigning\"` "),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.B,
			types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))),
		}, nil)).Description("array of `[valid, certs]`: if the input certificate chain could be verified then `valid` is `true` and `certs` is an array of X.509 certificates represented as objects; if the input certificate chain could not be verified then `valid` is `false` and `certs` is `[]`"),
	),
	canSkipBctx: true,
}

var CryptoX509ParseCertificateRequest = &Builtin{
	Name:        "crypto.x509.parse_certificate_request",
	Description: "Returns a PKCS #10 certificate signing request from the given PEM-encoded PKCS#10 certificate signing request.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("csr", types.S).Description("base64 string containing either a PEM encoded or DER CSR or a string containing a PEM CSR"),
		),
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("X.509 CSR represented as an object"),
	),
	canSkipBctx: true,
}

var CryptoX509ParseKeyPair = &Builtin{
	Name:        "crypto.x509.parse_keypair",
	Description: "Returns a valid key pair",
	Decl: types.NewFunction(
		types.Args(
			types.Named("cert", types.S).Description("string containing PEM or base64 encoded DER certificates"),
			types.Named("pem", types.S).Description("string containing PEM or base64 encoded DER keys"),
		),
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("if key pair is valid, returns the tls.certificate(https://pkg.go.dev/crypto/tls#Certificate) as an object. If the key pair is invalid, nil and an error are returned."),
	),
	canSkipBctx: true,
}
var CryptoX509ParseRSAPrivateKey = &Builtin{
	Name:        "crypto.x509.parse_rsa_private_key",
	Description: "Returns a JWK for signing a JWT from the given PEM-encoded RSA private key.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pem", types.S).Description("base64 string containing a PEM encoded RSA private key"),
		),
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).Description("JWK as an object"),
	),
	canSkipBctx: true,
}

var CryptoParsePrivateKeys = &Builtin{
	Name: "crypto.parse_private_keys",
	Description: `Returns zero or more private keys from the given encoded string containing DER certificate data.

If the input is empty, the function will return null. The input string should be a list of one or more concatenated PEM blocks. The whole input of concatenated PEM blocks can optionally be Base64 encoded.`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("keys", types.S).Description("PEM encoded data containing one or more private keys as concatenated blocks. Optionally Base64 encoded."),
		),
		types.Named("output", types.NewArray(nil, types.NewObject(nil, types.NewDynamicProperty(types.S, types.A)))).Description("parsed private keys represented as objects"),
	),
	canSkipBctx: true,
}

var CryptoMd5 = &Builtin{
	Name:        "crypto.md5",
	Description: "Returns a string representing the input string hashed with the MD5 function",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
		),
		types.Named("y", types.S).Description("MD5-hash of `x`"),
	),
	canSkipBctx: true,
}

var CryptoSha1 = &Builtin{
	Name:        "crypto.sha1",
	Description: "Returns a string representing the input string hashed with the SHA1 function",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
		),
		types.Named("y", types.S).Description("SHA1-hash of `x`"),
	),
	canSkipBctx: true,
}

var CryptoSha256 = &Builtin{
	Name:        "crypto.sha256",
	Description: "Returns a string representing the input string hashed with the SHA256 function",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
		),
		types.Named("y", types.S).Description("SHA256-hash of `x`"),
	),
	canSkipBctx: true,
}

var CryptoHmacMd5 = &Builtin{
	Name:        "crypto.hmac.md5",
	Description: "Returns a string representing the MD5 HMAC of the input message using the input key.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
			types.Named("key", types.S).Description("key to use"),
		),
		types.Named("y", types.S).Description("MD5-HMAC of `x`"),
	),
	canSkipBctx: true,
}

var CryptoHmacSha1 = &Builtin{
	Name:        "crypto.hmac.sha1",
	Description: "Returns a string representing the SHA1 HMAC of the input message using the input key.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
			types.Named("key", types.S).Description("key to use"),
		),
		types.Named("y", types.S).Description("SHA1-HMAC of `x`"),
	),
	canSkipBctx: true,
}

var CryptoHmacSha256 = &Builtin{
	Name:        "crypto.hmac.sha256",
	Description: "Returns a string representing the SHA256 HMAC of the input message using the input key.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
			types.Named("key", types.S).Description("key to use"),
		),
		types.Named("y", types.S).Description("SHA256-HMAC of `x`"),
	),
	canSkipBctx: true,
}

var CryptoHmacSha512 = &Builtin{
	Name:        "crypto.hmac.sha512",
	Description: "Returns a string representing the SHA512 HMAC of the input message using the input key.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.S).Description("input string"),
			types.Named("key", types.S).Description("key to use"),
		),
		types.Named("y", types.S).Description("SHA512-HMAC of `x`"),
	),
	canSkipBctx: true,
}

var CryptoHmacEqual = &Builtin{
	Name:        "crypto.hmac.equal",
	Description: "Returns a boolean representing the result of comparing two MACs for equality without leaking timing information.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("mac1", types.S).Description("mac1 to compare"),
			types.Named("mac2", types.S).Description("mac2 to compare"),
		),
		types.Named("result", types.B).Description("`true` if the MACs are equals, `false` otherwise"),
	),
	canSkipBctx: true,
}

/**
 * Graphs.
 */
var graphs = category("graph")

var WalkBuiltin = &Builtin{
	Name:        "walk",
	Relation:    true,
	Description: "Generates `[path, value]` tuples for all nested documents of `x` (recursively).  Queries can use `walk` to traverse documents nested under `x`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("value to walk"),
		),
		types.Named("output", types.NewArray(
			[]types.Type{
				types.NewArray(nil, types.A),
				types.A,
			},
			nil,
		)).Description("pairs of `path` and `value`: `path` is an array representing the pointer to `value` in `x`. If `path` is assigned a wildcard (`_`), the `walk` function will skip path creation entirely for faster evaluation."),
	),
	Categories:  graphs,
	canSkipBctx: true,
}

var ReachableBuiltin = &Builtin{
	Name:        "graph.reachable",
	Description: "Computes the set of reachable nodes in the graph from a set of starting nodes.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("graph", types.NewObject(
				nil,
				types.NewDynamicProperty(
					types.A,
					types.NewAny(
						types.SetOfAny,
						types.NewArray(nil, types.A)),
				)),
			).Description("object containing a set or array of neighboring vertices"),
			types.Named("initial", types.NewAny(types.SetOfAny, types.NewArray(nil, types.A))).Description("set or array of root vertices"),
		),
		types.Named("output", types.SetOfAny).Description("set of vertices reachable from the `initial` vertices in the directed `graph`"),
	),
	canSkipBctx: true,
}

var ReachablePathsBuiltin = &Builtin{
	Name:        "graph.reachable_paths",
	Description: "Computes the set of reachable paths in the graph from a set of starting nodes.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("graph", types.NewObject(
				nil,
				types.NewDynamicProperty(
					types.A,
					types.NewAny(
						types.SetOfAny,
						types.NewArray(nil, types.A)),
				)),
			).Description("object containing a set or array of root vertices"),
			types.Named("initial", types.NewAny(types.SetOfAny, types.NewArray(nil, types.A))).Description("initial paths"), // TODO(sr): copied. is that correct?
		),
		types.Named("output", types.NewSet(types.NewArray(nil, types.A))).Description("paths reachable from the `initial` vertices in the directed `graph`"),
	),
	canSkipBctx: true,
}

/**
 * Type
 */
var typesCat = category("types")

var IsNumber = &Builtin{
	Name:        "is_number",
	Description: "Returns `true` if the input value is a number.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is a number, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

var IsString = &Builtin{
	Name:        "is_string",
	Description: "Returns `true` if the input value is a string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is a string, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

var IsBoolean = &Builtin{
	Name:        "is_boolean",
	Description: "Returns `true` if the input value is a boolean.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is an boolean, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

var IsArray = &Builtin{
	Name:        "is_array",
	Description: "Returns `true` if the input value is an array.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is an array, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

var IsSet = &Builtin{
	Name:        "is_set",
	Description: "Returns `true` if the input value is a set.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is a set, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

var IsObject = &Builtin{
	Name:        "is_object",
	Description: "Returns true if the input value is an object",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is an object, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

var IsNull = &Builtin{
	Name:        "is_null",
	Description: "Returns `true` if the input value is null.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("result", types.B).Description("`true` if `x` is null, `false` otherwise."),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

/**
 * Type Name
 */

// TypeNameBuiltin returns the type of the input.
var TypeNameBuiltin = &Builtin{
	Name:        "type_name",
	Description: "Returns the type of its input value.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("x", types.A).Description("input value"),
		),
		types.Named("type", types.S).Description(`one of "null", "boolean", "number", "string", "array", "object", "set"`),
	),
	Categories:  typesCat,
	canSkipBctx: true,
}

/**
 * HTTP Request
 */

// Marked non-deterministic because HTTP request results can be non-deterministic.
var HTTPSend = &Builtin{
	Name:        "http.send",
	Description: "Returns a HTTP response to the given HTTP request.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("request", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).
				Description("the HTTP request object"),
		),
		types.Named("response", types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))).
			Description("the HTTP response object"),
	),
	Nondeterministic: true,
	canSkipBctx:      false,
}

/**
 * GraphQL
 */

// GraphQLParse returns a pair of AST objects from parsing/validation.
var GraphQLParse = &Builtin{
	Name:        "graphql.parse",
	Description: "Returns AST objects for a given GraphQL query and schema after validating the query against the schema. Returns undefined if errors were encountered during parsing or validation. The query and/or schema can be either GraphQL strings or AST objects from the other GraphQL builtin functions.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("query", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the GraphQL query"),
			types.Named("schema", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the GraphQL schema"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
		}, nil)).Description("`output` is of the form `[query_ast, schema_ast]`. If the GraphQL query is valid given the provided schema, then `query_ast` and `schema_ast` are objects describing the ASTs for the query and schema."),
	),
	canSkipBctx: false,
}

// GraphQLParseAndVerify returns a boolean and a pair of AST object from parsing/validation.
var GraphQLParseAndVerify = &Builtin{
	Name:        "graphql.parse_and_verify",
	Description: "Returns a boolean indicating success or failure alongside the parsed ASTs for a given GraphQL query and schema after validating the query against the schema. The query and/or schema can be either GraphQL strings or AST objects from the other GraphQL builtin functions.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("query", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the GraphQL query"),
			types.Named("schema", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the GraphQL schema"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.B,
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
			types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
		}, nil)).Description(" `output` is of the form `[valid, query_ast, schema_ast]`. If the query is valid given the provided schema, then `valid` is `true`, and `query_ast` and `schema_ast` are objects describing the ASTs for the GraphQL query and schema. Otherwise, `valid` is `false` and `query_ast` and `schema_ast` are `{}`."),
	),
	canSkipBctx: false,
}

// GraphQLParseQuery parses the input GraphQL query and returns a JSON
// representation of its AST.
var GraphQLParseQuery = &Builtin{
	Name:        "graphql.parse_query",
	Description: "Returns an AST object for a GraphQL query.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("query", types.S).Description("GraphQL query string"),
		),
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))).Description("AST object for the GraphQL query."),
	),
	canSkipBctx: true,
}

// GraphQLParseSchema parses the input GraphQL schema and returns a JSON
// representation of its AST.
var GraphQLParseSchema = &Builtin{
	Name:        "graphql.parse_schema",
	Description: "Returns an AST object for a GraphQL schema.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("schema", types.S).Description("GraphQL schema string"),
		),
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))).Description("AST object for the GraphQL schema."),
	),
	canSkipBctx: false,
}

// GraphQLIsValid returns true if a GraphQL query is valid with a given
// schema, and returns false for all other inputs.
var GraphQLIsValid = &Builtin{
	Name:        "graphql.is_valid",
	Description: "Checks that a GraphQL query is valid against a given schema. The query and/or schema can be either GraphQL strings or AST objects from the other GraphQL builtin functions.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("query", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the GraphQL query"),
			types.Named("schema", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the GraphQL schema"),
		),
		types.Named("output", types.B).Description("`true` if the query is valid under the given schema. `false` otherwise."),
	),
	canSkipBctx: false,
}

// GraphQLSchemaIsValid returns true if the input is valid GraphQL schema,
// and returns false for all other inputs.
var GraphQLSchemaIsValid = &Builtin{
	Name:        "graphql.schema_is_valid",
	Description: "Checks that the input is a valid GraphQL schema. The schema can be either a GraphQL string or an AST object from the other GraphQL builtin functions.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("schema", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the schema to verify"),
		),
		types.Named("output", types.B).Description("`true` if the schema is a valid GraphQL schema. `false` otherwise."),
	),
	canSkipBctx: false,
}

/**
 * JSON Schema
 */

// JSONSchemaVerify returns empty string if the input is valid JSON schema
// and returns error string for all other inputs.
var JSONSchemaVerify = &Builtin{
	Name:        "json.verify_schema",
	Description: "Checks that the input is a valid JSON schema object. The schema can be either a JSON string or an JSON object.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("schema", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("the schema to verify"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.B,
			types.NewAny(types.S, types.Null{}),
		}, nil)).
			Description("`output` is of the form `[valid, error]`. If the schema is valid, then `valid` is `true`, and `error` is `null`. Otherwise, `valid` is `false` and `error` is a string describing the error."),
	),
	Categories:  objectCat,
	canSkipBctx: true,
}

// JSONMatchSchema returns empty array if the document matches the JSON schema,
// and returns non-empty array with error objects otherwise.
var JSONMatchSchema = &Builtin{
	Name:        "json.match_schema",
	Description: "Checks that the document matches the JSON schema.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("document", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("document to verify by schema"),
			types.Named("schema", types.NewAny(types.S, types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)))).
				Description("schema to verify document by"),
		),
		types.Named("output", types.NewArray([]types.Type{
			types.B,
			types.NewArray(
				nil, types.NewObject(
					[]*types.StaticProperty{
						{Key: "error", Value: types.S},
						{Key: "type", Value: types.S},
						{Key: "field", Value: types.S},
						{Key: "desc", Value: types.S},
					},
					nil,
				),
			),
		}, nil)).
			Description("`output` is of the form `[match, errors]`. If the document is valid given the schema, then `match` is `true`, and `errors` is an empty array. Otherwise, `match` is `false` and `errors` is an array of objects describing the error(s)."),
	),
	Categories:  objectCat,
	canSkipBctx: false,
}

/**
 * Cloud Provider Helper Functions
 */
var providersAWSCat = category("providers.aws")

var ProvidersAWSSignReqObj = &Builtin{
	Name:        "providers.aws.sign_req",
	Description: "Signs an HTTP request object for Amazon Web Services. Currently implements [AWS Signature Version 4 request signing](https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html) by the `Authorization` header method.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("request", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).
				Description("HTTP request object"),
			types.Named("aws_config", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).
				Description("AWS configuration object"),
			types.Named("time_ns", types.N).Description("nanoseconds since the epoch"),
		),
		types.Named("signed_request", types.NewObject(nil, types.NewDynamicProperty(types.A, types.A))).
			Description("HTTP request object with `Authorization` header"),
	),
	Categories:  providersAWSCat,
	canSkipBctx: true,
}

/**
 * Rego
 */

var RegoParseModule = &Builtin{
	Name:        "rego.parse_module",
	Description: "Parses the input Rego string and returns an object representation of the AST.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("filename", types.S).Description("file name to attach to AST nodes' locations"),
			types.Named("rego", types.S).Description("Rego module"),
		),
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).
			Description("AST object for the Rego module"),
	),
	canSkipBctx: true,
}

var RegoMetadataChain = &Builtin{
	Name: "rego.metadata.chain",
	Description: `Returns the chain of metadata for the active rule.
Ordered starting at the active rule, going outward to the most distant node in its package ancestry.
A chain entry is a JSON document with two members: "path", an array representing the path of the node; and "annotations", a JSON document containing the annotations declared for the node.
The first entry in the chain always points to the active rule, even if it has no declared annotations (in which case the "annotations" member is not present).`,
	Decl: types.NewFunction(
		types.Args(),
		types.Named("chain", types.NewArray(nil, types.A)).Description("each array entry represents a node in the path ancestry (chain) of the active rule that also has declared annotations"),
	),
	canSkipBctx: true,
}

// RegoMetadataRule returns the metadata for the active rule
var RegoMetadataRule = &Builtin{
	Name:        "rego.metadata.rule",
	Description: "Returns annotations declared for the active rule and using the _rule_ scope.",
	Decl: types.NewFunction(
		types.Args(),
		types.Named("output", types.A).Description("\"rule\" scope annotations for this rule; empty object if no annotations exist"),
	),
	canSkipBctx: true,
}

/**
 * OPA
 */

// Marked non-deterministic because of unpredictable config/environment-dependent results.
var OPARuntime = &Builtin{
	Name:        "opa.runtime",
	Description: "Returns an object that describes the runtime environment where OPA is deployed.",
	Decl: types.NewFunction(
		nil,
		types.Named("output", types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))).
			Description("includes a `config` key if OPA was started with a configuration file; an `env` key containing the environment variables that the OPA process was started with; includes `version` and `commit` keys containing the version and build commit of OPA."),
	),
	Nondeterministic: true,
	canSkipBctx:      false,
}

/**
 * Trace
 */
var tracing = category("tracing")

var Trace = &Builtin{
	Name:        "trace",
	Description: "Emits `note` as a `Note` event in the query explanation. Query explanations show the exact expressions evaluated by OPA during policy execution. For example, `trace(\"Hello There!\")` includes `Note \"Hello There!\"` in the query explanation. To include variables in the message, use `sprintf`. For example, `person := \"Bob\"; trace(sprintf(\"Hello There! %v\", [person]))` will emit `Note \"Hello There! Bob\"` inside of the explanation.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("note", types.S).Description("the note to include"),
		),
		types.Named("result", types.B).Description("always `true`"),
	),
	Categories:  tracing,
	canSkipBctx: false,
}

/**
 * Glob
 */

var GlobMatch = &Builtin{
	Name:        "glob.match",
	Description: "Parses and matches strings against the glob notation. Not to be confused with `regex.globs_match`.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("glob pattern"),
			types.Named("delimiters", types.NewAny(
				types.NewArray(nil, types.S),
				types.Nl,
			)).Description("glob pattern delimiters, e.g. `[\".\", \":\"]`, defaults to `[\".\"]` if unset. If `delimiters` is `null`, glob match without delimiter."),
			types.Named("match", types.S).Description("string to match against `pattern`"),
		),
		types.Named("result", types.B).Description("true if `match` can be found in `pattern` which is separated by `delimiters`"),
	),
	canSkipBctx: false,
}

var GlobQuoteMeta = &Builtin{
	Name:        "glob.quote_meta",
	Description: "Returns a string which represents a version of the pattern where all asterisks have been escaped.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("pattern", types.S).Description("glob pattern"),
		),
		types.Named("output", types.S).Description("the escaped string of `pattern`"),
	),
	canSkipBctx: true,
	// TODO(sr): example for this was: Calling ``glob.quote_meta("*.github.com", output)`` returns ``\\*.github.com`` as ``output``.
}

/**
 * Networking
 */

var NetCIDRIntersects = &Builtin{
	Name:        "net.cidr_intersects",
	Description: "Checks if a CIDR intersects with another CIDR (e.g. `192.168.0.0/16` overlaps with `192.168.1.0/24`). Supports both IPv4 and IPv6 notations.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("cidr1", types.S).Description("first CIDR"),
			types.Named("cidr2", types.S).Description("second CIDR"),
		),
		types.Named("result", types.B).Description("`true` if `cidr1` intersects with `cidr2`"),
	),
	canSkipBctx: true,
}

var NetCIDRExpand = &Builtin{
	Name:        "net.cidr_expand",
	Description: "Expands CIDR to set of hosts  (e.g., `net.cidr_expand(\"192.168.0.0/30\")` generates 4 hosts: `{\"192.168.0.0\", \"192.168.0.1\", \"192.168.0.2\", \"192.168.0.3\"}`).",
	Decl: types.NewFunction(
		types.Args(
			types.Named("cidr", types.S).Description("CIDR to expand"),
		),
		types.Named("hosts", types.SetOfStr).Description("set of IP addresses the CIDR `cidr` expands to"),
	),
	canSkipBctx: false,
}

var NetCIDRContains = &Builtin{
	Name:        "net.cidr_contains",
	Description: "Checks if a CIDR or IP is contained within another CIDR. `output` is `true` if `cidr_or_ip` (e.g. `127.0.0.64/26` or `127.0.0.1`) is contained within `cidr` (e.g. `127.0.0.1/24`) and `false` otherwise. Supports both IPv4 and IPv6 notations.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("cidr", types.S).Description("CIDR to check against"),
			types.Named("cidr_or_ip", types.S).Description("CIDR or IP to check"),
		),
		types.Named("result", types.B).Description("`true` if `cidr_or_ip` is contained within `cidr`"),
	),
	canSkipBctx: true,
}

var NetCIDRContainsMatches = &Builtin{
	Name: "net.cidr_contains_matches",
	Description: "Checks if collections of cidrs or ips are contained within another collection of cidrs and returns matches. " +
		"This function is similar to `net.cidr_contains` except it allows callers to pass collections of CIDRs or IPs as arguments and returns the matches (as opposed to a boolean result indicating a match between two CIDRs/IPs).",
	Decl: types.NewFunction(
		types.Args(
			types.Named("cidrs", netCidrContainsMatchesOperandType).Description("CIDRs to check against"),
			types.Named("cidrs_or_ips", netCidrContainsMatchesOperandType).Description("CIDRs or IPs to check"),
		),
		types.Named("output", types.NewSet(types.NewArray([]types.Type{types.A, types.A}, nil))).Description("tuples identifying matches where `cidrs_or_ips` are contained within `cidrs`"),
	),
	canSkipBctx: true,
}

var NetCIDRMerge = &Builtin{
	Name: "net.cidr_merge",
	Description: "Merges IP addresses and subnets into the smallest possible list of CIDRs (e.g., `net.cidr_merge([\"192.0.128.0/24\", \"192.0.129.0/24\"])` generates `{\"192.0.128.0/23\"}`." +
		`This function merges adjacent subnets where possible, those contained within others and also removes any duplicates.
Supports both IPv4 and IPv6 notations. IPv6 inputs need a prefix length (e.g. "/128").`,
	Decl: types.NewFunction(
		types.Args(
			types.Named("addrs", types.NewAny(
				types.NewArray(nil, types.NewAny(types.S)),
				types.SetOfStr,
			)).Description("CIDRs or IP addresses"),
		),
		types.Named("output", types.SetOfStr).Description("smallest possible set of CIDRs obtained after merging the provided list of IP addresses and subnets in `addrs`"),
	),
	canSkipBctx: true,
}

var NetCIDRIsValid = &Builtin{
	Name:        "net.cidr_is_valid",
	Description: "Parses an IPv4/IPv6 CIDR and returns a boolean indicating if the provided CIDR is valid.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("cidr", types.S).Description("CIDR to validate"),
		),
		types.Named("result", types.B).Description("`true` if `cidr` is a valid CIDR"),
	),
	canSkipBctx: true,
}

var netCidrContainsMatchesOperandType = types.NewAny(
	types.S,
	types.NewArray(nil, types.NewAny(
		types.S,
		types.NewArray(nil, types.A),
	)),
	types.NewSet(types.NewAny(
		types.S,
		types.NewArray(nil, types.A),
	)),
	types.NewObject(nil, types.NewDynamicProperty(
		types.S,
		types.NewAny(
			types.S,
			types.NewArray(nil, types.A),
		),
	)),
)

// Marked non-deterministic because DNS resolution results can be non-deterministic.
var NetLookupIPAddr = &Builtin{
	Name:        "net.lookup_ip_addr",
	Description: "Returns the set of IP addresses (both v4 and v6) that the passed-in `name` resolves to using the standard name resolution mechanisms available.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("name", types.S).Description("domain name to resolve"),
		),
		types.Named("addrs", types.SetOfStr).Description("IP addresses (v4 and v6) that `name` resolves to"),
	),
	Nondeterministic: true,
	canSkipBctx:      false,
}

/**
 * Semantic Versions
 */

var SemVerIsValid = &Builtin{
	Name:        "semver.is_valid",
	Description: "Validates that the input is a valid SemVer string.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("vsn", types.A).Description("input to validate"),
		),
		types.Named("result", types.B).Description("`true` if `vsn` is a valid SemVer; `false` otherwise"),
	),
	canSkipBctx: true,
}

var SemVerCompare = &Builtin{
	Name:        "semver.compare",
	Description: "Compares valid SemVer formatted version strings.",
	Decl: types.NewFunction(
		types.Args(
			types.Named("a", types.S).Description("first version string"),
			types.Named("b", types.S).Description("second version string"),
		),
		types.Named("result", types.N).Description("`-1` if `a < b`; `1` if `a > b`; `0` if `a == b`"),
	),
	canSkipBctx: true,
}

/**
 * Printing
 */

// Print is a special built-in function that writes zero or more operands
// to a message buffer. The caller controls how the buffer is displayed. The
// operands may be of any type. Furthermore, unlike other built-in functions,
// undefined operands DO NOT cause the print() function to fail during
// evaluation.
var Print = &Builtin{
	Name: "print",
	Decl: types.NewVariadicFunction(nil, types.A, nil),
}

// InternalPrint represents the internal implementation of the print() function.
// The compiler rewrites print() calls to refer to the internal implementation.
var InternalPrint = &Builtin{
	Name: "internal.print",
	Decl: types.NewFunction([]types.Type{types.NewArray(nil, types.SetOfAny)}, nil),
}

var InternalTestCase = &Builtin{
	Name: "internal.test_case",
	Decl: types.NewFunction([]types.Type{types.NewArray(nil, types.A)}, nil),
}

/**
 * Deprecated built-ins.
 */

// SetDiff has been replaced by the minus built-in.
var SetDiff = &Builtin{
	Name: "set_diff",
	Decl: types.NewFunction(
		types.Args(
			types.SetOfAny,
			types.SetOfAny,
		),
		types.SetOfAny,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// NetCIDROverlap has been replaced by the `net.cidr_contains` built-in.
var NetCIDROverlap = &Builtin{
	Name: "net.cidr_overlap",
	Decl: types.NewFunction(
		types.Args(
			types.S,
			types.S,
		),
		types.B,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// CastArray checks the underlying type of the input. If it is array or set, an array
// containing the values is returned. If it is not an array, an error is thrown.
var CastArray = &Builtin{
	Name: "cast_array",
	Decl: types.NewFunction(
		types.Args(types.A),
		types.NewArray(nil, types.A),
	),
	deprecated:  true,
	canSkipBctx: true,
}

// CastSet checks the underlying type of the input.
// If it is a set, the set is returned.
// If it is an array, the array is returned in set form (all duplicates removed)
// If neither, an error is thrown
var CastSet = &Builtin{
	Name: "cast_set",
	Decl: types.NewFunction(
		types.Args(types.A),
		types.SetOfAny,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// CastString returns input if it is a string; if not returns error.
// For formatting variables, see sprintf
var CastString = &Builtin{
	Name: "cast_string",
	Decl: types.NewFunction(
		types.Args(types.A),
		types.S,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// CastBoolean returns input if it is a boolean; if not returns error.
var CastBoolean = &Builtin{
	Name: "cast_boolean",
	Decl: types.NewFunction(
		types.Args(types.A),
		types.B,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// CastNull returns null if input is null; if not returns error.
var CastNull = &Builtin{
	Name: "cast_null",
	Decl: types.NewFunction(
		types.Args(types.A),
		types.Nl,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// CastObject returns the given object if it is null; throws an error otherwise
var CastObject = &Builtin{
	Name: "cast_object",
	Decl: types.NewFunction(
		types.Args(types.A),
		types.NewObject(nil, types.NewDynamicProperty(types.A, types.A)),
	),
	deprecated:  true,
	canSkipBctx: true,
}

// RegexMatchDeprecated declares `re_match` which has been deprecated. Use `regex.match` instead.
var RegexMatchDeprecated = &Builtin{
	Name: "re_match",
	Decl: types.NewFunction(
		types.Args(
			types.S,
			types.S,
		),
		types.B,
	),
	deprecated:  true,
	canSkipBctx: false,
}

// All takes a list and returns true if all of the items
// are true. A collection of length 0 returns true.
var All = &Builtin{
	Name: "all",
	Decl: types.NewFunction(
		types.Args(
			types.NewAny(
				types.SetOfAny,
				types.NewArray(nil, types.A),
			),
		),
		types.B,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// Any takes a collection and returns true if any of the items
// is true. A collection of length 0 returns false.
var Any = &Builtin{
	Name: "any",
	Decl: types.NewFunction(
		types.Args(
			types.NewAny(
				types.SetOfAny,
				types.NewArray(nil, types.A),
			),
		),
		types.B,
	),
	deprecated:  true,
	canSkipBctx: true,
}

// Builtin represents a built-in function supported by OPA. Every built-in
// function is uniquely identified by a name.
type Builtin struct {
	Name        string `json:"name"`                  // Unique name of built-in function, e.g., <name>(arg1,arg2,...,argN)
	Description string `json:"description,omitempty"` // Description of what the built-in function does.

	// Categories of the built-in function. Omitted for namespaced
	// built-ins, i.e. "array.concat" is taken to be of the "array" category.
	// "minus" for example, is part of two categories: numbers and sets. (NOTE(sr): aspirational)
	Categories []string `json:"categories,omitempty"`

	Decl             *types.Function `json:"decl"`               // Built-in function type declaration.
	Infix            string          `json:"infix,omitempty"`    // Unique name of infix operator. Default should be unset.
	Relation         bool            `json:"relation,omitempty"` // Indicates if the built-in acts as a relation.
	deprecated       bool            // Indicates if the built-in has been deprecated.
	canSkipBctx      bool            // Built-in needs no data from the built-in context.
	Nondeterministic bool            `json:"nondeterministic,omitempty"` // Indicates if the built-in returns non-deterministic results.
}

// category is a helper for specifying a Builtin's Categories
func category(cs ...string) []string {
	return cs
}

// Minimal returns a shallow copy of b with the descriptions and categories and
// named arguments stripped out.
func (b *Builtin) Minimal() *Builtin {
	cpy := *b
	fargs := b.Decl.FuncArgs()
	if fargs.Variadic != nil {
		cpy.Decl = types.NewVariadicFunction(fargs.Args, fargs.Variadic, b.Decl.Result())
	} else {
		cpy.Decl = types.NewFunction(fargs.Args, b.Decl.Result())
	}
	cpy.Categories = nil
	cpy.Description = ""
	return &cpy
}

// IsDeprecated returns true if the Builtin function is deprecated and will be removed in a future release.
func (b *Builtin) IsDeprecated() bool {
	return b.deprecated
}

// IsDeterministic returns true if the Builtin function returns non-deterministic results.
func (b *Builtin) IsNondeterministic() bool {
	return b.Nondeterministic
}

// Expr creates a new expression for the built-in with the given operands.
func (b *Builtin) Expr(operands ...*Term) *Expr {
	ts := make([]*Term, len(operands)+1)
	ts[0] = NewTerm(b.Ref())
	for i := range operands {
		ts[i+1] = operands[i]
	}
	return &Expr{
		Terms: ts,
	}
}

// Call creates a new term for the built-in with the given operands.
func (b *Builtin) Call(operands ...*Term) *Term {
	call := make(Call, len(operands)+1)
	call[0] = NewTerm(b.Ref())
	for i := range operands {
		call[i+1] = operands[i]
	}
	return NewTerm(call)
}

// Ref returns a Ref that refers to the built-in function.
func (b *Builtin) Ref() Ref {
	parts := strings.Split(b.Name, ".")
	ref := make(Ref, len(parts))
	ref[0] = VarTerm(parts[0])
	for i := 1; i < len(parts); i++ {
		ref[i] = InternedTerm(parts[i])
	}
	return ref
}

// IsTargetPos returns true if a variable in the i-th position will be bound by
// evaluating the call expression.
func (b *Builtin) IsTargetPos(i int) bool {
	return b.Decl.Arity() == i
}

// NeedsBuiltInContext returns true if the built-in depends on the built-in context.
func (b *Builtin) NeedsBuiltInContext() bool {
	// Negated, so built-ins we don't know about (and who don't know about this option)
	// will get a built-in context provided to them.
	return !b.canSkipBctx
}

func init() {
	BuiltinMap = map[string]*Builtin{}
	for _, b := range &DefaultBuiltins {
		RegisterBuiltin(b)
	}
}
