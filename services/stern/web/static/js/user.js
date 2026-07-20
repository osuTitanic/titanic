// The active tab is derived from the url #<tab>, defaulting to the general tab.
// All other tabs are fetched on user interaction from their HTML partials.

var activeTab = window.location.hash !== "" ? window.location.hash.replace("#", "") : "general";
var rankGraphLoaded = false;
var playsGraphLoaded = false;
var viewsGraphLoaded = false;

function scrollToTab(id) {
    var heading = document.getElementById("tab-" + id);
    if (!heading) return;
    heading.scrollIntoView();

    if (window.location.hash !== "#" + id) {
        window.location.hash = "#" + id;
        return;
    }
}

function loadTab(id, onReady) {
    var content = document.getElementById(id);
    if (!content) {
        if (onReady) onReady();
        return;
    }

    var partial = content.getAttribute("data-partial");
    if (content.getAttribute("data-loaded") === "true" || !partial) {
        if (onReady) onReady();
        return;
    }

    // Mark the tab as loaded before the request, to prevent duplicate requests
    content.setAttribute("data-loaded", "true");
    $(content).load(partial + "?mode=" + mode, function () {
        renderTimeagoElements();
        if (onReady) onReady();
    });
}

function expandProfileTab(id, forceExpand) {
    var tab = document.getElementById(id);
    if (!tab) {
        expandProfileTab("general", forceExpand);
        return;
    }

    activeTab = id;
    var isExpanded = tab.className.indexOf("expanded") !== -1;

    // If forceExpand is false and the tab is already expanded, collapse it
    if (!forceExpand && isExpanded) {
        $(tab).removeClass("expanded");
        slideUp(tab);
        return;
    }

    // Slide the tab open once its content is loaded, then render its graphs
    loadTab(id, function () {
        loadTabGraphs(id);

        if (!$(tab).is(":hidden") && tab.style.height !== "0px") {
            // already expanded the tab
            return;
        }
        slideDown(tab);
        setTimeout(function () {
            tab.className += " expanded";
        }, 500);
    });

    if (forceExpand) {
        scrollToTab(activeTab);
    }
}

function loadTabGraphs(id) {
    if (id === "general") {
        loadPerformanceGraph(userId, modeName);
    } else if (id === "history") {
        loadPlaysGraph(userId, modeName);
        loadViewsGraph(userId, modeName);
    }
}

function loadMoreActivity(link, offset) {
    var row = $(link).closest("tr");
    row.text("Loading...");

    $.get("/partials/users/" + userId + "/activity?mode=" + mode + "&offset=" + offset, function (html) {
        row.replaceWith(html);
        renderTimeagoElements();
    });

    return false;
}

function loadMoreScores(link, section, offset) {
    var row = $(link).closest(".show-more");
    row.find("b").text("Loading...");

    $.get(
        "/partials/users/" + userId + "/scores?section=" + section + "&offset=" + offset + "&mode=" + mode,
        function (html) {
            row.replaceWith(html);
            renderTimeagoElements();
        }
    );

    return false;
}

function scoreClicked(event, scoreId) {
    event = event || window.event;
    var target = event.target || event.srcElement;

    // Ignore clicks on the inner links, replay button and pin icons
    if ($(target).closest("a")[0]) return true;
    if ($(target).closest(".score-replay")[0]) return true;
    if ($(target).closest(".score-pin-icon")[0]) return true;
    if ($(target).closest(".score-pinned-icon")[0]) return true;

    window.location.href = "/scores/" + scoreId;
    return false;
}

function stopScoreClick(event) {
    event = event || window.event;
    if (!event) return;

    if (event.stopPropagation) {
        event.stopPropagation();
    }
    event.cancelBubble = true;
}

function togglePin(event, icon, scoreId, pinned) {
    stopScoreClick(event);

    if (!isLoggedIn()) return false;

    var method = pinned ? "DELETE" : "POST";

    performApiRequest(method, "/users/" + userId + "/pinned", { score_id: scoreId }, function () {
        // Reload the whole tab so the pinned section reflects the change
        var content = document.getElementById("leader");
        $(content).load("/partials/users/" + userId + "/leader?mode=" + mode, function () {
            renderTimeagoElements();
        });
    });

    return false;
}

function toggleBeatmapContainer(section) {
    var container = $(section).find(".profile-beatmaps-container").first()[0];
    var beatmapsSection = document.getElementById("beatmaps");

    if (!container) return;

    if (container.style.display === "none") {
        container.style.display = "block";
        $(beatmapsSection).slideDown(500);
    } else {
        container.style.display = "none";
        $(beatmapsSection).css("height", "auto");
    }
}

