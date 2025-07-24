CREATE DATABASE miniRustpbxgo;

CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，自增主键',
                                     username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，唯一',
                                     password VARCHAR(255) NOT NULL COMMENT '密码（建议存储加密后的密码）',
                                     email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱，唯一',
                                     phone VARCHAR(20) UNIQUE COMMENT '手机号，可选，唯一',
                                     nickname VARCHAR(50) COMMENT '昵称',
                                     status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
                                     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS robotKeys (
                                    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '密钥ID，自增主键',
                                    user_id INT NOT NULL COMMENT '关联的用户ID',
                                    name VARCHAR(100) COMMENT '密钥名称',

    -- 大模型（LLM）配置
                                    llm_provider VARCHAR(100) COMMENT '大模型提供商',
                                    llm_api_key VARCHAR(255) COMMENT '大模型API密钥',
                                    llm_api_url VARCHAR(255) COMMENT '大模型API地址',

    -- 语音识别（ASR）配置
                                    asr_provider VARCHAR(100) COMMENT '语音识别提供商',
                                    asr_app_id VARCHAR(100) COMMENT '语音识别App ID',
                                    asr_secret_id VARCHAR(255) COMMENT '语音识别Secret ID',
                                    asr_secret_key VARCHAR(255) COMMENT '语音识别Secret Key',
                                    asr_language VARCHAR(20) DEFAULT 'zh' COMMENT '语音识别语言',

    -- 语音合成（TTS）配置
                                    tts_provider VARCHAR(100) COMMENT '语音合成提供商',
                                    tts_app_id VARCHAR(100) COMMENT '语音合成App ID',
                                    tts_secret_id VARCHAR(255) COMMENT '语音合成Secret ID',
                                    tts_secret_key VARCHAR(255) COMMENT '语音合成Secret Key',

    -- 新增的API密钥字段
                                    api_key VARCHAR(255) COMMENT 'API密钥',
                                    api_secret VARCHAR(255) COMMENT 'API密钥的Secret',

                                    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 外键约束，关联到users表的id字段
                                    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='机器人密钥配置表';

CREATE TABLE IF NOT EXISTS robots (
                                      id INT AUTO_INCREMENT PRIMARY KEY COMMENT '机器人配置ID，自增主键',
                                      user_id INT NOT NULL COMMENT '关联的用户ID，外键关联users表',
                                      name VARCHAR(100) COMMENT '名称',
                                      speed FLOAT COMMENT '语音语速，浮点型（支持0.5-2.0）',
                                      volume INT COMMENT '语音音量，整数型（数值越大音量越高，通常范围0-10）',
                                      speaker VARCHAR(50) COMMENT '根据枚举类查找腾讯服务商具体对应信息',
                                      emotion VARCHAR(50) COMMENT '语音情感（如"happy"、"sad"、"neutral"等情感类型）',
                                      system_prompt TEXT COMMENT '系统提示词，用于定义机器人的行为模式或角色设定',
                                      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    -- 外键约束，关联users表的id字段，级联删除
                                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '机器人语音及行为配置表';
