-- 创建 go_role_user_rel 表  role_id、 user_id的索引
CREATE INDEX index_role_id  ON go_role_user_rel (role_id);
CREATE INDEX index_user_id  ON go_role_user_rel (user_id);
-- 查看是否创建成功
show  CREATE TABLE go_role_user_rel
-- 查看索引是否生效
EXPLAIN SELECT * from go_role_user_rel

-- 创建 go_role_actione_rel 表  role_id、 action_id的索引
alter table go_role_actione_rel add INDEX index_role_id (role_id);
alter table go_role_actione_rel add INDEX index_action_id (action_id);


-- 创建 go_role_menu_rel 表  role_id、 menu_id的索引
alter table go_role_menu_rel add INDEX index_role_id (role_id);
alter table go_role_menu_rel add INDEX index_menu_id (menu_id);

-- 删除索引
DROP index index_logic_delete ON go_role_menu_rel;