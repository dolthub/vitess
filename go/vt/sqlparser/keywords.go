// Copyright 2023 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlparser

// keywords is a map of mysql keywords that fall into two categories:
// 1) keywords considered reserved by MySQL
// 2) keywords for us to handle specially in sql.y
//
// Those marked as UNUSED are likely reserved keywords. We add them here so that
// when rewriting queries we can properly backtick quote them so they don't cause issues
//
// NOTE: If you add new keywords, add them also to the reserved_keywords or
// non_reserved_keywords grammar in sql.y -- this will allow the keyword to be used
// in identifiers. See the docs for each grammar to determine which one to put it into.
// Everything in this list will be escaped when used as identifier in printing.
var keywords = map[string]int{
	"_armscii8":                     UNDERSCORE_ARMSCII8,
	"_ascii":                        UNDERSCORE_ASCII,
	"_big5":                         UNDERSCORE_BIG5,
	"_binary":                       UNDERSCORE_BINARY,
	"_cp1250":                       UNDERSCORE_CP1250,
	"_cp1251":                       UNDERSCORE_CP1251,
	"_cp1256":                       UNDERSCORE_CP1256,
	"_cp1257":                       UNDERSCORE_CP1257,
	"_cp850":                        UNDERSCORE_CP850,
	"_cp852":                        UNDERSCORE_CP852,
	"_cp866":                        UNDERSCORE_CP866,
	"_cp932":                        UNDERSCORE_CP932,
	"_dec8":                         UNDERSCORE_DEC8,
	"_eucjpms":                      UNDERSCORE_EUCJPMS,
	"_euckr":                        UNDERSCORE_EUCKR,
	"_gb18030":                      UNDERSCORE_GB18030,
	"_gb2312":                       UNDERSCORE_GB2312,
	"_gbk":                          UNDERSCORE_GBK,
	"_geostd8":                      UNDERSCORE_GEOSTD8,
	"_greek":                        UNDERSCORE_GREEK,
	"_hebrew":                       UNDERSCORE_HEBREW,
	"_hp8":                          UNDERSCORE_HP8,
	"_keybcs2":                      UNDERSCORE_KEYBCS2,
	"_koi8r":                        UNDERSCORE_KOI8R,
	"_koi8u":                        UNDERSCORE_KOI8U,
	"_latin1":                       UNDERSCORE_LATIN1,
	"_latin2":                       UNDERSCORE_LATIN2,
	"_latin5":                       UNDERSCORE_LATIN5,
	"_latin7":                       UNDERSCORE_LATIN7,
	"_macce":                        UNDERSCORE_MACCE,
	"_macroman":                     UNDERSCORE_MACROMAN,
	"_sjis":                         UNDERSCORE_SJIS,
	"_swe7":                         UNDERSCORE_SWE7,
	"_tis620":                       UNDERSCORE_TIS620,
	"_ucs2":                         UNDERSCORE_UCS2,
	"_ujis":                         UNDERSCORE_UJIS,
	"_utf16":                        UNDERSCORE_UTF16,
	"_utf16le":                      UNDERSCORE_UTF16LE,
	"_utf32":                        UNDERSCORE_UTF32,
	"_utf8":                         UNDERSCORE_UTF8,
	"_utf8mb3":                      UNDERSCORE_UTF8MB3,
	"_utf8mb4":                      UNDERSCORE_UTF8MB4,
	"accessible":                    ACCESSIBLE,
	"account":                       ACCOUNT,
	"action":                        ACTION,
	"add":                           ADD,
	"admin":                         ADMIN,
	"after":                         AFTER,
	"against":                       AGAINST,
	"algorithm":                     ALGORITHM,
	"all":                           ALL,
	"alter":                         ALTER,
	"always":                        ALWAYS,
	"analyze":                       ANALYZE,
	"and":                           AND,
	"application_password_admin":    APPLICATION_PASSWORD_ADMIN,
	"array":                         ARRAY,
	"as":                            AS,
	"asc":                           ASC,
	"asensitive":                    ASENSITIVE,
	"at":                            AT,
	"attribute":                     ATTRIBUTE,
	"audit_abort_exempt":            AUDIT_ABORT_EXEMPT,
	"audit_admin":                   AUDIT_ADMIN,
	"authentication":                AUTHENTICATION,
	"authentication_policy_admin":   AUTHENTICATION_POLICY_ADMIN,
	"autoextend_size":               AUTOEXTEND_SIZE,
	"auto_increment":                AUTO_INCREMENT,
	"avg":                           AVG,
	"avg_row_length":                AVG_ROW_LENGTH,
	"backup_admin":                  BACKUP_ADMIN,
	"before":                        BEFORE,
	"begin":                         BEGIN,
	"between":                       BETWEEN,
	"bigint":                        BIGINT,
	"binary":                        BINARY,
	"binlog_admin":                  BINLOG_ADMIN,
	"binlog_encryption_admin":       BINLOG_ENCRYPTION_ADMIN,
	"bit":                           BIT,
	"bit_and":                       BIT_AND,
	"bit_or":                        BIT_OR,
	"bit_xor":                       BIT_XOR,
	"blob":                          BLOB,
	"bool":                          BOOL,
	"boolean":                       BOOLEAN,
	"both":                          BOTH,
	"by":                            BY,
	"call":                          CALL,
	"cascade":                       CASCADE,
	"cascaded":                      CASCADED,
	"case":                          CASE,
	"cast":                          CAST,
	"catalog_name":                  CATALOG_NAME,
	"chain":                         CHAIN,
	"change":                        CHANGE,
	"channel":                       CHANNEL,
	"char":                          CHAR,
	"character":                     CHARACTER,
	"charset":                       CHARSET,
	"check":                         CHECK,
	"checksum":                      CHECKSUM,
	"cipher":                        CIPHER,
	"class_origin":                  CLASS_ORIGIN,
	"client":                        CLIENT,
	"clone_admin":                   CLONE_ADMIN,
	"close":                         CLOSE,
	"coalesce":                      COALESCE,
	"collate":                       COLLATE,
	"collation":                     COLLATION,
	"column":                        COLUMN,
	"column_name":                   COLUMN_NAME,
	"columns":                       COLUMNS,
	"comment":                       COMMENT_KEYWORD,
	"commit":                        COMMIT,
	"committed":                     COMMITTED,
	"compact":                       COMPACT,
	"completion":                    COMPLETION,
	"compressed":                    COMPRESSED,
	"compression":                   COMPRESSION,
	"condition":                     CONDITION,
	"connection":                    CONNECTION,
	"connection_admin":              CONNECTION_ADMIN,
	"consistent":                    CONSISTENT,
	"constraint":                    CONSTRAINT,
	"constraint_catalog":            CONSTRAINT_CATALOG,
	"constraint_name":               CONSTRAINT_NAME,
	"constraint_schema":             CONSTRAINT_SCHEMA,
	"contains":                      CONTAINS,
	"contained":                     CONTAINED,
	"continue":                      CONTINUE,
	"convert":                       CONVERT,
	"copy":                          COPY,
	"count":                         COUNT,
	"create":                        CREATE,
	"cross":                         CROSS,
	"cube":                          CUBE,
	"cume_dist":                     CUME_DIST,
	"current":                       CURRENT,
	"current_date":                  CURRENT_DATE,
	"current_time":                  CURRENT_TIME,
	"current_timestamp":             CURRENT_TIMESTAMP,
	"current_user":                  CURRENT_USER,
	"cursor":                        CURSOR,
	"cursor_name":                   CURSOR_NAME,
	"data":                          DATA,
	"database":                      DATABASE,
	"databases":                     DATABASES,
	"date":                          DATE,
	"datetime":                      DATETIME,
	"day":                           DAY,
	"day_hour":                      DAY_HOUR,
	"day_microsecond":               DAY_MICROSECOND,
	"day_minute":                    DAY_MINUTE,
	"day_second":                    DAY_SECOND,
	"deallocate":                    DEALLOCATE,
	"dec":                           DEC,
	"decimal":                       DECIMAL,
	"declare":                       DECLARE,
	"default":                       DEFAULT,
	"definer":                       DEFINER,
	"definition":                    DEFINITION,
	"delay_key_write":               DELAY_KEY_WRITE,
	"delayed":                       DELAYED,
	"delete":                        DELETE,
	"dense_rank":                    DENSE_RANK,
	"desc":                          DESC,
	"describe":                      DESCRIBE,
	"description":                   DESCRIPTION,
	"deterministic":                 DETERMINISTIC,
	"directory":                     DIRECTORY,
	"disable":                       DISABLE,
	"discard":                       DISCARD,
	"disk":                          DISK,
	"distinct":                      DISTINCT,
	"distinctrow":                   DISTINCTROW,
	"div":                           DIV,
	"do":                            DO,
	"double":                        DOUBLE,
	"drop":                          DROP,
	"dual":                          DUAL,
	"dumpfile":                      DUMPFILE,
	"duplicate":                     DUPLICATE,
	"dynamic":                       DYNAMIC,
	"each":                          EACH,
	"else":                          ELSE,
	"elseif":                        ELSEIF,
	"empty":                         EMPTY,
	"enable":                        ENABLE,
	"enclosed":                      ENCLOSED,
	"encrypted":                     ENCRYPTED,
	"encryption":                    ENCRYPTION,
	"encryption_key_admin":          ENCRYPTION_KEY_ADMIN,
	"encryption_key_id":             ENCRYPTION_KEY_ID,
	"end":                           END,
	"ends":                          ENDS,
	"enforced":                      ENFORCED,
	"engine":                        ENGINE,
	"engine_attribute":              ENGINE_ATTRIBUTE,
	"engines":                       ENGINES,
	"enum":                          ENUM,
	"error":                         ERROR,
	"errors":                        ERRORS,
	"escape":                        ESCAPE,
	"escaped":                       ESCAPED,
	"event":                         EVENT,
	"events":                        EVENTS,
	"every":                         EVERY,
	"except":                        EXCEPT,
	"exchange":                      EXCHANGE,
	"exclusive":                     EXCLUSIVE,
	"execute":                       EXECUTE,
	"exists":                        EXISTS,
	"exit":                          EXIT,
	"expansion":                     EXPANSION,
	"expire":                        EXPIRE,
	"explain":                       EXPLAIN,
	"extended":                      EXTENDED,
	"extract":                       EXTRACT,
	"failed_login_attempts":         FAILED_LOGIN_ATTEMPTS,
	"false":                         FALSE,
	"fetch":                         FETCH,
	"fields":                        FIELDS,
	"file":                          FILE,
	"filter":                        FILTER,
	"firewall_admin":                FIREWALL_ADMIN,
	"firewall_exempt":               FIREWALL_EXEMPT,
	"firewall_user":                 FIREWALL_USER,
	"first":                         FIRST,
	"first_value":                   FIRST_VALUE,
	"fixed":                         FIXED,
	"float":                         FLOAT_TYPE,
	"float4":                        FLOAT4,
	"float8":                        FLOAT8,
	"flush":                         FLUSH,
	"flush_optimizer_costs":         FLUSH_OPTIMIZER_COSTS,
	"flush_status":                  FLUSH_STATUS,
	"flush_tables":                  FLUSH_TABLES,
	"flush_user_resources":          FLUSH_USER_RESOURCES,
	"following":                     FOLLOWING,
	"follows":                       FOLLOWS,
	"for":                           FOR,
	"force":                         FORCE,
	"foreign":                       FOREIGN,
	"format":                        FORMAT,
	"found":                         FOUND,
	"from":                          FROM,
	"full":                          FULL,
	"fulltext":                      FULLTEXT,
	"function":                      FUNCTION,
	"general":                       GENERAL,
	"generated":                     GENERATED,
	"geometry":                      GEOMETRY,
	"geometrycollection":            GEOMETRYCOLLECTION,
	"get":                           GET,
	"get_format":                    GET_FORMAT,
	"global":                        GLOBAL,
	"grant":                         GRANT,
	"grants":                        GRANTS,
	"group":                         GROUP,
	"group_concat":                  GROUP_CONCAT,
	"group_replication_admin":       GROUP_REPLICATION_ADMIN,
	"group_replication_stream":      GROUP_REPLICATION_STREAM,
	"grouping":                      GROUPING,
	"groups":                        GROUPS,
	"handler":                       HANDLER,
	"hash":                          HASH,
	"having":                        HAVING,
	"high_priority":                 HIGH_PRIORITY,
	"histogram":                     HISTOGRAM,
	"history":                       HISTORY,
	"hosts":                         HOSTS,
	"hour":                          HOUR,
	"hour_microsecond":              HOUR_MICROSECOND,
	"hour_minute":                   HOUR_MINUTE,
	"hour_second":                   HOUR_SECOND,
	"identified":                    IDENTIFIED,
	"if":                            IF,
	"ignore":                        IGNORE,
	"import":                        IMPORT,
	"in":                            IN,
	"index":                         INDEX,
	"indexes":                       INDEXES,
	"infile":                        INFILE,
	"initial":                       INITIAL,
	"inner":                         INNER,
	"innodb_redo_log_archive":       INNODB_REDO_LOG_ARCHIVE,
	"innodb_redo_log_enable":        INNODB_REDO_LOG_ENABLE,
	"inout":                         INOUT,
	"insensitive":                   INSENSITIVE,
	"insert":                        INSERT,
	"insert_method":                 INSERT_METHOD,
	"inplace":                       INPLACE,
	"instant":                       INSTANT,
	"int":                           INT,
	"int1":                          INT1,
	"int2":                          INT2,
	"int3":                          INT3,
	"int4":                          INT4,
	"int8":                          INT8,
	"integer":                       INTEGER,
	"interval":                      INTERVAL,
	"intersect":                     INTERSECT,
	"into":                          INTO,
	"invisible":                     INVISIBLE,
	"invoker":                       INVOKER,
	"io_after_gtids":                IO_AFTER_GTIDS,
	"io_before_gtids":               IO_BEFORE_GTIDS,
	"io_thread":                     IO_THREAD,
	"is":                            IS,
	"isolation":                     ISOLATION,
	"issuer":                        ISSUER,
	"itef_quotes":                   ITEF_QUOTES,
	"iterate":                       ITERATE,
	"join":                          JOIN,
	"json":                          JSON,
	"json_arrayagg":                 JSON_ARRAYAGG,
	"json_objectagg":                JSON_OBJECTAGG,
	"json_table":                    JSON_TABLE,
	"key":                           KEY,
	"key_block_size":                KEY_BLOCK_SIZE,
	"keys":                          KEYS,
	"kill":                          KILL,
	"lag":                           LAG,
	"language":                      LANGUAGE,
	"last":                          LAST,
	"last_insert_id":                LAST_INSERT_ID,
	"last_value":                    LAST_VALUE,
	"lateral":                       LATERAL,
	"lead":                          LEAD,
	"leading":                       LEADING,
	"leave":                         LEAVE,
	"left":                          LEFT,
	"less":                          LESS,
	"level":                         LEVEL,
	"like":                          LIKE,
	"limit":                         LIMIT,
	"linear":                        LINEAR,
	"lines":                         LINES,
	"linestring":                    LINESTRING,
	"list":                          LIST,
	"load":                          LOAD,
	"local":                         LOCAL,
	"localtime":                     LOCALTIME,
	"localtimestamp":                LOCALTIMESTAMP,
	"lock":                          LOCK,
	"locked":                        LOCKED,
	"log":                           LOG,
	"logs":                          LOGS,
	"long":                          LONG,
	"longblob":                      LONGBLOB,
	"longtext":                      LONGTEXT,
	"loop":                          LOOP,
	"low_priority":                  LOW_PRIORITY,
	"master":                        MASTER,
	"master_bind":                   MASTER_BIND,
	"master_ssl_verify_server_cert": MASTER_SSL_VERIFY_SERVER_CERT,
	"match":                         MATCH,
	"max":                           MAX,
	"max_connections_per_hour":      MAX_CONNECTIONS_PER_HOUR,
	"max_queries_per_hour":          MAX_QUERIES_PER_HOUR,
	"max_rows":                      MAX_ROWS,
	"max_updates_per_hour":          MAX_UPDATES_PER_HOUR,
	"max_user_connections":          MAX_USER_CONNECTIONS,
	"maxvalue":                      MAXVALUE,
	"mediumblob":                    MEDIUMBLOB,
	"mediumint":                     MEDIUMINT,
	"mediumtext":                    MEDIUMTEXT,
	"memory":                        MEMORY,
	"merge":                         MERGE,
	"message_text":                  MESSAGE_TEXT,
	"middleint":                     MIDDLEINT,
	"microsecond":                   MICROSECOND,
	"min":                           MIN,
	"min_rows":                      MIN_ROWS,
	"minute":                        MINUTE,
	"minute_microsecond":            MINUTE_MICROSECOND,
	"minute_second":                 MINUTE_SECOND,
	"mod":                           MOD,
	"mode":                          MODE,
	"modifies":                      MODIFIES,
	"modify":                        MODIFY,
	"month":                         MONTH,
	"multilinestring":               MULTILINESTRING,
	"multipoint":                    MULTIPOINT,
	"multipolygon":                  MULTIPOLYGON,
	"mysql_errno":                   MYSQL_ERRNO,
	"name":                          NAME,
	"names":                         NAMES,
	"national":                      NATIONAL,
	"natural":                       NATURAL,
	"nested":                        NESTED,
	"nchar":                         NCHAR,
	"ndb_stored_user":               NDB_STORED_USER,
	"never":                         NEVER,
	"next":                          NEXT,
	"no":                            NO,
	"no_write_to_binlog":            NO_WRITE_TO_BINLOG,
	"none":                          NONE,
	"not":                           NOT,
	"now":                           NOW,
	"nowait":                        NOWAIT,
	"nth_value":                     NTH_VALUE,
	"ntile":                         NTILE,
	"null":                          NULL,
	"numeric":                       NUMERIC,
	"nvarchar":                      NVARCHAR,
	"of":                            OF,
	"off":                           OFF,
	"offset":                        OFFSET,
	"on":                            ON,
	"only":                          ONLY,
	"open":                          OPEN,
	"optimize":                      OPTIMIZE,
	"optimizer_costs":               OPTIMIZER_COSTS,
	"option":                        OPTION,
	"optional":                      OPTIONAL,
	"optionally":                    OPTIONALLY,
	"or":                            OR,
	"order":                         ORDER,
	"ordinality":                    ORDINALITY,
	"organization":                  ORGANIZATION,
	"out":                           OUT,
	"outer":                         OUTER,
	"outfile":                       OUTFILE,
	"over":                          OVER,
	"pack_keys":                     PACK_KEYS,
	"page_checksum":                 PAGE_CHECKSUM,
	"page_compressed":               PAGE_COMPRESSED,
	"page_compression_level":        PAGE_COMPRESSION_LEVEL,
	"partition":                     PARTITION,
	"partitioning":                  PARTITIONING,
	"partitions":                    PARTITIONS,
	"password":                      PASSWORD,
	"password_lock_time":            PASSWORD_LOCK_TIME,
	"passwordless_user_admin":       PASSWORDLESS_USER_ADMIN,
	"path":                          PATH,
	"percent_rank":                  PERCENT_RANK,
	"persist":                       PERSIST,
	"persist_only":                  PERSIST_ONLY,
	"persist_ro_variables_admin":    PERSIST_RO_VARIABLES_ADMIN,
	"plan":                          PLAN,
	"plugins":                       PLUGINS,
	"point":                         POINT,
	"polygon":                       POLYGON,
	"position":                      POSITION,
	"precedes":                      PRECEDES,
	"preceding":                     PRECEDING,
	"precision":                     PRECISION,
	"prepare":                       PREPARE,
	"preserve":                      PRESERVE,
	"primary":                       PRIMARY,
	"privileges":                    PRIVILEGES,
	"procedure":                     PROCEDURE,
	"process":                       PROCESS,
	"processlist":                   PROCESSLIST,
	"proxy":                         PROXY,
	"purge":                         PURGE,
	"quarter":                       QUARTER,
	"query":                         QUERY,
	"random":                        RANDOM,
	"range":                         RANGE,
	"rank":                          RANK,
	"rebuild":                       REBUILD,
	"read":                          READ,
	"read_write":                    READ_WRITE,
	"reads":                         READS,
	"real":                          REAL,
	"recursive":                     RECURSIVE,
	"redundant":                     REDUNDANT,
	"reference":                     REFERENCE,
	"references":                    REFERENCES,
	"regexp":                        REGEXP,
	"relay":                         RELAY,
	"release":                       RELEASE,
	"reload":                        RELOAD,
	"remove":                        REMOVE,
	"rename":                        RENAME,
	"reorganize":                    REORGANIZE,
	"repair":                        REPAIR,
	"repeat":                        REPEAT,
	"repeatable":                    REPEATABLE,
	"replace":                       REPLACE,
	"replica":                       REPLICA,
	"replicas":                      REPLICAS,
	"replicate_do_table":            REPLICATE_DO_TABLE,
	"replicate_ignore_table":        REPLICATE_IGNORE_TABLE,
	"replication":                   REPLICATION,
	"replication_applier":           REPLICATION_APPLIER,
	"replication_slave_admin":       REPLICATION_SLAVE_ADMIN,
	"require":                       REQUIRE,
	"reset":                         RESET,
	"resignal":                      RESIGNAL,
	"resource_group_admin":          RESOURCE_GROUP_ADMIN,
	"resource_group_user":           RESOURCE_GROUP_USER,
	"restrict":                      RESTRICT,
	"return":                        RETURN,
	"returning":                     RETURNING,
	"reuse":                         REUSE,
	"revoke":                        REVOKE,
	"right":                         RIGHT,
	"rlike":                         REGEXP,
	"role":                          ROLE,
	"role_admin":                    ROLE_ADMIN,
	"rollback":                      ROLLBACK,
	"routine":                       ROUTINE,
	"row":                           ROW,
	"row_format":                    ROW_FORMAT,
	"row_number":                    ROW_NUMBER,
	"rows":                          ROWS,
	"savepoint":                     SAVEPOINT,
	"schedule":                      SCHEDULE,
	"schema":                        SCHEMA,
	"schema_name":                   SCHEMA_NAME,
	"schemas":                       SCHEMAS,
	"second":                        SECOND,
	"second_microsecond":            SECOND_MICROSECOND,
	"secondary_engine":              SECONDARY_ENGINE,
	"secondary_engine_attribute":    SECONDARY_ENGINE_ATTRIBUTE,
	"security":                      SECURITY,
	"select":                        SELECT,
	"sensitive":                     SENSITIVE,
	"sensitive_variables_observer":  SENSITIVE_VARIABLES_OBSERVER,
	"separator":                     SEPARATOR,
	"sequence":                      SEQUENCE,
	"serial":                        SERIAL,
	"serializable":                  SERIALIZABLE,
	"session":                       SESSION,
	"session_variables_admin":       SESSION_VARIABLES_ADMIN,
	"set":                           SET,
	"set_user_id":                   SET_USER_ID,
	"share":                         SHARE,
	"shared":                        SHARED,
	"show":                          SHOW,
	"show_routine":                  SHOW_ROUTINE,
	"shutdown":                      SHUTDOWN,
	"signal":                        SIGNAL,
	"signed":                        SIGNED,
	"skip":                          SKIP,
	"skip_query_rewrite":            SKIP_QUERY_REWRITE,
	"slave":                         SLAVE,
	"slow":                          SLOW,
	"smallint":                      SMALLINT,
	"snapshot":                      SNAPSHOT,
	"source":                        SOURCE,
	"source_auto_position":          SOURCE_AUTO_POSITION,
	"source_connect_retry":          SOURCE_CONNECT_RETRY,
	"source_host":                   SOURCE_HOST,
	"source_ssl":                    SOURCE_SSL,
	"source_password":               SOURCE_PASSWORD,
	"source_port":                   SOURCE_PORT,
	"source_retry_count":            SOURCE_RETRY_COUNT,
	"source_user":                   SOURCE_USER,
	"spatial":                       SPATIAL,
	"specific":                      SPECIFIC,
	"sql":                           SQL,
	"sql_big_result":                SQL_BIG_RESULT,
	"sql_cache":                     SQL_CACHE,
	"sql_calc_found_rows":           SQL_CALC_FOUND_ROWS,
	"sql_no_cache":                  SQL_NO_CACHE,
	"sql_small_result":              SQL_SMALL_RESULT,
	"sqlexception":                  SQLEXCEPTION,
	"sqlstate":                      SQLSTATE,
	"sql_thread":                    SQL_THREAD,
	"sqlwarning":                    SQLWARNING,
	"srid":                          SRID,
	"ssl":                           SSL,
	"start":                         START,
	"starts":                        STARTS,
	"starting":                      STARTING,
	"stats_auto_recalc":             STATS_AUTO_RECALC,
	"stats_persistent":              STATS_PERSISTENT,
	"stats_sample_pages":            STATS_SAMPLE_PAGES,
	"status":                        STATUS,
	"std":                           STD,
	"stddev":                        STDDEV,
	"stddev_pop":                    STDDEV_POP,
	"stddev_samp":                   STDDEV_SAMP,
	"stop":                          STOP,
	"storage":                       STORAGE,
	"stored":                        STORED,
	"straight_join":                 STRAIGHT_JOIN,
	"stream":                        STREAM,
	"subclass_origin":               SUBCLASS_ORIGIN,
	"subject":                       SUBJECT,
	"subpartition":                  SUBPARTITION,
	"subpartitions":                 SUBPARTITIONS,
	"substr":                        SUBSTR,
	"substring":                     SUBSTRING,
	"sum":                           SUM,
	"super":                         SUPER,
	"system":                        SYSTEM,
	"system_time":                   SYSTEM_TIME,
	"system_variables_admin":        SYSTEM_VARIABLES_ADMIN,
	"table":                         TABLE,
	"table_checksum":                TABLE_CHECKSUM,
	"table_encryption_admin":        TABLE_ENCRYPTION_ADMIN,
	"table_name":                    TABLE_NAME,
	"tables":                        TABLES,
	"tablespace":                    TABLESPACE,
	"temporary":                     TEMPORARY,
	"temptable":                     TEMPTABLE,
	"terminated":                    TERMINATED,
	"text":                          TEXT,
	"than":                          THAN,
	"then":                          THEN,
	"time":                          TIME,
	"timestamp":                     TIMESTAMP,
	"timestampadd":                  TIMESTAMPADD,
	"timestampdiff":                 TIMESTAMPDIFF,
	"tinyblob":                      TINYBLOB,
	"tinyint":                       TINYINT,
	"tinytext":                      TINYTEXT,
	"to":                            TO,
	"tp_connection_admin":           TP_CONNECTION_ADMIN,
	"trailing":                      TRAILING,
	"transaction":                   TRANSACTION,
	"transactional":                 TRANSACTIONAL,
	"trigger":                       TRIGGER,
	"triggers":                      TRIGGERS,
	"trim":                          TRIM,
	"true":                          TRUE,
	"truncate":                      TRUNCATE,
	"unbounded":                     UNBOUNDED,
	"uncommitted":                   UNCOMMITTED,
	"undefined":                     UNDEFINED,
	"undo":                          UNDO,
	"union":                         UNION,
	"unique":                        UNIQUE,
	"unlock":                        UNLOCK,
	"unknown":                       UNKNOWN,
	"unsigned":                      UNSIGNED,
	"until":                         UNTIL,
	"update":                        UPDATE,
	"usage":                         USAGE,
	"use":                           USE,
	"user":                          USER,
	"user_resources":                USER_RESOURCES,
	"using":                         USING,
	"utc_date":                      UTC_DATE,
	"utc_time":                      UTC_TIME,
	"utc_timestamp":                 UTC_TIMESTAMP,
	"validation":                    VALIDATION,
	"value":                         VALUE,
	"values":                        VALUES,
	"var_pop":                       VAR_POP,
	"var_samp":                      VAR_SAMP,
	"varbinary":                     VARBINARY,
	"varchar":                       VARCHAR,
	"varcharacter":                  VARCHARACTER,
	"variables":                     VARIABLES,
	"variance":                      VARIANCE,
	"varying":                       VARYING,
	"vector":                        VECTOR,
	"version":                       VERSION,
	"versioning":                    VERSIONING,
	"versions":                      VERSIONS,
	"version_token_admin":           VERSION_TOKEN_ADMIN,
	"view":                          VIEW,
	"virtual":                       VIRTUAL,
	"visible":                       VISIBLE,
	"warnings":                      WARNINGS,
	"week":                          WEEK,
	"when":                          WHEN,
	"where":                         WHERE,
	"while":                         WHILE,
	"window":                        WINDOW,
	"with":                          WITH,
	"without":                       WITHOUT,
	"work":                          WORK,
	"write":                         WRITE,
	"x509":                          X509,
	"xa_recover_admin":              XA_RECOVER_ADMIN,
	"xor":                           XOR,
	"year":                          YEAR,
	"year_month":                    YEAR_MONTH,
	"yes":                           YES,
	"zerofill":                      ZEROFILL,
}

// keywordStrings contains the reverse mapping of token to keyword strings
var keywordStrings = map[int]string{}

func init() {
	for str, id := range keywords {
		if id == UNUSED {
			continue
		}
		keywordStrings[id] = str
	}
}
