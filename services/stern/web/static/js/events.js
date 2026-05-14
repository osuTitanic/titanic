var EventTypes = {
    RanksGained: 1,
    NumberOne: 2,
    BeatmapLeaderboardRank: 3,
    LostFirstPlace: 4,
    PPRecord: 5,
    TopPlay: 6,
    AchievementUnlocked: 7,
    ScoreSubmitted: 8,
    BeatmapUploaded: 9,
    BeatmapUpdated: 10,
    BeatmapRevived: 11,
    BeatmapFavouriteAdded: 12,
    BeatmapFavouriteRemoved: 13,
    BeatmapRated: 14,
    BeatmapCommented: 15,
    BeatmapDownloaded: 16,
    BeatmapStatusUpdated: 17,
    BeatmapNominated: 18,
    ForumTopicCreated: 19,
    ForumPostCreated: 20,
    ForumSubscribed: 21,
    ForumUnsubscribed: 22,
    ForumBookmarked: 23,
    ForumUnbookmarked: 24,
    OsuCoinsReceived: 25,
    OsuCoinsUsed: 26,
    FriendAdded: 27,
    FriendRemoved: 28,
    ReplayWatched: 29,
    ScreenshotUploaded: 30,
    UserRegistration: 31,
    UserLogin: 32,
    UserChatMessage: 33,
    UserMatchCreated: 34,
    UserMatchJoined: 35,
    UserMatchLeft: 36,
    BeatmapNuked: 37
};

var LoginClient = {
    osu: "osu!",
    irc: "IRC"
};

