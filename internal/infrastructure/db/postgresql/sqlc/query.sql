-- name: CreateAttendance :one
INSERT INTO attendances (student_id, schedule_id, here, counter)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAttendanceByID :one
SELECT * FROM attendances WHERE id = $1;

-- name: UpdateAttendance :one
UPDATE attendances
SET student_id = $2,
    schedule_id = $3,
    here = $4,
    counter = $5
WHERE id = $1
RETURNING *;

-- name: DeleteAttendance :exec
DELETE FROM attendances WHERE id = $1;

-- name: GetAttendanceByStudentID :many
SELECT * FROM attendances WHERE student_id = $1;

-- name: GetAttendanceByScheduleID :many
SELECT * FROM attendances WHERE schedule_id = $1;






-- name: CreateHomework :one
INSERT INTO homeworks (teacher_id, lesson_id, class_id, title, content, due_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetHomeworkByID :one
SELECT * FROM homeworks WHERE id = $1;

-- name: UpdateHomework :one
UPDATE homeworks
SET teacher_id = $2,
    lesson_id = $3,
    class_id = $4,
    title = $5,
    content = $6,
    due_date = $7
WHERE id = $1
RETURNING *;

-- name: DeleteHomework :exec
DELETE FROM homeworks WHERE id = $1;

-- name: GetAllHomeworks :many
SELECT * FROM homeworks;

-- name: GetHomeworksByTeacherID :many
SELECT * FROM homeworks WHERE teacher_id = $1;

-- name: GetHomeworksByLessonID :many
SELECT * FROM homeworks WHERE lesson_id = $1;

-- name: GetHomeworksByClassID :many
SELECT * FROM homeworks WHERE class_id = $1;



-- name: CreateLesson :one
INSERT INTO lessons (lesson_name)
VALUES ($1)
RETURNING *;

-- name: GetLessonByID :one
SELECT * FROM lessons WHERE id = $1;

-- name: UpdateLesson :one
UPDATE lessons
SET lesson_name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteLesson :exec
DELETE FROM lessons WHERE id = $1;

-- name: GetAllLessons :many
SELECT * FROM lessons;




-- name: CreateSchedule :one
INSERT INTO schedules (date, time, teacher_id, lesson_id, class_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetScheduleByID :one
SELECT * FROM schedules WHERE id = $1;

-- name: UpdateSchedule :one
UPDATE schedules
SET date = $2,
    time = $3,
    teacher_id = $4,
    lesson_id = $5,
    class_id = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedules WHERE id = $1;

-- name: GetAllSchedules :many
SELECT * FROM schedules;

-- name: GetSchedulesByTeacherID :many
SELECT * FROM schedules WHERE teacher_id = $1;

-- name: GetSchedulesByClassID :many
SELECT * FROM schedules WHERE class_id = $1;
