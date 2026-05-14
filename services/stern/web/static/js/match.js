var matchId = window.location.pathname.split("/")[2];
var lastEventTime = undefined;
var refreshRate = 6000;

var ScoringType = {
    0: "Score",
    1: "Accuracy",
    2: "Combo"
};

var TeamType = {
    0: "Head to Head",
    1: "Tag Co-op",
    2: "Team VS",
    3: "Tag Team VS"
};

var Team = {
    0: "None",
    1: "Blue",
    2: "Red"
};

function generateResultsTable(results, matchMods) {
    var table = document.createElement("table");
    var headerWrapper = document.createElement("thead");
    var header = document.createElement("tr");

    var headerPlace = document.createElement("th");
    headerPlace.style.width = "25px";
    var headerPlayer = document.createElement("th");
    headerPlayer.innerHTML = "Player";
    var headerScore = document.createElement("th");
    headerScore.innerHTML = "Score";
    var headerAccuracy = document.createElement("th");
    headerAccuracy.innerHTML = "Accuracy";
    var headerCombo = document.createElement("th");
    headerCombo.innerHTML = "Combo";
    var headerMods = document.createElement("th");
    headerMods.innerHTML = "Mods";
    var c300 = document.createElement("th");
    c300.innerHTML = "300s";
    var c100 = document.createElement("th");
    c100.innerHTML = "100s";
    var c50 = document.createElement("th");
    c50.innerHTML = "50s";
    var cMiss = document.createElement("th");
    cMiss.innerHTML = "Misses";

    header.appendChild(headerPlace);
    header.appendChild(headerPlayer);
    header.appendChild(headerScore);
    header.appendChild(headerAccuracy);
    header.appendChild(headerCombo);
    header.appendChild(headerMods);
    header.appendChild(c300);
    header.appendChild(c100);
    header.appendChild(c50);
    header.appendChild(cMiss);
    headerWrapper.appendChild(header);
    table.appendChild(headerWrapper);

    var tableBody = document.createElement("tbody");

    results.forEach(function (result) {
        var row = document.createElement("tr");

        if (result.score.failed) row.classList.add("fail");

        if (results.indexOf(result) % 2 == 0) row.classList.add("light-row");
        else row.classList.add("dark-row");

        var place = document.createElement("td");
        place.innerHTML = result.place;

        if (result.player.team != 0) place.classList.add(result.player.team == 1 ? "team-blue" : "team-red");

        var playerLink = document.createElement("a");
        playerLink.innerText = result.player.name;
        playerLink.href = "/u/" + result.player.id;

        var playerFlag = document.createElement("img");
        playerFlag.src = "/images/flags/" + result.player.country.toLowerCase() + ".gif";
        playerFlag.classList.add("flag");

        var player = document.createElement("td");
        player.appendChild(playerFlag);
        player.appendChild(playerLink);

        var score = document.createElement("td");
        score.innerHTML = result.score.score.toLocaleString();

        if (result.score.failed) {
            var failed = document.createElement("span");
            failed.style.color = "#ff0000";
            failed.innerHTML = " FAIL";
            score.style.fontWeight = "bold";
            score.appendChild(failed);
        }

        var c300 = document.createElement("td");
        c300.innerHTML = result.score.c300.toLocaleString();

        var c100 = document.createElement("td");
        c100.innerHTML = result.score.c100.toLocaleString();

        var c50 = document.createElement("td");
        c50.innerHTML = result.score.c50.toLocaleString();

        var cMiss = document.createElement("td");
        cMiss.innerHTML = result.score.cMiss.toLocaleString();

        var accuracy = document.createElement("td");
        accuracy.innerHTML = result.score.accuracy + "%";

        var combo = document.createElement("td");
        combo.innerHTML = result.score.max_combo;

        try {
            var mods = document.createElement("td");
            mods.innerHTML = Mods.getString(result.player.mods + matchMods);
        } catch (e) {
            console.warn("Failed to parse mods: " + result.player.mods);
            mods.innerHTML = "??";
        }

        row.appendChild(place);
        row.appendChild(player);
        row.appendChild(score);
        row.appendChild(accuracy);
        row.appendChild(combo);
        row.appendChild(mods);
        row.appendChild(c300);
        row.appendChild(c100);
        row.appendChild(c50);
        row.appendChild(cMiss);
        tableBody.appendChild(row);
    });

    table.appendChild(tableBody);

    return table;
}

