var activeTab = window.location.hash !== "" ? window.location.hash.replace("#", "") : "general";
var beatmapsSectionLoaded = false;
var topLeaderOffset = 0;
var topScoreOffset = 0;

function pinScore(scoreId, userId) {
    var url = "/users/" + currentUser + "/pinned";
    performApiRequest("POST", url, { score_id: scoreId }, function () {
        loadPinnedScores(userId, modeName);
    });
}

function unpinScore(scoreId, userId) {
    var url = "/users/" + currentUser + "/pinned";
    performApiRequest("DELETE", url, { score_id: scoreId }, function () {
        loadPinnedScores(userId, modeName);
    });
}

function expandProfileTab(id, forceExpand) {
    var tab = document.getElementById(id);
    activeTab = id;

    if (!tab) {
        expandProfileTab(id.indexOf("score-top") === 0 ? "leader" : "general", forceExpand);
        return;
    }

    var isExpanded = tab.className.indexOf("expanded") !== -1;

    // If forceExpand is true, always expand; otherwise toggle
    if (forceExpand || !isExpanded) {
        if ($(tab).is(":hidden") || tab.style.height === "0px") {
            slideDown(tab);
            setTimeout(function () {
                tab.className += " expanded";
            }, 500);
        }

        if (forceExpand) {
            window.location.hash = "#" + activeTab;
        }
    } else {
        // Collapse the tab
        $(tab).removeClass("expanded");
        slideUp(tab);
    }

    if (activeTab === "general") {
        loadPerformanceGraph(userId, modeName);
    } else if (activeTab === "history") {
        loadPlaysGraph(userId, modeName);
        loadViewsGraph(userId, modeName);
    } else if (activeTab === "beatmaps") {
        loadBeatmapsTab(userId);
    }
}

function expandRecentActivity() {
    var recentPreview = document.getElementById("profile-recent-preview");
    var recentFull = document.getElementById("profile-recent-full");
    var general = document.getElementById("general");

    $(recentPreview).hide();
    $(recentFull).show();
    $(tab).slideDown(500);
}

