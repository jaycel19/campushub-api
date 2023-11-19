package services

import (
	"context"
	"time"
)

type Profile struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	ProfilePic        string    `json:"profile_pic"`
	Age               string    `json:"age"`
	Program           string    `json:"program"`
	Year              string    `json:"Year"`
	ProfileBackground string    `json:"profile_background"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ProfileRequest struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Age               string    `json:"age"`
	Program           string    `json:"program"`
	Year              string    `json:"Year"`
	ProfileBackground string    `json:"profile_background"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (p *Profile) GetAllProfiles() ([]*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from profiles`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var profiles []*Profile
	for rows.Next() {
		var profile Profile
		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Username,
			&profile.ProfilePic,
			&profile.Age,
			&profile.Program,
			&profile.Year,
			&profile.ProfileBackground,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, &profile)
	}
	return profiles, nil
}

func (p *Profile) GetProfileByUser(username string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	SELECT * FROM profiles WHERE username = $1
	`
	var profile Profile

	row := db.QueryRowContext(ctx, query, username)
	err := row.Scan(
		&profile.ID,
		&profile.Name,
		&profile.Username,
		&profile.ProfilePic,
		&profile.Age,
		&profile.Program,
		&profile.Year,
		&profile.ProfileBackground,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// TODO: check if all fields is populated
func (p *Profile) UpdateProfile(username string, body Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE profiles
		SET
			name = $1,
			username = $2,
			age = $3,
			program = $4,
			year = $5,
			updated_at = $6
		WHERE username=$7
	`

	_, err := db.ExecContext(
		ctx,
		query,
		body.Name,
		body.Username,
		body.Age,
		body.Program,
		body.Year,
		time.Now(),
		username,
	)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (p *Profile) ProfileChangeBackground(username string, background string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE profiles
		SET
			profile_background = $1,
			updated_at = $2
		WHERE username=$3
	`

	_, err := db.ExecContext(
		ctx,
		query,
		background,
		time.Now(),
		username,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Profile) DeleteProfile(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM profiles WHERE username=$1`
	_, err := db.ExecContext(ctx, query, username)
	if err != nil {
		return err
	}
	return nil
}

func (p *Profile) CreateProfile(profile Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO profiles (name, username, age, program, year, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning *
	`

	_, err := db.ExecContext(
		ctx,
		query,
		profile.Name,
		profile.Username,
		profile.Age,
		profile.Program,
		profile.Year,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
