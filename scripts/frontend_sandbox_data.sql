
USE `psjudge_frontend`;
INSERT INTO contest (title, start_time, end_time) VALUES ('Test Contest', '2018-10-10 10:00:00', '2018-10-10 12:00:00');
SELECT @contest_id := id FROM contest WHERE title='Test Contest';
INSERT INTO assignment (contest_id, title, article) VALUES (@contest_id, 'A+B Problem', 'Solve A+B problem using C++ or PASCAL');
INSERT INTO assignment (contest_id, title, article) VALUES (@contest_id, 'A*B Problem', 'Solve A*B problem using C++ or PASCAL');
INSERT INTO user (username, password, roles, active_contest_id) VALUES ('Martin', '2018', 'student', @contest_id);
INSERT INTO user (username, password, roles, active_contest_id) VALUES ('Piter', '2018', 'judge', @contest_id);
