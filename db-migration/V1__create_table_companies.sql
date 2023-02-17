CREATE TABLE companies (
    id UUID NOT NULL,
    name VARCHAR(15) NOT NULL UNIQUE,
    description TEXT,
    amount_of_employees INTEGER NOT NULL,
    registered BOOLEAN NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship')),
    PRIMARY KEY (id)
);