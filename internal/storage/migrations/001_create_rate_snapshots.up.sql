CREATE TABLE IF NOT EXISTS rate_snapshots (
    id          BIGSERIAL   PRIMARY KEY,
    ask_price   NUMERIC     NOT NULL,
    bid_price   NUMERIC     NOT NULL,
    top_n       NUMERIC     NOT NULL,
    avg_nm      NUMERIC     NOT NULL,
    fetched_at  TIMESTAMPTZ NOT NULL
);
