-- lessons
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_name VARCHAR(255) NOT NULL
);

-- schedules
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    time TIME NOT NULL,
    teacher_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    class_id UUID NOT NULL,
    CONSTRAINT fk_lesson FOREIGN KEY(lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);

-- attendances
CREATE TABLE attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    schedule_id UUID NOT NULL,
    here BOOLEAN NOT NULL DEFAULT FALSE,
    counter INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_schedule FOREIGN KEY(schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);

-- homeworks
CREATE TABLE homeworks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    class_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    due_date TIMESTAMP NOT NULL,
    CONSTRAINT fk_lesson FOREIGN KEY(lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);


