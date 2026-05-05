USE employee_management;
-- departments
CREATE TABLE departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP 
        DEFAULT CURRENT_TIMESTAMP 
        ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- employees
CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    position VARCHAR(255),
    department_id INT,
    salary DECIMAL(10,2),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP 
        DEFAULT CURRENT_TIMESTAMP 
        ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    CONSTRAINT fk_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE SET NULL
);

-- index (quan trọng cho performance)
CREATE INDEX idx_emp_name ON employees(name);
CREATE INDEX idx_emp_department ON employees(department_id);
CREATE INDEX idx_emp_deleted_at ON employees(deleted_at);

-- Insert sample data
-- INSERT INTO departments (name) VALUES ('Engineering');

-- INSERT INTO employees (name, age, position, salary, department_id) VALUES
-- ('Nguyen Van A', 28, 'Backend Engineer', 2000, 1)
