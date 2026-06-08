INSERT INTO
  "user" (display_name, email, handle, cognito_user_id)
VALUES
  (
    'Andrew Brown',
    'legal@legal.com',
    'andrewbrown',
    'MOCK'
  ),
  (
    'Andrew Bayko',
    'luiz@luiz.com.br',
    'bayko',
    'MOCK'
  );

INSERT INTO
  activity (user_id, message, expires_at)
VALUES
  (
    (
      SELECT
        id
      FROM
        public."user"
      WHERE
        "user".handle = 'andrewbrown'
      LIMIT
        1
    ),
    'This was imported as seed data!',
    current_timestamp + interval '10 day'
  );