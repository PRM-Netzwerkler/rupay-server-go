package schemas

import "github.com/guregu/null/v5"

type SavaResponse struct {
	Success bool        `json:"succes"`           // Note the field name matches the typo in your sample ("succes" not "success")
	Result  null.Float  `json:"result,omitempty"` // Use a pointer to distinguish between empty and missing fields
	Error   null.String `json:"error,omitempty"`  // Use a pointer for optional fields
}

type SavaUser struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
