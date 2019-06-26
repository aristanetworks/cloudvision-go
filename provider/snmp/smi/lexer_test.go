package smi

import (
	"io/ioutil"
	"testing"

	"modernc.org/golex/lex"
)

type lexerTestCase struct {
	name           string
	input          string
	expectedOutput []Token
}

func runLexerTest(t *testing.T, tc lexerTestCase) {
	l, err := newLexer(tc.input)
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for {
		tok := l.scan()
		if tok.char.Rune == lex.RuneEOF {
			if i != len(tc.expectedOutput) {
				t.Fatal("Expected more tokens")
			}
			break
		}
		if i >= len(tc.expectedOutput) {
			t.Fatal("More tokens than expected")
		}
		if tok.tokType != tc.expectedOutput[i].tokType {
			t.Fatalf("On token %d, expected tokType '%s', got '%s'",
				i, tokstr(tc.expectedOutput[i].tokType), tokstr(tok.tokType))
		}
		if tok.literal != tc.expectedOutput[i].literal {
			t.Fatalf("On token %d, expected '%s', got '%s'",
				i, tc.expectedOutput[i].literal, tok.literal)
		}
		i++
	}
}

func newToken(t int, l string) Token {
	return Token{
		tokType: t,
		literal: l,
	}
}

func tokenTestCase(t int, l string) lexerTestCase {
	return lexerTestCase{
		name:           tokstr(t),
		input:          l,
		expectedOutput: []Token{newToken(t, l)},
	}
}

