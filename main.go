package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

// Create a db Handle Object
var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	cfg := mysql.Config{
		// For now passing the values, need to modify later like
		// User:   os.Getenv("DBUser"),
		// Passwd: os.Getenv("DBPass"),
		User:   "root",
		Passwd: "12345678",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	// Getting a db Handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// Check if connection is alive
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("DB Connected!")

	albums, err := albumsByArtist("iniko")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albums)

	albumByIdO, err := albumById(3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albumByIdO)

	newAlbum := Album{
		Title:  "company",
		Artist: "emiway",
		Price:  80,
	}

	res, err := addAlbum(newAlbum)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Album inserted: %v\n", res)

	// db.Close()

}

// Get all the albums by a artist
func albumsByArtist(artist string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM ALBUM WHERE artist= ?", artist)
	if err != nil {
		return albums, fmt.Errorf("albumsByArtist %q: %v", artist, err)
	}

	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return albums, fmt.Errorf("albumsByArtist %q: %v", artist, err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return albums, fmt.Errorf("albumsByArtist %q: %v", artist, err)
	}
	return albums, err
}

// Getting a single row
func albumById(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album where id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById: No album with the requested Id %q", id)
		}
		return alb, fmt.Errorf("getArtistById %q: %v", id, err)
	}
	return alb, nil
}

// Put data into the table

func addAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist,price) VALUES (?,?,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum %v", err)
	}
	return id, err
}
