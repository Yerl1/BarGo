CREATE TABLE IF NOT EXISTS stores (
  store_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  address TEXT,
  photo TEXT,
  description TEXT,
  owner_id UUID REFERENCES users(user_id),
  coord UUID REFERENCES coords(coord_id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);