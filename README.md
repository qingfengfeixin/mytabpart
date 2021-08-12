# mytabpart

管理mysql分区表
1. 通过 CFG_AUTO_MAN_TAB_PART_INFO 进行分表配置
2. 可以选择直接进行数据库执行 或者 输出到sql脚本
3. 可以选择是否执行drop分区命令
4. 可以根据小时、天、周、月进行分区
5. 选择是否按照最后时间创建分区，还是最近一个月开始创建分区


init脚本
1. 创建 C_AUTO_MAN_TAB_PART_INFO 表

 drop table C_AUTO_MAN_TAB_PART_INFO;
 create table C_AUTO_MAN_TAB_PART_INFO(
 TABLE_NAME      varchar(50) not null COMMENT '表名',
 INTER_VAL       int  COMMENT '间隔',
 RETENTION_HOUR  int COMMENT '保留时长（小时）',
 IS_PART int  default 1  COMMENT '是否分区表 1 是',
 PRIMARY KEY (TABLE_NAME)
 )ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='分区配置表';
   
测试

drop table t1;
create table t1(
id int,
hiredate datetime
)partition by range columns(hiredate)(
PARTITION t1_20210101 values less than ('2021-01-01 01')
);

delete from C_AUTO_MAN_TAB_PART_INFO;
insert into C_AUTO_MAN_TAB_PART_INFO values('t1',24,168,1);

