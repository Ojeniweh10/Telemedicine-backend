-- USERS TABLE
CREATE TABLE users (
    usertag VARCHAR(50) PRIMARY KEY,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    phone_no VARCHAR(20),
    gender VARCHAR(10),
    date_of_birth DATE,
    password TEXT,
    transaction_pin VARCHAR(6),
    otp VARCHAR(10),
    otp_expiry TIMESTAMP,
    state VARCHAR(100),
    delivery_address TEXT,
    profile_pic_url TEXT
);

-- HOSPITALS TABLE
CREATE TABLE hospitals (
    hospital_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address TEXT,
    country VARCHAR(100),
    state VARCHAR(100),
    profile_pic_url TEXT,
    about TEXT
);

-- DOCTORS TABLE
CREATE TABLE doctors (
    doctortag VARCHAR(50) PRIMARY KEY,
    fullname VARCHAR(200),
    date_of_birth DATE,
    phone_number VARCHAR(20),
    gender VARCHAR(10),
    specialization VARCHAR(100),
    country VARCHAR(100),
    city VARCHAR(100),
    yrs_of_experience INTEGER,
    price_per_session NUMERIC(10, 2),
    about TEXT,
    password TEXT NOT NULL,
    hospital_id INTEGER,
    availability JSONB, -- e.g. ["2025-08-01T10:00:00", "2025-08-02T14:00:00"]
    profile_pic_url TEXT,
    FOREIGN KEY (hospital_id) REFERENCES hospitals(hospital_id) ON DELETE SET NULL
);

-- INVENTORY (MEDICATIONS) TABLE
CREATE TABLE inventory (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    milligram VARCHAR(50),
    price NUMERIC(10, 2),
    quantity INTEGER DEFAULT 0,
    product_image_url TEXT
);

-- PHARMACY TABLE
CREATE TABLE pharmacies (
    pharmacy_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address TEXT,
    country VARCHAR(100),
    state VARCHAR(100),
    about TEXT,
    pharmacy_picture_url TEXT
);

-- ORDERS TABLE
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    usertag VARCHAR(50),
    total NUMERIC(10, 2),
    status VARCHAR(20) CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled')),
    items JSONB, -- Stores array of {product_id, item_name, price, quantity}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE
);

-- TEST CENTRE TABLE
CREATE TABLE test_centres (
    center_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    address TEXT,
    country VARCHAR(100),
    state VARCHAR(100),
    daily_capacity INTEGER,
    about TEXT,
    availability JSONB, -- e.g. [{"date": "2025-08-01", "slots": ["10:00", "11:00"]}]
    test_types TEXT[], -- array of test types
    price_per_test NUMERIC(10, 2)
);

-- PRESCRIPTION TABLE
CREATE TABLE prescriptions (
    id SERIAL PRIMARY KEY,
    prescription TEXT,
    doctor_notes TEXT,
    usertag VARCHAR(50),
    doctortag VARCHAR(50),
    prescription_date DATE,
    doctor_note_date DATE,
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE,
    FOREIGN KEY (doctortag) REFERENCES doctors(doctortag) ON DELETE CASCADE
);

-- REVIEWS TABLE
CREATE TABLE reviews (
    review_id SERIAL PRIMARY KEY,
    usertag VARCHAR(50),
    doctortag VARCHAR(50),
    review TEXT,
    star_rating INTEGER CHECK (star_rating BETWEEN 1 AND 5),
    status VARCHAR(20) CHECK (status IN ('approved', 'pending')),
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE,
    FOREIGN KEY (doctortag) REFERENCES doctors(doctortag) ON DELETE CASCADE
);

-- ADMINS TABLE
CREATE TABLE admins (
    admin_id SERIAL PRIMARY KEY,
    admintag VARCHAR(50) UNIQUE NOT NULL,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    password TEXT NOT NULL,
    otp VARCHAR(10),
    otp_expiry TIMESTAMP,
    profile_pic_url TEXT,
    role VARCHAR(50)
);

-- APPOINTMENTS TABLE
CREATE TABLE appointments (
    appointment_id SERIAL PRIMARY KEY,
    patient_tag VARCHAR(50),
    doctor_tag VARCHAR(50),
    scheduled_at TIMESTAMP,
    reason TEXT,
    file_url TEXT,
    status VARCHAR(20) CHECK (status IN ('pending', 'confirmed', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_tag) REFERENCES users(usertag) ON DELETE CASCADE,
    FOREIGN KEY (doctor_tag) REFERENCES doctors(doctortag) ON DELETE CASCADE
);

--CARTS
CREATE TABLE carts (
    cart_id SERIAL PRIMARY KEY,
    usertag VARCHAR(50),
    product_id INTEGER,
    quantity INTEGER,
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES inventory(product_id) ON DELETE RESTRICT
    ADD CONSTRAINT unique_usertag_product UNIQUE (usertag, product_id)
);

--BILLING DETAILS
CREATE TABLE billing_details (
    usertag VARCHAR(50) PRIMARY KEY,
    fullname VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    phone_no VARCHAR(20),
    state VARCHAR(100),
    delivery_address TEXT
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE
);

--wallet system
CREATE TABLE wallets (
    usertag VARCHAR(50) PRIMARY KEY,
    balance NUMERIC(10, 2) DEFAULT 0.00 NOT NULL CHECK (balance >= 0),
    pending_balance FLOAT DEFAULT 0.00 CHECK (pending_balance >= 0),
    wallet_status VARCHAR(20) CHECK (wallet_status IN ('active', 'inactive')) DEFAULT 'active',
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE
);

CREATE TABLE wallet_transactions (
    transaction_id SERIAL PRIMARY KEY,
    usertag VARCHAR(50),
    amount NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
    transaction_type VARCHAR(20) CHECK (transaction_type IN ('credit', 'debit')),
    transaction_reference VARCHAR(100) UNIQUE,
    paystack_reference VARCHAR(100) UNIQUE,
    transfer_code VARCHAR(100),
    access_code VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('pending', 'success', 'failed', 'reversed', 'disputed' )),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE
);

CREATE TABLE payout_accounts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  usertag VARCHAR(50) NOT NULL,
  recipient_code VARCHAR(50) UNIQUE NOT NULL,
  account_number VARCHAR(20) NOT NULL,       
  bank_code VARCHAR(10) NOT NULL,
  bank_name VARCHAR(100) NOT NULL, 
  account_name VARCHAR(100) NOT NULL,
  currency VARCHAR(3) DEFAULT 'NGN',
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (usertag) REFERENCES users(usertag) ON DELETE CASCADE
);