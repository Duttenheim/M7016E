/*
 * GET servers page.
 */

exports.index = function(req, res){
  res.render('servers', { title: 'Servers' })
};