function getTeamWinner(results, condition) {
    var teamResults = document.createElement("div");
    teamResults.classList.add("team-results");

    switch (condition) {
        case 0:
            // Score
            var blueScore = 0;
            var redScore = 0;

            results.forEach(function (result) {
                if (result.player.team == 1) blueScore += result.score.score;
                else if (result.player.team == 2) redScore += result.score.score;
            });

            var teamBlue = document.createElement("span");
            teamBlue.style.color = "#0000ff";
            teamBlue.innerHTML = blueScore.toLocaleString();

            var teamRed = document.createElement("span");
            teamRed.style.color = "#ff0000";
            teamRed.innerHTML = redScore.toLocaleString();

            teamResults.appendChild(teamBlue);
            teamResults.appendChild(document.createTextNode(" vs. "));
            teamResults.appendChild(teamRed);

            var scoreDifference = Math.abs(blueScore - redScore);
            var winner = document.createElement("span");
            winner.classList.add("winner");
            winner.style.fontWeight = "bold";

            if (blueScore > redScore) {
                winner.style.color = "#0000ff";
                winner.innerHTML = "Blue";
            } else if (redScore > blueScore) {
                winner.style.color = "#ff0000";
                winner.innerHTML = "Red";
            } else {
                winner.style.color = "#000000";
                winner.innerHTML = "Draw";
            }

            teamResults.appendChild(document.createElement("br"));
            teamResults.appendChild(winner);
            teamResults.appendChild(document.createTextNode(" wins by "));
            teamResults.appendChild(document.createTextNode(scoreDifference.toLocaleString()));
            teamResults.appendChild(document.createTextNode(" points!"));
            return teamResults;

        case 1:
            // Average Accuracy
            var blueAccs = [];
            var redAccs = [];

            results.forEach(function (result) {
                if (result.player.team == 1) blueAccs.push(result.score.accuracy);
                else if (result.player.team == 2) redAccs.push(result.score.accuracy);
            });

            var blueAcc =
                blueAccs.length > 0
                    ? blueAccs.reduce(function (a, b) {
                          return a + b;
                      }, 0) / blueAccs.length
                    : 0;
            var redAcc =
                redAccs.length > 0
                    ? redAccs.reduce(function (a, b) {
                          return a + b;
                      }, 0) / redAccs.length
                    : 0;

            var teamBlue = document.createElement("span");
            teamBlue.style.color = "#0000ff";
            teamBlue.innerHTML = blueAcc.toFixed(2) + "%";

            var teamRed = document.createElement("span");
            teamRed.style.color = "#ff0000";
            teamRed.innerHTML = redAcc.toFixed(2) + "%";

            teamResults.appendChild(teamBlue);
            teamResults.appendChild(document.createTextNode(" vs. "));
            teamResults.appendChild(teamRed);

            var accDifference = Math.abs(blueAcc - redAcc);
            var winner = document.createElement("span");
            winner.classList.add("winner");
            winner.style.fontWeight = "bold";

            if (blueAcc > redAcc) {
                winner.style.color = "#0000ff";
                winner.innerHTML = "Blue";
            } else if (redAcc > blueAcc) {
                winner.style.color = "#ff0000";
                winner.innerHTML = "Red";
            } else {
                winner.style.color = "#000000";
                winner.innerHTML = "Draw";
            }

            teamResults.appendChild(document.createElement("br"));
            teamResults.appendChild(winner);
            teamResults.appendChild(document.createTextNode(" wins by "));
            teamResults.appendChild(document.createTextNode(accDifference.toFixed(2)));
            teamResults.appendChild(document.createTextNode("%!"));
            return teamResults;

        case 2:
            // Combo
            var blueCombo = 0;
            var redCombo = 0;

            results.forEach(function (result) {
                if (result.player.team == 1) blueCombo += result.score.max_combo;
                else if (result.player.team == 2) redCombo += result.score.max_combo;
            });

            var teamBlue = document.createElement("span");
            teamBlue.style.color = "#0000ff";
            teamBlue.innerHTML = blueCombo.toLocaleString();

            var teamRed = document.createElement("span");
            teamRed.style.color = "#ff0000";
            teamRed.innerHTML = redCombo.toLocaleString();

            teamResults.appendChild(teamBlue);
            teamResults.appendChild(document.createTextNode(" vs. "));
            teamResults.appendChild(teamRed);

            var comboDifference = Math.abs(blueCombo - redCombo);
            var winner = document.createElement("span");
            winner.classList.add("winner");
            winner.style.fontWeight = "bold";

            if (blueCombo > redCombo) {
                winner.style.color = "#0000ff";
                winner.innerHTML = "Blue";
            } else if (redCombo > blueCombo) {
                winner.style.color = "#ff0000";
                winner.innerHTML = "Red";
            } else {
                winner.style.color = "#000000";
                winner.innerHTML = "Draw";
            }

            teamResults.appendChild(document.createElement("br"));
            teamResults.appendChild(winner);
            teamResults.appendChild(document.createTextNode(" wins by "));
            teamResults.appendChild(document.createTextNode(comboDifference.toLocaleString()));
            teamResults.appendChild(document.createTextNode(" combo!"));
            return teamResults;
    }
}

