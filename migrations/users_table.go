package main

import (
	"strings"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras" // Your provided package
	"github.com/rs/zerolog/log"
)

func main() {
	config := configs.Get()

	mysqlConn := infras.ProvideMySQLConn(config)

	query := `
		CREATE TABLE IF NOT EXISTS ums_dept (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL,
			created_by VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(255) NOT NULL,
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS ums_placement (
			id VARCHAR(50) PRIMARY KEY,
			city VARCHAR(50) UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL,
			created_by VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(255) NOT NULL,
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS ums_profiles (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255),
			gender ENUM('male', 'female'),
			dob VARCHAR(10),
			education VARCHAR(50),
			address VARCHAR(255),
			city VARCHAR(50),
			province VARCHAR(50),
			phone_number VARCHAR(50),
			created_at TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP,
			updated_by VARCHAR(255),
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS ums_status (
			id VARCHAR(36) PRIMARY KEY,
			status VARCHAR(50),
			job_role VARCHAR(50),
			created_at TIMESTAMP,
			created_by VARCHAR(255),
			updated_at TIMESTAMP,
			updated_by VARCHAR(255),
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS ums_users (
			id VARCHAR(36) PRIMARY KEY,
			profile_id VARCHAR(36) UNIQUE NOT NULL,
			status_id VARCHAR(36) UNIQUE NOT NULL,
			dept_id VARCHAR(50),
			placement_id VARCHAR(50),
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARBINARY(255) NOT NULL,
			role ENUM('admin', 'trainee') NOT NULL,
			created_at TIMESTAMP NOT NULL,
			created_by VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(255) NOT NULL,
			deleted_at TIMESTAMP,
			deleted_by VARCHAR(255),
			FOREIGN KEY (profile_id) REFERENCES ums_profiles(id),
			FOREIGN KEY (status_id) REFERENCES ums_status(id),
			FOREIGN KEY (dept_id) REFERENCES ums_dept(id),
			FOREIGN KEY (placement_id) REFERENCES ums_placement(id)
		);

		-- Create indexes
		CREATE INDEX idx_users_username ON ums_users (username);
		CREATE INDEX idx_profiles_address ON ums_profiles (address);
		CREATE INDEX idx_profiles_name ON ums_profiles (name);
		CREATE INDEX idx_status_role ON ums_status (job_role);
		CREATE INDEX idx_status ON ums_status (status);
	`

	// Split the query into separate statements
	statements := strings.Split(query, ";")

	// Remove empty statements
	var validStatements []string
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			validStatements = append(validStatements, stmt)
		}
	}

	// Execute each statement
	for _, stmt := range validStatements {
		_, err := mysqlConn.Write.Exec(stmt)
		if err != nil {
			log.Error().Err(err).Msg("Error executing statement")
			return
		}
	}
}
