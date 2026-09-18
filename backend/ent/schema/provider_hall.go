package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
)

type ProviderHallConfig struct{ ent.Schema }

func (ProviderHallConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "provider_hall_config", Checks: map[string]string{
		"provider_hall_config_singleton":          "id = 1",
		"provider_hall_config_budget_nonnegative": "daily_budget >= 0",
	}}}
}

func (ProviderHallConfig) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int64("id").Default(1).Immutable().Annotations(entsql.Annotation{Incremental: new(false)}),
		field.Bool("collection_enabled").Default(false),
		field.Bool("display_enabled").Default(false),
		field.Bool("tasks_enabled").Default(false),
		field.Bool("auto_schedule_enabled").Default(false),
		field.String("default_model").Default("").MaxLen(200),
		field.Enum("default_protocol").Values("responses", "chat_completions", "messages").Default("responses"),
		field.Enum("default_range").Values("6h", "24h", "7d", "30d").Default("6h"),
		field.String("gateway_origin").Default("").MaxLen(2048),
		field.Int64("operator_user_id").Optional().Nillable().Positive(),
		field.String("daily_budget").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.Postgres: "numeric(20,8)"}).DefaultFunc(func() decimal.Decimal { return decimal.Zero }),
		field.JSON("expected_nodes", []string{}).Default([]string{}),
	}, providerHallVersionFields()...)
}

// ProviderHallGroup uses the existing group ID as its primary key. The SQL
// migration owns its FK; no original group is created by a hall configuration.
type ProviderHallGroup struct{ ent.Schema }

func (ProviderHallGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "provider_hall_groups"}}
}

func (ProviderHallGroup) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int64("id").StorageKey("group_id").Immutable().Annotations(entsql.Annotation{Incremental: new(false)}),
		field.Bool("listed").Default(false),
		field.String("display_name").Default("").MaxRuneLen(100),
		field.String("description").Default("").MaxRuneLen(2000),
		field.Int("display_order").Default(0),
	}, providerHallVersionFields()...)
}

type ProviderHallProfile struct{ ent.Schema }

func (ProviderHallProfile) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "provider_hall_profiles", Checks: map[string]string{
		"provider_hall_profile_output_limit": "output_limit BETWEEN 1 AND 1024",
		"provider_hall_profile_input_price":  "reference_input_price >= 0",
		"provider_hall_profile_cache_price":  "reference_cache_price >= 0",
		"provider_hall_profile_cache_rate":   "reference_cache_rate BETWEEN 0 AND 1",
	}}}
}

func (ProviderHallProfile) Fields() []ent.Field {
	return append([]ent.Field{
		field.String("model").NotEmpty().MaxLen(200),
		field.Enum("protocol").Values("responses", "chat_completions", "messages"),
		field.Bool("supports_tools").Default(false),
		field.Int("output_limit").Default(256).Min(1).Max(1024),
		field.JSON("model_aliases", []string{}).Default([]string{}),
		field.String("reference_input_price").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.Postgres: "numeric(24,10)"}).Optional().Nillable(),
		field.String("reference_cache_price").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.Postgres: "numeric(24,10)"}).Optional().Nillable(),
		field.String("reference_cache_rate").GoType(decimal.Decimal{}).SchemaType(map[string]string{dialect.Postgres: "numeric(11,10)"}).Optional().Nillable(),
		field.Time("reference_confirmed_at").Optional().Nillable().SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}, providerHallVersionFields()...)
}

func (ProviderHallProfile) Indexes() []ent.Index {
	return []ent.Index{index.Fields("model", "protocol").Unique()}
}

type ProviderHallTarget struct{ ent.Schema }

func (ProviderHallTarget) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "provider_hall_targets", Checks: map[string]string{
		"provider_hall_target_probe_interval":        "probe_interval_seconds BETWEEN 60 AND 86400",
		"provider_hall_target_verification_interval": "verification_interval_seconds BETWEEN 3600 AND 604800",
		"provider_hall_target_enabled_key":           "NOT enabled OR probe_key_id IS NOT NULL",
	}}}
}

func (ProviderHallTarget) Fields() []ent.Field {
	return append([]ent.Field{
		field.Int64("group_id").Positive().Immutable(),
		field.Int64("profile_id").Positive().Immutable(),
		field.Int64("probe_key_id").Optional().Nillable().Positive(),
		field.Bool("enabled").Default(false),
		field.Bool("auto_schedule_enabled").Default(false),
		field.Int("probe_interval_seconds").Default(300).Min(60).Max(86400),
		field.Int("verification_interval_seconds").Default(86400).Min(3600).Max(604800),
	}, providerHallVersionFields()...)
}

func (ProviderHallTarget) Indexes() []ent.Index {
	return []ent.Index{index.Fields("group_id", "profile_id").Unique(), index.Fields("probe_key_id")}
}

// Permanent provenance contains IDs only. Intentionally no API key, user or
// group FK: removing the source entity must not erase its probe classification.
type ProviderHallProbeKey struct{ ent.Schema }

func (ProviderHallProbeKey) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Table: "provider_hall_probe_keys"}}
}

func (ProviderHallProbeKey) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("api_key_id").Immutable().Annotations(entsql.Annotation{Incremental: new(false)}),
		field.Int64("operator_user_id").Positive().Immutable(),
		field.Int64("group_id").Positive().Immutable(),
		field.Int64("registered_by").Positive().Immutable(),
		field.Time("registered_at").Immutable().Default(time.Now).SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}

func providerHallVersionFields() []ent.Field {
	return []ent.Field{
		field.Int64("version").Default(1).Positive(),
		field.Int64("updated_by").Optional().Nillable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).SchemaType(map[string]string{dialect.Postgres: "timestamptz"}),
	}
}