function reloadBeatmapsTab() {
    var content = document.getElementById("beatmaps");
    if (!content) return;
    $(content).load("/partials/users/" + userId + "/beatmaps", function () {
        renderTimeagoElements();
    });
}

function removeFavourite(setId) {
    if (!isLoggedIn()) return false;

    performApiRequest("DELETE", "/users/" + userId + "/favourites/" + setId, null, reloadBeatmapsTab);
    return false;
}

function deleteBeatmap(setId) {
    if (!isLoggedIn()) return false;
    if (!confirm("Are you sure you want to delete this beatmap?")) return false;

    performApiRequest(
        "DELETE",
        "/users/" + userId + "/beatmapsets/" + setId,
        null,
        reloadBeatmapsTab,
        apiErrorAlert
    );
    return false;
}

function reviveBeatmap(setId) {
    if (!isLoggedIn()) return false;

    performApiRequest(
        "POST",
        "/users/" + userId + "/beatmapsets/" + setId + "/revive",
        null,
        reloadBeatmapsTab,
        apiErrorAlert
    );
    return false;
}

function updatePlaystyleElement(element) {
    var nowUsing = $(element).toggleClass("playstyle-using").hasClass("playstyle-using");
    performApiRequest(nowUsing ? "POST" : "DELETE", "/users/" + userId + "/playstyle", { playstyle: element.id });
}

function updateFriendStatus(xhr, currentAdded) {
    var data = $.parseJSON(xhr.responseText);
    var targetAdded = data.status === "mutual" || superFriendly;
    $("#friend-status").attr("class", "friend-current-" + currentAdded + "-target-" + targetAdded);
}

function addFriendFromProfile() {
    return addFriend(
        userId,
        function (xhr) {
            updateFriendStatus(xhr, true);
        },
        function (xhr) {
            apiErrorAlert(xhr, "The user could not be added as a friend.");
        }
    );
}

function removeFriendFromProfile() {
    return removeFriend(
        userId,
        function (xhr) {
            updateFriendStatus(xhr, false);
        },
        function (xhr) {
            apiErrorAlert(xhr, "The user could not be removed from your friends.");
        }
    );
}

function processRankEntries(entries, rankingType) {
    var bestEntryByDate = [];
    var bestWasLast = false;
    var best = null;
    var entry = null;

    for (var i = 0; i < entries.length; i++) {
        entry = {
            daysAgo: getDaysAgo(new Date(entries[i].time)),
            value: entries[i][rankingType]
        };

        if (i == 0) {
            best = entry;
        } else if (entry.daysAgo == best.daysAgo) {
            if (entry.value < best.value) {
                best = entry;
            }
        } else {
            bestEntryByDate.push(best);
            bestWasLast = i == entries.length - 1;

            if (!bestWasLast) {
                best = entry;
            }
        }
    }

    if (best != null && !bestWasLast) {
        bestEntryByDate.push(entry);
    }

    return $.map(bestEntryByDate, function (entry, i) {
        var daysAgo = entry.daysAgo == 0 ? 0 : -entry.daysAgo;
        return {
            x: daysAgo,
            y: -entry.value
        };
    });
}

function processRankHistory(entries) {
    var globalRankValues = processRankEntries(entries, "global_rank");
    var scoreRankValues = processRankEntries(entries, "score_rank");
    var countryRankValues = processRankEntries(entries, "country_rank");
    var ppv1RankValues = processRankEntries(entries, "ppv1_rank");

    if (entries.length > 0) {
        countryRankValues.unshift({ x: 0, y: -countryRank });
        globalRankValues.unshift({ x: 0, y: -globalRank });
        scoreRankValues.unshift({ x: 0, y: -scoreRank });
        ppv1RankValues.unshift({ x: 0, y: -ppv1Rank });
    }

    ppv1RankValues = $.grep(ppv1RankValues, function (e, i) {
        return e.y != 0;
    });
    scoreRankValues = $.grep(scoreRankValues, function (e, i) {
        return e.y != 0;
    });
    globalRankValues = $.grep(globalRankValues, function (e, i) {
        return e.y != 0;
    });
    countryRankValues = $.grep(countryRankValues, function (e, i) {
        return e.y != 0;
    });

    ppv1RankValues = ppv1RankValues.reverse();
    scoreRankValues = scoreRankValues.reverse();
    globalRankValues = globalRankValues.reverse();
    countryRankValues = countryRankValues.reverse();

    return [
        { values: globalRankValues, key: "Global Rank", color: "#ff7f0e" },
        { values: countryRankValues, key: "Country Rank", color: "#0ec7ff", disabled: true },
        { values: scoreRankValues, key: "Score Rank", color: "#d30eff", disabled: true },
        { values: ppv1RankValues, key: "PPv1 Rank", color: "#51f542", disabled: true }
    ];
}