function createScoreElement(score, index, type) {
    var scoreDiv = document.createElement("div");
    scoreDiv.id = "score-" + type + "-" + index;
    scoreDiv.className = "score";
    scoreDiv.onclick = function (e) {
        e = e || window.event;
        var target = e.target || e.srcElement;
        if ($(target).closest("a")[0]) return;
        if ($(target).closest(".score-pin-icon")[0]) return;
        if ($(target).closest(".score-pinned-icon")[0]) return;
        if ($(target).closest(".score-replay")[0]) return;
        window.location.href = "/scores/" + score.id;
    };

    var scoreTable = document.createElement("table");
    var tableBody = document.createElement("tbody");
    var tableRow = document.createElement("tr");

    var leftColumn = document.createElement("td");
    leftColumn.className = "score-left";

    var scoreGrade = document.createElement("img");
    scoreGrade.src = "/images/grades/" + score.grade + "_small.png";
    scoreGrade.loading = "lazy";

    var beatmapInfo = document.createElement("a");
    beatmapInfo.href = "/b/" + score.beatmap.id + "?mode=" + score.mode;
    $(beatmapInfo).text(
        score.beatmap.beatmapset.artist + " - " + score.beatmap.beatmapset.title + " [" + score.beatmap.version + "]"
    );

    var modsText = document.createElement("b");
    if (score.mods > 0) {
        $(modsText).text("+" + Mods.getString(score.mods));
    }

    var scoreInfo = document.createElement("b");
    scoreInfo.appendChild(beatmapInfo);
    scoreInfo.appendChild(modsText);

    var accuracyText = document.createTextNode("(" + (score.acc * 100).toFixed(2) + "%)");

    // Parse date to a format that timeago can understand
    var scoreDateString = formatDateTimeForTitle(score.submitted_at);

    var rightColumn = document.createElement("td");
    rightColumn.className = "score-right";

    var ppText = document.createElement("b");
    $(ppText).text(score.pp.toFixed(0) + "pp");

    var ppDisplay = document.createElement("div");
    ppDisplay.className = "pp-display";
    ppDisplay.appendChild(ppText);

    var ppWeightPercent = document.createElement("b");
    $(ppWeightPercent).text((Math.pow(0.95, index + topScoreOffset) * 100).toFixed(0) + "%");

    var ppWeight = document.createElement("div");
    ppWeight.className = "pp-display-weight";

    if (type == "top") {
        ppWeight.appendChild(document.createTextNode("weighted "));
        ppWeight.appendChild(ppWeightPercent);
        ppWeight.appendChild(
            document.createTextNode(" (" + (score.pp * Math.pow(0.95, index + topScoreOffset)).toFixed(0) + "pp)")
        );
    }

    if (!approvedRewards && score.beatmap.status > 2) {
        // Display heart icon for loved maps
        ppText.innerHTML = '<i class="icon-heart"></i>';
        ppText.title = score.pp.toFixed(0) + "pp (if ranked)";
        // Reset pp weight text
        ppWeight.innerHTML = type == "top" ? "weighted <b>0%</b> (0pp)" : "";
    }

    var iconContainer = document.createElement("div");
    iconContainer.className = "score-icon-container";

    var scoreInfoDiv = document.createElement("div");
    scoreInfoDiv.appendChild(scoreGrade);
    scoreInfoDiv.appendChild(scoreInfo);
    scoreInfoDiv.appendChild(accuracyText);

    var scoreBottomDiv = document.createElement("div");
    scoreBottomDiv.className = "score-bottom";

    // Score's Date
    var dateText = document.createElement("time");
    dateText.setAttribute("datetime", score.submitted_at);
    $(dateText).text(score.submitted_at);
    dateText.title = scoreDateString;
    dateText.className = "timeago";

    scoreBottomDiv.appendChild(dateText);

    // Score's Client Version
    versionText = score.client_string ? score.client_string : undefined;

    if (versionText != undefined) {
        var clientText = document.createElement("div");
        clientText.className = "score-version";
        clientText.innerHTML += " &mdash; on ";

        var clientTextVersion = document.createElement("span");
        clientTextVersion.className = "score-version-number";
        $(clientTextVersion).text(versionText);

        clientText.appendChild(clientTextVersion);
        scoreBottomDiv.appendChild(clientText);
    }

    var replayLink = document.createElement("a");
    replayLink.href = "/scores/" + score.id + "/download";
    replayLink.className = "score-replay";
    replayLink.title = "Download Replay";
    replayLink.target = "_blank";

    var replayIcon = document.createElement("i");
    replayIcon.className = "icon-download-alt";
    replayLink.appendChild(replayIcon);

    iconContainer.appendChild(replayLink);

    if (currentUser === userId) {
        var pinIcon = document.createElement("i");
        pinIcon.className = "icon-star score-pin-" + score.id;

        if (!score.pinned) {
            pinIcon.className += " score-pin-icon";
            pinIcon.title = "Pin Score";
            pinIcon.onclick = function () {
                var icons = $(document)
                    .find(".score-pin-" + score.id)
                    .toArray();
                for (var j = 0; j < icons.length; j++) {
                    $(icons[j]).removeClass("score-pin-icon");
                    $(icons[j]).addClass("score-pinned-icon");
                    icons[j].title = "Unpin Score";
                }
                pinScore(score.id, userId);
                pinIcon.onclick = function () {
                    unpinScore(score.id, userId);
                    $(pinIcon).removeClass("score-pinned-icon");
                    $(pinIcon).addClass("score-pin-icon");
                    pinIcon.title = "Pin Score";
                };
            };
        } else {
            pinIcon.className += " score-pinned-icon";
            pinIcon.title = "Unpin Score";
            pinIcon.onclick = function () {
                var icons = $(document)
                    .find(".score-pin-" + score.id)
                    .toArray();
                for (var j = 0; j < icons.length; j++) {
                    $(icons[j]).removeClass("score-pinned-icon");
                    $(icons[j]).addClass("score-pin-icon");
                    icons[j].title = "Pin Score";
                }
                unpinScore(score.id, userId);
                pinIcon.onclick = function () {
                    pinScore(score.id, userId);
                    $(pinIcon).removeClass("score-pin-icon");
                    $(pinIcon).addClass("score-pinned-icon");
                    pinIcon.title = "Unpin Score";
                };
            };
        }

        iconContainer.appendChild(pinIcon);
    }

    ppWeight.appendChild(iconContainer);
    leftColumn.appendChild(scoreInfoDiv);
    leftColumn.appendChild(scoreBottomDiv);
    rightColumn.appendChild(ppDisplay);
    rightColumn.appendChild(ppWeight);
    tableRow.appendChild(leftColumn);
    tableRow.appendChild(rightColumn);
    tableBody.appendChild(tableRow);
    scoreTable.appendChild(tableBody);
    scoreDiv.appendChild(scoreTable);
    return scoreDiv;
}

function loadPinnedScores(userId, mode) {
    var url = "/users/" + userId + "/pinned/" + mode;
    var scoreContainer = document.getElementById("pinned-scores");

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var data = JSON.parse(xhr.responseText);
            var loadingText = document.getElementById("pinned-scores-loading");
            var scores = data.scores;

            if (loadingText) {
                $($(loadingText).parent()[0]).removeClass("score");
                $(loadingText).remove();
            }

            // Reset container
            scoreContainer.innerHTML = "<h2>Pinned Scores</h2>";

            if (data.total <= 0) {
                scoreContainer.appendChild(document.createTextNode("This player has not pinned any scores yet :("));
                return;
            }

            // Update total score count
            var heading = document.getElementById("pinned-scores").getElementsByTagName("h2")[0];
            heading.innerHTML = "Pinned Scores (" + data.total.toLocaleString() + ")";

            for (var index = 0; index < scores.length; index++) {
                var score = scores[index];
                var scoreDiv = createScoreElement(score, index, "pinned");
                scoreContainer.appendChild(scoreDiv);
            }

            // Render timeago elements
            renderTimeagoElements();
        },
        function (xhr) {
            var errorText = document.createElement("p");
            $(errorText).text("Failed to load pinned scores.");
            $(errorText).addClass("score");
            scoreContainer.appendChild(errorText);

            var loadingText = document.getElementById("pinned-scores-loading");
            if (loadingText) {
                $($(loadingText).parent()[0]).removeClass("score");
                $(loadingText).remove();
            }
        }
    );
}

