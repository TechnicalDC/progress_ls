package progress

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