function resetPerformanceGraph() {
    var $rankGraph = $("#rank-graph svg");
    if ($rankGraph.length == 0) {
        return;
    }

    // nuke all child nodes of the graph
    while ($rankGraph[0].firstChild) $rankGraph[0].removeChild($rankGraph[0].firstChild);
}

function getPerformanceGraphRange(backupEntries, backupEntriesKey) {
    var rankMin = null;
    var rankMax = null;

    if (backupEntries == undefined) {
        var legendData = d3.selectAll(".nv-series").data();
        for (var i = 0; i < legendData.length; i++) {
            var legendElement = legendData[i];
            if (legendElement.disabled) {
                continue;
            }

            for (var j = 0; j < legendElement.values.length; j++) {
                var rankValue = Math.abs(legendElement.values[j].y);

                rankMin = rankMin == null ? rankValue : Math.min(rankMin, rankValue);
                rankMax = rankMax == null ? rankValue : Math.max(rankMax, rankValue);
            }
        }
    } else {
        for (var i = 0; i < backupEntries.length; i++) {
            var value = backupEntries[i][backupEntriesKey];

            rankMin = rankMin == null ? value : Math.min(rankMin, value);
            rankMax = rankMax == null ? value : Math.max(rankMax, value);
        }
    }

    return [rankMin, rankMax];
}

function updatePerformanceGraphYAxis(chart, range) {
    var userDigits = range[1].toString().length - 1;

    var minRankDigits = "1" + (userDigits > 0 ? userDigits * "0" : "");
    var relativeMinRank = Math.round(range[0] / minRankDigits) * minRankDigits;

    var maxRankDigits = "1" + userDigits * "0";
    var relativeMaxRank = Math.round(range[1] / maxRankDigits) * maxRankDigits;

    var betweenRank = relativeMaxRank - relativeMaxRank / 2;

    chart.yScale(d3.scale.linear().domain([-relativeMinRank - 1, -relativeMaxRank]));

    // Only display certain tick values
    chart.xAxis.tickValues([-90, -60, -30, 0]);
    chart.yAxis.tickValues([-relativeMaxRank, -betweenRank, -relativeMinRank]);

    // Force chart to show range between min/max ranks
    chart.forceY([-range[1] + 1, -range[0] - 1]);
}

function loadPerformanceGraph(userId, mode) {
    if (typeof nv === "undefined" || nv === undefined || !nv.addGraph) return;
    if (typeof d3 === "undefined" || d3 === undefined || !d3.selectAll) return;
    if (rankGraphLoaded) return;
    rankGraphLoaded = true;

    resetPerformanceGraph();
    var url = "/users/" + userId + "/history/rank/" + mode;

    performApiRequest("GET", url, null, function (xhr) {
        var entries = JSON.parse(xhr.responseText);

        if (!entries || entries.length <= 0) {
            // Allow the graph to be reloaded if there is no data yet
            rankGraphLoaded = false;
            return;
        }

        var rankData = processRankHistory(entries);

        // Used for initial y-axis calculation
        var defaultSelectedRankType = "global_rank";

        nv.addGraph(function () {
            var chart = nv.models
                .lineChart()
                .margin({ left: 80, bottom: 20, right: 50 })
                .useInteractiveGuideline(true)
                .transitionDuration(250)
                .interpolate("linear")
                .showLegend(true)
                .showYAxis(true)
                .showXAxis(true);

            chart.xAxis.axisLabel("Days").tickFormat(function (days) {
                if (days == 0) return "now";
                if (days > 0) return days != 1 ? "In " + days + " days" : "In " + days + " day";
                return days != -1 ? -days + " days ago" : -days + " day ago";
            });

            chart.yAxis.axisLabel("Rank").tickFormat(function (rank) {
                rank = Math.round(rank);
                if (rank >= 0) return "";
                return "#" + -rank;
            });

            chart.legend.dispatch.on("legendClick", function (state) {
                setTimeout(function () {
                    updatePerformanceGraphYAxis(chart, getPerformanceGraphRange());
                    chart.update();
                }, 0);
            });

            // Calculate the range of the y axis
            var range = getPerformanceGraphRange(entries, defaultSelectedRankType);

            // Update the y axis with the calculated range
            updatePerformanceGraphYAxis(chart, range);

            d3.select("#rank-graph svg").datum(rankData).call(chart);

            nv.utils.windowResize(function () {
                chart.update();
            });

            // Reset "dy" value
            var noDataElements = document.querySelectorAll(".nv-noData");
            for (var i = 0; i < noDataElements.length; i++) {
                noDataElements[i].setAttribute("dy", 0);
            }

            return chart;
        });
    });
}