function loadTopPlays(userId, mode, limit, offset) {
    var url = "/users/" + userId + "/top/" + mode + "?limit=" + limit + "&offset=" + offset;
    var scoreContainer = document.getElementById("top-scores");

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var data = JSON.parse(xhr.responseText);
            var loadingText = document.getElementById("top-scores-loading");
            var scores = data.scores;

            if (loadingText) {
                $($(loadingText).parent()[0]).removeClass("score");
                $(loadingText).remove();
            }

            if (data.total <= 0) {
                var noScoresText = document.createElement("p");
                $(noScoresText).text("No awesome performance records yet :(");
                scoreContainer.appendChild(noScoresText);
                return;
            }

            // Update total score count
            var heading = document.getElementById("top-scores").getElementsByTagName("h2")[0];
            heading.innerHTML = "Best Performance (" + data.total.toLocaleString() + ")";

            for (var index = 0; index < scores.length; index++) {
                var score = scores[index];
                if (score.beatmap.status > 2 && !approvedRewards) {
                    continue;
                }

                var scoreDiv = createScoreElement(score, index, "top");
                scoreContainer.appendChild(scoreDiv);
            }
            topScoreOffset += scores.length;

            // Render timeago elements
            renderTimeagoElements();

            if (scores.length >= limit) {
                // Create show more text
                var showMoreText = document.createElement("b");
                $(showMoreText).text("Show me more!");

                // Add onclick event
                var showMoreHref = document.createElement("a");
                showMoreHref.href = "#score-top-" + scores.length;
                showMoreHref.id = "show-more-top";
                showMoreHref.appendChild(showMoreText);
                showMoreHref.onclick = function () {
                    var loadingText = document.createElement("p");
                    $(loadingText).text("Loading...");
                    loadingText.id = "top-scores-loading";

                    var showMore = document.getElementById("show-more-top");
                    $(showMore).parent()[0].appendChild(loadingText);
                    $(showMore).remove();

                    loadTopPlays(userId, modeName, 50, topScoreOffset);
                };

                // Create wrapper that contains styling
                var showMore = document.createElement("div");
                showMore.className = "score show-more";
                showMore.appendChild(showMoreHref);

                // Append show more text to container
                scoreContainer.appendChild(showMore);
            }
        },
        function (xhr) {
            var errorText = document.createElement("p");
            $(errorText).text("Failed to load top plays.");
            $(errorText).addClass("score");
            scoreContainer.appendChild(errorText);

            var loadingText = document.getElementById("top-scores-loading");

            if (loadingText) {
                $($(loadingText).parent()[0]).removeClass("score");
                $(loadingText).remove();
            }
        }
    );

    return false;
}

function loadLeaderScores(userId, mode, limit, offset) {
    var url = "/users/" + userId + "/first/" + mode + "?limit=" + limit + "&offset=" + offset;
    var scoreContainer = document.getElementById("leader-scores");

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var data = JSON.parse(xhr.responseText);
            var loadingText = document.getElementById("leader-scores-loading");
            var scores = data.scores;

            if (loadingText) {
                $($(loadingText).parent()[0]).removeClass("score");
                $(loadingText).remove();
            }

            if (data.total <= 0) {
                var noScoresText = document.createElement("p");
                $(noScoresText).text("No first place records currently :(");
                scoreContainer.appendChild(noScoresText);
                return;
            }

            // Update count & rank
            var heading = document.getElementById("leader-scores").getElementsByTagName("h2")[0];
            var leaderScoresText = "First Place Ranks";
            var total = " (" + data.total.toLocaleString() + ")";
            var rank = " - #" + firstsRank.toLocaleString();
            if (data.total <= 0) {
                rank = "";
            }

            $(heading).text(leaderScoresText + total + rank);

            for (var i = 0; i < scores.length; i++) {
                var scoreDiv = createScoreElement(scores[i], i, "leader");
                scoreContainer.appendChild(scoreDiv);
            }
            topLeaderOffset += scores.length;

            // Render timeago elements
            renderTimeagoElements();

            if (scores.length >= limit) {
                var showMoreText = document.createElement("b");
                $(showMoreText).text("Show me more!");

                // Add onclick event
                var showMoreHref = document.createElement("a");
                showMoreHref.href = "#score-leader-" + scores.length;
                showMoreHref.id = "show-more-leader";
                showMoreHref.appendChild(showMoreText);
                showMoreHref.onclick = function () {
                    var loadingText = document.createElement("p");
                    $(loadingText).text("Loading...");
                    loadingText.id = "leader-scores-loading";

                    var showMore = document.getElementById("show-more-leader");
                    $(showMore).parent()[0].appendChild(loadingText);
                    $(showMore).remove();

                    loadLeaderScores(userId, modeName, 50, topLeaderOffset);
                };

                // Create wrapper that contains styling
                var showMore = document.createElement("div");
                showMore.className = "score show-more";
                showMore.appendChild(showMoreHref);

                // Append show more text to container
                scoreContainer.appendChild(showMore);
            }
        },
        function (xhr) {
            var errorText = document.createElement("p");
            $(errorText).text("Failed to load first place ranks.");
            $(errorText).addClass("score");
            scoreContainer.appendChild(errorText);

            var loadingText = document.getElementById("leader-scores-loading");
            if (loadingText) {
                $($(loadingText).parent()[0]).removeClass("score");
                $(loadingText).remove();
            }
        }
    );

    return false;
}