var EventRenderers = {
    [EventTypes.RanksGained]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            ` has gained ${event.data.ranks_gained} rank${event.data.ranks_gained !== 1 ? "s" : ""},`,
            ` now placed `,
            renderBoldElement(`#${event.data.rank}`),
            ` in ${event.data.mode}`
        ];
    },
    [EventTypes.NumberOne]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            ` has taken the lead as the top-ranked ${event.data.mode} player!`
        ];
    },
    [EventTypes.BeatmapLeaderboardRank]: (event) => {
        var modsText = event.data.mods ? ` with ${event.data.mods}` : "";
        var ppText = event.data.pp ? ` (${event.data.pp}pp)` : "";
        var modeText = ` <${event.data.mode}>`;

        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            ` achieved rank #${event.data.beatmap_rank} on `,
            renderBeatmap(event.data.beatmap, event.data.beatmap_id),
            modsText,
            modeText,
            ppText
        ];
    },
    [EventTypes.LostFirstPlace]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            ` has lost first place on `,
            renderBeatmap(event.data.beatmap, event.data.beatmap_id),
            ` <${event.data.mode}>`
        ];
    },
    [EventTypes.PPRecord]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            ` has set the new pp record on `,
            renderBeatmap(event.data.beatmap, event.data.beatmap_id),
            ` with ${event.data.pp}pp <${event.data.mode}>`
        ];
    },
    [EventTypes.TopPlay]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            ` achieved a new top play on `,
            renderBeatmap(event.data.beatmap, event.data.beatmap_id),
            ` with ${event.data.pp}pp <${event.data.mode}>`
        ];
    },
    [EventTypes.AchievementUnlocked]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " unlocked an achievement: ",
            renderBoldElement(event.data.achievement)
        ];
    },
    [EventTypes.ScoreSubmitted]: (event) => {
        var modsText = event.data.mods ? ` with ${event.data.mods}` : "";
        var ppText = event.data.pp ? ` (${event.data.pp}pp)` : "";
        var modeText = ` <${event.data.mode}>`;

        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            " submitted a score on ",
            renderBeatmap(event.data.beatmap, event.data.beatmap_id),
            modsText,
            modeText,
            ppText
        ];
    },
    [EventTypes.BeatmapUploaded]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            " uploaded a new beatmap: ",
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id)
        ];
    },
    [EventTypes.BeatmapUpdated]: (event) => {
        return [
            renderProfileWithMode(event.data.username, event.user_id, event.mode),
            " updated a beatmap: ",
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id)
        ];
    },
    [EventTypes.BeatmapRevived]: (event) => {
        return [
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id),
            " has been revived from eternal slumber by ",
            renderProfileWithMode(event.data.username, event.user_id, event.mode)
        ];
    },
    [EventTypes.BeatmapFavouriteAdded]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " added a beatmap to their favourites: ",
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id)
        ];
    },
    [EventTypes.BeatmapFavouriteRemoved]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " removed a beatmap from their favourites: ",
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id)
        ];
    },
    [EventTypes.BeatmapRated]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            ` rated a beatmap with ${event.data.rating}/10 `,
            "(",
            renderBeatmap(event.data.beatmap_name, event.data.beatmap_id),
            ")"
        ];
    },
    [EventTypes.BeatmapCommented]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " commented on a beatmap: ",
            renderBeatmap(event.data.beatmap_name, event.data.beatmap_id)
        ];
    },
    [EventTypes.BeatmapStatusUpdated]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " updated the status of a beatmap to ",
            renderBoldElement(BeatmapStatus.toString(event.data.status)),
            ": ",
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id)
        ];
    },
    [EventTypes.BeatmapNominated]: (event) => {
        var profileLink = renderProfile(event.data.username, event.user_id);
        var beatmapsetLink = renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id);

        if (event.data.type == "reset") return [profileLink, " popped the bubble for ", beatmapsetLink];

        return [profileLink, " nominated a beatmap: ", beatmapsetLink];
    },
    [EventTypes.BeatmapNuked]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " nuked a beatmap: ",
            renderBeatmapset(event.data.beatmapset_name, event.data.beatmapset_id)
        ];
    },
    [EventTypes.ForumTopicCreated]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " created a new topic: ",
            renderTopic(event.data.topic_name, event.data.topic_id)
        ];
    },
    [EventTypes.ForumPostCreated]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " posted in a topic: ",
            renderPost(event.data.topic_name, event.data.post_id)
        ];
    },
    [EventTypes.ForumSubscribed]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " subscribed to a topic: ",
            renderTopic(event.data.topic_name, event.data.topic_id)
        ];
    },
    [EventTypes.ForumUnsubscribed]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " unsubscribed from a topic: ",
            renderTopic(event.data.topic_name, event.data.topic_id)
        ];
    },
    [EventTypes.ForumBookmarked]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " bookmarked a topic: ",
            renderTopic(event.data.topic_name, event.data.topic_id)
        ];
    },
    [EventTypes.ForumUnbookmarked]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " removed a topic from their bookmarks: ",
            renderTopic(event.data.topic_name, event.data.topic_id)
        ];
    },
    [EventTypes.OsuCoinsReceived]: (event) => {
        if (event.data.amount === 0) return;

        return [
            renderProfile(event.data.username, event.user_id),
            ` received ${event.data.amount} osu!coin${event.data.amount !== 1 ? "s" : ""},`,
            ` now standing on ${event.data.coins} coins in total`
        ];
    },
    [EventTypes.OsuCoinsUsed]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            ` used ${event.data.amount} osu!coin${event.data.amount !== 1 ? "s" : ""},`,
            ` now standing on ${event.data.coins} coins in total`
        ];
    },
    [EventTypes.FriendAdded]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " is now following ",
            renderProfile(event.data.target_username, event.data.target_id)
        ];
    },
    [EventTypes.FriendRemoved]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " is no longer following ",
            renderProfile(event.data.target_username, event.data.target_id)
        ];
    },
    [EventTypes.ReplayWatched]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " downloaded a replay on ",
            renderBeatmap(event.data.beatmap_name, event.data.beatmap_id),
            " (",
            renderScoreLink("download", event.data.score_id),
            ")"
        ];
    },
    [EventTypes.ScreenshotUploaded]: (event) => {
        return [renderProfile(event.data.username, event.user_id), " uploaded a screenshot"];
    },
    [EventTypes.UserRegistration]: (event) => {
        return [renderProfile(event.data.username, event.user_id), " just registered!"];
    },
    [EventTypes.UserLogin]: (event) => {
        var profileLink = renderProfile(event.data.username, event.user_id);
        var text = " logged in to the website";

        if (event.data.location == "bancho") {
            var clientType = LoginClient[event.data.client] || event.data.client;
            text = ` logged in to Bancho using ${clientType}`;

            if (event.data.version !== undefined) {
                text += ` (${event.data.version})`;
            }
        }

        return [profileLink, text];
    },
    [EventTypes.UserMatchCreated]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " created a new match: ",
            renderMatchLink(event.data.match_name, event.data.match_id)
        ];
    },
    [EventTypes.UserMatchJoined]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " joined the match: ",
            renderMatchLink(event.data.match_name, event.data.match_id)
        ];
    },
    [EventTypes.UserMatchLeft]: (event) => {
        return [
            renderProfile(event.data.username, event.user_id),
            " left the match: ",
            renderMatchLink(event.data.match_name, event.data.match_id)
        ];
    }
};

