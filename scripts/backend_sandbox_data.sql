
USE `psjudge_frontend`;

-- testPasswordHash keeps hash of '2018' string.
SET @testPasswordHash = '4476d6a3edee189e699ca3c2cfd80905abc8d999954a08d1c504e6ae437cc28dd4194a1051d84bf5cb2cfc19e09e339ce7f2ff83fec56b07ee39ec2205c2adba';

-- adminPasswordHash keeps hash of 'ej19g72d' string.
SET @adminPasswordHash = 'ee17183b0ffeccf0d2dbb0e608c8383383fe31d8a00e69505107be86d1a063d18164c39010ec9d3ebeff369653e0d631249923d7d4a9c8a724c6e4bfacff5111';

INSERT INTO contest (title, max_reviews) VALUES ('Test Contest', 2);
SELECT @contest_id := id FROM contest WHERE title='Test Contest';
INSERT INTO `group` (`name`) VALUES ('test_group');
SELECT @group_id := id FROM `group` WHERE name='test_group';
INSERT INTO appointment (group_id, contest_id, start_time, end_time) VALUES (@group_id, @contest_id, '2016-10-10 10:00:00', '2020-10-10 11:00:00');

INSERT INTO user (username, password, roles, active_contest_id) VALUES ('psjudge', @adminPasswordHash, 'admin', 1);
INSERT INTO user (username, password, roles, active_contest_id) VALUES ('test_judge', @testPasswordHash, 'judge', 2);
INSERT INTO user (username, password, roles, active_contest_id) VALUES ('test_student', @testPasswordHash, 'student', 3);
SELECT @student_id := id FROM user WHERE username='test_student';
INSERT INTO group_relation (user_id, group_id) VALUES (@student_id, @group_id);

INSERT INTO assignment (contest_id, uuid, title, article) VALUES (@contest_id, 'd31fc051b882448798705a5b016141af', 'A+B Problem', 'Write program that reads two integer numbers from input and writes their sum to output');
INSERT INTO assignment (contest_id, uuid, title, article) VALUES (@contest_id, '4fec83ac54424d909a5a343c3cf36aa9', 'A*B Problem', 'Write program that reads two integer numbers from input and writes their product to output');
