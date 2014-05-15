/*
 * GET manage page.
 */

exports.index = function(req, res){
  res.render('manage', { title: 'Manage' })
};