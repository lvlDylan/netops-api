CREATE TABLE IF NOT EXISTS device (
    id VARCHAR(36) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    os VARCHAR(50) NOT NULL,
    platform VARCHAR(255) NOT NULL,
    arch VARCHAR(20) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS history (
    id BIGINT NOT NULL AUTO_INCREMENT,
    send_at DATETIME(6) NOT NULL,
    device_id VARCHAR(36) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    agent_version VARCHAR(20) NOT NULL,
    uptime_seconds BIGINT NOT NULL,
    mem_usage_percent DOUBLE NOT NULL,
    mem_used_bytes BIGINT NOT NULL,
    mem_total_bytes BIGINT NOT NULL,
    cpu_usage_percent DOUBLE NOT NULL,
    cpu_cores INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (device_id) REFERENCES device(id) ON DELETE CASCADE,
    INDEX idx_device_send_at (device_id, send_at)
);

CREATE TABLE IF NOT EXISTS disk (
    id BIGINT NOT NULL AUTO_INCREMENT,
    history_id BIGINT NOT NULL,
    mountpoint VARCHAR(255) NOT NULL,
    usage_percent DOUBLE NOT NULL,
    used_bytes BIGINT NOT NULL,
    total_bytes BIGINT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (history_id) REFERENCES history(id) ON DELETE CASCADE
);
