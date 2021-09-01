package model

type EsUser struct {
	ID       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nickname"`
	Gender   int32  `json:"gender"`
}

func (EsUser) GetIndexName() string {
	return "user"
}

func (EsUser) GetMapping() string {
	userMapping := `
	{
		"mappings" : {
			"properties" : {
				"Mobile" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"Mobile" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"NickName" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"gender" : {
					"type" : "integer"
				}
			}
		}
	}`
	return userMapping
}
