package logger

type (
	Log struct {
		Level   int    `json:"level"`
		Line    int    `json:"line"`
		File    string `json:"file"`
		Func    string `json:"func"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	User struct {
		Name string `json:"name"`
		Id   any    `json:"id"`
		Role any    `json:"role"`
	}

	Access struct {
		Ip        string `json:"ip"`
		Route     string `json:"route"`
		Method    string `json:"method"`
		UserAgent string `json:"user_agent"`
		User      *User  `json:"user"`
	}
)
