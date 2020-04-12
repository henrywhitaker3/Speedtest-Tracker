module.exports = process.env.CONCAT_COV
  ? require('./lib-cov/concatenate')
  : require('./lib/concatenate');