function processPlayHistory(entries) {
    var currentDate = new Date();
    var currentYear = currentDate.getFullYear();
    var currentMonth = currentDate.getMonth();

    var values = $.map(entries, function (entry, i) {
        var elapsedMonths = (currentYear - entry.year) * 12 + (currentMonth - (entry.month - 1));
        return { x: -elapsedMonths, y: entry.plays };
    });

    values.sort(function (a, b) {
        return a.x - b.x;
    });

    return [{ values: values, key: "Plays", color: "#f5f242", area: true }];
}

function loadPlaysGraph(userId, mode) {
    if (typeof nv === "undefined" || !nv.addGraph) return;
    if (playsGraphLoaded) return;
    playsGraphLoaded = true;

    var url = "/users/" + userId + "/history/plays/" + mode;

    performApiRequest("GET", url, null, function (xhr) {
        var entries = JSON.parse(xhr.responseText);
        var playData = processPlayHistory(entries);

        nv.addGraph(function () {
            var chart = nv.models
                .lineChart()
                .margin({ left: 80, bottom: 20, right: 50 })
                .useInteractiveGuideline(true)
                .transitionDuration(250)
                .interpolate("linear")
                .showLegend(false)
                .showYAxis(true)
                .showXAxis(true);

            chart.xAxis.axisLabel("Months").tickFormat(monthTickFormat);

            chart.yAxis.axisLabel("Plays").tickFormat(integerTickFormat);

            // Force Y-axis to start at 0 and use nice round numbers
            chart.forceY([0, graphYAxisMax(playData)]);

            d3.select("#play-graph svg").datum(playData).call(chart);

            nv.utils.windowResize(function () {
                chart.update();
            });

            resetNoDataOffset();
            return chart;
        });
    });
}

function processViewsHistory(entries) {
    var currentDate = new Date();
    var currentYear = currentDate.getFullYear();
    var currentMonth = currentDate.getMonth();

    var values = $.map(entries, function (entry, i) {
        var elapsedMonths = (currentYear - entry.year) * 12 + (currentMonth - (entry.month - 1));
        return { x: -elapsedMonths, y: entry.replay_views };
    });

    values.sort(function (a, b) {
        return a.x - b.x;
    });

    return [{ values: values, key: "Replay Views", color: "#f78e25", area: true }];
}

function loadViewsGraph(userId, mode) {
    if (typeof nv === "undefined" || !nv.addGraph) return;
    if (viewsGraphLoaded) return;
    viewsGraphLoaded = true;

    var url = "/users/" + userId + "/history/views/" + mode;

    performApiRequest("GET", url, null, function (xhr) {
        var entries = JSON.parse(xhr.responseText);
        var viewsData = processViewsHistory(entries);

        nv.addGraph(function () {
            var chart = nv.models
                .lineChart()
                .margin({ left: 80, bottom: 20, right: 50 })
                .useInteractiveGuideline(true)
                .transitionDuration(250)
                .interpolate("linear")
                .showLegend(false)
                .showYAxis(true)
                .showXAxis(true);

            chart.xAxis.axisLabel("Months").tickFormat(monthTickFormat);

            chart.yAxis.axisLabel("Views").tickFormat(integerTickFormat);

            // Force Y-axis to start at 0 and use nice round numbers
            chart.forceY([0, graphYAxisMax(viewsData)]);

            d3.select("#replay-graph svg").datum(viewsData).call(chart);

            nv.utils.windowResize(function () {
                chart.update();
            });

            resetNoDataOffset();
            return chart;
        });
    });
}

function monthTickFormat(month) {
    if (month % 1 !== 0) return "";
    if (month == 0) return "This Month";
    if (month > 0) return month != 1 ? "In " + month + " months" : "In " + month + " month";
    return month != -1 ? -month + " months ago" : -month + " month ago";
}

function integerTickFormat(value) {
    var rounded = Math.round(value);
    if (value !== rounded) return "";
    return rounded.toString();
}

function graphYAxisMax(data) {
    var max = 0;
    var multiplier = 2;

    if (data[0] && data[0].values) {
        for (var i = 0; i < data[0].values.length; i++) {
            if (data[0].values[i].y > max) {
                max = data[0].values[i].y;
            }
        }
        // Makes the graph look better if the user has a single entry
        multiplier = 2 / data[0].values.length;
    }

    return max > 0 ? max * multiplier : 10;
}

function resetNoDataOffset() {
    var noDataElements = document.querySelectorAll(".nv-noData");
    for (var i = 0; i < noDataElements.length; i++) {
        noDataElements[i].setAttribute("dy", 0);
    }
}

$(document).ready(function () {
    expandProfileTab(activeTab);

    if (window.location.hash !== "") {
        scrollToTab(activeTab);
    }
});
