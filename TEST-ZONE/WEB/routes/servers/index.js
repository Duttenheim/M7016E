/*
 * GET servers page.
 */

exports.index = function(req, res)
{
  var serverList = { server1: '130.240.134.117:2020', 
		  				server2: '130.240.134.120:2020', 
		  				server5: '127.0.0.1:2020'};
  res.render('servers', { title: 'SuperNode',  servers: serverList })
  
};