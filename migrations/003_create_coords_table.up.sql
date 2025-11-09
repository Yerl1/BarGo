CREATE TABLE IF NOT EXISTS coords (
  coord_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  entity_id UUID NOT NULL,
  entity_type VARCHAR(50) NOT NULL,
  latitude DECIMAL(10, 8) NOT NULL CHECK (
    latitude >= -90
    AND latitude <= 90
  ),
  longitude DECIMAL(11, 8) NOT NULL CHECK (
    longitude >= -180
    AND longitude <= 180
  ),
  address TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CREATE EXTENSION IF NOT EXISTS postgis;
CREATE TABLE IF NOT EXISTS coords (
  coord_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  entity_id UUID NOT NULL,
  entity_type VARCHAR(50) NOT NULL,
  latitude DECIMAL(10, 8) NOT NULL CHECK (
    latitude >= -90
    AND latitude <= 90
  ),
  longitude DECIMAL(11, 8) NOT NULL CHECK (
    longitude >= -180
    AND longitude <= 180
  ),
  address TEXT,
  geom geometry(Point, 4326) GENERATED ALWAYS AS (
    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
  ) STORED,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CURRENT_TIMESTAMP
);