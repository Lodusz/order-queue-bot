-- Создаем main table orders

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для поиска заказов по пользователю 
CREATE INDEX idx_orders_user_id ON orders(user_id);


CREATE INDEX idx_orders_status ON orders(status) 
WHERE status IN ('new', 'in_progress');