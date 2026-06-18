-- =====================================================
-- 计件薪资核算系统 数据库脚本
-- =====================================================
CREATE DATABASE IF NOT EXISTS piece_wage DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE piece_wage;

-- =====================================================
-- 1. 用户表（工人 + 管理员）
-- =====================================================
DROP TABLE IF EXISTS sys_user;
CREATE TABLE sys_user (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    username    VARCHAR(50)  NOT NULL COMMENT '登录账号',
    password    VARCHAR(255) NOT NULL COMMENT '密码(BCrypt加密)',
    real_name   VARCHAR(50)  NOT NULL COMMENT '真实姓名',
    phone       VARCHAR(20)  DEFAULT NULL COMMENT '手机号',
    role        TINYINT      NOT NULL DEFAULT 1 COMMENT '角色:1-工人,2-核算员,9-管理员',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '状态:0-禁用,1-启用',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';

-- =====================================================
-- 2. 班组表
-- =====================================================
DROP TABLE IF EXISTS work_team;
CREATE TABLE work_team (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    team_name   VARCHAR(100) NOT NULL COMMENT '班组名称',
    team_code   VARCHAR(50)  NOT NULL COMMENT '班组编号',
    leader_id   BIGINT UNSIGNED DEFAULT NULL COMMENT '班组长ID',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '状态:0-停用,1-启用',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_team_code (team_code),
    KEY idx_leader (leader_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='班组表';

-- =====================================================
-- 3. 班组成员关联表
-- =====================================================
DROP TABLE IF EXISTS team_member;
CREATE TABLE team_member (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    team_id     BIGINT UNSIGNED NOT NULL COMMENT '班组ID',
    user_id     BIGINT UNSIGNED NOT NULL COMMENT '工人ID',
    join_date   DATE         NOT NULL COMMENT '入组日期',
    leave_date  DATE         DEFAULT NULL COMMENT '离组日期',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_team_user (team_id, user_id),
    KEY idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='班组成员表';

-- =====================================================
-- 4. 产品表
-- =====================================================
DROP TABLE IF EXISTS product;
CREATE TABLE product (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    product_code VARCHAR(50)  NOT NULL COMMENT '产品编号',
    product_name VARCHAR(100) NOT NULL COMMENT '产品名称',
    spec         VARCHAR(200) DEFAULT NULL COMMENT '规格型号',
    status       TINYINT      NOT NULL DEFAULT 1 COMMENT '状态:0-停用,1-启用',
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_product_code (product_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品表';

-- =====================================================
-- 5. 工序表
-- =====================================================
DROP TABLE IF EXISTS process_step;
CREATE TABLE process_step (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    process_code  VARCHAR(50)  NOT NULL COMMENT '工序编号',
    process_name  VARCHAR(100) NOT NULL COMMENT '工序名称',
    product_id    BIGINT UNSIGNED NOT NULL COMMENT '所属产品ID',
    difficulty    TINYINT      NOT NULL DEFAULT 1 COMMENT '难度等级:1-简单,2-普通,3-复杂,4-高难度',
    description   VARCHAR(500) DEFAULT NULL COMMENT '工序描述',
    is_shared     TINYINT      NOT NULL DEFAULT 0 COMMENT '是否班组共享工序:0-否,1-是',
    status        TINYINT      NOT NULL DEFAULT 1 COMMENT '状态:0-停用,1-启用',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_process_code (process_code),
    KEY idx_product (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工序表';

-- =====================================================
-- 6. 工序单价版本表（核心：单价带版本，按生效时间管理）
-- =====================================================
DROP TABLE IF EXISTS process_price;
CREATE TABLE process_price (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    process_id    BIGINT UNSIGNED NOT NULL COMMENT '工序ID',
    version_no    INT          NOT NULL COMMENT '版本号',
    grade_level   VARCHAR(20)  NOT NULL DEFAULT 'STD' COMMENT '等级:STD-标准,PRE-特级,ADV-高级',
    unit_price    DECIMAL(12,4) NOT NULL COMMENT '单件工价(元)',
    effective_date DATE        NOT NULL COMMENT '生效日期',
    expiry_date   DATE         DEFAULT NULL COMMENT '失效日期(NULL表示永不过期)',
    remark        VARCHAR(500) DEFAULT NULL COMMENT '调整说明',
    created_by    BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_process_version (process_id, version_no, grade_level),
    KEY idx_effective (process_id, effective_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工序单价版本表';

-- =====================================================
-- 7. 报工单表（核心：记录报工当时的单价快照，避免后续调价影响历史数据）
-- =====================================================
DROP TABLE IF EXISTS production_report;
CREATE TABLE production_report (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    report_no      VARCHAR(32)  NOT NULL COMMENT '报工单号',
    worker_id      BIGINT UNSIGNED NOT NULL COMMENT '报工工人ID',
    team_id        BIGINT UNSIGNED DEFAULT NULL COMMENT '所属班组ID(共享工序时使用)',
    process_id     BIGINT UNSIGNED NOT NULL COMMENT '工序ID',
    price_id       BIGINT UNSIGNED NOT NULL COMMENT '报工当时生效的单价版本ID(快照)',
    unit_price     DECIMAL(12,4) NOT NULL COMMENT '报工当时单价快照(元)',
    grade_level    VARCHAR(20)  NOT NULL DEFAULT 'STD' COMMENT '报工等级',
    report_date    DATE         NOT NULL COMMENT '报工日期',
    qty_good       INT          NOT NULL DEFAULT 0 COMMENT '合格数量',
    qty_defect     INT          NOT NULL DEFAULT 0 COMMENT '不良/返工数量(扣减)',
    qty_total      INT          NOT NULL DEFAULT 0 COMMENT '总完工数量(=合格+不良)',
    unit_defect    DECIMAL(12,4) NOT NULL DEFAULT 0.5 COMMENT '不良品扣款比例(0.5=不良按半价扣)',
    gross_amount   DECIMAL(14,2) NOT NULL DEFAULT 0 COMMENT '计件总额(合格数量×单价)',
    defect_amount  DECIMAL(14,2) NOT NULL DEFAULT 0 COMMENT '不良扣款(不良数量×单价×扣款比例)',
    net_amount     DECIMAL(14,2) NOT NULL DEFAULT 0 COMMENT '实发金额(=总额-扣款)',
    work_order_no  VARCHAR(50)  DEFAULT NULL COMMENT '关联生产工单',
    remark         VARCHAR(500) DEFAULT NULL COMMENT '备注',
    status         TINYINT      NOT NULL DEFAULT 1 COMMENT '状态:0-作废,1-已生效',
    created_by     BIGINT UNSIGNED NOT NULL COMMENT '录入人ID',
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_report_no (report_no),
    KEY idx_worker_date (worker_id, report_date),
    KEY idx_team_date (team_id, report_date),
    KEY idx_process (process_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报工单表';

-- =====================================================
-- 8. 班组工资分配明细表（共享工序按比例分配）
-- =====================================================
DROP TABLE IF EXISTS team_wage_allocation;
CREATE TABLE team_wage_allocation (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    report_id       BIGINT UNSIGNED NOT NULL COMMENT '报工单ID',
    team_id         BIGINT UNSIGNED NOT NULL COMMENT '班组ID',
    total_amount    DECIMAL(14,2) NOT NULL COMMENT '待分配总金额',
    allocation_rule TINYINT      NOT NULL DEFAULT 1 COMMENT '分配规则:1-平均分,2-按比例分,3-按工时权重分',
    created_by      BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_report (report_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='班组工资分配主表';

-- =====================================================
-- 9. 班组工资分配人员明细表
-- =====================================================
DROP TABLE IF EXISTS team_wage_allocation_item;
CREATE TABLE team_wage_allocation_item (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    allocation_id  BIGINT UNSIGNED NOT NULL COMMENT '分配主表ID',
    worker_id      BIGINT UNSIGNED NOT NULL COMMENT '工人ID',
    weight_ratio   DECIMAL(8,4)  NOT NULL COMMENT '分配权重/比例',
    work_hours     DECIMAL(8,2)  DEFAULT NULL COMMENT '工时(按工时分配时使用)',
    allocated_amt  DECIMAL(14,2) NOT NULL COMMENT '分配金额',
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_allocation (allocation_id),
    KEY idx_worker (worker_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='班组工资分配人员明细表';

-- =====================================================
-- 10. 工人计件工资明细表（按报工单汇总+班组分配金额）
-- =====================================================
DROP TABLE IF EXISTS worker_wage_detail;
CREATE TABLE worker_wage_detail (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    worker_id      BIGINT UNSIGNED NOT NULL COMMENT '工人ID',
    report_id      BIGINT UNSIGNED DEFAULT NULL COMMENT '关联报工单ID(直接报工)',
    allocation_id  BIGINT UNSIGNED DEFAULT NULL COMMENT '关联分配明细ID(班组分配)',
    detail_type    TINYINT      NOT NULL COMMENT '类型:1-直接计件,2-班组分配,3-其他调整',
    process_id     BIGINT UNSIGNED DEFAULT NULL COMMENT '工序ID',
    wage_date      DATE         NOT NULL COMMENT '工资归属日期',
    qty_good       INT          NOT NULL DEFAULT 0 COMMENT '合格数量',
    qty_defect     INT          NOT NULL DEFAULT 0 COMMENT '不良数量',
    unit_price     DECIMAL(12,4) NOT NULL DEFAULT 0 COMMENT '单价',
    amount         DECIMAL(14,2) NOT NULL COMMENT '金额(正负)',
    remark         VARCHAR(500) DEFAULT NULL COMMENT '备注',
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_worker_date (worker_id, wage_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工人工资明细表';

-- =====================================================
-- 11. 工人月度工资汇总表
-- =====================================================
DROP TABLE IF EXISTS worker_wage_summary;
CREATE TABLE worker_wage_summary (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    worker_id       BIGINT UNSIGNED NOT NULL COMMENT '工人ID',
    summary_month   VARCHAR(7)     NOT NULL COMMENT '汇总月份(YYYY-MM)',
    total_qty_good  INT            NOT NULL DEFAULT 0 COMMENT '当月总合格数量',
    total_qty_defect INT           NOT NULL DEFAULT 0 COMMENT '当月总不良数量',
    gross_amount    DECIMAL(14,2)  NOT NULL DEFAULT 0 COMMENT '计件总额',
    defect_amount   DECIMAL(14,2)  NOT NULL DEFAULT 0 COMMENT '不良扣款总额',
    allocation_amt  DECIMAL(14,2)  NOT NULL DEFAULT 0 COMMENT '班组分配金额',
    adjust_amount   DECIMAL(14,2)  NOT NULL DEFAULT 0 COMMENT '其他调整',
    net_amount      DECIMAL(14,2)  NOT NULL DEFAULT 0 COMMENT '应发合计',
    calc_status     TINYINT        NOT NULL DEFAULT 0 COMMENT '计算状态:0-未结算,1-已结算',
    last_calc_time  DATETIME       DEFAULT NULL COMMENT '最后计算时间',
    created_at      DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at      DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_worker_month (worker_id, summary_month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工人月度工资汇总表';

-- =====================================================
-- 初始化基础数据
-- =====================================================
-- 默认管理员: admin / 123456 (BCrypt哈希)
INSERT INTO sys_user (username, password, real_name, phone, role, status) VALUES
('admin', '$2a$10$IwGYKHak313aoht4KrDgPuw70H7xqwx.1DkUKA5YD5qf0WbznW1pq', '系统管理员', '13800000000', 9, 1),
('accountant01', '$2a$10$IwGYKHak313aoht4KrDgPuw70H7xqwx.1DkUKA5YD5qf0WbznW1pq', '张核算员', '13800000001', 2, 1),
('worker01', '$2a$10$IwGYKHak313aoht4KrDgPuw70H7xqwx.1DkUKA5YD5qf0WbznW1pq', '李工人', '13800000011', 1, 1),
('worker02', '$2a$10$IwGYKHak313aoht4KrDgPuw70H7xqwx.1DkUKA5YD5qf0WbznW1pq', '王工人', '13800000012', 1, 1),
('worker03', '$2a$10$IwGYKHak313aoht4KrDgPuw70H7xqwx.1DkUKA5YD5qf0WbznW1pq', '赵工人', '13800000013', 1, 1),
('worker04', '$2a$10$IwGYKHak313aoht4KrDgPuw70H7xqwx.1DkUKA5YD5qf0WbznW1pq', '钱工人', '13800000014', 1, 1);

-- 班组
INSERT INTO work_team (team_name, team_code, leader_id) VALUES
('冲压一组', 'TM001', 3),
('装配一组', 'TM002', 4);

-- 班组成员
INSERT INTO team_member (team_id, user_id, join_date) VALUES
(1, 3, '2024-01-01'),
(1, 5, '2024-01-15'),
(2, 4, '2024-01-01'),
(2, 6, '2024-02-01');

-- 产品
INSERT INTO product (product_code, product_name, spec, status) VALUES
('P-1001', '精密齿轮A', '直径50mm', 1),
('P-1002', '精密齿轮B', '直径80mm', 1),
('P-2001', '电机外壳', '铝合金', 1);

-- 工序
INSERT INTO process_step (process_code, process_name, product_id, difficulty, description, is_shared) VALUES
('PRC-1001-01', '落料', 1, 1, '按图纸下料', 0),
('PRC-1001-02', '粗车', 1, 2, '车床粗加工', 0),
('PRC-1001-03', '精车', 1, 3, '车床精加工', 0),
('PRC-1001-04', '热处理', 1, 2, '淬火回火', 1),
('PRC-1002-01', '锻造', 2, 2, '毛坯锻造', 1),
('PRC-1002-02', '滚齿', 2, 3, '滚齿机加工', 0),
('PRC-2001-01', '压铸', 3, 2, '压铸成型', 1),
('PRC-2001-02', '去毛刺', 3, 1, '人工去毛刺', 0);

-- 工序单价版本（注意：同一个工序有多版本，按生效日期区分）
INSERT INTO process_price (process_id, version_no, grade_level, unit_price, effective_date, expiry_date, remark, created_by) VALUES
(1, 1, 'STD', 0.8000, '2024-01-01', '2024-03-31', '初始定价', 1),
(1, 2, 'STD', 0.8500, '2024-04-01', NULL, '原材料上涨调整+5分', 1),
(2, 1, 'STD', 1.5000, '2024-01-01', NULL, '初始定价', 1),
(2, 1, 'PRE', 1.8000, '2024-01-01', NULL, '特级工价', 1),
(3, 1, 'STD', 3.2000, '2024-01-01', NULL, '初始定价', 1),
(3, 1, 'ADV', 3.8000, '2024-01-01', NULL, '高级工价', 1),
(4, 1, 'STD', 0.6000, '2024-01-01', NULL, '热处理班组共享', 1),
(5, 1, 'STD', 5.0000, '2024-01-01', NULL, '锻造班组共享', 1),
(6, 1, 'STD', 2.5000, '2024-01-01', NULL, '初始定价', 1),
(7, 1, 'STD', 4.0000, '2024-01-01', NULL, '压铸班组共享', 1),
(8, 1, 'STD', 0.5000, '2024-01-01', NULL, '初始定价', 1);

-- 报工示例（2024年1月的报工，此时落料工序使用版本1的单价0.80）
INSERT INTO production_report (report_no, worker_id, team_id, process_id, price_id, unit_price, grade_level,
    report_date, qty_good, qty_defect, qty_total, unit_defect,
    gross_amount, defect_amount, net_amount, work_order_no, remark, status, created_by) VALUES
('RPT20240115001', 3, NULL, 1, 1, 0.8000, 'STD', '2024-01-15', 1000, 5, 1005, 0.5,
 800.00, 2.00, 798.00, 'WO202401001', '正常报工', 1, 2),
('RPT20240115002', 3, NULL, 2, 3, 1.5000, 'STD', '2024-01-15', 500, 2, 502, 0.5,
 750.00, 1.50, 748.50, 'WO202401001', NULL, 1, 2),
('RPT20240115003', 4, NULL, 3, 5, 3.2000, 'STD', '2024-01-15', 300, 1, 301, 0.5,
 960.00, 1.60, 958.40, 'WO202401002', NULL, 1, 2),
('RPT20240410001', 3, NULL, 1, 2, 0.8500, 'STD', '2024-04-10', 800, 0, 800, 0.5,
 680.00, 0.00, 680.00, 'WO202404001', '新单价生效', 1, 2);
