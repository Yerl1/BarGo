-- 007_insert_astana_stores.up.sql
-- Make sure PostGIS is enabled
CREATE EXTENSION IF NOT EXISTS postgis;
-- Insert coordinates (now with geom column auto-generated)
INSERT INTO coords (
    coord_id,
    entity_id,
    entity_type,
    latitude,
    longitude,
    address
  )
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'aaaaaaaa-0000-0000-0000-000000000001',
    'store',
    51.1280,
    71.4305,
    'Kabanbay Batyr Ave 20'
  ),
  (
    '22222222-2222-2222-2222-222222222222',
    'aaaaaaaa-0000-0000-0000-000000000002',
    'store',
    51.1338,
    71.4220,
    'Tauelsizdik Ave 34'
  ),
  (
    '33333333-3333-3333-3333-333333333333',
    'aaaaaaaa-0000-0000-0000-000000000003',
    'store',
    51.1650,
    71.4195,
    'Qabanbay Batyr Ave 62'
  ),
  (
    '44444444-4444-4444-4444-444444444444',
    'aaaaaaaa-0000-0000-0000-000000000004',
    'store',
    51.1605,
    71.4704,
    'Syganak St 37'
  ),
  (
    '55555555-5555-5555-5555-555555555555',
    'aaaaaaaa-0000-0000-0000-000000000005',
    'store',
    51.1460,
    71.4325,
    'Turan Ave 18'
  ),
  (
    '66666666-6666-6666-6666-666666666666',
    'aaaaaaaa-0000-0000-0000-000000000006',
    'store',
    51.1692,
    71.4033,
    'Baitursynov St 15'
  ),
  (
    '77777777-7777-7777-7777-777777777777',
    'aaaaaaaa-0000-0000-0000-000000000007',
    'store',
    51.1565,
    71.4400,
    'Zhenis Ave 65'
  ),
  (
    '88888888-8888-8888-8888-888888888888',
    'aaaaaaaa-0000-0000-0000-000000000008',
    'store',
    51.1711,
    71.4515,
    'Eurasia Ave 12A'
  ),
  (
    '99999999-9999-9999-9999-999999999999',
    'aaaaaaaa-0000-0000-0000-000000000009',
    'store',
    51.1508,
    71.4700,
    'Abai Ave 10'
  ),
  (
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'aaaaaaaa-0000-0000-0000-000000000010',
    'store',
    51.1635,
    71.4750,
    'Mangilik El Ave 25'
  );
-- Insert stores referencing those coords
INSERT INTO stores (store_id, name, address, coord)
VALUES (
    'aaaaaaaa-0000-0000-0000-000000000001',
    'Green Market',
    'Kabanbay Batyr Ave 20',
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000002',
    'Astana Mall',
    'Tauelsizdik Ave 34',
    '22222222-2222-2222-2222-222222222222'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000003',
    'Mega Silk Way',
    'Qabanbay Batyr Ave 62',
    '33333333-3333-3333-3333-333333333333'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000004',
    'Keruen City',
    'Syganak St 37',
    '44444444-4444-4444-4444-444444444444'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000005',
    'Aru Supermarket',
    'Turan Ave 18',
    '55555555-5555-5555-5555-555555555555'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000006',
    'Nomad Market',
    'Baitursynov St 15',
    '66666666-6666-6666-6666-666666666666'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000007',
    'Asia Park',
    'Zhenis Ave 65',
    '77777777-7777-7777-7777-777777777777'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000008',
    'Small Store 24/7',
    'Eurasia Ave 12A',
    '88888888-8888-8888-8888-888888888888'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000009',
    'Dina Grocery',
    'Abai Ave 10',
    '99999999-9999-9999-9999-999999999999'
  ),
  (
    'aaaaaaaa-0000-0000-0000-000000000010',
    'KazMart',
    'Mangilik El Ave 25',
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'
  );