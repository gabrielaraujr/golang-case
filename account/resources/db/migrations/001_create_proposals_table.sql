CREATE TABLE IF NOT EXISTS proposals (
    id UUID PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    salary DECIMAL(10, 2) NOT NULL,
    birthdate DATE NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address_street VARCHAR(255) NOT NULL,
    address_city VARCHAR(100) NOT NULL,
    address_state VARCHAR(100) NOT NULL,
    address_zip VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    CONSTRAINT unique_cpf UNIQUE (cpf)
);

CREATE INDEX idx_proposals_cpf ON proposals(cpf);
CREATE INDEX idx_proposals_status ON proposals(status);
