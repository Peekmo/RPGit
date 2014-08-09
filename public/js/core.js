jQuery(function($) {
	$('.tab-ranking a').click(function (e) {
	  e.preventDefault()
	  $(this).tab('show')
	});


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
