CREATE TABLE attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,     -- Keycloak user id
    schedule_id UUID NOT NULL,    -- Schedule tablosu ile bağlantı
    here BOOLEAN NOT NULL DEFAULT FALSE,
    counter INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_schedule FOREIGN KEY(schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
);


CREATE TABLE homeworks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    teacher_id UUID NOT NULL,       -- Keycloak teacher user ID
    lesson_id UUID NOT NULL,        -- Lesson tablosu ile bağlantı
    class_id UUID NOT NULL,         -- Class tablosu ile bağlantı
    title VARCHAR(255) NOT NULL,
    content TEXT,
    due_date TIMESTAMP NOT NULL,
    CONSTRAINT fk_lesson FOREIGN KEY(lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT fk_class FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);



CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_name VARCHAR(255) NOT NULL
);



CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL,
    time TIME NOT NULL,
    teacher_id UUID NOT NULL,      -- Keycloak teacher user ID
    lesson_id UUID NOT NULL,       -- Lesson tablosu ile bağlantı
    class_id UUID NOT NULL,        -- Class tablosu ile bağlantı
    CONSTRAINT fk_lesson FOREIGN KEY(lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT fk_class FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);
