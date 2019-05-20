%yyc c
%yyn c = l.Next()
%yym l.Mark()
%yyt l.sc

%{
package smi

import (
    "bytes"
    "fmt"
    "go/token"
    "strings"

    "modernc.org/golex/lex"
)

var tokMap = map[int]string{
	ACCESS: "ACCESS",
	AGENT_CAPABILITIES: "AGENT_CAPABILITIES",
	APPLICATION: "APPLICATION",
	AUGMENTS: "AUGMENTS",
	BEGIN: "BEGIN",
	BIN_STRING: "BIN_STRING",
	BITS: "BITS",
	CHOICE: "CHOICE",
	COLON_COLON_EQUAL: "COLON_COLON_EQUAL",
	COMMENT: "COMMENT",
	CONTACT_INFO: "CONTACT_INFO",
	CREATION_REQUIRES: "CREATION_REQUIRES",
	COUNTER32: "COUNTER32",
	COUNTER64: "COUNTER64",
	DEFINITIONS: "DEFINITIONS",
	DEFVAL: "DEFVAL",
	DESCRIPTION: "DESCRIPTION",
	DISPLAY_HINT: "DISPLAY_HINT",
	DOT_DOT: "DOT_DOT",
	END: "END",
	ENTERPRISE: "ENTERPRISE",
	EXPORTS: "EXPORTS",
	EXTENDS: "EXTENDS",
	FROM: "FROM",
	GROUP: "GROUP",
	GAUGE32: "GAUGE32",
	HEX_STRING: "HEX_STRING",
	IDENTIFIER: "IDENTIFIER",
	IMPLICIT: "IMPLICIT",
	IMPLIED: "IMPLIED",
	IMPORTS: "IMPORTS",
	INCLUDES: "INCLUDES",
	INDEX: "INDEX",
	INSTALL_ERRORS: "INSTALL_ERRORS",
	INTEGER: "INTEGER",
	INTEGER32: "INTEGER32",
	INTEGER64: "INTEGER64",
	IPADDRESS: "IPADDRESS",
	LAST_UPDATED: "LAST_UPDATED",
	LOWERCASE_IDENTIFIER: "LOWERCASE_IDENTIFIER",
	MACRO: "MACRO",
	MANDATORY_GROUPS: "MANDATORY_GROUPS",
	MAX_ACCESS: "MAX_ACCESS",
	MIN_ACCESS: "MIN_ACCESS",
	MODULE: "MODULE",
	MODULE_COMPLIANCE: "MODULE_COMPLIANCE",
	MODULE_IDENTITY: "MODULE_IDENTITY",
	NEGATIVE_NUMBER: "NEGATIVE_NUMBER",
	NOTIFICATION_GROUP: "NOTIFICATION_GROUP",
	NOTIFICATION_TYPE: "NOTIFICATION_TYPE",
	NOTIFICATIONS: "NOTIFICATIONS",
	NUMBER: "NUMBER",
	OBJECT: "OBJECT",
	OBJECT_GROUP: "OBJECT_GROUP",
	OBJECT_IDENTITY: "OBJECT_IDENTITY",
	OBJECT_TYPE: "OBJECT_TYPE",
	OBJECTS: "OBJECTS",
	OCTET: "OCTET",
	OF: "OF",
	ORGANIZATION: "ORGANIZATION",
	OPAQUE: "Opaque",
	PIB_ACCESS: "PIB_ACCESS",
	PIB_DEFINITIONS: "PIB_DEFINITIONS",
	PIB_INDEX: "PIB_INDEX",
	PIB_MIN_ACCESS: "PIB_MIN_ACCESS",
	PIB_REFERENCES: "PIB_REFERENCES",
	PIB_TAG: "PIB_TAG",
	POLICY_ACCESS: "POLICY_ACCESS",
	PRODUCT_RELEASE: "PRODUCT_RELEASE",
	QUOTED_STRING: "QUOTED_STRING",
	REFERENCE: "REFERENCE",
	REVISION: "REVISION",
	SEQUENCE: "SEQUENCE",
	SIZE: "SIZE",
	STATUS: "STATUS",
	STRING: "STRING",
	SUBJECT_CATEGORIES: "SUBJECT_CATEGORIES",
	SUPPORTS: "SUPPORTS",
	SYNTAX: "SYNTAX",
	TEXTUAL_CONVENTION: "TEXTUAL_CONVENTION",
	TIMETICKS: "TIMETICKS",
	TRAP_TYPE: "TRAP_TYPE",
	UNIQUENESS: "UNIQUENESS",
	UNITS: "UNITS",
	UNIVERSAL: "UNIVERSAL",
	UNSIGNED32: "UNSIGNED32",
	UNSIGNED64: "UNSIGNED64",
	UPPERCASE_IDENTIFIER: "UPPERCASE_IDENTIFIER",
	VALUE: "VALUE",
	VARIABLES: "VARIABLES",
	VARIATION: "VARIATION",
	WRITE_SYNTAX: "WRITE_SYNTAX",
}

func tokstr(t int) string {
    if t < 127 {
        return string(t)
    }
    return tokMap[t]
}

func sanitized(s string) string {
    f := strings.Fields(s)
    j := strings.Join(f, " ")
    return strings.Replace(j, "\"", "", -1)
}

type lexer struct {
	*lex.Lexer
    sc int
    modules map[string]*Module
}

type Token struct {
    tokType int
    literal string
    char lex.Char
}

func (t *Token) String() string {
    return fmt.Sprintf("%s %s %v", tokstr(t.tokType), t.literal, t.char)
}

func (l *lexer) char(r int) lex.Char {
	return lex.NewChar(l.First.Pos(), rune(r))
}

func (l *lexer) token(t int) Token {
    return Token{
        tokType: t,
        literal: string(l.TokenBytes(nil)),
        char: l.char(t),
    }
}

func (l *lexer) tokenWithLiteral(t int, s string) Token {
    return Token{
        tokType: t,
        literal: s,
        char: l.char(t),
    }
}

func (l *lexer) begin(cond int) {
    l.sc = cond
}

func (l *lexer) scan() Token {
	c := l.Enter()

    const (
        INITIAL = iota
        Macro
        Exports
        Choice
    )
%}

%x Macro
%x Exports
%x Choice

binString \'[01]*\'[bB]
colonColonEqual ::=
comment --.*$
dotDot ".."
hexString \'[0123456789AaBbCcDdEeFf]*\'[hH]
leadingZero 0+[0-9]
lowercaseIdentifier [a-z](-?[a-zA-Z0-9_]+)*-?
negNumber -[0-9]+
number [0-9]+
quotedString \"[^\"]*\"
uppercaseIdentifier [A-Z](-?[a-zA-Z0-9_]+)*-?
%%

    c = l.Rule0()

[ \t\r\n]+

MACRO
    l.begin(Macro)
    return l.token(MACRO)

<Macro>[ \t\r\n]+

<Macro>.

<Macro>END
    l.begin(INITIAL)
    return l.token(END)

EXPORTS
    l.begin(Exports)
    return l.token(EXPORTS)

<Exports>[ \t\r\n]+

<Exports>[^\;]*

<Exports>\;
    l.begin(INITIAL)
    return l.token(';')

CHOICE
    l.begin(Choice)
    return l.token(CHOICE)

<Choice>[ \t\r\n]+

<Choice>[^\}]*

