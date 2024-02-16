/* added */ 
CREATE TABLE user_profile (
  user_id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_name VARCHAR(255),
  user_email VARCHAR(255) UNIQUE,
);

/* added */ 
CREATE TABLE contest (
  contest_id INTEGER PRIMARY KEY AUTOINCREMENT,
  contest_title VARCHAR(255) UNIQUE,
  contest_description TEXT,
  contest_start_time DATETIME,
  contest_end_time DATETIME,
  is_public BOOLEAN,
  creator_id INT,
  FOREIGN KEY (creator_id) REFERENCES user_profile(user_id)
);

/* added */ 
CREATE TABLE distribute_problems_to_contest (
    contest_id INT,
    problem_id INT,
    FOREIGN KEY(contest_id) REFERENCES contest(contest_id),
    FOREIGN KEY(problem_id) REFERENCES problem(problem_id)
);

/* added */ 
CREATE TABLE testcase (
    testcase_id INTEGER PRIMARY KEY AUTOINCREMENT,
    testcase_input BLOB,
    testcase_output BLOB,
    problem_id INT,
    FOREIGN KEY (problem_id) REFERENCES problem(problem_id)
);

/* added */ 
CREATE TABLE problem (
  problem_id INTEGER PRIMARY KEY AUTOINCREMENT,
  problem_title VARCHAR(255) UNIQUE,
  problem_description TEXT,
  constraints_desc TEXT,
  creator_id INT,
  is_private BOOLEAN,
  FOREIGN KEY (creator_id) REFERENCES user_profile(user_id)
); 


/* added */ 
CREATE TABLE allowed_list (
  contest_id INT,
  language VARCHAR(255),
  time_limit INT,
  memory_limit INT,
  PRIMARY KEY (contest_id, language)
);

/* added */ 
CREATE TABLE submission (
  submission_id INTEGER PRIMARY KEY AUTOINCREMENT,
  problem_id INT,
  user_id INT,
  code TEXT,
  language VARCHAR(255),
  result VARCHAR(255),
  execution_time INT,
  memory_used INT,
  submission_date_time DATETIME,
  FOREIGN KEY (problem_id) REFERENCES problem(problem_id),
  FOREIGN KEY (user_id) REFERENCES user_profile(user_id)
);

/* added */ 
CREATE TABLE contest_user (
  contest_id INT,
  user_id INT,
  PRIMARY KEY (contest_id, user_id),
  FOREIGN KEY (contest_id) REFERENCES contest(contest_id),
  FOREIGN KEY (user_id) REFERENCES user_profile(user_id)
);

/* added */ 
CREATE TABLE virtual_contest (
  virutal_contest_id INT PRIMARY KEY,
  contest_id INT,
  user_id INT,
  FOREIGN KEY (contest_id) REFERENCES contest(contest_id),
  FOREIGN KEY (user_id) REFERENCES user_profile(user_id)
);
