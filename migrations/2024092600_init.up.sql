CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

create table users (
    id serial primary key,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    first_name text not null,
    last_name text not null,
    email text unique not null,
    password_hash text not null,
    created_at timestamptz default now()
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

create table generated_images (
    id serial primary key,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    url text not null
);

CREATE TRIGGER update_generated_images_updated_at
BEFORE UPDATE ON generated_images
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
