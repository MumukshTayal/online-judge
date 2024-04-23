CREATE TABLE user_profile (
  user_email VARCHAR(255) PRIMARY KEY,
  user_name VARCHAR(255)
);

CREATE TABLE contest (
  contest_id INTEGER PRIMARY KEY AUTOINCREMENT,
  contest_title VARCHAR(255) UNIQUE,
  contest_description TEXT,
  contest_start_time DATETIME,
  contest_end_time DATETIME,
  is_public BOOLEAN,
  creator_email VARCHAR(255),
  FOREIGN KEY (creator_email) REFERENCES user_profile(user_email)
);

CREATE TABLE distribute_problems_to_contest (
    contest_id INT,
    problem_id INT,
    FOREIGN KEY(contest_id) REFERENCES contest(contest_id),
    FOREIGN KEY(problem_id) REFERENCES problem(problem_id)
);

CREATE TABLE testcase (
    testcase_id INTEGER PRIMARY KEY AUTOINCREMENT,
    testcase_input BLOB,
    testcase_output BLOB,
    problem_id INT,
    hidden BOOLEAN,
    FOREIGN KEY (problem_id) REFERENCES problem(problem_id)
); 

CREATE TABLE problem (
  problem_id INTEGER PRIMARY KEY AUTOINCREMENT,
  problem_title VARCHAR(255) UNIQUE,
  problem_description TEXT,
  constraints_desc TEXT,
  input_format TEXT,
  output_format TEXT,
  sample_input TEXT,
  sample_output TEXT,
  creator_email VARCHAR(255),
  is_private BOOLEAN,
  FOREIGN KEY (creator_email) REFERENCES user_profile(user_email)
); 

CREATE TABLE allowed_list (
  contest_id INT,
  language VARCHAR(255),
  time_limit INT,
  memory_limit INT,
  PRIMARY KEY (contest_id, language)
); 

CREATE TABLE submission (
  submission_id INTEGER PRIMARY KEY AUTOINCREMENT,
  problem_id INT,
  user_email VARCHAR(255),
  code TEXT,
  language VARCHAR(255),
  result VARCHAR(255),
  execution_time INT,
  memory_used INT,
  submission_date_time DATETIME,
  FOREIGN KEY (problem_id) REFERENCES problem(problem_id),
  FOREIGN KEY (user_email) REFERENCES user_profile(user_email)
);

CREATE TABLE contest_user (
  contest_id INT,
  user_email VARCHAR(255),
  PRIMARY KEY (contest_id, user_email),
  FOREIGN KEY (contest_id) REFERENCES contest(contest_id),
  FOREIGN KEY (user_email) REFERENCES user_profile(user_email)
);

CREATE TABLE virtual_contest (
  virutal_contest_id INT PRIMARY KEY,
  contest_id INT,
  user_email VARCHAR(255),
  FOREIGN KEY (contest_id) REFERENCES contest(contest_id),
  FOREIGN KEY (user_email) REFERENCES user_profile(user_email)
);
