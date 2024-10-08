
create or replace function update_updated_at_column()
returns trigger as $$
begin
   new.updated_at = now();
   return new;
end;
$$ language 'plpgsql';

create table users (
    id uuid primary key default gen_random_uuid(),
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    first_name text not null,
    last_name text not null,
    email text unique not null,
    password_hash text not null
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();


create table generated_images (
    id uuid primary key default gen_random_uuid(),
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    url text not null,
    user_id uuid not null,
    constraint fk_user
        foreign key(user_id) 
        references users(id)
);

CREATE TRIGGER update_generated_images_updated_at
BEFORE UPDATE ON generated_images
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