function webSocketApiResolver() {
    return window.WebSocket || window.MozWebSocket;
}

function supportsWebSockets() {
    try {
        return Boolean(webSocketApiResolver());
    } catch (e) {
        return false;
    }
}

function setupWebSocket() {
    var activityContainer = document.getElementById("activity-feed-container");
    var statusText = document.getElementById("status-text");

    if (!activityContainer || !statusText) {
        console.error("Activity feed container or status text not found!");
        return;
    }

    var webSocket = webSocketApiResolver();
    var webSocketUrl = activityContainer.dataset.wsEndpoint;
    console.info("Connecting to websocket at:", webSocketUrl);
    statusText.textContent = "Connecting to websocket...";

    var socket = new webSocket(webSocketUrl);
    socket.onopen = onWebSocketOpen;
    socket.onmessage = onWebSocketMessage;
    socket.onerror = onWebSocketError;
    socket.onclose = onWebSocketClose;
}

function onWebSocketOpen(event) {
    console.info("Websocket connection established:", event);

    var statusText = document.getElementById("status-text");
    if (statusText) {
        statusText.textContent = "Connected! Waiting for data...";
    }
}

function onWebSocketMessage(event) {
    var activityContainer = document.getElementById("activity-feed-container");
    var statusText = document.getElementById("status-text");

    if (statusText.style.opacity !== "0") {
        statusText.textContent = "Connected!";
        statusText.style.opacity = "0";
    }

    try {
        var data = JSON.parse(event.data);
        console.info("Websocket message received:", data);

        var eventElement = renderEvent(data);
        if (!eventElement) return;

        activityContainer.insertBefore(eventElement, activityContainer.firstChild);
        setTimeout(function () {
            removeEvent(eventElement);
        }, 10000);
    } catch (e) {
        if (statusText) {
            statusText.textContent = "Error processing message.";
            statusText.style.opacity = "1";
        }
        console.error("Failed to parse websocket message:", event.data, e);
        return;
    }
}

function onWebSocketError(event) {
    console.error("Websocket error:", event);

    var statusText = document.getElementById("status-text");
    if (statusText) {
        statusText.textContent = "Something went wrong :(";
        statusText.style.opacity = "1";
    }
}

function onWebSocketClose(event) {
    console.warn("Websocket connection closed:", event);

    var statusText = document.getElementById("status-text");
    if (statusText) {
        statusText.textContent = "Connection closed. Reconnecting...";
        statusText.style.opacity = "1";
    }

    setTimeout(setupWebSocket, 5000);
}

