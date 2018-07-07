
USE `psjudge_frontend`;
-- passwordHash keeps hash of '2018' string.
SET @passwordHash = '4476d6a3edee189e699ca3c2cfd80905abc8d999954a08d1c504e6ae437cc28dd4194a1051d84bf5cb2cfc19e09e339ce7f2ff83fec56b07ee39ec2205c2adba';
INSERT INTO contest (title, start_time, end_time) VALUES ('Test Contest', '2018-10-10 10:00:00', '2018-10-10 12:00:00');
SELECT @contest_id := id FROM contest WHERE title='Test Contest';
INSERT INTO assignment (contest_id, uuid, title, article) VALUES (@contest_id, 'd31fc051b882448798705a5b016141af', 'A+B Problem', 'Solve A+B problem using C++ or PASCAL');
INSERT INTO assignment (contest_id, uuid, title, article) VALUES (@contest_id, '4fec83ac54424d909a5a343c3cf36aa9', 'A*B Problem', 'Solve A*B problem using C++ or PASCAL');
INSERT INTO user (username, password, roles, active_contest_id) VALUES ('Martin', @passwordHash, 'student', @contest_id);
INSERT INTO user (username, password, roles, active_contest_id) VALUES ('Piter', @passwordHash, 'judge', @contest_id);
