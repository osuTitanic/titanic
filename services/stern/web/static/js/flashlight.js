// Modified from https://gist.github.com/peppy/2276367
var fl = null;
var modToggle = null;

$("document").ready(function () {
    // Disable on mobile entirely
    var isMobile = window.innerWidth <= 768;
    if (isMobile) {
        return;
    }

    var today = new Date();
    var isAprilFools = today.getMonth() === 3 && today.getDate() === 1;

    // Remove settings if it's not april fools anymore
    if (!isAprilFools) {
        localStorage.removeItem("aprilFoolsSettings");
        return;
    }

    var settings = JSON.parse(localStorage.getItem("aprilFoolsSettings"));
    if (!settings) {
        // The user has not seen the april fools joke yet, so we
        // randomly decide whether to enable it for them or not
        var randomMode = Math.floor(Math.random() * 10);
        if (randomMode !== 0 && randomMode !== 1 && randomMode !== 3) {
            return;
        }
        settings = {
            flEnabled: true,
            hdEnabled: true
        };
        localStorage.setItem("aprilFoolsSettings", JSON.stringify(settings));
    }

    function saveSettings() {
        localStorage.setItem("aprilFoolsSettings", JSON.stringify(settings));
    }

    $("body").append("<div class='flashlight'></div><div class='modtoggle fl'></div>");
    fl = $(".flashlight");

    if (!settings.flEnabled) {
        fl.hide();
    }

    $("body").mousemove(function (e) {
        if (fl.is(":visible")) {
            fl.css("background-position", e.pageX - 1280 + "px " + (e.pageY - 720 - $(document).scrollTop()) + "px");
        }
    });

    $(".fl").click(function (e) {
        fl.toggle();
        settings.flEnabled = fl.is(":visible");
        saveSettings();
    });

    $("body").append("<div class='modtoggle hd'></div>");

    if (settings.hdEnabled) {
        $("a").addClass("hiddenA");
    }

    $(".hd").click(function () {
        $("a").toggleClass("hiddenA");
        settings.hdEnabled = $("a").hasClass("hiddenA");
        saveSettings();
    });

    modToggle = $(".modtoggle");
    if (modToggle.length > 0) {
        modToggle.css("opacity", 0.9);
        if (fl != null) {
            fl.css("opacity", 0.9);
        }
    }
});
