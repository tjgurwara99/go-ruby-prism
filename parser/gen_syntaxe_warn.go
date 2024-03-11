/******************************************************************************/
/* This file is generated by the templates/template.rb script and should not  */
/* be modified manually. See                                                  */
/* templates/gen_syntaxe_warn.go.erb                                          */
/* if you are looking to modify the                                           */
/* template                                                                   */
/******************************************************************************/

package parser

type SyntaxWarningLevel int
type SyntaxWarningType int

const (
	SyntaxWarningDefault SyntaxWarningLevel = iota
	SyntaxWarningVerbose
)

const (
	AMBIGUOUS_FIRST_ARGUMENT_MINUS SyntaxWarningType = 0
	AMBIGUOUS_FIRST_ARGUMENT_PLUS  SyntaxWarningType = 1
	AMBIGUOUS_PREFIX_STAR          SyntaxWarningType = 2
	AMBIGUOUS_SLASH                SyntaxWarningType = 3
	DOT_DOT_DOT_EOL                SyntaxWarningType = 4
	EQUAL_IN_CONDITIONAL           SyntaxWarningType = 5
	END_IN_METHOD                  SyntaxWarningType = 6
	DUPLICATED_HASH_KEY            SyntaxWarningType = 7
	DUPLICATED_WHEN_CLAUSE         SyntaxWarningType = 8
	FLOAT_OUT_OF_RANGE             SyntaxWarningType = 9
	INTEGER_IN_FLIP_FLOP           SyntaxWarningType = 10
	KEYWORD_EOL                    SyntaxWarningType = 11
)

var SyntaxWarningLevels = []SyntaxWarningLevel{SyntaxWarningDefault, SyntaxWarningVerbose}
var SyntaxWarningTypes = []SyntaxWarningType{
	AMBIGUOUS_FIRST_ARGUMENT_MINUS,
	AMBIGUOUS_FIRST_ARGUMENT_PLUS,
	AMBIGUOUS_PREFIX_STAR,
	AMBIGUOUS_SLASH,
	DOT_DOT_DOT_EOL,
	EQUAL_IN_CONDITIONAL,
	END_IN_METHOD,
	DUPLICATED_HASH_KEY,
	DUPLICATED_WHEN_CLAUSE,
	FLOAT_OUT_OF_RANGE,
	INTEGER_IN_FLIP_FLOP,
	KEYWORD_EOL,
}

type SyntaxWarning struct {
	Message  string
	Location *Location
	Level    SyntaxWarningLevel
	Type     SyntaxWarningType
}

func NewSyntaxWarning(
	message string,
	location *Location,
	level SyntaxWarningLevel,
	warnType SyntaxWarningType,
) *SyntaxWarning {
	return &SyntaxWarning{
		Message:  message,
		Location: location,
		Level:    level,
		Type:     warnType,
	}
}
