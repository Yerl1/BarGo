DO $$ BEGIN IF NOT EXISTS (
  SELECT 1
  FROM pg_type
  WHERE typname = 'subscription_plan'
) THEN CREATE TYPE subscription_plan AS ENUM ('none', 'basic', 'premium');
END IF;
IF NOT EXISTS (
  SELECT 1
  FROM pg_type
  WHERE typname = 'roles'
) THEN CREATE TYPE roles AS ENUM ('user', 'admin');
END IF;
END $$;