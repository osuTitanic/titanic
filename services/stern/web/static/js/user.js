// The active tab is derived from the url #<tab>, defaulting to the general tab.
// All other tabs are fetched on user interaction from their HTML partials.

var activeTab = window.location.hash !== "" ? window.location.hash.replace("#", "") : "general";
var rankGraphLoaded = false;

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

    // Otherwise, expand it
    if (id === "general") {
        loadPerformanceGraph(userId, modeName);
    }

    // Slide the tab open once its content is loaded
    loadTab(id, function () {
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
        window.location.hash = "#" + activeTab;
    }
}

function loadMoreActivity(link, offset) {
    var row = $(link).closest("tr");

    $.get("/partials/users/" + userId + "/activity?mode=" + mode + "&offset=" + offset, function (html) {
        row.replaceWith(html);
        renderTimeagoElements();
    });

    return false;
}

function loadMoreScores(link, section, offset) {
    var row = $(link).closest(".show-more");
    row.closest("b").innerText = "Loading...";

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
    if ($(target).closest("a")[0]) return false;
    if ($(target).closest(".score-replay")[0]) return false;
    if ($(target).closest(".score-pin-icon")[0]) return false;
    if ($(target).closest(".score-pinned-icon")[0]) return false;

    window.location.href = "/scores/" + scoreId;
    return false;
}

function togglePin(icon, scoreId, pinned) {
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

function updatePlaystyleElement(element) {
    var nowUsing = $(element).toggleClass("playstyle-using").hasClass("playstyle-using");
    performApiRequest(nowUsing ? "POST" : "DELETE", "/users/" + userId + "/playstyle", { playstyle: element.id });
}

function addFriend() {
    if (!isLoggedIn()) return false;

    performApiRequest("POST", "/account/friends?id=" + userId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var targetAdded = data.status === "mutual" || superFriendly;
        document.getElementById("friend-status").className = "friend-current-true-target-" + targetAdded;
    });

    return false;
}

function removeFriend() {
    if (!isLoggedIn()) return false;

    performApiRequest("DELETE", "/account/friends?id=" + userId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var targetAdded = data.status === "mutual" || superFriendly;
        document.getElementById("friend-status").className = "friend-current-false-target-" + targetAdded;
    });

    return false;
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
    if (typeof nv === "undefined" || !nv.addGraph) return;
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

$(document).ready(function () {
    expandProfileTab(activeTab);
});
