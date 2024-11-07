package repository

const (
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS user (
		discord_user_id CHAR(30) PRIMARY KEY,
	    name CHAR(20) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS scrum (
		user_id CHAR(30) NOT NULL,
		goal TEXT NOT NULL,
		commitment TEXT,
		feel_score INTEGER,
		feel_reason TEXT,
		created_at TIMESTAMP DEFAULT (datetime('now', '+09:00')),
		
		PRIMARY KEY (user_id, created_at),
		FOREIGN KEY (user_id) REFERENCES user(discord_user_id)
	);
	`
	insertUserQuery       = `INSERT INTO user (name, discord_user_id) VALUES (?, ?)`
	insertScrumQuery      = `INSERT INTO scrum (user_id, goal, commitment, feel_score, feel_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	existScrumQuery       = `SELECT COUNT(user_id) FROM scrum WHERE user_id = ? and created_at = ?`
	existUserQuery        = `SELECT EXISTS (SELECT 1 FROM user WHERE discord_user_id = ?)`
	selectTodayScrumQuery = `SELECT u.name, goal, commitment, feel_score, feel_reason FROM scrum as s JOIN user as u ON s.user_id = u.discord_user_id WHERE created_at = ?`
)
