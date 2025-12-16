DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'semester_category') THEN
        CREATE TYPE semester_category AS ENUM ('Ganjil', 'Genap');
    END IF;
END $$;

-- Create schedule_from_ticket table
CREATE TABLE schedule_ticket (
    id_schedule SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    kategori ticket_category NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create semester_unblocking table
CREATE TABLE unblocking (
    id SERIAL PRIMARY KEY,
    tahun INTEGER NOT NULL,
    semester semester_category NOT NULL, -- Ganjil/Genap
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    start_date TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    end_date TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create schedule_reguler table
CREATE TABLE schedule_reguler (
    id_schedule SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_schedule_from_ticket_ticket_id ON schedule_ticket(ticket_id);
CREATE INDEX idx_schedule_from_ticket_user_id ON schedule_ticket(user_id);
CREATE INDEX idx_schedule_from_ticket_start_date ON schedule_ticket(start_date);
CREATE INDEX idx_schedule_from_ticket_end_date ON schedule_ticket(end_date);

CREATE INDEX idx_semester_unblocking_user_id ON unblocking(user_id);
CREATE INDEX idx_semester_unblocking_tahun ON unblocking(tahun);
CREATE INDEX idx_semester_unblocking_semester ON unblocking(semester);

CREATE INDEX idx_schedule_reguler_user_id ON schedule_reguler(user_id);
CREATE INDEX idx_schedule_reguler_start_date ON schedule_reguler(start_date);
CREATE INDEX idx_schedule_reguler_end_date ON schedule_reguler(end_date);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_schedule_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at timestamp
CREATE TRIGGER trigger_schedule_from_ticket_updated_at
    BEFORE UPDATE ON schedule_ticket
    FOR EACH ROW
    EXECUTE FUNCTION update_schedule_updated_at_column();

CREATE TRIGGER trigger_semester_unblocking_updated_at
    BEFORE UPDATE ON unblocking
    FOR EACH ROW
    EXECUTE FUNCTION update_schedule_updated_at_column();

CREATE TRIGGER trigger_schedule_reguler_updated_at
    BEFORE UPDATE ON schedule_reguler
    FOR EACH ROW
    EXECUTE FUNCTION update_schedule_updated_at_column();