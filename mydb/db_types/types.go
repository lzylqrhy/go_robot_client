package db_types

type DBRow map[string]interface{}

func (r DBRow)GetString(field string) string {
	if f, isOK := r[field]; isOK {
		switch f.(type) {
		case []byte:
			return string(f.([]byte))
		case string:
			return f.(string)
		case nil:
			return ""
		}
	}
	return ""
}