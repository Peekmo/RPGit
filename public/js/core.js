jQuery(function($) {
	if ($('.nav-menu .item').length) {
		$('.nav-menu .item')
		.tab()
		.tab('activate tab', 'first')
		.tab('activate navigation', 'first')
	;
	}

	$('#search-btn').click(function(e) {
		var name = $('#search-field').val();
		if (name != "") {
			document.location.href = '/users/' + name;
		}
	});	

	$('#search-field').keyup(function(e) {
		if (e.keyCode == 13) {
			var name = $(this).val();
			if (name != "") {
				document.location.href = '/users/' + name;
			}
		}
	});

	//var Home = {
	//	init: function() {
	//		this.compile();
	//		this.render();
	//	},

	//	compile: function() {
	//		this.dailyTpl = Handlebars.compile($("#overview-ranking").html());
	//		this.$daily = jQuery("#ranking-daily");
	//	},

	//	render: function() {
	//		this.$daily.html(this.dailyTpl(test));
	//	}
	//};

	//Home.init();
});
