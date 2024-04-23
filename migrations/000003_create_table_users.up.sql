--Inserting 10 sample users
INSERT INTO users (id, first_name, last_name, birth_date, phone_number, password, gender, refresh_token, created_at)
VALUES
    ('123e4567-e89b-12d3-a456-426614174001', 'John', 'Doe', '1990-05-15', '1234567890', 'password123', 'male', 'random_refresh_token_1', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174002', 'Jane', 'Doe', '1992-08-20', '2345678901', 'password456', 'female', 'random_refresh_token_2', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174003', 'Alice', 'Smith', '1985-03-10', '3456789012', 'password789', 'female', 'random_refresh_token_3', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174004', 'Bob', 'Johnson', '1988-11-25', '4567890123', 'passwordabc', 'male', 'random_refresh_token_4', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174005', 'Emily', 'Brown', '1995-07-05', '5678901234', 'passworddef', 'female', 'random_refresh_token_5', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174006', 'Michael', 'Wilson', '1983-09-30', '6789012345', 'passwordghi', 'male', 'random_refresh_token_6', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174007', 'Sarah', 'Martinez', '1993-01-18', '7890123456', 'passwordjkl', 'female', 'random_refresh_token_7', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174008', 'David', 'Taylor', '1980-12-08', '8901234567', 'passwordmno', 'male', 'random_refresh_token_8', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174009', 'Jennifer', 'Lopez', '1977-06-22', '9012345678', 'passwordpqr', 'female', 'random_refresh_token_9', CURRENT_TIMESTAMP),
    ('123e4567-e89b-12d3-a456-426614174010', 'Christopher', 'Lee', '1970-04-12', '0123456789', 'passwordstu', 'male', 'random_refresh_token_10', CURRENT_TIMESTAMP);
