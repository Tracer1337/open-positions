module.exports = ({ env }) => ({
  connection: {
    client: 'mysql',
    connection: {
      host: env('DATABASE_HOST', 'localhost'),
      port: parseInt(env('DATABASE_PORT', 3306)),
      database: env('DATABASE_NAME', 'open-positions'),
      user: env('DATABASE_USER', 'root'),
      password: env('DATABASE_PASSWORD', 'root'),
    },
    useNullAsDefault: true,
  },
});
