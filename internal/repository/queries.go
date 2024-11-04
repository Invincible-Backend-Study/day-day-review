package repository

const (
	createTableQuery = `
	CREATE TABLE IF NOT EXISTS user (
		discord_user_id CHAR(30) PRIMARY KEY,
	    name CHAR(20) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS scrum (
		user_id INTEGER NOT NULL,
		goal TEXT NOT NULL,
		commitment TEXT,
		feel_score INTEGER,
		feel_reason TEXT,
		created_at TIMESTAMP DEFAULT (datetime('now', '+09:00')),
		
		PRIMARY KEY (user_id, created_at)
	);
	`
	insertUserQuery  = `INSERT INTO User (name, discord_user_id) VALUES (?, ?)`
	insertScrumQuery = `INSERT INTO Scrum (user_id, goal, commitment, feel_score, feel_reason, created_at) VALUES (?, ?, ?, ?, ?, ?)`
)
