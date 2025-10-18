-- Create items_category table (parent)
CREATE TABLE items_category (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL,
    specification TEXT,
    quantity INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create items table (child)
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    year INTEGER,
    kondisi VARCHAR(100),
    note TEXT,
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES items_category(id) ON DELETE CASCADE
);

-- Create index for better performance on foreign key lookups
CREATE INDEX idx_items_category_id ON items(category_id);

-- Function to update quantity in items_category
CREATE OR REPLACE FUNCTION update_category_quantity()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Increment quantity when new item is added
        UPDATE items_category 
        SET quantity = quantity + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = NEW.category_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Decrement quantity when item is deleted
        UPDATE items_category 
        SET quantity = quantity - 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = OLD.category_id;
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Handle category change
        IF OLD.category_id != NEW.category_id THEN
            -- Decrement from old category
            UPDATE items_category 
            SET quantity = quantity - 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = OLD.category_id;
            
            -- Increment to new category
            UPDATE items_category 
            SET quantity = quantity + 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.category_id;
        END IF;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update quantity
CREATE TRIGGER trigger_update_category_quantity
    AFTER INSERT OR UPDATE OR DELETE ON items
    FOR EACH ROW
    EXECUTE FUNCTION update_category_quantity();

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at timestamp
CREATE TRIGGER trigger_items_category_updated_at
    BEFORE UPDATE ON items_category
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_items_updated_at
    BEFORE UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();