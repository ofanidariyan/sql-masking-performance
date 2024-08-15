package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bxcodec/faker"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
)

const (
	encryptionKey = "0123456789abcdef0123456789abcdef" // Use lowercase for constants
	batchSize     = 10000
)

// User represents a user with name and email.
type User struct {
	Name  string `faker:"name"`
	Email string `faker:"email"`
}

func main() {
	// Establish database connection
	db, err := sql.Open("mysql", "homecare:homecaredb123@tcp(127.0.0.1:3308)/homecare-db")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Check for command-line arguments to determine which function to run
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <masking-golang|masking-sql> <num_records>")
		return
	}

	command := os.Args[1]
	numRecords, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid number of records:", err)
		return
	}

	switch command {
	case "masking-golang":
		insertWithEncryptionInGo(db, numRecords)
		listWithDecryptionInGo(db)
	case "masking-sql":
		insertWithEncryptionInSQL(db, numRecords)
		listWithDecryptionInSQL(db)
	default:
		fmt.Println("Invalid command. Use 'masking-golang' or 'masking-sql'")
	}
}

// Insert user data with encryption performed directly in the SQL query.
func insertWithEncryptionInSQL(db *sql.DB, numRecords int) {
	start := time.Now()

	for i := 0; i < numRecords; i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > numRecords {
			batchEnd = numRecords
		}

		// Construct the SQL query with placeholders for batch insertion
		valueStrings := make([]string, 0, batchEnd-i)
		valueArgs := make([]interface{}, 0, 2*(batchEnd-i)) // 2 placeholders per user
		for j := i; j < batchEnd; j++ {
			user := generateFakeUser()
			valueStrings = append(valueStrings, "(AES_ENCRYPT(?, ?), AES_ENCRYPT(?, ?))")
			valueArgs = append(valueArgs, user.Name, encryptionKey, user.Email, encryptionKey)
		}
		stmt := fmt.Sprintf("INSERT INTO users (name, email) VALUES %s", strings.Join(valueStrings, ","))

		// Execute the batch insertion
		_, err := db.Exec(stmt, valueArgs...)
		if err != nil {
			log.Fatal("Error inserting data with SQL encryption:", err)
		}
	}

	elapsed := time.Since(start)
	color.Green("1. Insert with encryption in SQL (%d records): %s\n", numRecords, elapsed)
}

// Insert user data with encryption performed in Go before sending to the database.
func insertWithEncryptionInGo(db *sql.DB, numRecords int) {
	start := time.Now()

	for i := 0; i < numRecords; i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > numRecords {
			batchEnd = numRecords
		}

		// Construct the SQL query with placeholders for batch insertion
		valueStrings := make([]string, 0, batchEnd-i)
		valueArgs := make([]interface{}, 0, 2*(batchEnd-i)) // 2 placeholders per user
		for j := i; j < batchEnd; j++ {
			user := generateFakeUser()
			encryptedName := encryptAES(user.Name, encryptionKey)
			encryptedEmail := encryptAES(user.Email, encryptionKey)
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, encryptedName, encryptedEmail)
		}
		stmt := fmt.Sprintf("INSERT INTO users (name, email) VALUES %s", strings.Join(valueStrings, ","))

		// Execute the batch insertion
		_, err := db.Exec(stmt, valueArgs...)
		if err != nil {
			log.Fatal("Error inserting data with Go encryption:", err)
		}
	}

	elapsed := time.Since(start)
	color.Green("1. Insert with encryption in Go (%d records): %s\n", numRecords, elapsed)
}

// Retrieve and decrypt user data using MySQL's AES_DECRYPT function
func listWithDecryptionInSQL(db *sql.DB) {
	start := time.Now()

	rows, err := db.Query(`
		SELECT id, 
		       AES_DECRYPT(name, ?), 
		       AES_DECRYPT(email, ?) 
		FROM users
	`, encryptionKey, encryptionKey)
	if err != nil {
		log.Fatal("Error retrieving data with SQL decryption:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Fatal("Error scanning row:", err)
		}
		// fmt.Println(id, name, email) // Uncomment to print the decrypted data
	}

	elapsed := time.Since(start)
	color.Green("2. List with decryption in SQL: %s", elapsed)
}

// Retrieve user data and perform decryption in Go
func listWithDecryptionInGo(db *sql.DB) {
	start := time.Now()

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal("Error retrieving data:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var encryptedName, encryptedEmail string
		if err := rows.Scan(&id, &encryptedName, &encryptedEmail); err != nil {
			log.Fatal("Error scanning row:", err)
		}

		// Decrypt data in Go
		_ = decryptAES(encryptedName, encryptionKey)
		_ = decryptAES(encryptedEmail, encryptionKey)

		//fmt.Println(id, name, email)
	}

	elapsed := time.Since(start)
	color.Green("2. List with decryption in Go: %s", elapsed)
}

// Helper function to generate a fake user
func generateFakeUser() User {
	var user User
	if err := faker.FakeData(&user); err != nil {
		log.Fatal("Error generating fake user:", err)
	}
	return user
}

// AES encryption function
func encryptAES(plaintext, keyString string) string {
	key := []byte(keyString)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext)
}

// AES decryption function
func decryptAES(ciphertextB64, keyString string) string {
	key := []byte(keyString)
	ciphertext, err := base64.URLEncoding.DecodeString(ciphertextB64)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}