<Choice>\}
    l.begin(INITIAL)
    return l.token('}')

{comment}

{dotDot} return l.token(DOT_DOT)
{colonColonEqual} return l.token(COLON_COLON_EQUAL)

ACCESS return l.token(ACCESS)
AGENT-CAPABILITIES return l.token(AGENT_CAPABILITIES)
APPLICATION return l.token(APPLICATION)
AUGMENTS return l.token(AUGMENTS)
BEGIN return l.token(BEGIN)
BITS return l.token(BITS)
CONTACT-INFO return l.token(CONTACT_INFO)
CREATION-REQUIRES return l.token(CREATION_REQUIRES)
Counter32 return l.token(COUNTER32)
Counter64 return l.token(COUNTER64)
DEFINITIONS return l.token(DEFINITIONS)
DEFVAL return l.token(DEFVAL)
DESCRIPTION return l.token(DESCRIPTION)
DISPLAY-HINT return l.token(DISPLAY_HINT)
END return l.token(END)
ENTERPRISE return l.token(ENTERPRISE)
EXTENDS return l.token(EXTENDS)
FROM return l.token(FROM)
GROUP return l.token(GROUP)
Gauge32 return l.token(GAUGE32)
IDENTIFIER return l.token(IDENTIFIER)
IMPLICIT return l.token(IMPLICIT)
IMPLIED return l.token(IMPLIED)
IMPORTS return l.token(IMPORTS)
INCLUDES return l.token(INCLUDES)
INDEX return l.token(INDEX)
INSTALL-ERRORS return l.token(INSTALL_ERRORS)
INTEGER return l.token(INTEGER)
Integer32 return l.token(INTEGER32)
Integer64 return l.token(INTEGER64)
IpAddress return l.token(IPADDRESS)
LAST-UPDATED return l.token(LAST_UPDATED)
MANDATORY-GROUPS return l.token(MANDATORY_GROUPS)
MAX-ACCESS return l.token(MAX_ACCESS)
MIN-ACCESS return l.token(MIN_ACCESS)
MODULE return l.token(MODULE)
MODULE-COMPLIANCE return l.token(MODULE_COMPLIANCE)
MODULE-IDENTITY return l.token(MODULE_IDENTITY)
NOTIFICATION-GROUP return l.token(NOTIFICATION_GROUP)
NOTIFICATION-TYPE return l.token(NOTIFICATION_TYPE)
NOTIFICATIONS return l.token(NOTIFICATIONS)
OBJECT return l.token(OBJECT)
OBJECT-GROUP return l.token(OBJECT_GROUP)
OBJECT-IDENTITY return l.token(OBJECT_IDENTITY)
OBJECT-TYPE return l.token(OBJECT_TYPE)
OBJECTS return l.token(OBJECTS)
OCTET return l.token(OCTET)
OF return l.token(OF)
ORGANIZATION return l.token(ORGANIZATION)
Opaque return l.token(OPAQUE)
PIB-ACCESS return l.token(PIB_ACCESS)
PIB-DEFINITIONS return l.token(PIB_DEFINITIONS)
PIB-INDEX return l.token(PIB_INDEX)
PIB-MIN-ACCESS return l.token(PIB_MIN_ACCESS)
PIB-REFERENCES return l.token(PIB_REFERENCES)
PIB-TAG return l.token(PIB_TAG)
POLICY-ACCESS return l.token(POLICY_ACCESS)
PRODUCT-RELEASE return l.token(PRODUCT_RELEASE)
REFERENCE return l.token(REFERENCE)
REVISION return l.token(REVISION)
SEQUENCE return l.token(SEQUENCE)
SIZE return l.token(SIZE)
STATUS return l.token(STATUS)
STRING return l.token(STRING)
SUBJECT-CATEGORIES return l.token(SUBJECT_CATEGORIES)
SUPPORTS return l.token(SUPPORTS)
SYNTAX return l.token(SYNTAX)
TEXTUAL-CONVENTION return l.token(TEXTUAL_CONVENTION)
TimeTicks return l.token(TIMETICKS)
TRAP-TYPE return l.token(TRAP_TYPE)
UNIQUENESS return l.token(UNIQUENESS)
UNITS return l.token(UNITS)
UNIVERSAL return l.token(UNIVERSAL)
Unsigned32 return l.token(UNSIGNED32)
Unsigned64 return l.token(UNSIGNED64)
VALUE return l.token(VALUE)
VARIABLES return l.token(VARIABLES)
VARIATION return l.token(VARIATION)
WRITE-SYNTAX return l.token(WRITE_SYNTAX)

{uppercaseIdentifier} return l.token(UPPERCASE_IDENTIFIER)
{lowercaseIdentifier} return l.token(LOWERCASE_IDENTIFIER)
{leadingZero} // XXX TODO Error
{number} return l.token(NUMBER)
{negNumber} return l.token(NEGATIVE_NUMBER) // XXX TODO 64-bit?
{binString} return l.token(BIN_STRING)
{hexString} return l.token(HEX_STRING)
{quotedString}
    return l.tokenWithLiteral(QUOTED_STRING,
        sanitized(string(l.TokenBytes(nil))))

%%

    if c, ok := l.Abort(); ok {
		return l.token(c)
	}

    goto yyAction
}

func (l *lexer) Lex(lval *yySymType) int {
    t := l.scan()
    lval.token = t
    return t.tokType
}

func newLexer(text string) (*lexer, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("whatever.txt", -1, len(text))
	src := bytes.NewBufferString(text)
	lx, err := lex.New(file, src)
	if err != nil {
		return nil, err
	}
	return &lexer{lx, 0, make(map[string]*Module)}, nil
}
