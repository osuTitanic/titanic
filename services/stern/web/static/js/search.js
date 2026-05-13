function bindBeatmapAudio() {
    $(".beatmap-image .icon-play, .beatmap-image .icon-pause").on("click", function (event) {
        event.preventDefault();

        var $icon = $(this);
        var $beatmapset = $icon.closest(".beatmapset");
        var audio = $beatmapset.find('audio[id^="beatmap-preview-"]')[0];

        if (!audio) {
            return;
        }

        $('audio[id^="beatmap-preview-"]')
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
        $icon.toggleClass("icon-play", audio.paused).toggleClass("icon-pause", !audio.paused);
    });

    $('audio[id^="beatmap-preview-"]').on("ended", function () {
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

$(function () {
    bindBeatmapAudio();
    bindFavourites();
});
