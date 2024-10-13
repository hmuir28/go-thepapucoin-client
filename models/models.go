package models

type Transaction struct {
	To			*string						`json:"to"			validate:"required"`
	From		*string						`json:"from" 		validate:"required"`
	Money		int							`json:"money"		validate:"required"`
}
