
ALTER TABLE `plugin_config_interface_parameters` ADD COLUMN `multiple` VARCHAR(10) NULL;
ALTER TABLE `plugin_config_interface_parameters` ADD COLUMN `ref_object_name` VARCHAR(60) NULL;
ALTER TABLE `plugin_object_property_meta` DROP COLUMN `ref_type`;
ALTER TABLE `plugin_object_property_meta` CHANGE COLUMN `ref_name` `ref_object_name` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_bin' NULL DEFAULT NULL ;
ALTER TABLE `plugin_object_property_meta` ADD COLUMN `multiple` VARCHAR(10) NULL;

ALTER TABLE `plugin_object_property_var` RENAME TO `core_object_property_var`;
ALTER TABLE `plugin_object_var` RENAME TO `core_object_var`;

ALTER TABLE `core_re_task_node_def_info` ADD COLUMN `prev_ctx_node_ids` VARCHAR(255) COLLATE utf8_bin DEFAULT NULL;
ALTER TABLE `core_re_proc_def_info` ADD COLUMN `tags` VARCHAR(255) NULL;