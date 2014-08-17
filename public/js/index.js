function filter(element) {
	Home.language = element;
	Home.load('pushevent');
	$('#reset').fadeIn();
}

function reset() {
	Home.language = undefined;
	Home.load('pushevent');
	$('#reset').fadeOut();
}

$(".show-ranking").click(function(e) {
	var type = $(this).attr("data-type");
	var icon = $(this).find('i');

	$(".show-ranking").each(function(i, e) {
		if ($(e).attr('data-type') != type) {
			$(e).find('i').removeClass('minus sign');
			$(e).find('i').addClass('add sign');
			$(e).removeClass('red');
			$(e).addClass('green');
		}
	});

	$(".handle-ranking").fadeOut();

	if (icon.hasClass('add sign')) {
		icon.removeClass('add sign');
		icon.addClass('minus sign');
		$(this).removeClass('green');
		$(this).addClass('red');

		$("#ranking-" + type).fadeIn();
	} else {
		icon.removeClass('minus sign');
		icon.addClass('add sign');
		$(this).removeClass('red');
		$(this).addClass('green');
	}
});

var Home = {
	init: function() {
		this.compile();
		this.load('pushevent');
	},

	compile: function() {
		this.language = undefined;
		this.data = {};
		this.tpl = Handlebars.compile($("#overview-ranking").html());
	},

	render: function(data) {
		for (var kind in data) {
			for (var type in data[kind]) {
				$("#" + kind + "-" + type).html(this.tpl(data[kind][type]));
				
				if (data[kind][type].length > 0) {
					$("#" + kind + "-top-" + type + "-name").html(data[kind][type][0].key);

					var value = data[kind][type][0].value;
					var value = type == "experience" ? value + " XP" : value + " pushes";
					$("#" + kind + "-top-" + type + "-value").html(value);

					if (type != "language") {
						$("#" + kind + "-top-" + type + "-link").attr("href", "/users/" + data[kind][type][0].key);
					}
				}
			}
		}
	},

	load: function(type) {
		var route = "/api/ranking/home/" + type;

		if (this.language !== undefined) {
			route = route + "/" + this.language;
		}

		$.get(route, function(data) {
			Home.render(data);
		});
	}
};

Home.init();
