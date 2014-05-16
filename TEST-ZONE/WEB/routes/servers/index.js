/*
 * GET servers page.
 */

exports.index = function(req, res)
{
  var servers = { server1: '130.240.134.116', server2: '130.240.134.117', server3: '130.240.134.118', server4: '130.240.134.119' };
  res.render('servers', { title: 'Servers',  servers: servers })
};