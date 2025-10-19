// contains the model of data sent through/from endpoint
package model

// represents albums in our program
type Album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

