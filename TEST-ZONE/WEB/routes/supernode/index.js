/*
 * GET supernode page.
 */

exports.index = function(req, res)
{
  res.render('supernode', { title: 'EdgeNodes' })
};