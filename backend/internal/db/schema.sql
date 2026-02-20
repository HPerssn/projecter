CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE routes (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       text,
    grade      text,
    photo_url  text,
    created_at timestamptz DEFAULT now()
);

CREATE TABLE holds (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id   uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    x_pct      float NOT NULL,
    y_pct      float NOT NULL,
    note       text,
    seq_order  int
);
