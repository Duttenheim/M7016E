/*
 * GET images page.
 */

exports.index = function(req, res){
  res.render('images', { title: 'Private repository images' })
};