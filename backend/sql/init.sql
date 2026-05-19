create database my_blog
default character set utf8mb4
default collate utf8mb4_unicode_ci;

use my_blog;

-- 管理员
CREATE TABLE admin_user(
                           id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '管理员ID',
                           username varchar(50) NOT NULL COMMENT '用户名',
                           password varchar(50) not null comment '加密后的密码',
                           nickname varchar(255) not null comment '昵称',
                           create_time datetime not null default current_timestamp comment'创建时间',
                           update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',

                           unique key uk_username(username)
)engine=Innodb default charset=utf8mb4 comment='管理员表';

-- 文章分类
create table category(
                         id BIGINT primary key auto_increment comment '分类ID',
                         category_name varchar(50) not null  comment'分类名称',
                         slug varchar (80) default null comment '分类别名，用于UDL',
                         sort INT not null default 0 comment '排序值(越大越靠前)',
                         status tinyint not null default 1 comment'状态：1正常，0禁用',
                         article_count int not null default 0 comment '文章数量',
                         create_time datetime not null default current_timestamp comment '创建时间',
                         update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',

                         unique key uk_name(category_name),
                         unique key uk_slug(slug),
                         key idx_status_sort(status,sort)
)engine=InnoDB default charset=utf8mb4 comment='文章分类表';

-- 标签
create table tag (
                     id BIGINT primary key auto_increment comment '标签ID',
                     tag_name varchar(50) not null comment '标签名称',
                     slug varchar(80) default null comment '标签别名，用于URL',
                     article_count int not null default 0 comment'文章数量',
                     create_time datetime not null default current_timestamp comment '创建时间',
                     update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',

                     unique key uk_name(tag_name),
                     unique key uk_slug (slug)
)engine=InnoDB default charset =utf8mb4 comment='标签表';

-- 文章
create table article (
                         id bigint primary key auto_increment comment '文章ID',
                         category_id bigint not null comment '分类ID',
                         title varchar(150) not null comment '文章标题',
                         summary varchar(500) default null comment '文章摘要',
                         content mediumint not null comment '文章正文，Markdown内容',
                         cover_url varchar(255) default null comment '封面图URL',

                         status tinyint not null default 0 comment '状态：0草稿，1发布，2下架',
                         is_top tinyint not null default 0 comment '是否置顶：1是，0否',
                         is_deleted tinyint not null default 0 comment  '逻辑删除“1删除，0未删除',

    -- view_count int not null default 0 comment '浏览量',
                         comment_count int not null default 0 comment '评论数',

                         published_time datetime default null comment '发布时间',
                         create_time datetime not null default current_timestamp comment'创建时间',
                         update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',

                         key idx_category_status_time(category_id,status,is_deleted,published_time),
                         key idx_status_top_time (status,is_deleted,is_top,published_time),
                         key idx_status_time (status,is_deleted,published_time),
                         key idx_title (title)
)engine=Innodb default charset =utf8mb4 comment='文章表';

-- 文章标签
create table article_tag(
                            id bigint primary key auto_increment comment '主键ID',
                            article_id bigint not null comment '文章ID',
                            tag_id bigint not null comment '标签ID',
                            create_time datetime not null default current_timestamp comment '创建时间',

                            unique key uk_article_tag (article_id,tag_id),
                            key idx_tag_article (tag_id,article_id)
)engine=innodb default charset =utf8mb4 comment='文章标签关联表';

-- 评论表
create table comment (
                         id bigint primary key auto_increment comment '评论ID',

                         target_type tinyint not null comment '评论目标类型：1文章，2留言板',
                         target_id bigint default 0 comment '目标ID：文章ID;留言板固定为0',

    -- is_top tinyint not null default 0 comment '是否置顶：1是，0否',

                         root_id bigint not null default 0 comment '根评论ID',
                         parent_id bigint not null default 0 comment '父评论ID，0表示一级评论',

                         nickname varchar(50) not null comment '游客昵称',
                         email varchar(100) default null comment '游客邮箱，不公开展示',
                         website varchar (255) default null comment '游客网站',
                         avatar varchar(255) default null comment '头像URL，可根据邮箱生成',

                         content varchar(1000) not null comment '评论内容',
                         ip varchar (64) default null comment '评论者IP',
                         user_agent varchar(500) default null comment '浏览器UA',

                         is_admin_reply tinyint not null default 0 comment '是否管理员回复：1是，0否',
                         is_deleted tinyint not null default 0 comment '逻辑删除：1删除，0未删除',

                         create_time datetime not null default current_timestamp comment '创建时间',
                         update_time datetime not null default current_timestamp on update current_timestamp,

                         key idx_target_time (target_type,target_id,is_deleted,create_time),
                         key idx_parent_time (parent_id,create_time),
                         key idx_root_time (root_id,create_time),
                         key idx_time (create_time),
                         key idx_email(email)

) engine=InnoDB default charset=utf8mb4 comment='评论表';

-- 友链表
create table friend_link(
                            id bigint primary key auto_increment comment '友链ID',
                            name varchar (100) not null comment '站点名称',
                            url varchar (255) not null comment '站点地址',
                            logo varchar(255) default null comment '站点描述',
                            sort int not null default 0 comment '排序值，越大越靠前',
                            status tinyint not null default 1 comment '状态：1正常，0隐藏',
                            create_time datetime not null default current_timestamp comment '创建时间',
                            update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',

                            unique key uk_url (url),
                            key idx_status_sort (status,sort)
)engine=InnoDB default charset=utf8mb4 comment='友链表';

-- 网站配置
create table site_config(
                            id bigint primary key auto_increment comment '配置ID',
                            config_key varchar (100) not null comment '配置键',
                            config_value text comment '配置值',
                            description varchar(255) default  null comment '配置说明',
                            create_time datetime not null default  current_timestamp comment '创建时间',
                            update_time datetime not null default current_timestamp on update current_timestamp comment '更新时间',

                            unique key uk_config_key(config_key)
)engine =InnoDB default charset = utf8mb4 comment ='网站配置表'