func TestLexer(t *testing.T) {
	for _, tc := range []lexerTestCase{
		{
			name: "IF-MIB sample",
			input: `
IF-MIB DEFINITIONS ::= BEGIN
-- the Interfaces table

-- The Interfaces table contains information on the entity's

ifTable OBJECT-TYPE
    SYNTAX      SEQUENCE OF IfEntry
    MAX-ACCESS  not-accessible
    STATUS      current
    DESCRIPTION
            "A list of interface entries.  The number of entries is
            given by the value of ifNumber."
    ::= { interfaces 2 }

ifEntry OBJECT-TYPE
    SYNTAX      IfEntry
    MAX-ACCESS  not-accessible
    STATUS      current
    DESCRIPTION
            "An entry containing management information applicable to a
            particular interface."
    INDEX   { ifIndex }
    ::= { ifTable 1 }
END
`,
			expectedOutput: []Token{
				newToken(UPPERCASE_IDENTIFIER, "IF-MIB"),
				newToken(DEFINITIONS, "DEFINITIONS"),
				newToken(COLON_COLON_EQUAL, "::="),
				newToken(BEGIN, "BEGIN"),
				newToken(LOWERCASE_IDENTIFIER, "ifTable"),
				newToken(OBJECT_TYPE, "OBJECT-TYPE"),
				newToken(SYNTAX, "SYNTAX"),
				newToken(SEQUENCE, "SEQUENCE"),
				newToken(OF, "OF"),
				newToken(UPPERCASE_IDENTIFIER, "IfEntry"),
				newToken(MAX_ACCESS, "MAX-ACCESS"),
				newToken(LOWERCASE_IDENTIFIER, "not-accessible"),
				newToken(STATUS, "STATUS"),
				newToken(LOWERCASE_IDENTIFIER, "current"),
				newToken(DESCRIPTION, "DESCRIPTION"),
				newToken(QUOTED_STRING, "A list of interface entries. The "+
					"number of entries is given by the value of ifNumber."),
				newToken(COLON_COLON_EQUAL, "::="),
				newToken('{', "{"),
				newToken(LOWERCASE_IDENTIFIER, "interfaces"),
				newToken(NUMBER, "2"),
				newToken('}', "}"),
				newToken(LOWERCASE_IDENTIFIER, "ifEntry"),
				newToken(OBJECT_TYPE, "OBJECT-TYPE"),
				newToken(SYNTAX, "SYNTAX"),
				newToken(UPPERCASE_IDENTIFIER, "IfEntry"),
				newToken(MAX_ACCESS, "MAX-ACCESS"),
				newToken(LOWERCASE_IDENTIFIER, "not-accessible"),
				newToken(STATUS, "STATUS"),
				newToken(LOWERCASE_IDENTIFIER, "current"),
				newToken(DESCRIPTION, "DESCRIPTION"),
				newToken(QUOTED_STRING, "An entry containing management "+
					"information applicable to a particular interface."),
				newToken(INDEX, "INDEX"),
				newToken('{', "{"),
				newToken(LOWERCASE_IDENTIFIER, "ifIndex"),
				newToken('}', "}"),
				newToken(COLON_COLON_EQUAL, "::="),
				newToken('{', "{"),
				newToken(LOWERCASE_IDENTIFIER, "ifTable"),
				newToken(NUMBER, "1"),
				newToken('}', "}"),
				newToken(END, "END"),
			},
		},
		tokenTestCase(ACCESS, "ACCESS"),
		tokenTestCase(AGENT_CAPABILITIES, "AGENT-CAPABILITIES"),
		tokenTestCase(APPLICATION, "APPLICATION"),
		tokenTestCase(AUGMENTS, "AUGMENTS"),
		tokenTestCase(BEGIN, "BEGIN"),
		tokenTestCase(BITS, "BITS"),
		tokenTestCase(CONTACT_INFO, "CONTACT-INFO"),
		tokenTestCase(CREATION_REQUIRES, "CREATION-REQUIRES"),
		tokenTestCase(COUNTER32, "Counter32"),
		tokenTestCase(COUNTER64, "Counter64"),
		tokenTestCase(DEFINITIONS, "DEFINITIONS"),
		tokenTestCase(DEFVAL, "DEFVAL"),
		tokenTestCase(DESCRIPTION, "DESCRIPTION"),
		tokenTestCase(DISPLAY_HINT, "DISPLAY-HINT"),
		tokenTestCase(END, "END"),
		tokenTestCase(ENTERPRISE, "ENTERPRISE"),
		tokenTestCase(EXTENDS, "EXTENDS"),
		tokenTestCase(FROM, "FROM"),
		tokenTestCase(GROUP, "GROUP"),
		tokenTestCase(GAUGE32, "Gauge32"),
		tokenTestCase(IDENTIFIER, "IDENTIFIER"),
		tokenTestCase(IMPLICIT, "IMPLICIT"),
		tokenTestCase(IMPLIED, "IMPLIED"),
		tokenTestCase(IMPORTS, "IMPORTS"),
		tokenTestCase(INCLUDES, "INCLUDES"),
		tokenTestCase(INDEX, "INDEX"),
		tokenTestCase(INSTALL_ERRORS, "INSTALL-ERRORS"),
		tokenTestCase(INTEGER, "INTEGER"),
		tokenTestCase(INTEGER32, "Integer32"),
		tokenTestCase(INTEGER64, "Integer64"),
		tokenTestCase(IPADDRESS, "IpAddress"),
		tokenTestCase(LAST_UPDATED, "LAST-UPDATED"),
		tokenTestCase(MANDATORY_GROUPS, "MANDATORY-GROUPS"),
		tokenTestCase(MAX_ACCESS, "MAX-ACCESS"),
		tokenTestCase(MIN_ACCESS, "MIN-ACCESS"),
		tokenTestCase(MODULE, "MODULE"),
		tokenTestCase(MODULE_COMPLIANCE, "MODULE-COMPLIANCE"),
		tokenTestCase(MODULE_IDENTITY, "MODULE-IDENTITY"),
		tokenTestCase(NOTIFICATION_GROUP, "NOTIFICATION-GROUP"),
		tokenTestCase(NOTIFICATION_TYPE, "NOTIFICATION-TYPE"),
		tokenTestCase(NOTIFICATIONS, "NOTIFICATIONS"),
		tokenTestCase(OBJECT, "OBJECT"),
		tokenTestCase(OBJECT_GROUP, "OBJECT-GROUP"),
		tokenTestCase(OBJECT_IDENTITY, "OBJECT-IDENTITY"),
		tokenTestCase(OBJECT_TYPE, "OBJECT-TYPE"),
		tokenTestCase(OBJECTS, "OBJECTS"),
		tokenTestCase(OCTET, "OCTET"),
		tokenTestCase(OF, "OF"),
		tokenTestCase(ORGANIZATION, "ORGANIZATION"),
		tokenTestCase(OPAQUE, "Opaque"),
		tokenTestCase(PIB_ACCESS, "PIB-ACCESS"),
		tokenTestCase(PIB_DEFINITIONS, "PIB-DEFINITIONS"),
		tokenTestCase(PIB_INDEX, "PIB-INDEX"),
		tokenTestCase(PIB_MIN_ACCESS, "PIB-MIN-ACCESS"),
		tokenTestCase(PIB_REFERENCES, "PIB-REFERENCES"),
		tokenTestCase(PIB_TAG, "PIB-TAG"),
		tokenTestCase(POLICY_ACCESS, "POLICY-ACCESS"),
		tokenTestCase(PRODUCT_RELEASE, "PRODUCT-RELEASE"),
		tokenTestCase(REFERENCE, "REFERENCE"),
		tokenTestCase(REVISION, "REVISION"),
		tokenTestCase(SEQUENCE, "SEQUENCE"),
		tokenTestCase(SIZE, "SIZE"),
		tokenTestCase(STATUS, "STATUS"),
		tokenTestCase(STRING, "STRING"),
		tokenTestCase(SUBJECT_CATEGORIES, "SUBJECT-CATEGORIES"),
		tokenTestCase(SUPPORTS, "SUPPORTS"),
		tokenTestCase(SYNTAX, "SYNTAX"),
		tokenTestCase(TEXTUAL_CONVENTION, "TEXTUAL-CONVENTION"),
		tokenTestCase(TIMETICKS, "TimeTicks"),
		tokenTestCase(TRAP_TYPE, "TRAP-TYPE"),
		tokenTestCase(UNIQUENESS, "UNIQUENESS"),
		tokenTestCase(UNITS, "UNITS"),
		tokenTestCase(UNIVERSAL, "UNIVERSAL"),
		tokenTestCase(UNSIGNED32, "Unsigned32"),
		tokenTestCase(UNSIGNED64, "Unsigned64"),
		tokenTestCase(VALUE, "VALUE"),
		tokenTestCase(VARIABLES, "VARIABLES"),
		tokenTestCase(VARIATION, "VARIATION"),
		tokenTestCase(WRITE_SYNTAX, "WRITE-SYNTAX"),
		tokenTestCase('[', "["),
		tokenTestCase(']', "]"),
		tokenTestCase('{', "{"),
		tokenTestCase('}', "}"),
		tokenTestCase('(', "("),
		tokenTestCase(')', ")"),
		tokenTestCase(':', ":"),
		tokenTestCase(';', ";"),
		tokenTestCase(',', ","),
		tokenTestCase('-', "-"),
		tokenTestCase('.', "."),
		tokenTestCase('|', "|"),
		tokenTestCase(DOT_DOT, ".."),
		tokenTestCase(COLON_COLON_EQUAL, "::="),
		{
			name:  "comment",
			input: "OF -- comment stuff junk\nOF",
			expectedOutput: []Token{
				newToken(OF, "OF"),
				newToken(OF, "OF"),
			},
		},
		{
			name:           "tab",
			input:          "\t",
			expectedOutput: []Token{},
		},
		tokenTestCase(UPPERCASE_IDENTIFIER, "ETC-ETC-STUFF"),
		tokenTestCase(LOWERCASE_IDENTIFIER, "eTC-ETC-STUFF"),
		tokenTestCase(NUMBER, "123"),
		tokenTestCase(NEGATIVE_NUMBER, "-123"),
		tokenTestCase(BIN_STRING, "'010101'B"),
		tokenTestCase(HEX_STRING, "'e8d53b'h"),
		{
			name:  "simple quoted string",
			input: "\"mumbo jumbo etc., and so on\"",
			expectedOutput: []Token{
				newToken(QUOTED_STRING, "mumbo jumbo etc., and so on"),
			},
		},
		{
			name:  "long quoted string",
			input: "\"the quick  brown\tfox\njumped\"",
			expectedOutput: []Token{
				newToken(QUOTED_STRING, "the quick brown fox jumped"),
			},
		},
		{
			name: "MACRO",
			input: "MODULE-IDENTITY MACRO ::=\n" +
				"BEGIN\n" +
				"TYPE NOTATION ::=\n" +
				"    \"LAST-UPDATED\" value(Update ExtUTCTime)" +
				"    \"ORGANIZATION\" Text" +
				"VALUE NOTATION ::=\n" +
				"    stuff\n" +
				"END\n",
			expectedOutput: []Token{
				newToken(MODULE_IDENTITY, "MODULE-IDENTITY"),
				newToken(MACRO, "MACRO"),
				newToken(END, "END"),
			},
		},
		{
			name: "EXPORTS",
			input: "EXPORTS -- EVERYTHING\n" +
				"internet, directory, mgmt,\n" +
				"etc;\n",
			expectedOutput: []Token{
				newToken(EXPORTS, "EXPORTS"),
				newToken(';', ";"),
			},
		},
		{
			name: "CHOICE",
			input: "ObjectSyntax ::=\n" +
				"    CHOICE {\n" +
				"        stuff\n" +
				"		 etc.\n" +
				"    }\n" +
				"UNITS",
			expectedOutput: []Token{
				newToken(UPPERCASE_IDENTIFIER, "ObjectSyntax"),
				newToken(COLON_COLON_EQUAL, "::="),
				newToken(CHOICE, "CHOICE"),
				newToken('}', "}"),
				newToken(UNITS, "UNITS"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runLexerTest(t, tc)
		})
	}
}

func TestLexerOnMIBs(t *testing.T) {
	b, err := ioutil.ReadFile("mibs/IF-MIB")
	if err != nil {
		t.Fatal(err)
	}
	lx, err := newLexer(string(b))
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		i++
		tok := lx.scan()
		if tok.char.Rune == lex.RuneEOF {
			break
		}
	}
}