function renderEvent(eventData) {
    var eventType = eventData.type;
    var renderer = EventRenderers[eventType];

    if (!renderer) {
        console.warn("No renderer found for event type:", eventType);
        return;
    }

    var processedEventData = renderer(eventData);
    if (!processedEventData) {
        console.warn("Renderer returned no content for event type:", eventType);
        return;
    }

    // If renderer returns a single element, wrap it in an array
    if (!Array.isArray(processedEventData)) {
        processedEventData = [processedEventData];
    }

    var element = document.createElement("div");
    element.className = "event-item";
    element.style.opacity = "1";

    // Append event data elements
    processedEventData.forEach(function (e) {
        if (typeof e === "string") {
            var textNode = document.createTextNode(e);
            element.appendChild(textNode);
        } else {
            element.appendChild(e);
        }
    });

    return element;
}

function removeEvent(eventElement) {
    if (totalEvents() <= 10) {
        // Wait until enough events are present before removing
        setTimeout(function () {
            removeEvent(eventElement);
        }, 5000);
        return;
    }
    eventElement.style.opacity = "0";
    setTimeout(function () {
        eventElement.remove();
    }, 1000);
}

function totalEvents() {
    var activityContainer = document.getElementById("activity-feed-container");
    if (!activityContainer) return 0;
    return activityContainer.children.length;
}

function renderProfile(username, userId) {
    var profileLink = document.createElement("a");
    profileLink.textContent = username;
    profileLink.href = `/u/${userId}`;
    profileLink.className = "username-link";
    return profileLink;
}

function renderProfileWithMode(username, userId, mode) {
    var profileLink = document.createElement("a");
    profileLink.textContent = username;
    profileLink.href = `/u/${userId}?mode=${mode}`;
    profileLink.className = "username-link";
    return profileLink;
}

function renderBeatmap(beatmapName, beatmapId) {
    var beatmapLink = document.createElement("a");
    beatmapLink.textContent = beatmapName;
    beatmapLink.href = `/b/${beatmapId}`;
    beatmapLink.className = "beatmap-link";
    return beatmapLink;
}

function renderBeatmapset(beatmapsetName, beatmapsetId) {
    var beatmapsetLink = document.createElement("a");
    beatmapsetLink.textContent = beatmapsetName;
    beatmapsetLink.href = `/s/${beatmapsetId}`;
    beatmapsetLink.className = "beatmapset-link";
    return beatmapsetLink;
}

function renderScoreLink(linkText, scoreId) {
    var replayLink = document.createElement("a");
    replayLink.textContent = linkText;
    replayLink.href = `/scores/${scoreId}/download`;
    replayLink.className = "replay-link";
    return replayLink;
}

function renderMatchLink(matchName, matchId) {
    var matchLink = document.createElement("a");
    matchLink.textContent = matchName;
    matchLink.href = `/mp/${matchId}`;
    matchLink.className = "match-link";
    return matchLink;
}

function renderTopic(topicName, topicId) {
    var topicLink = document.createElement("a");
    topicLink.textContent = topicName;
    topicLink.href = `/forum/t/${topicId}`;
    topicLink.className = "topic-link";
    return topicLink;
}

function renderPost(postContent, postId) {
    var postLink = document.createElement("a");
    postLink.textContent = postContent;
    postLink.href = `/forum/p/${postId}`;
    postLink.className = "post-link";
    return postLink;
}

function renderBoldElement(text) {
    var boldElement = document.createElement("strong");
    boldElement.textContent = text;
    return boldElement;
}

function renderItalicElement(text) {
    var italicElement = document.createElement("em");
    italicElement.textContent = text;
    return italicElement;
}

$(document).ready(function () {
    var activityContainer = document.getElementById("activity-feed-container");
    var statusText = document.getElementById("status-text");

    if (!activityContainer) {
        statusText.textContent = "Activity feed container not found.";
        console.error("Activity feed container not found!");
        return;
    }

    if (!supportsWebSockets()) {
        statusText.textContent = "WebSockets are not supported by your browser.";
        console.warn("WebSockets are not supported by this browser.");
        return;
    }

    setupWebSocket();
});
