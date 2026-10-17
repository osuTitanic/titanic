function bindBeatmapAudio() {
    $(".beatmap-image .icon-play, .beatmap-image .icon-pause").on("click", function (event) {
        event.preventDefault();

        var $icon = $(this);
        var $beatmapset = $icon.closest(".beatmapset");
        var audio = $beatmapset.find(".beatmap-preview-audio")[0];

        if (!audio) {
            return;
        }
        if (!audio.getAttribute("src")) {
            var source = audio.getAttribute("data-src");
            if (!source) {
                return;
            }
            audio.setAttribute("src", source);
        }

        $(".beatmap-preview-audio")
            .not(audio)
            .each(function () {
                if (this.paused) {
                    return;
                }

                this.pause();
                this.currentTime = 0;
                $(this).siblings(".beatmap-image").find("i").removeClass("icon-pause").addClass("icon-play");
            });

        resetOrPlayAudio(audio.id);
        audio.volume = getSavedVolume();
        $icon.toggleClass("icon-play", audio.paused).toggleClass("icon-pause", !audio.paused);
    });

    $(".beatmap-preview-audio").on("ended", function () {
        $(this).siblings(".beatmap-image").find("i").removeClass("icon-pause").addClass("icon-play");
    });
}

function bindFavourites() {
    $(".beatmap-favourite-link").each(function () {
        var beatmapsetId = $(this).closest(".beatmapset").attr("id");
        if (beatmapsetId) {
            markFavourite(beatmapsetId.replace("beatmapset-", ""), false);
        }
    });
}

function addFavorite(beatmapsetId) {
    var url = "/users/" + currentUser + "/favourites";

    performApiRequest(
        "POST",
        url,
        { set_id: beatmapsetId },
        function () {
            markFavourite(beatmapsetId, true);
        },
        function () {
            markFavourite(beatmapsetId, true);
        }
    );
}

function removeFavorite(beatmapsetId) {
    var url = "/users/" + currentUser + "/favourites/" + beatmapsetId;

    performApiRequest("DELETE", url, null, function () {
        markFavourite(beatmapsetId, false);
    });
}

function markFavourite(beatmapsetId, favourited) {
    var $link = $("#beatmapset-" + beatmapsetId + " .beatmap-favourite-link");
    if (!$link.length) {
        return;
    }

    $link
        .css("color", favourited ? "red" : "")
        .off("click.searchFavourite")
        .on("click.searchFavourite", function (event) {
            event.preventDefault();
            if (!currentUser) {
                showLoginForm();
                return;
            }
            if (favourited) {
                removeFavorite(beatmapsetId);
            } else {
                addFavorite(beatmapsetId);
            }
        });
}

function setPreviewVolume(volume) {
    saveVolume(volume);

    $(".beatmap-preview-audio").each(function () {
        this.volume = volume;
    });
    var percent = Math.round(volume * 100);

    $("#beatmap-volume-slider").val(percent);
    $("#beatmap-volume-control").attr("title", "Preview Volume: " + percent + "%");
}

function bindVolumeSlider() {
    var control = $("#beatmap-volume-control");
    var slider = $("#beatmap-volume-slider");
    if (!control.length || !slider.length) {
        // volume control not available
        return;
    }

    var initialVolume = getSavedVolume();
    setPreviewVolume(initialVolume);

    slider.on("input change", function () {
        var volume = parseInt(this.value, 10) / 100;
        setPreviewVolume(volume);
    });

    slider.on("mousedown touchstart", function () {
        control.addClass("active");
    });
    $(document).on("mouseup touchend", function () {
        control.removeClass("active");
    });

    // show & hide based on mouse hover
    control.hover(
        function () {
            $(this).stop().fadeTo(300, 1);
        },
        function () {
            if (!$(this).hasClass("active") && !$(this).find(":focus").length) {
                $(this).stop().fadeTo(300, 0);
            }
        }
    );
}

function getSavedVolume() {
    try {
        var saved = localStorage.getItem("beatmapPreviewVolume");
        if (saved === null) {
            return 0.35;
        }

        var parsed = parseFloat(saved);
        if (isNaN(parsed)) {
            return 0.35;
        }
        if (parsed < 0 || parsed > 1) {
            return 0.35;
        }
        return parsed;
    } catch (e) {
        // local storage not available
    }
    return 0.35;
}

function saveVolume(vol) {
    try {
        localStorage.setItem("beatmapPreviewVolume", vol);
    } catch (e) {
        // local storage not available
    }
}

$(function () {
    bindBeatmapAudio();
    bindFavourites();
    bindVolumeSlider();
});
