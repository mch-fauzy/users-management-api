package users

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	GetData(filter UserFilter, page, size int) ([]UserView, error)
	CountTotalData(filter UserFilter) (int, error)
	GetProfile(uuid string) (*ProfileView, error)
	UpdateProfile(uuid string, profile *UpdateProfile) (*UpdateProfile, error)
	DeleteUserByID(uuid string) error
}

type UserRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideUserRepositoryMySQL(db *infras.MySQLConn) *UserRepositoryMySQL {
	return &UserRepositoryMySQL{
		DB: db,
	}
}

func (r *UserRepositoryMySQL) GetData(filter UserFilter, page, size int) ([]UserView, error) {
	query := `
		SELECT 
			u.username,
		 	p.name,
			u.role,
			p.gender,
			p.dob,
			p.education,
			p.city,
			p.province,
			p.address,
			p.phone_number,
			s.job_role,
			s.status,
			pl.city AS placement,
			d.name AS department_name
		FROM 
			ums_users AS u
		LEFT JOIN
			ums_profiles AS p
				ON u.profile_id = p.id
		LEFT JOIN
			ums_status AS s
				ON u.status_id = s.id
		LEFT JOIN
			ums_placement AS pl
				ON u.placement_id = pl.id
		LEFT JOIN
			ums_dept AS d
				ON u.dept_id = d.id
	`

	// Add filters
	args := []interface{}{}
	if filter.Name != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " p.name LIKE ?"
		args = append(args, "%"+filter.Name+"%")
	}

	if filter.City != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " p.city LIKE ?"
		args = append(args, "%"+filter.City+"%")
	}

	if filter.Province != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " p.province LIKE ?"
		args = append(args, "%"+filter.Province+"%")
	}

	if filter.JobRole != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " s.job_role LIKE ?"
		args = append(args, "%"+filter.JobRole+"%")
	}

	if filter.Status != "" {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " s.status LIKE ?"
		args = append(args, "%"+filter.Status+"%")
	}

	// Add pagination
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 5
	}

	query += " LIMIT ? OFFSET ?"
	offset := (page - 1) * size
	args = append(args, size, offset)

	var users []UserView
	err := r.DB.Read.Select(&users, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read users from db")
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryMySQL) CountTotalData(filter UserFilter) (int, error) {
	totalDataQuery := `
		SELECT 
			COUNT(*) 
		FROM 
			ums_users AS u
		LEFT JOIN 
			ums_profiles AS p 
				ON u.profile_id = p.id
		LEFT JOIN 
			ums_status AS s 
				ON u.status_id = s.id
		LEFT JOIN 
			ums_placement AS pl 
				ON u.placement_id = pl.id
		LEFT JOIN 
			ums_dept AS d 
				ON u.dept_id = d.id
	`

	// Add filters
	argsTotalData := []interface{}{}
	if filter.Name != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " p.name LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.Name+"%")
	}

	if filter.City != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " p.city LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.City+"%")
	}

	if filter.Province != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " p.province LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.Province+"%")
	}

	if filter.JobRole != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " s.job_role LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.JobRole+"%")
	}

	if filter.Status != "" {
		if len(argsTotalData) > 0 {
			totalDataQuery += " AND"
		} else {
			totalDataQuery += " WHERE"
		}
		totalDataQuery += " s.status LIKE ?"
		argsTotalData = append(argsTotalData, "%"+filter.Status+"%")
	}

	var totalData int
	err := r.DB.Read.Get(&totalData, totalDataQuery, argsTotalData...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get total data")
		return 0, err
	}

	return totalData, nil
}

func (r *UserRepositoryMySQL) DeleteUserByID(uuid string) error {
	query := `
		DELETE FROM ums_users
		WHERE id = ? AND role != 'admin'
	`

	result, err := r.DB.Write.Exec(query, uuid)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Failed to check affected rows")
		return err
	}

	if rowsAffected == 0 {
		return errors.New("User with 'admin' role cannot be deleted or user not found")
	}

	return nil
}

func (r *UserRepositoryMySQL) GetProfile(uuid string) (*ProfileView, error) {
	query := `
	SELECT 
		p.name,
		u.role,
		p.gender,
		p.dob,
		p.education,
		p.city,
		p.province,
		p.address,
		p.phone_number,
		s.job_role,
		s.status,
		pl.city AS placement_city,
		d.name AS department_name
	FROM 
		ums_users AS u
	LEFT JOIN
		ums_profiles AS p
			ON u.profile_id = p.id
	LEFT JOIN
		ums_status AS s
			ON u.status_id = s.id
	LEFT JOIN
		ums_placement AS pl
			ON u.placement_id = pl.id
	LEFT JOIN
		ums_dept AS d
			ON u.dept_id = d.id
	WHERE u.id = ?
	`

	var profile ProfileView
	err := r.DB.Read.Get(&profile, query, uuid)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get profile")
		return nil, err
	}
	return &profile, nil
}

func (r *UserRepositoryMySQL) UpdateProfile(uuid string, profile *UpdateProfile) (*UpdateProfile, error) {
	// Build the SET clause for the update query
	setClauses := []string{
		"p.name = COALESCE(?, p.name)",
		"p.gender = COALESCE(?, p.gender)",
		"p.dob = COALESCE(?, p.dob)",
		"p.education = COALESCE(?, p.education)",
		"p.address = COALESCE(?, p.address)",
		"p.city = COALESCE(?, p.city)",
		"p.province = COALESCE(?, p.province)",
		"p.phone_number = COALESCE(?, p.phone_number)",
		"p.updated_at = ?, p.updated_by = ?",
		"u.updated_at = ?, u.updated_by = ?",
	}

	// Construct the SQL query
	query := fmt.Sprintf(`
		UPDATE ums_profiles AS p
		INNER JOIN ums_users AS u ON p.id = u.profile_id
		SET %s
		WHERE u.id = ?`,
		strings.Join(setClauses, ", "))

	profile.UpdatedAt = time.Now()
	values := []interface{}{
		lowercaseOrNil(profile.Name),
		lowercaseOrNil(profile.Gender),
		profile.DoB,
		lowercaseOrNil(profile.Education),
		lowercaseOrNil(profile.Address),
		lowercaseOrNil(profile.City),
		lowercaseOrNil(profile.Province),
		profile.PhoneNumber,
		profile.UpdatedAt,
		profile.UpdatedBy,
		profile.UpdatedAt,
		profile.UpdatedBy,
		uuid,
	}

	// Execute the update query
	_, err := r.DB.Write.Exec(query, values...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update profile")
		return nil, err
	}

	return profile, nil
}

// lowercaseOrNil converts a string pointer to lowercase if it's not nil.
func lowercaseOrNil(s *string) interface{} {
	if s != nil {
		return strings.ToLower(*s)
	}
	return nil
}
