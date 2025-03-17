-- Table creations --
-- Make sure to create table under the same user that will query it
-- This avoids further complexion of DB privilege management

CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    rate INT,
    model VARCHAR(255)
);


-- Seeds for testing purposes --

INSERT INTO devices(rate, model) VALUES(10000, 'Mid-end PC');
