/*
 * GET manage page.
 */

var request = require("request");

exports.index = function(req, res){
	var IP = '130.240.134.118';	
	var url = 'http://'+IP+':5000/v1/search?';
	request({
			url: url,
			json: true
		}, function (error, response, body) {
			var dataList = " ";
			if (!error && response.statusCode == 200) {
				dataList = body;
			}
			res.render('manage', { title: 'Manage', serv_addr: req.query.addr, node_id: req.query.node , datalist: dataList});
		})
};
