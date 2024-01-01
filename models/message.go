package models

type Message struct {
    Id      int    `json:"Id"`
    ShowId  string `json:"ShowId"`
    JuryId  string `json:"JuryId"`
    Rating  int    `json:"Rating"`
    Comment string `json:"Comment"`
}