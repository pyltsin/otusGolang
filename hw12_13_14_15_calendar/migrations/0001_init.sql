-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id varchar(255) NOT NULL,
	title varchar(255) NOT NULL,
	date_event datetime NOT NULL,
	latency int(16) NOT NULL,
	note text,
	userID int(16),
	notify int(16)
  PRIMARY KEY (`id`)
);;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