function loadMostPlayed(userId, limit, offset) {
    var loadingText = document.getElementById("plays-loading");

    if (!loadingText) return;

    var url = "/users/" + userId + "/plays" + "?limit=" + limit + "&offset=" + offset;
    var playsContainer = document.getElementById("plays-container");

    performApiRequest("GET", url, null, function (xhr) {
        var plays = JSON.parse(xhr.responseText);
        if (plays.length <= 0) {
            playsContainer.appendChild(document.createTextNode("This player has not played anything yet :("));
            return;
        }

        for (var index = 0; index < plays.length; index++) {
            var item = plays[index];
            var beatmapLink = document.createElement("a");
            $(beatmapLink).text(
                item.beatmap.beatmapset.artist +
                    " - " +
                    item.beatmap.beatmapset.title +
                    " [" +
                    item.beatmap.version +
                    "]"
            );
            beatmapLink.href = "/b/" + item.beatmap.id;

            var playsDiv = document.createElement("div");
            playsDiv.style.fontSize = 180 * Math.pow(0.95, index + 1) + "%";
            playsDiv.style.margin = "2.5px";
            playsDiv.appendChild(document.createTextNode(item.count + " plays - "));
            playsDiv.appendChild(beatmapLink);

            playsContainer.appendChild(playsDiv);
        }
    });

    $(loadingText).remove();
}

function loadRecentPlays(userId, mode) {
    var loadingText = document.getElementById("recent-loading");

    if (!loadingText) return;

    var url = "/users/" + userId + "/recent/" + mode;
    var playsContainer = document.getElementById("recent-container");

    performApiRequest("GET", url, null, function (xhr) {
        var response = JSON.parse(xhr.responseText);
        var scores = response.scores;
        if (response.total <= 0) {
            playsContainer.appendChild(document.createTextNode("No recent scores set by this player :("));
            return;
        }

        for (var i = 0; i < scores.length; i++) {
            var score = scores[i];

            // Parse date to a format that timeago can understand
            var scoreDateString = formatDateTimeForTitle(score.submitted_at);

            var dateText = document.createElement("time");
            dateText.setAttribute("datetime", score.submitted_at);
            $(dateText).text(score.submitted_at);
            dateText.title = scoreDateString;
            dateText.className += " timeago";

            var beatmapLink = document.createElement("a");
            $(beatmapLink).text(
                score.beatmap.beatmapset.artist +
                    " - " +
                    score.beatmap.beatmapset.title +
                    " [" +
                    score.beatmap.version +
                    "]"
            );
            beatmapLink.href = "/b/" + score.beatmap.id;

            var modsText = "";

            if (score.mods > 0) modsText = "+" + Mods.getString(score.mods);

            var scoreDiv = document.createElement("div");
            scoreDiv.appendChild(dateText);
            scoreDiv.appendChild(document.createTextNode(" - "));
            scoreDiv.appendChild(beatmapLink);
            scoreDiv.appendChild(
                document.createTextNode(
                    " " + Math.round(score.pp).toLocaleString("en-US") + "pp (" + score.grade + ") " + modsText
                )
            );
            scoreDiv.style.margin = "2.5px";

            playsContainer.appendChild(scoreDiv);
        }

        // Render timeago elements
        renderTimeagoElements();
    });

    loadingText.parentNode.removeChild(loadingText);
}

function getDaysAgo(date) {
    var currentDate = new Date();
    currentDate.setHours(0, 0, 0, 0);
    date.setHours(0, 0, 0, 0);
    return Math.ceil((currentDate.valueOf() - date.valueOf()) / 86400 / 1000);
}

