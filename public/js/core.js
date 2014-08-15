jQuery(function($) {
	$('.nav-menu .item')
		.tab()
		.tab('activate tab', 'first')
		.tab('activate navigation', 'first')
	;

	var Home = {
		init: function() {
			this.compile();
			this.render();
		},

		compile: function() {
			this.dailyTpl = Handlebars.compile($("#overview-ranking").html());
			this.$daily = jQuery("#ranking-daily");
		},

		render: function() {
			this.$daily.html(this.dailyTpl(test));
		}
	};

	Home.init();
});
