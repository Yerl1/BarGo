CREATE TABLE IF NOT EXISTS comments (
  comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(user_id),
  store_id UUID REFERENCES stores(store_id),
  content TEXT NOT NULL,
  rating INT CHECK (
    rating >= 1
    AND rating <= 5
  ),
  helpful_votes INT DEFAULT 0,
  is_toxic BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);