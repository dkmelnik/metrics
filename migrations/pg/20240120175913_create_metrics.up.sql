CREATE TABLE IF NOT EXISTS metrics (
   id VARCHAR(60) not null default gen_random_uuid(),
   name VARCHAR(60) NOT NULL,
   type VARCHAR(14) NOT NULL,
   delta BIGINT,
   value DOUBLE PRECISION,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP DEFAULT NOW()
)