function processRankEntries(entries, rankingType) {
    var bestEntryByDate = [];
    var bestWasLast = false;
    var best = null;

    for (var i = 0; i < entries.length; i++) {
        var entry = {
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
        countryRankValues.unshift({
            x: 0,
            y: -countryRank
        });
        globalRankValues.unshift({
            x: 0,
            y: -globalRank
        });
        scoreRankValues.unshift({
            x: 0,
            y: -scoreRank
        });
        ppv1RankValues.unshift({
            x: 0,
            y: -ppv1Rank
        });
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
        {
            values: globalRankValues,
            key: "Global Rank",
            color: "#ff7f0e"
        },
        {
            values: countryRankValues,
            key: "Country Rank",
            color: "#0ec7ff",
            disabled: true
        },
        {
            values: scoreRankValues,
            key: "Score Rank",
            color: "#d30eff",
            disabled: true
        },
        {
            values: ppv1RankValues,
            key: "PPv1 Rank",
            color: "#51f542",
            disabled: true
        }
    ];
}

function resetPerformanceGraph() {
    var $rankGraph = $("#rank-graph svg");
    if ($rankGraph.length == 0) {
        return;
    }

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

    resetPerformanceGraph();
    var url = "/users/" + userId + "/history/rank/" + mode;

    performApiRequest("GET", url, null, function (xhr) {
        var entries = JSON.parse(xhr.responseText);
        var rankData = processRankHistory(entries);

        if (!entries || entries.length <= 0) {
            // Handle no data case
            return;
        }

        // Used for initial y-axis calculation
        var defaultSelectedRankType = "global_rank";

        nv.addGraph(function () {
            var chart = nv.models
                .lineChart()
                .margin({
                    left: 80,
                    bottom: 20,
                    right: 50
                })
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
        return {
            x: -elapsedMonths,
            y: entry.plays
        };
    });

    values.sort(function (a, b) {
        return a.x - b.x;
    });

    return [
        {
            values: values,
            key: "Plays",
            color: "#f5f242",
            area: true
        }
    ];
}

function loadPlaysGraph(userId, mode) {
    if (typeof nv === "undefined" || !nv.addGraph) return;

    var url = "/users/" + userId + "/history/plays/" + mode;

    performApiRequest("GET", url, null, function (xhr) {
        var entries = JSON.parse(xhr.responseText);
        var playData = processPlayHistory(entries);

        nv.addGraph(function () {
            var chart = nv.models
                .lineChart()
                .margin({
                    left: 80,
                    bottom: 20,
                    right: 50
                })
                .useInteractiveGuideline(true)
                .transitionDuration(250)
                .interpolate("linear")
                .showLegend(false)
                .showYAxis(true)
                .showXAxis(true);

            chart.xAxis.axisLabel("Months").tickFormat(function (month) {
                if (month % 1 !== 0) return "";
                if (month == 0) return "This Month";
                if (month > 0) return month != 1 ? "In " + month + " months" : "In " + month + " month";
                return month != -1 ? -month + " months ago" : -month + " month ago";
            });

            chart.yAxis.axisLabel("Plays").tickFormat(function (plays) {
                var rounded = Math.round(plays);
                if (plays !== rounded) return "";
                return rounded.toString();
            });

            // Calculate proper Y-axis range
            var playMultiplier = 2;
            var maxPlays = 0;

            if (playData[0] && playData[0].values) {
                for (var i = 0; i < playData[0].values.length; i++) {
                    if (playData[0].values[i].y > maxPlays) {
                        maxPlays = playData[0].values[i].y;
                    }
                }

                // Makes graph look better if user has a single "play" entry
                playMultiplier = 2 / playData[0].values.length;
            }

            // Force Y-axis to start at 0 and use nice round numbers
            chart.forceY([0, maxPlays > 0 ? maxPlays * playMultiplier : 10]);

            d3.select("#play-graph svg").datum(playData).call(chart);

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

function processViewsHistory(entries) {
    var currentDate = new Date();
    var currentYear = currentDate.getFullYear();
    var currentMonth = currentDate.getMonth();

    var values = $.map(entries, function (entry, i) {
        var elapsedMonths = (currentYear - entry.year) * 12 + (currentMonth - (entry.month - 1));
        return {
            x: -elapsedMonths,
            y: entry.replay_views
        };
    });

    values.sort(function (a, b) {
        return a.x - b.x;
    });

    return [
        {
            values: values,
            key: "Replay Views",
            color: "#f78e25",
            area: true
        }
    ];
}

function loadViewsGraph(userId, mode) {
    if (typeof nv === "undefined" || !nv.addGraph) return;

    var url = "/users/" + userId + "/history/views/" + mode;

    performApiRequest("GET", url, null, function (xhr) {
        var entries = JSON.parse(xhr.responseText);
        var viewsData = processViewsHistory(entries);

        nv.addGraph(function () {
            var chart = nv.models
                .lineChart()
                .margin({
                    left: 80,
                    bottom: 20,
                    right: 50
                })
                .useInteractiveGuideline(true)
                .transitionDuration(250)
                .interpolate("linear")
                .showLegend(false)
                .showYAxis(true)
                .showXAxis(true);

            chart.xAxis.axisLabel("Months").tickFormat(function (month) {
                if (month % 1 !== 0) return "";
                if (month == 0) return "This Month";
                if (month > 0) return month != 1 ? "In " + month + " months" : "In " + month + " month";
                return month != -1 ? -month + " months ago" : -month + " month ago";
            });

            chart.yAxis.axisLabel("Views").tickFormat(function (views) {
                var rounded = Math.round(views);
                if (views !== rounded) return "";
                return rounded.toString();
            });

            // Calculate proper Y-axis range
            var viewMultiplier = 2;
            var maxViews = 0;

            if (viewsData[0] && viewsData[0].values) {
                for (var i = 0; i < viewsData[0].values.length; i++) {
                    if (viewsData[0].values[i].y > maxViews) {
                        maxViews = viewsData[0].values[i].y;
                    }
                }

                // Makes graph look better if user has a single "view" entry
                viewMultiplier = 2 / viewsData[0].values.length;
            }

            // Force Y-axis to start at 0 and use nice round numbers
            chart.forceY([0, maxViews > 0 ? maxViews * viewMultiplier : 10]);

            d3.select("#replay-graph svg").datum(viewsData).call(chart);

            nv.utils.windowResize(function () {
                chart.update();
            });

            // Reset "dy" value
            var textElements = document.querySelectorAll(".nv-noData");
            for (var i = 0; i < textElements.length; i++) {
                textElements[i].setAttribute("dy", 0);
            }

            return chart;
        });
    });
}

function updatePlaystyleElement(element) {
    performApiRequest(
        $(element).toggleClass("playstyle-using") ? "POST" : "DELETE",
        "/users/" + userId + "/playstyle",
        {
            playstyle: element.id
        }
    );
}

function addFriend() {
    if (!isLoggedIn()) return;

    performApiRequest("POST", "/account/friends?id=" + userId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var targetAdded = data.status === "mutual" || superFriendly;
        var friendStatus = document.getElementById("friend-status");

        friendStatus.className = "friend-current-true-target-" + targetAdded;
    });

    return false;
}

function removeFriend() {
    if (!isLoggedIn()) return;

    performApiRequest("DELETE", "/account/friends?id=" + userId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var targetAdded = data.status === "mutual" || superFriendly;
        var friendStatus = document.getElementById("friend-status");

        friendStatus.className = "friend-current-false-target-" + targetAdded;
    });

    return false;
}

function removeFavourite(setId) {
    var url = "/users/" + currentUser + "/favourites/" + setId;

    performApiRequest("DELETE", url, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var favouritesCount = document.getElementById("beatmaps-favourites-count");
        var favouritesListContainer = document.getElementById("beatmaps-favourites-container");
        var container = document.getElementById("favourite-" + setId);
        var beatmapsContainer = document.getElementById("beatmaps");

        if (!container) return;

        container.style.opacity = 0;

        setTimeout(function () {
            $(container).remove();

            if (favouritesCount && data && data.length != null) $(favouritesCount).text(data.length.toString());

            if (data.length == 0) {
                // User has no favourite beatmaps anymore
                if (favouritesListContainer) {
                    favouritesListContainer.innerHTML = "";
                    var textElement = document.createElement("p");
                    textElement.style.margin = "5px";
                    textElement.innerHTML = "This player has no favourite beatmaps :(";
                    favouritesListContainer.appendChild(textElement);
                }
            }
        }, 350);

        setTimeout(function () {
            slideUp(beatmapsContainer);
            setTimeout(function () {
                slideDown(beatmapsContainer);
            }, 150);
        }, 400);
    });

    return false;
}

function deleteBeatmap(setId) {
    var proceed = confirm("Are you sure you want to delete this beatmap?");
    var url = "/users/" + currentUser + "/beatmapsets/" + setId;

    if (!proceed) return;

    performApiRequest(
        "DELETE",
        url,
        null,
        function (xhr) {
            reloadPageSoon(250, "/u/" + currentUser + "#beatmaps");
        },
        function (xhr) {
            try {
                var data = JSON.parse(xhr.responseText);
                alert(data.details);
            } catch (e) {
                console.error("Failed to delete beatmap:", e);
                alert("Failed to delete beatmap.");
            }
        }
    );
}

function reviveBeatmap(setId) {
    var url = "/users/" + currentUser + "/beatmapsets/" + setId + "/revive";
    performApiRequest(
        "POST",
        url,
        null,
        function (xhr) {
            reloadPageSoon(250, "/u/" + currentUser + "#beatmaps");
        },
        function (xhr) {
            try {
                var data = JSON.parse(xhr.responseText);
                alert(data.details);
            } catch (e) {
                console.error("Failed to revive beatmap:", e);
                alert("Failed to revive beatmap.");
            }
        }
    );
}

function toggleBeatmapContainer(section) {
    var container = $(section).find(".profile-beatmaps-container").first()[0];
    var beatmapsSection = document.getElementById("beatmaps");

    if (container.style.display === "none") {
        container.style.display = "block";
        $(beatmapsSection).slideDown(500);
    } else {
        container.style.display = "none";
        $(beatmapsSection).css("height", "auto");
    }
}

function loadBeatmapsTab(userId) {
    if (beatmapsSectionLoaded) return;

    beatmapsSectionLoaded = true;
    loadUserCreatedBeatmapsets(userId);
    loadUserNominatedBeatmapsets(userId);
    loadUserFavouriteBeatmapsets(userId);
}

function loadUserCreatedBeatmapsets(userId) {
    var url = "/users/" + userId + "/beatmapsets";
    var beatmapSections = document.getElementById("user-beatmaps-sections");
    var beatmapLoadingText = document.getElementById("user-beatmaps-loading");

    if (!beatmapSections) return;

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var beatmapsets = JSON.parse(xhr.responseText) || [];
            beatmapSections.innerHTML = "";

            var categories = {
                Ranked: [],
                Loved: [],
                Qualified: [],
                Pending: [],
                WIP: [],
                Graveyarded: []
            };

            for (var i = 0; i < beatmapsets.length; i++) {
                var bm = beatmapsets[i];

                if (!bm) continue;

                if (bm.status === BeatmapStatus.Inactive) continue;

                if (bm.status === BeatmapStatus.Loved) {
                    categories["Loved"].push(bm);
                } else if (bm.status === BeatmapStatus.Qualified) {
                    categories["Qualified"].push(bm);
                } else if (bm.status === BeatmapStatus.Ranked || bm.status === BeatmapStatus.Approved) {
                    categories["Ranked"].push(bm);
                } else if (bm.status === BeatmapStatus.Pending) {
                    categories["Pending"].push(bm);
                } else if (bm.status === BeatmapStatus.WIP) {
                    categories["WIP"].push(bm);
                } else {
                    categories["Graveyarded"].push(bm);
                }
            }

            var order = ["Ranked", "Loved", "Qualified", "Pending", "WIP", "Graveyarded"];

            var unrankedStatuses = ["Pending", "WIP", "Graveyarded"];

            for (var j = 0; j < order.length; j++) {
                var category = order[j];
                var items = categories[category] || [];
                if (items.length === 0) continue;

                var section = createBeatmapsSectionElement(
                    " " + category + " Beatmaps",
                    items.length,
                    "user-beatmaps-" + category.toLowerCase()
                );
                var container = $(section).find(".profile-beatmaps-container").first()[0];

                for (var k = 0; k < items.length; k++) {
                    var set = items[k];
                    var showOwnerControls = currentUser === userId && $.inArray(category, unrankedStatuses) !== -1;
                    var table = createBeatmapsetTable(set, {
                        showOwnerControls: showOwnerControls,
                        showRevive: showOwnerControls && category === "Graveyarded"
                    });
                    container.appendChild(table);
                }

                var clear = document.createElement("div");
                clear.style.clear = "both";
                container.appendChild(clear);
                beatmapSections.appendChild(section);
            }

            if (beatmapSections.children.length === 0) {
                beatmapSections.innerHTML = "";
            }
        },
        function (xhr) {
            if (beatmapLoadingText) $(beatmapLoadingText).remove();
        }
    );
}

function loadUserFavouriteBeatmapsets(userId) {
    var url = "/users/" + userId + "/favourites";
    var favouritesCount = document.getElementById("beatmaps-favourites-count");
    var favouritesContainer = document.getElementById("beatmaps-favourites-container");

    if (!favouritesContainer) return;

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var favourites = JSON.parse(xhr.responseText) || [];

            if (favouritesCount) $(favouritesCount).text(favourites.length.toString());

            favouritesContainer.innerHTML = "";

            if (favourites.length === 0) {
                var empty = document.createElement("p");
                empty.style.margin = "5px";
                empty.innerHTML = "This player has no favourite beatmaps :(";
                favouritesContainer.appendChild(empty);
                return;
            }

            favourites.sort(function (a, b) {
                var at = new Date(a.created_at).valueOf();
                var bt = new Date(b.created_at).valueOf();
                return at - bt;
            });

            for (var i = 0; i < favourites.length; i++) {
                var fav = favourites[i];
                if (!fav || !fav.beatmapset) continue;

                var setId = fav.beatmapset.id;
                var table = createBeatmapsetTable(fav.beatmapset, {
                    tableId: "favourite-" + setId
                });

                if (currentUser === userId) {
                    var leftCell = $(table).find("td").first()[0];
                    if (leftCell) {
                        leftCell.appendChild(document.createElement("br"));
                        var deleteWrap = document.createElement("div");
                        deleteWrap.style.cssFloat = "right";
                        deleteWrap.style.textAlign = "right";
                        deleteWrap.style.paddingRight = "2.5px";

                        var deleteLink = document.createElement("a");
                        deleteLink.onclick = function () {
                            return removeFavourite(setId);
                        };
                        $(deleteLink).text("delete");
                        deleteWrap.appendChild(deleteLink);
                        leftCell.appendChild(deleteWrap);
                    }
                }

                favouritesContainer.appendChild(table);
            }

            var clear = document.createElement("div");
            clear.style.clear = "both";
            favouritesContainer.appendChild(clear);
        },
        function (xhr) {
            favouritesContainer.innerHTML = '<p style="margin:5px">Failed to load favourite beatmaps.</p>';
            if (favouritesCount) $(favouritesCount).text("0");
        }
    );
}

function loadUserNominatedBeatmapsets(userId) {
    var url = "/users/" + userId + "/nominations";
    var nominationsCount = document.getElementById("beatmaps-nominations-count");
    var nominationsSection = document.getElementById("beatmaps-nominations-section");
    var nominationsContainer = document.getElementById("beatmaps-nominations-container");

    if (!nominationsContainer || !nominationsSection) return;

    performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var nominations = JSON.parse(xhr.responseText) || [];

            if (nominations.length === 0) {
                $(nominationsSection).remove();
                return;
            }

            if (nominationsCount) $(nominationsCount).text(nominations.length.toString());

            nominationsContainer.innerHTML = "";

            for (var i = 0; i < nominations.length; i++) {
                var nomination = nominations[i];
                if (!nomination || !nomination.beatmapset) continue;

                var table = createBeatmapsetTable(nomination.beatmapset, {});
                nominationsContainer.appendChild(table);
            }

            var clear = document.createElement("div");
            clear.style.clear = "both";
            nominationsContainer.appendChild(clear);
        },
        function (xhr) {
            if (xhr && xhr.status === 404) {
                $(nominationsSection).remove();
                return;
            }

            nominationsContainer.innerHTML = '<p style="margin:5px">Failed to load nominated beatmaps.</p>';
            if (nominationsCount) $(nominationsCount).text("0");
        }
    );
}

