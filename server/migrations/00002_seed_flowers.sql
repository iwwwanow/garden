-- +goose Up

INSERT INTO flowers (season, image_path) VALUES
    (1, 'flowers/season1_rose.png'),
    (1, 'flowers/season1_tulip.png'),
    (1, 'flowers/season1_daisy.png');

-- +goose Down

DELETE FROM flowers WHERE season = 1;
