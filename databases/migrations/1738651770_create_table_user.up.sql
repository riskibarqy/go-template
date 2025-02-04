CREATE TABLE public."user"
(
    "id" SERIAL NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(100) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "token" TEXT,
    "token_expired_at" INT NULL,
    "created_at" INT NOT NULL,
    "updated_at" INT NOT NULL,
    "deleted_at" INT,
    CONSTRAINT user_pkey PRIMARY KEY ("id")
);

-- Add an index on email for faster lookups
CREATE INDEX user_email_idx ON public."user"("email");

-- Add an index on deleted_at to optimize soft delete queries
CREATE INDEX user_deleted_at_idx ON public."user"("deleted_at");