function createBeatmapsSectionElement(titleText, count, containerId) {
    var section = document.createElement("div");
    var header = document.createElement("h2");
    header.className = "profile-stats-header";

    var link = document.createElement("a");
    link.onclick = function () {
        toggleBeatmapContainer(section);
    };
    link.innerHTML = titleText + ' (<span class="beatmaps-count">' + count + "</span>)";
    header.appendChild(link);

    var container = document.createElement("div");
    container.className = "profile-beatmaps-container";
    container.style.display = "none";
    container.id = containerId;

    section.appendChild(header);
    section.appendChild(container);
    return section;
}

function createBeatmapsetTable(beatmapset, options) {
    options = options || {};

    var table = document.createElement("table");
    table.className = "profile-beatmap";

    if (options.tableId) table.id = options.tableId;

    var tbody = document.createElement("tbody");
    var tr = document.createElement("tr");
    var tdLeft = document.createElement("td");

    var descriptionLink = document.createElement("a");
    descriptionLink.href = "/s/" + beatmapset.id;
    descriptionLink.className = "profile-beatmap-description";

    var bold = document.createElement("b");
    $(bold).text((beatmapset.artist || "Unknown") + " - " + (beatmapset.title || "Unknown"));
    descriptionLink.appendChild(bold);
    tdLeft.appendChild(descriptionLink);
    tdLeft.appendChild(document.createElement("br"));

    var creatorName = beatmapset.creator;
    var creatorId = beatmapset.creator_id;
    var server = beatmapset.server;

    if (creatorName) {
        if (server === 0) {
            var creatorLink = document.createElement("a");
            creatorLink.href = "https://osu.ppy.sh/u/" + creatorName;
            $(creatorLink).text("mapped by " + creatorName);
            tdLeft.appendChild(creatorLink);
        } else {
            var creatorLink = document.createElement("a");
            creatorLink.href = "/u/" + creatorId;
            $(creatorLink).text("mapped by " + creatorName);
            tdLeft.appendChild(creatorLink);
        }
    }

    if (options.showOwnerControls) {
        tdLeft.appendChild(document.createElement("br"));

        var controls = document.createElement("div");
        var deleteDiv = document.createElement("div");
        deleteDiv.style.cssFloat = "right";
        deleteDiv.style.textAlign = "right";
        deleteDiv.style.paddingRight = "2.5px";

        var deleteLink = document.createElement("a");
        deleteLink.onclick = function () {
            return deleteBeatmap(beatmapset.id);
        };
        $(deleteLink).text("delete");
        deleteDiv.appendChild(deleteLink);
        controls.appendChild(deleteDiv);

        if (options.showRevive) {
            var reviveDiv = document.createElement("div");
            reviveDiv.style.cssFloat = "right";
            reviveDiv.style.textAlign = "right";
            reviveDiv.style.marginRight = "5px";

            var reviveLink = document.createElement("a");
            reviveLink.onclick = function () {
                return reviveBeatmap(beatmapset.id);
            };
            $(reviveLink).text("revive");
            reviveDiv.appendChild(reviveLink);
            controls.appendChild(reviveDiv);
        }

        tdLeft.appendChild(controls);
    }

    var tdRight = document.createElement("td");
    tdRight.style.width = "80px";
    tdRight.style.height = "60px";

    var image = document.createElement("img");
    image.src = "/mt/" + beatmapset.id;
    image.alt = "";
    image.loading = "lazy";
    tdRight.appendChild(image);

    tr.appendChild(tdLeft);
    tr.appendChild(tdRight);
    tbody.appendChild(tr);
    table.appendChild(tbody);
    return table;
}

$(document).ready(function (event) {
    var beatmapContainers = $(".profile-beatmaps-container");
    for (var i = 0; i < beatmapContainers.length; i++) {
        beatmapContainers[i].style.display = "none";
    }

    expandProfileTab(activeTab);
    loadPinnedScores(userId, modeName);
    loadTopPlays(userId, modeName, 5, 0);
    loadLeaderScores(userId, modeName, 5, 0);
    loadRecentPlays(userId, modeName);
    loadMostPlayed(userId, 15, 0);
});
