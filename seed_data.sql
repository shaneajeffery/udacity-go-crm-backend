CREATE TABLE customers (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name text,
    role text,
    email text,
    phone text,
    contacted boolean
);

CREATE UNIQUE INDEX customers_pkey ON customers(id uuid_ops);

INSERT INTO "public"."customers"("id","name","role","email","phone","contacted")
VALUES
(E'abe92bf9-253b-411c-bae3-c660e6a0f485',E'Customer 1',E'Customer',E'customer1@customer.com',E'555-555-5555',TRUE),
(E'bc8a1c91-59f8-46e2-ad2d-a1d81201805d',E'Customer 2',E'Customer',E'customer2@customer.com',E'555-555-5556',FALSE),
(E'33a42be1-f4c0-4391-ab94-edc4d706d219',E'Customer 3',E'Customer',E'customer3@customer.com',E'555-555-5557',FALSE);