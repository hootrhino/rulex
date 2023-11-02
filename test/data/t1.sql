--
-- SQLiteStudio v3.4.4 生成的文件，周四 11月 2 21:33:20 2023
--
-- 所用的文本编码：System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- 表：m_ai_bases
DROP TABLE IF EXISTS m_ai_bases;

CREATE TABLE IF NOT EXISTS m_ai_bases (
    id          INTEGER,
    created_at  DATETIME,
    uuid        TEXT     NOT NULL,
    name        TEXT     NOT NULL,
    type        TEXT     NOT NULL,
    is_build_in NUMERIC  NOT NULL,
    version     TEXT     NOT NULL,
    filepath    TEXT     NOT NULL,
    description TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);


-- 表：m_apps
DROP TABLE IF EXISTS m_apps;

CREATE TABLE IF NOT EXISTS m_apps (
    id          INTEGER,
    created_at  DATETIME,
    uuid        TEXT     NOT NULL,
    name        TEXT     NOT NULL,
    version     TEXT     NOT NULL,
    auto_start  NUMERIC  NOT NULL,
    lua_source  TEXT     NOT NULL,
    description TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_apps (
                       id,
                       created_at,
                       uuid,
                       name,
                       version,
                       auto_start,
                       lua_source,
                       description
                   )
                   VALUES (
                       1,
                       '2023-11-02 21:11:00.3877632+08:00',
                       'APPEPAKFK',
                       'Hello',
                       '1.0.0',
                       1,
                       '
--
-- App use lua syntax, goto https://hootrhino.github.io for more document
-- APPID: APPEPAKFK
--
AppNAME = "Hello"
AppVERSION = "1.0.0"
AppDESCRIPTION = ""
--
-- Main
--

function Main(arg)
	applib:Debug("Hello World:" .. applib:Time())
	return 0
end

',
                       ''
                   );


-- 表：m_cron_results
DROP TABLE IF EXISTS m_cron_results;

CREATE TABLE IF NOT EXISTS m_cron_results (
    id         INTEGER,
    created_at DATETIME,
    task_id    INTEGER,
    status     TEXT,
    exit_code  TEXT,
    log_path   TEXT,
    start_time DATETIME,
    end_time   DATETIME,
    PRIMARY KEY (
        id
    )
);


-- 表：m_cron_tasks
DROP TABLE IF EXISTS m_cron_tasks;

CREATE TABLE IF NOT EXISTS m_cron_tasks (
    id         INTEGER,
    created_at DATETIME,
    name       TEXT     NOT NULL,
    cron_expr  TEXT,
    enable     TEXT,
    task_type  INTEGER,
    command    TEXT,
    args       TEXT,
    is_root    TEXT,
    work_dir   TEXT,
    env        TEXT,
    script     TEXT,
    updated_at DATETIME,
    PRIMARY KEY (
        id
    )
);


-- 表：m_data_schemas
DROP TABLE IF EXISTS m_data_schemas;

CREATE TABLE IF NOT EXISTS m_data_schemas (
    id         INTEGER,
    created_at DATETIME,
    uuid       TEXT     NOT NULL,
    name       TEXT     NOT NULL,
    type       TEXT     NOT NULL,
    schema     TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);


-- 表：m_devices
DROP TABLE IF EXISTS m_devices;

CREATE TABLE IF NOT EXISTS m_devices (
    id          INTEGER,
    created_at  DATETIME,
    uuid        TEXT     NOT NULL,
    name        TEXT     NOT NULL,
    type        TEXT     NOT NULL,
    config      TEXT,
    bind_rules  TEXT,
    description TEXT,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_devices (
                          id,
                          created_at,
                          uuid,
                          name,
                          type,
                          config,
                          bind_rules,
                          description
                      )
                      VALUES (
                          1,
                          '2023-10-26 23:14:18.7475785+08:00',
                          'DEVICE7NPKXH',
                          'GENERIC_CAMERA',
                          'GENERIC_CAMERA',
                          '{"device":"USB2.0 PC CAMERA","inputMode":"LOCAL","outputMode":"H264_STREAM","rtspUrl":""}',
                          '[]',
                          'GENERIC_CAMERA'
                      );

INSERT INTO m_devices (
                          id,
                          created_at,
                          uuid,
                          name,
                          type,
                          config,
                          bind_rules,
                          description
                      )
                      VALUES (
                          2,
                          '2023-11-02 21:10:33.4760578+08:00',
                          'DEVICEKRO89A',
                          '温湿度传感器',
                          'GENERIC_MODBUS',
                          '{"commonConfig":{"autoRequest":true,"frequency":3000,"mode":"RTU","retryTime":5,"timeout":3000,"transport":"rawserial"},"deviceConfig":{},"registers":[{"address":0,"alias":"d2","function":3,"initValue":0,"quantity":2,"slaverId":2,"tag":"d2","value":"","weight":1},{"address":0,"alias":"d1","function":3,"initValue":0,"quantity":2,"slaverId":1,"tag":"d1","value":"","weight":1}],"rtuConfig":{"baudRate":4800,"dataBits":8,"parity":"N","stopBits":1,"timeout":3000,"uart":"COM3"}}',
                          '["RULEJVVODN"]',
                          '温湿度传感器'
                      );


-- 表：m_generic_group_relations
DROP TABLE IF EXISTS m_generic_group_relations;

CREATE TABLE IF NOT EXISTS m_generic_group_relations (
    id         INTEGER,
    created_at DATETIME,
    uuid       TEXT     NOT NULL,
    gid        TEXT     NOT NULL,
    rid        TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);


-- 表：m_generic_groups
DROP TABLE IF EXISTS m_generic_groups;

CREATE TABLE IF NOT EXISTS m_generic_groups (
    id         INTEGER,
    created_at DATETIME,
    uuid       TEXT     NOT NULL,
    name       TEXT     NOT NULL,
    type       TEXT     NOT NULL,
    parent     TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_generic_groups (
                                 id,
                                 created_at,
                                 uuid,
                                 name,
                                 type,
                                 parent
                             )
                             VALUES (
                                 1,
                                 '2023-10-26 23:11:57.1575319+08:00',
                                 'VROOT',
                                 '默认分组',
                                 'VISUAL',
                                 'NULL'
                             );

INSERT INTO m_generic_groups (
                                 id,
                                 created_at,
                                 uuid,
                                 name,
                                 type,
                                 parent
                             )
                             VALUES (
                                 2,
                                 '2023-10-26 23:11:57.1612891+08:00',
                                 'DROOT',
                                 '默认分组',
                                 'DEVICE',
                                 'NULL'
                             );


-- 表：m_goods
DROP TABLE IF EXISTS m_goods;

CREATE TABLE IF NOT EXISTS m_goods (
    id           INTEGER,
    created_at   DATETIME,
    uuid         TEXT     NOT NULL,
    local_path   TEXT     NOT NULL,
    goods_type   TEXT     NOT NULL,
    execute_type TEXT     NOT NULL,
    auto_start   NUMERIC  NOT NULL,
    net_addr     TEXT     NOT NULL,
    args         TEXT     NOT NULL,
    description  TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_goods (
                        id,
                        created_at,
                        uuid,
                        local_path,
                        goods_type,
                        execute_type,
                        auto_start,
                        net_addr,
                        args,
                        description
                    )
                    VALUES (
                        1,
                        '2023-11-02 21:12:04.9252762+08:00',
                        'GOODSFV2RO7',
                        './upload/TrailerGoods/goods_1698930724917722.exe',
                        'LOCAL',
                        'EXE',
                        0,
                        '127.0.0.1:7701',
                        '-port 7701',
                        '简单的RPC示例'
                    );

INSERT INTO m_goods (
                        id,
                        created_at,
                        uuid,
                        local_path,
                        goods_type,
                        execute_type,
                        auto_start,
                        net_addr,
                        args,
                        description
                    )
                    VALUES (
                        2,
                        '2023-11-02 21:12:13.0439456+08:00',
                        'GOODS8LRD3A',
                        './upload/TrailerGoods/goods_1698930733035913.exe',
                        'LOCAL',
                        'EXE',
                        0,
                        '127.0.0.1:7702',
                        '-port 7702',
                        '简单的RPC示例2'
                    );

INSERT INTO m_goods (
                        id,
                        created_at,
                        uuid,
                        local_path,
                        goods_type,
                        execute_type,
                        auto_start,
                        net_addr,
                        args,
                        description
                    )
                    VALUES (
                        3,
                        '2023-11-02 21:12:18.5218802+08:00',
                        'GOODSFEFBF4',
                        './upload/TrailerGoods/goods_1698930738513917.exe',
                        'LOCAL',
                        'EXE',
                        0,
                        '127.0.0.1:7703',
                        '-port 7703',
                        '简单的RPC示例3'
                    );


-- 表：m_in_ends
DROP TABLE IF EXISTS m_in_ends;

CREATE TABLE IF NOT EXISTS m_in_ends (
    id            INTEGER,
    created_at    DATETIME,
    uuid          TEXT     NOT NULL,
    type          TEXT     NOT NULL,
    name          TEXT     NOT NULL,
    bind_rules    TEXT,
    description   TEXT,
    config        TEXT,
    x_data_models TEXT,
    PRIMARY KEY (
        id
    )
);


-- 表：m_ip_routes
DROP TABLE IF EXISTS m_ip_routes;

CREATE TABLE IF NOT EXISTS m_ip_routes (
    id            INTEGER,
    created_at    DATETIME,
    uuid          TEXT     NOT NULL,
    iface         TEXT     NOT NULL,
    ip            TEXT     NOT NULL,
    network       TEXT     NOT NULL,
    gateway       TEXT     NOT NULL,
    netmask       TEXT     NOT NULL,
    ip_pool_begin TEXT     NOT NULL,
    ip_pool_end   TEXT     NOT NULL,
    iface_from    TEXT     NOT NULL,
    iface_to      TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_ip_routes (
                            id,
                            created_at,
                            uuid,
                            iface,
                            ip,
                            network,
                            gateway,
                            netmask,
                            ip_pool_begin,
                            ip_pool_end,
                            iface_from,
                            iface_to
                        )
                        VALUES (
                            1,
                            '2023-10-26 23:11:57.007169+08:00',
                            '0',
                            'eth1',
                            '192.168.64.100',
                            '',
                            '192.168.64.100',
                            '255.255.255.0',
                            '192.168.64.100',
                            '192.168.64.150',
                            'eth1',
                            'usb0'
                        );


-- 表：m_modbus_point_positions
DROP TABLE IF EXISTS m_modbus_point_positions;

CREATE TABLE IF NOT EXISTS m_modbus_point_positions (
    id            INTEGER,
    created_at    DATETIME,
    device_uuid   TEXT     NOT NULL,
    tag           TEXT     NOT NULL,
    function      INTEGER  NOT NULL,
    slaver_id     INTEGER  NOT NULL,
    start_address INTEGER  NOT NULL,
    quality       INTEGER  NOT NULL,
    PRIMARY KEY (
        id
    )
);


-- 表：m_network_configs
DROP TABLE IF EXISTS m_network_configs;

CREATE TABLE IF NOT EXISTS m_network_configs (
    id           INTEGER,
    created_at   DATETIME,
    type         TEXT     NOT NULL,
    interface    TEXT     NOT NULL,
    address      TEXT     NOT NULL,
    netmask      TEXT     NOT NULL,
    gateway      TEXT     NOT NULL,
    dns          TEXT     NOT NULL,
    dhcp_enabled NUMERIC  NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_network_configs (
                                  id,
                                  created_at,
                                  type,
                                  interface,
                                  address,
                                  netmask,
                                  gateway,
                                  dns,
                                  dhcp_enabled
                              )
                              VALUES (
                                  1,
                                  '2023-10-26 23:11:56.9868179+08:00',
                                  '',
                                  'eth0',
                                  '192.168.1.100',
                                  '255.255.255.0',
                                  '192.168.1.1',
                                  '["8.8.8.8","114.114.114.114"]',
                                  1
                              );

INSERT INTO m_network_configs (
                                  id,
                                  created_at,
                                  type,
                                  interface,
                                  address,
                                  netmask,
                                  gateway,
                                  dns,
                                  dhcp_enabled
                              )
                              VALUES (
                                  2,
                                  '2023-10-26 23:11:56.9921209+08:00',
                                  '',
                                  'eth1',
                                  '192.168.64.100',
                                  '255.255.255.0',
                                  '192.168.64.1',
                                  '["8.8.8.8","114.114.114.114"]',
                                  0
                              );


-- 表：m_out_ends
DROP TABLE IF EXISTS m_out_ends;

CREATE TABLE IF NOT EXISTS m_out_ends (
    id          INTEGER,
    created_at  DATETIME,
    uuid        TEXT     NOT NULL,
    type        TEXT     NOT NULL,
    name        TEXT     NOT NULL,
    description TEXT,
    config      TEXT,
    PRIMARY KEY (
        id
    )
);


-- 表：m_protocol_apps
DROP TABLE IF EXISTS m_protocol_apps;

CREATE TABLE IF NOT EXISTS m_protocol_apps (
    id         INTEGER,
    created_at DATETIME,
    uuid       TEXT     NOT NULL,
    name       TEXT     NOT NULL,
    type       TEXT     NOT NULL,
    content    TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);


-- 表：m_rules
DROP TABLE IF EXISTS m_rules;

CREATE TABLE IF NOT EXISTS m_rules (
    id          INTEGER,
    created_at  DATETIME,
    uuid        TEXT     NOT NULL,
    name        TEXT     NOT NULL,
    type        TEXT,
    from_source TEXT,
    from_device TEXT,
    expression  TEXT     NOT NULL,
    actions     TEXT     NOT NULL,
    success     TEXT     NOT NULL,
    failed      TEXT     NOT NULL,
    description TEXT,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_rules (
                        id,
                        created_at,
                        uuid,
                        name,
                        type,
                        from_source,
                        from_device,
                        expression,
                        actions,
                        success,
                        failed,
                        description
                    )
                    VALUES (
                        1,
                        '2023-11-02 21:13:18.5380365+08:00',
                        'RULEJVVODN',
                        '日志输出而已',
                        'lua',
                        '[]',
                        '["DEVICEKRO89A"]',
                        '',
                        'Actions = {
	function(args)
		rulexlib:Debug(args)
		return true, args
	end
}',
                        'function Success()
--rulexlib:log("success")
end',
                        'function Failed(error)
rulexlib:log(error)
end',
                        ''
                    );


-- 表：m_site_configs
DROP TABLE IF EXISTS m_site_configs;

CREATE TABLE IF NOT EXISTS m_site_configs (
    id         INTEGER,
    created_at DATETIME,
    uuid       TEXT     NOT NULL,
    site_name  TEXT     NOT NULL,
    logo       TEXT     NOT NULL,
    app_name   TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_site_configs (
                               id,
                               created_at,
                               uuid,
                               site_name,
                               logo,
                               app_name
                           )
                           VALUES (
                               1,
                               '2023-10-26 23:11:57.0018409+08:00',
                               '0',
                               'RhinoEEKIT',
                               'data:image/png;base64,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',
                               'RhinoEEKIT'
                           );


-- 表：m_users
DROP TABLE IF EXISTS m_users;

CREATE TABLE IF NOT EXISTS m_users (
    id          INTEGER,
    created_at  DATETIME,
    role        TEXT     NOT NULL,
    username    TEXT     NOT NULL,
    password    TEXT     NOT NULL,
    description TEXT,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_users (
                        id,
                        created_at,
                        role,
                        username,
                        password,
                        description
                    )
                    VALUES (
                        1,
                        '2023-11-02 21:09:09.0899151+08:00',
                        'admin',
                        'hootrhino',
                        '25d55ad283aa400af464c76d713c07ad',
                        ''
                    );


-- 表：m_visuals
DROP TABLE IF EXISTS m_visuals;

CREATE TABLE IF NOT EXISTS m_visuals (
    id         INTEGER,
    created_at DATETIME,
    uuid       TEXT     NOT NULL,
    name       TEXT     NOT NULL,
    type       TEXT     NOT NULL,
    status     NUMERIC  NOT NULL,
    content    TEXT     NOT NULL,
    thumbnail  TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);


-- 表：m_wifi_configs
DROP TABLE IF EXISTS m_wifi_configs;

CREATE TABLE IF NOT EXISTS m_wifi_configs (
    id         INTEGER,
    created_at DATETIME,
    interface  TEXT     NOT NULL,
    ss_id      TEXT     NOT NULL,
    password   TEXT     NOT NULL,
    security   TEXT     NOT NULL,
    PRIMARY KEY (
        id
    )
);

INSERT INTO m_wifi_configs (
                               id,
                               created_at,
                               interface,
                               ss_id,
                               password,
                               security
                           )
                           VALUES (
                               1,
                               '2023-10-26 23:11:56.9980347+08:00',
                               'wlan0',
                               'example.net',
                               '123456',
                               'wpa2-psk'
                           );


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
