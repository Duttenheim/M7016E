/*
 * GET manage page.
 */

exports.index = function(req, res){
  res.render('manage', { title: 'Manage', serv_addr: req.query.addr, node_id: req.query.node })
};