function loadMatchEvents(id, after) {
    var statusText = document.getElementById("status-text");
    var container = document.getElementById("match-events");
    var args = "";

    if (after != undefined) {
        args = "?after=" + after;
    }

    performApiRequest(
        "GET",
        "/multiplayer/" + id + "/events" + args,
        null,
        function (xhr) {
            var events = JSON.parse(xhr.responseText);
            statusText.innerHTML = "";

            if (events.length > 0) {
                lastEventTime = events[events.length - 1].time;
            }

            events.forEach(function (event) {
                var eventDate = new Date(event.time);
                var eventElement = document.createElement("div");
                eventElement.classList.add("event");

                var timeElement = document.createElement("span");
                timeElement.classList.add("event-time");
                timeElement.innerHTML = eventDate.getHours() + ":" + eventDate.getMinutes();

                switch (event.type) {
                    case 0:
                        if (!event.data.name) throw new Error("Invalid api response: " + event.data);

                        var userElement = document.createElement("a");
                        userElement.innerText = event.data.name;
                        userElement.href = "/u/" + event.data.user_id;
                        var descriptionElement = document.createElement("span");
                        descriptionElement.classList.add("event-description");
                        descriptionElement.appendChild(userElement);
                        descriptionElement.appendChild(document.createTextNode(" has joined the match."));
                        eventElement.appendChild(timeElement);
                        eventElement.appendChild(descriptionElement);
                        break;

                    case 1:
                        if (!event.data.name) throw new Error("Invalid api response: " + event.data);

                        var userElement = document.createElement("a");
                        userElement.innerText = event.data.name;
                        userElement.href = "/u/" + event.data.user_id;
                        var descriptionElement = document.createElement("span");
                        descriptionElement.classList.add("event-description");
                        descriptionElement.appendChild(userElement);
                        descriptionElement.appendChild(document.createTextNode(" has left the match."));
                        eventElement.appendChild(timeElement);
                        eventElement.appendChild(descriptionElement);
                        break;

                    case 2:
                        if (!event.data.name) throw new Error("Invalid api response: " + event.data);

                        var userElement = document.createElement("a");
                        userElement.innerText = event.data.name;
                        userElement.href = "/u/" + event.data.user_id;
                        var descriptionElement = document.createElement("span");
                        descriptionElement.classList.add("event-description");
                        descriptionElement.appendChild(userElement);
                        descriptionElement.appendChild(document.createTextNode(" was kicked from the match."));
                        eventElement.appendChild(timeElement);
                        eventElement.appendChild(descriptionElement);
                        break;

                    case 3:
                        if (!event.data["new"]) throw new Error("Invalid api response: " + event.data);

                        var userElement = document.createElement("a");
                        userElement.innerText = event.data["new"].name;
                        userElement.href = "/u/" + event.data["new"].id;
                        var descriptionElement = document.createElement("span");
                        descriptionElement.className = "event-description";
                        descriptionElement.appendChild(userElement);
                        descriptionElement.appendChild(document.createTextNode(" has become the host."));
                        eventElement.appendChild(timeElement);
                        eventElement.appendChild(descriptionElement);
                        break;

                    case 4:
                        var descriptionElement = document.createElement("span");
                        descriptionElement.classList.add("event-description");
                        descriptionElement.appendChild(document.createTextNode("The match was disbanded."));
                        eventElement.appendChild(timeElement);
                        eventElement.appendChild(descriptionElement);
                        break;

                    case 5:
                        // Match was started
                        // TODO: Display this?
                        break;

                    case 6:
                        var startTime = new Date(event.data.start_time);
                        var endTime = new Date(event.data.end_time);
                        var duration = endTime - startTime;
                        var durationMinutes = Math.floor(duration / 1000 / 60);
                        var durationSecondsRemainder = Math.floor((duration / 1000) % 60);
                        var durationString = durationMinutes + "m " + durationSecondsRemainder + "s";

                        var teamType = TeamType[event.data.team_mode];
                        var scoringType = ScoringType[event.data.scoring_mode];
                        var mode = Mode[event.data.mode];

                        var matchDetails = document.createElement("div");
                        matchDetails.classList.add("match-details");

                        var durationElement = document.createElement("div");
                        var durationTitle = document.createElement("strong");
                        durationTitle.innerHTML = "Duration: ";
                        durationElement.appendChild(durationTitle);
                        durationElement.appendChild(document.createTextNode(durationString));

                        var gameModeElement = document.createElement("div");
                        var gameModeTitle = document.createElement("strong");
                        gameModeTitle.innerHTML = "Game Mode: ";
                        gameModeElement.appendChild(gameModeTitle);
                        gameModeElement.appendChild(document.createTextNode(mode + " (" + teamType + ")"));

                        var scoringTypeElement = document.createElement("div");
                        var scoringTypeTitle = document.createElement("strong");
                        scoringTypeTitle.innerHTML = "Scoring Type: ";
                        scoringTypeElement.appendChild(scoringTypeTitle);
                        scoringTypeElement.appendChild(document.createTextNode(scoringType));

                        matchDetails.appendChild(durationElement);
                        matchDetails.appendChild(gameModeElement);
                        matchDetails.appendChild(scoringTypeElement);

                        var beatmapDetails = document.createElement("div");
                        var beatmapTitle = document.createElement("strong");
                        beatmapTitle.innerHTML = "Beatmap: ";
                        beatmapDetails.appendChild(beatmapTitle);

                        if (event.data.beatmap_id != 0) {
                            var beatmapLink = document.createElement("a");
                            beatmapLink.href = "/b/" + event.data.beatmap_id;
                            beatmapLink.appendChild(beatmapDetails);
                            beatmapLink.innerText = event.data.beatmap_text;
                            beatmapDetails.appendChild(beatmapLink);
                        } else {
                            beatmapDetails.appendChild(document.createTextNode(event.data.beatmap_text));
                        }

                        var beatmapInfo = document.createElement("div");
                        beatmapInfo.classList.add("beatmap-info");
                        beatmapInfo.appendChild(beatmapDetails);

                        eventElement.appendChild(matchDetails);
                        eventElement.appendChild(beatmapInfo);

                        var matchResults = document.createElement("div");
                        matchResults.classList.add("match-results");

                        var resultsTable = generateResultsTable(event.data.results, event.data.mods);

                        matchResults.appendChild(resultsTable);
                        eventElement.appendChild(matchResults);

                        if (event.data.team_mode == 2 || event.data.team_mode == 3) {
                            var teamResults = getTeamWinner(event.data.results, event.data.scoring_mode);
                            eventElement.appendChild(teamResults);
                        }

                        eventElement.classList.add("game");
                        eventElement.classList.remove("event");
                        break;

                    case 7:
                        var descriptionElement = document.createElement("span");
                        descriptionElement.classList.add("event-description");
                        descriptionElement.appendChild(document.createTextNode("The match was aborted!"));
                        descriptionElement.style.color = "#ff0000";
                        eventElement.appendChild(timeElement);
                        eventElement.appendChild(descriptionElement);
                }

                container.appendChild(eventElement);
            });
        },
        function (xhr) {
            document.querySelectorAll(".event").forEach(function (element) {
                element.remove();
            });
            document.querySelectorAll(".game").forEach(function (element) {
                element.remove();
            });
            statusText.innerHTML = "Failed to load match. Please try again!";
        }
    );
}

function loadMatchEventsLoop() {
    setTimeout(function () {
        var events = document.getElementById("match-events").innerHTML;

        if (events.includes("The match was disbanded.")) return;

        loadMatchEvents(matchId, lastEventTime);
        loadMatchEventsLoop();
    }, refreshRate);
}

// TODO: Add option for displaying chat

$(document).ready(function (event) {
    loadMatchEvents(matchId, undefined);
    loadMatchEventsLoop();
});
