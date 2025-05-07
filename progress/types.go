package progress

import "strings"

const (
	ProgressPublic    string = "public"
	ProgressProtected string = "protected"
	ProgressPrivate   string = "string"
)

const (
	ProgressBlob       string = "blob"
	ProgressCharacter  string = "character"
	ProgressClob       string = "clob"
	ProgressInteger    string = "integer"
	ProgressDecimal    string = "decimal"
	ProgressLogical    string = "logical"
	ProgressDate       string = "date"
	ProgressDateTime   string = "datetime"
	ProgressDateTimeTz string = "datetime-tz"
	ProgressInt64      string = "int64"
	ProgressLongChar   string = "longchar"
	ProgressMemPtr     string = "memptr"
	ProgressRaw        string = "raw"
	ProgressRecid      string = "recid"
	ProgressRowid      string = "rowid"
)

type ProgressTypes struct {
	URI        string
	Classes    []string
	Methods    []string
	Properties []string
}

const (
	ProgressLVC string = "lvc_"
	ProgressLVL string = "lvl_"
	ProgressLVI string = "lvi_"
	ProgressLVD string = "lvd_"
	ProgressIVC string = "ivc_"
	ProgressIVL string = "ivl_"
	ProgressIVD string = "ivi_"
	ProgressIVI string = "ivd_"
	ProgressOVC string = "ovc_"
	ProgressOVL string = "ovl_"
	ProgressOVD string = "ovi_"
	ProgressOVI string = "ovd_"
)

func IndexOfRestrictedText(text string) int {
	var idx int
	if idx <= 0 {idx = strings.Index(text, ProgressLVC)}
	if idx <= 0 {idx = strings.Index(text, ProgressLVL)}
	if idx <= 0 {idx = strings.Index(text, ProgressLVI)}
	if idx <= 0 {idx = strings.Index(text, ProgressLVD)}
	if idx <= 0 {idx = strings.Index(text, ProgressIVC)}
	if idx <= 0 {idx = strings.Index(text, ProgressIVL)}
	if idx <= 0 {idx = strings.Index(text, ProgressIVD)}
	if idx <= 0 {idx = strings.Index(text, ProgressIVI)}
	if idx <= 0 {idx = strings.Index(text, ProgressOVC)}
	if idx <= 0 {idx = strings.Index(text, ProgressOVL)}
	if idx <= 0 {idx = strings.Index(text, ProgressOVD)}
	if idx <= 0 {idx = strings.Index(text, ProgressOVI)}
	return idx
}

func IsDefaultDataType(text string) bool {
	switch text {
		case ProgressBlob:       return true
		case ProgressCharacter:  return true
		case ProgressClob:       return true
		case ProgressInteger:    return true
		case ProgressInt64:      return true
		case ProgressDate:       return true
		case ProgressDateTime:   return true
		case ProgressDateTimeTz: return true
		case ProgressDecimal:    return true
		case ProgressLongChar:   return true
		case ProgressLogical:    return true
		case ProgressMemPtr:     return true
		case ProgressRaw:        return true
		case ProgressRecid:      return true
		case ProgressRowid:      return true
	}
	return false
}
