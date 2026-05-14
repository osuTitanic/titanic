var mainChannelId = 2; // #osu
var socket = null;

var connected = false;
var isLoadingHistory = false;
var loadingHistoryFor = null;

var hasMoreMessages = {};
var messageHistory = {};
var channels = {};
var users = {};

var activeChannel = null;
var activeDM = null;

// Target channel/DM to join after connection
var pendingTarget = null;

// Compiled regexes for message link parsing
var osuLinkRegex = /\[((?:https?:\/\/)[^\s\]]+)\s+(.+?)\]/g;
var urlRegex = /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)/g;

// Disallowed channels to load message histories from
var disallowedChannels = ["#multi_", "#spec_"];

// Message type handlers, used for "msg" events
var messageHandlers = {
    part: handleUserPart,
    quit: handleUserQuit,
    whois: handleWhoIsResponse,
    message: handleChannelMessage
};

function getQueryParameter(name) {
    var urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
}

function initializeSocket(username, password) {
    if (!username || !password) {
        console.error("Username and password are required");
        return;
    }

    socket = io(loungeBackend, { transports: ["polling"] });

    socket.onAny((event, ...args) => {
        console.debug("Incoming event:", event, args);
    });

    socket.onAnyOutgoing((event, ...args) => {
        console.debug("Outgoing event:", event, args);
    });

    socket.on("connect", function () {
        console.log("Connected to IRC backend");
        updateStatusText("Connected to chat server.");
    });

    socket.on("disconnect", function () {
        console.log("Disconnected from IRC backend");
        updateStatusText("Disconnected from chat server. Please refresh the page to reconnect!");
        disableChatInput();
        connected = false;
    });

    socket.on("init", function () {
        channels = {};
        users = {};
    });

    socket.on("configuration", function () {
        socket.emit("network:new", {
            nick: username,
            username: username,
            realname: username,
            leaveMessage: "Leaving...",
            join: "#osu",
            password: password
        });
    });

    socket.on("msg", function (data) {
        if (data.msg.type in messageHandlers) {
            messageHandlers[data.msg.type](data);
        }
    });

    socket.on("network:status", function (data) {
        console.log("Network status:", data);
        if (!data.connected) {
            updateStatusText("Disconnected from chat server. Please refresh the page to reconnect!");
            disableChatInput();
            connected = false;
        } else {
            updateStatusText("Connecting...");
            connected = true;
        }
    });

    socket.on("network", onNetworkConfiguration);
    socket.on("join", onChannelJoin);
    socket.on("part", onChannelPart);
    socket.on("topic", onChannelTopic);
    socket.on("users", onChannelUsers);
    socket.on("names", onChannelNames);
}

function resetConnection() {
    if (socket) {
        socket.disconnect();
        socket = null;
    }
    hasMoreMessages = {};
    messageHistory = {};
    channels = {};
    users = {};
}

function onNetworkConfiguration(data) {
    // Assuming we only have one network
    var network = data.networks[0];

    // Initialize channels
    for (var i = 0; i < network.channels.length; i++) {
        var channel = network.channels[i];

        // Check if "channel" is actually a channel
        // There's also a "lobby" type, which is used for error messages
        if (channel.type !== "channel") {
            continue;
        }

        channels[channel.id] = channel;
        mainChannelId = channel.id;
    }

    connected = true;
    populateChannels();
    populateDMs();

    // Handle pending target if set
    if (pendingTarget) {
        handlePendingTarget();
        return;
    }

    // Auto-switch to main channel only if no target is specified (#osu)
    if (channels[mainChannelId]) {
        switchToChannel(mainChannelId);
    }
}

function onChannelJoin(data) {
    channels[data.chan.id] = data.chan;
}

function onChannelPart(data) {
    var channelId = data.chan;
    delete channels[channelId];
    populateChannels();

    // If the current channel is the one we parted from, clear UI
    if (activeChannel && activeChannel.id === channelId) {
        activeChannel = null;
        clearChatLog();
        disableChatInput();
        updateChatTitle("Select a channel");
    }
}

function onChannelTopic(data) {
    var channel = channels[data.chan];
    if (!channel) {
        console.error("Received topic for unknown channel:", data.chan);
        return;
    }
    channel.topic = data.topic;
    populateChannels();
}

function onChannelUsers(data) {
    var channel = channels[data.chan];
    if (!channel) {
        console.error("Received user listing request for unknown channel:", data.chan);
        return;
    }
    socket.emit("names", { target: channel.id });
}

function onChannelNames(data) {
    var channel = channels[data.id];
    if (!channel) {
        console.error("Received names for unknown channel:", data.id);
        return;
    }
    channel.users = data.users;

    if (channel.id != mainChannelId) {
        // We are not in the main channel
        return;
    }

    if (channel.users.length == users.length) {
        // No change in user count
        return;
    }

    // Request whois info about each user to resolve their user ID
    for (var i = 0; i < data.users.length; i++) {
        // Check if user already exists
        if (getUserByName(data.users[i].nick)) {
            continue;
        }
        sendWhoIs(data.users[i].nick);
    }
}

function handleUserPart(data) {
    var channel = channels[data.chan];
    if (!channel) {
        console.error("Received user part for unknown channel:", data.chan);
        return;
    }

    // Remove user from channel's user list
    for (var nick in channel.users) {
        if (channel.users[nick].nick === data.msg.from.nick) {
            delete channel.users[nick];
        }
    }
}

function handleUserQuit(data) {
    for (var id in users) {
        if (users[id].nick === data.msg.from.nick) {
            delete users[id];
            var dmElement = document.getElementById("dm-" + id);
            if (dmElement) {
                dmElement.dataset.isOnline = "false";
            }
        }
    }
}

function handlePendingTarget() {
    if (!pendingTarget) {
        return;
    }

    var target = pendingTarget;

    // Clear target to avoid repeated attempts
    pendingTarget = null;

    if (!isNaN(parseInt(target))) {
        // We got a user ID
        handleStartDMById(parseInt(target));
        return;
    }

    if (!target.startsWith("#")) {
        // We got a username
        handleStartDMByName(target);
        return;
    }

    // We got a channel
    var channel = getChannelByName(target);
    if (channel) {
        // Already in this channel, just switch to it
        switchToChannel(channel.id);
        return;
    }

    // Join the channel
    joinChannel(target);

    // Wait a moment for the join to complete, then try to switch
    setTimeout(function () {
        var joinedChannel = getChannelByName(target);
        if (joinedChannel) {
            switchToChannel(joinedChannel.id);
        } else {
            console.error("Failed to join channel:", target);
            switchToChannel(mainChannelId);
        }
    }, 1000);
}

function handleWhoIsResponse(data) {
    var whois = data.msg.whois;
    if (!whois || !whois.nick) {
        console.error("Invalid whois data:", data);
        return;
    }

    // Parse user ID from real name
    // "http://osu.titanic.sh/u/12345"
    var identParts = whois.real_name.split("/");
    var userId = parseInt(identParts[identParts.length - 1]);
    users[userId] = whois;
    users[userId].id = userId;
    users[userId].status = null;

    // Update any DM entries for this user
    var dmElement = document.getElementById("dm-" + userId);
    if (dmElement) {
        dmElement.dataset.isOnline = "true";
    }
}

function handleChannelMessage(data) {
    if (data.chan == 1) {
        // System message, log to console
        console.log("[System]:", data.msg.text);
        return;
    }

    var channel = channels[data.chan];
    if (!channel) {
        console.error("Received message for unknown channel:", data.chan);
        return;
    }

    var sender = getUserByName(data.msg.from.nick);
    if (!sender) {
        sender = { nick: data.msg.from.nick };
    }

    var message = data.msg.text;
    var highlight = data.msg.highlight;
    var time = data.msg.time || new Date();

    // Store message in appropriate cache
    var historyKey = null;

    if (channel.type === "channel") {
        historyKey = getChannelHistoryKey(channel);
    } else if (channel.type === "query") {
        historyKey = getDMHistoryKey(sender.id);
        highlight = false;
    }

    if (historyKey) {
        if (!messageHistory[historyKey]) {
            messageHistory[historyKey] = [];
        }
        messageHistory[historyKey].push({
            sender: sender,
            text: message,
            highlight: highlight,
            time: time
        });
    }

    if (channel.type === "channel") {
        console.log("[" + channel.name + "] " + sender.nick + ":", message);

        // Display the message if we're in this channel
        if (activeChannel && activeChannel.id === data.chan) {
            displayMessage(sender, message, highlight, time);
            return;
        }

        markChannelAsUnread(data.chan);
    }

    if (channel.type === "query") {
        console.log("[DM]", sender.nick + ":", message);

        // Display the message if we're in this DM
        if (activeDM && sender.id === activeDM) {
            displayMessage(sender, message, false, time);
            return;
        }

        markDmAsUnread(sender.id);
    }
}

function sendWhoIs(username) {
    sendInput(2, "/whois " + username);
}

function sendWhoIsMany(usernames) {
    if (!usernames || usernames.length === 0) {
        console.error("Usernames array is empty");
        return;
    }
    sendInput(2, "/whois " + usernames.join(" "));
}

function sendChannelMessage(channel, message) {
    if (!channel || !message) {
        console.error("Channel and message are required to send a channel message");
        return;
    }
    sendInput(channel, message);
}

function sendDirectMessage(username, message) {
    if (!username || !message) {
        console.error("Username and message are required to send a direct message");
        return;
    }
    var channel = getChannelByName(username);
    sendInput(channel.id, message);
}

function sendInput(channel, message) {
    if (!socket) {
        console.error("Socket not initialized");
        return;
    }
    if (channel == undefined || !message) {
        console.error("Invalid arguments for sendInput");
        return;
    }
    socket.emit("input", { target: channel, text: message });
}

function joinChannel(channelName) {
    sendInput(1, "/join " + channelName);
}

function leaveChannel(channelName) {
    var channel = getChannelByName(channelName);
    if (!channel) return;
    sendInput(channel.id, "/close");
    onChannelPart({ chan: channel.id });
}

function getChannelByName(name) {
    for (var id in channels) {
        if (channels[id].name === name) {
            return channels[id];
        }
    }
    return null;
}

function getChannelById(channelId) {
    return channels[channelId] || null;
}

function getUserByName(username) {
    for (var id in users) {
        if (users[id].nick === username) {
            return users[id];
        }
    }
    return null;
}

function getUserById(userId) {
    return users[userId] || null;
}

function getChannelHistoryKey(channel) {
    return "channel_" + channel.id;
}

function getDMHistoryKey(userId) {
    return "dm_" + userId;
}

function getHistoryMessages(key) {
    return messageHistory[key] || [];
}

function hasHistoryCache(key) {
    return messageHistory[key] && messageHistory[key].length > 0;
}

function storeHistoryMessages(key, messages, append) {
    if (!messageHistory[key]) {
        messageHistory[key] = [];
    }

    if (append) {
        // Prepend older messages at the beginning
        messageHistory[key] = messages.concat(messageHistory[key]);
    } else {
        messageHistory[key] = messages;
    }
}

function fetchChannelMessageHistory(channel, offset, limit, onSuccess, onFailure) {
    var url = "/chat/channels/" + encodeURIComponent(channel) + "/messages";

    if (offset) {
        url += "?offset=" + encodeURIComponent(offset);
    }
    if (limit) {
        url += (offset ? "&" : "?") + "limit=" + encodeURIComponent(limit);
    }

    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to fetch channel message history:", xhr);
            onFailure(xhr);
        }
    );
}

function fetchDirectMessageHistory(userId, offset, limit, onSuccess, onFailure) {
    var url = "/chat/dms/" + userId + "/messages";

    if (offset) {
        url += "?offset=" + encodeURIComponent(offset);
    }
    if (limit) {
        url += (offset ? "&" : "?") + "limit=" + encodeURIComponent(limit);
    }

    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to fetch DM message history:", xhr);
            onFailure(xhr);
        }
    );
}

function fetchUserById(userId, onSuccess, onFailure) {
    var url = "/users/" + userId;
    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to fetch user by ID:", xhr);
            onFailure(xhr);
        }
    );
}

function fetchUserByName(username, onSuccess, onFailure) {
    var url = "/users/lookup/" + encodeURIComponent(username);
    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to fetch user by name:", xhr);
            onFailure(xhr);
        }
    );
}

function fetchUserStatus(userId, onSuccess, onFailure) {
    var url = "/users/" + userId + "/status";
    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to fetch user status:", xhr);
            onFailure(xhr);
        }
    );
}

function fetchDirectMessageSelection(onSuccess, onFailure) {
    var url = "/chat/dms";

    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to fetch DM selection:", xhr);
            onFailure(xhr);
        }
    );
}

function postDirectMessage(userId, message, onSuccess, onFailure) {
    var url = "/chat/dms/" + userId + "/messages";
    var data = { message: message };

    return performApiRequest(
        "POST",
        url,
        data,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to post direct message:", xhr);
            onFailure(xhr);
        }
    );
}

function markDmAsRead(targetId, messageId, onSuccess, onFailure) {
    var url = "/chat/dms/" + targetId + "/messages/" + messageId + "/read";
    return performApiRequest(
        "POST",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to mark DM as read:", xhr);
            onFailure(xhr);
        }
    );
}

function markAllDmsAsRead(targetId, onSuccess, onFailure) {
    var url = "/chat/dms/" + targetId + "/messages/read";
    return performApiRequest(
        "POST",
        url,
        null,
        function (xhr) {
            var response = JSON.parse(xhr.responseText);
            if (!response) {
                console.error("Invalid response format:", response);
                onFailure(xhr);
                return;
            }
            if (onSuccess) onSuccess(response);
        },
        function (xhr) {
            console.error("Failed to mark all DMs as read:", xhr);
            onFailure(xhr);
        }
    );
}

function markChannelAsUnread(channelId) {
    var channelContainer = document.getElementById("channel-container");
    if (!channelContainer) {
        return;
    }

    var channelElement = channelContainer.querySelector('[data-channel-id="' + channelId + '"]');
    if (channelElement && !channelElement.classList.contains("active")) {
        channelElement.classList.add("unread");
        channelElement.style.fontWeight = "bold";
    }
}

function markDmAsUnread(userId) {
    var dmContainer = document.getElementById("dm-container");
    if (!dmContainer) {
        return;
    }

    var dmElement = dmContainer.querySelector('[data-user-id="' + userId + '"]');
    if (dmElement && !dmElement.classList.contains("active")) {
        dmElement.classList.add("unread");
        dmElement.style.fontWeight = "bold";
    }
}

function clearChannelUnread(channelId) {
    var channelContainer = document.getElementById("channel-container");
    if (!channelContainer) {
        return;
    }

    var channelElement = channelContainer.querySelector('[data-channel-id="' + channelId + '"]');
    if (channelElement) {
        channelElement.classList.remove("unread");
        channelElement.style.fontWeight = "normal";
    }
}

function clearDMUnread(userId) {
    var dmContainer = document.getElementById("dm-container");
    if (!dmContainer) {
        return;
    }

    var dmElement = dmContainer.querySelector('[data-user-id="' + userId + '"]');
    if (dmElement) {
        dmElement.classList.remove("unread");
        dmElement.style.fontWeight = "normal";
    }
}

function populateChannels() {
    var channelContainer = document.getElementById("channel-container");

    // Nuke everything first
    while (channelContainer.firstChild) {
        channelContainer.removeChild(channelContainer.firstChild);
    }

    for (var id in channels) {
        var channel = channels[id];
        if (channel.type !== "channel") continue;

        var channelElement = document.createElement("div");
        channelElement.className = "channel-entry";
        channelElement.dataset.channelId = id;
        channelElement.textContent = channel.name;
        channelElement.onclick = function () {
            switchToChannel(this.dataset.channelId);
        };

        var closeButton = document.createElement("span");
        closeButton.className = "close-channel";
        closeButton.textContent = "x";
        closeButton.onclick = function (e) {
            e.stopPropagation();
            var channelId = this.parentElement.dataset.channelId;
            var channel = getChannelById(channelId);
            if (channel) leaveChannel(channel.name);
        };

        if (channel.name !== "#osu") {
            // Only add close button for non-forced channels
            channelElement.appendChild(closeButton);
        }

        // Mark active channel
        if (activeChannel && activeChannel.id === channel.id) {
            channelElement.classList.add("active");
        }

        channelContainer.appendChild(channelElement);
    }
}

function populateDMs() {
    var dmContainer = document.getElementById("dm-container");

    fetchDirectMessageSelection(
        function (dms) {
            if (!dms || dms.length === 0) {
                // No DMs available
                return;
            }

            // Nuke everything first
            while (dmContainer.firstChild) {
                dmContainer.removeChild(dmContainer.firstChild);
            }

            for (var i = 0; i < dms.length; i++) {
                var dm = dms[i];
                var dmElement = document.createElement("div");
                dmElement.className = "dm-entry";
                dmElement.id = "dm-" + dm.user.id;
                dmElement.textContent = dm.user.name;
                dmElement.dataset.name = dm.user.name;
                dmElement.dataset.userId = dm.user.id;
                dmElement.dataset.isOnline = getUserById(dm.user.id) ? "true" : "false";
                dmElement.addEventListener("click", function () {
                    switchToDM(parseInt(this.dataset.userId));
                });

                // Mark active DM
                if (activeDM && activeDM === dm.user.id) {
                    dmElement.classList.add("active");
                }

                dmContainer.appendChild(dmElement);

                if (!dm.last_message) {
                    continue;
                }

                // Mark as unread depending on last message's state
                if (dm.last_message.read == false && dm.last_message.sender_id == dm.user.id) {
                    console.log("Unread conversation with", dm.user.name);
                    markDmAsUnread(dm.user.id);
                }
            }
        },
        function (xhr) {
            updateStatusText("Failed to load DMs.");
            disableChatInput();
        }
    );
}

function updateActiveChannel() {
    var channelContainer = document.getElementById("channel-container");

    for (var i = 0; i < channelContainer.children.length; i++) {
        var channelElement = channelContainer.children[i];
        channelElement.classList.remove("active");

        if (activeChannel && parseInt(channelElement.dataset.channelId) === activeChannel.id) {
            channelElement.classList.add("active");
        }
    }
}

function updateActiveDM() {
    var dmContainer = document.getElementById("dm-container");

    for (var i = 0; i < dmContainer.children.length; i++) {
        var dmElement = dmContainer.children[i];
        dmElement.classList.remove("active");

        if (activeDM && parseInt(dmElement.dataset.userId) === activeDM) {
            dmElement.classList.add("active");
        }
    }
}

function switchToChannel(channelId) {
    var channel = channels[channelId];
    if (!channel) {
        console.error("Channel not found:", channelId);
        return;
    }

    activeChannel = channel;
    activeDM = null;

    // Clear unread marker
    clearChannelUnread(channelId);

    // Update UI to show active channel
    updateActiveChannel();
    updateActiveDM();
    updateChatTitle(channel.name);

    // Clear and reload chat log
    clearChatLog();
    loadChannelHistory(channel);
    enableChatInput();
}

function switchToDM(userId) {
    activeDM = userId;
    activeChannel = null;

    // Clear unread marker
    clearDMUnread(userId);

    // Update UI to show active DM
    updateActiveChannel();
    updateActiveDM();

    fetchUserById(
        userId,
        function (user) {
            updateChatTitle("Direct Message with " + user.name);

            // Clear and reload chat log
            clearChatLog();
            loadDMHistory(user);
            enableChatInput();
            markAllDmsAsRead(userId);
        },
        function (xhr) {
            console.error("Failed to fetch user for DM:", xhr);
            updateStatusText("Failed to load DM");
        }
    );
}

function loadChannelHistory(channel) {
    if (disallowedChannels.some((prefix) => channel.name.startsWith(prefix))) return;

    var historyKey = getChannelHistoryKey(channel);

    // Check if we have cached messages
    if (hasHistoryCache(historyKey)) {
        console.debug("Loading messages from cache for", channel.name);
        var cachedMessages = getHistoryMessages(historyKey);

        // Display cached messages
        for (var i = 0; i < cachedMessages.length; i++) {
            var msg = cachedMessages[i];
            displayMessage(msg.sender, msg.text, msg.highlight, msg.time);
        }

        updateStatusText("Type a message...");
        scrollChatToBottom();
        return;
    }

    // No cache, fetch from API
    updateStatusText("Loading messages...");
    isLoadingHistory = true;
    loadingHistoryFor = historyKey;

    fetchChannelMessageHistory(
        channel.name,
        0,
        50,
        function (messages) {
            isLoadingHistory = false;

            // Check if we're still viewing this channel
            if (loadingHistoryFor !== historyKey) {
                console.debug("Channel changed during load, ignoring results");
                return;
            }
            loadingHistoryFor = null;

            if (!messages || messages.length === 0) {
                updateStatusText("No messages in this channel yet.");
                hasMoreMessages[historyKey] = false;
                return;
            }

            // Store messages in cache
            var historicalMessages = [];
            for (var i = messages.length - 1; i >= 0; i--) {
                var msg = messages[i];
                historicalMessages.push({
                    sender: { nick: msg.sender.name, id: msg.sender.id },
                    text: msg.message,
                    highlight: false,
                    time: new Date(msg.time)
                });
                displayHistoricalMessage(msg);
            }

            storeHistoryMessages(historyKey, historicalMessages, false);
            hasMoreMessages[historyKey] = messages.length >= 50;

            updateStatusText("Type a message...");
            scrollChatToBottom();
        },
        function (xhr) {
            isLoadingHistory = false;
            loadingHistoryFor = null;
            console.error("Failed to load channel history:", xhr);
            updateStatusText("Failed to load message history.");
        }
    );
}

function loadDMHistory(user) {
    var historyKey = getDMHistoryKey(user.id);

    // Check if we have cached messages
    if (hasHistoryCache(historyKey)) {
        console.debug("Loading DM messages from cache for", user.name);
        var cachedMessages = getHistoryMessages(historyKey);

        // Display cached messages
        for (var i = 0; i < cachedMessages.length; i++) {
            var msg = cachedMessages[i];
            displayMessage(msg.sender, msg.text, msg.highlight, msg.time);
        }

        updateStatusText("Type a message...");
        scrollChatToBottom();
        return;
    }

    // No cache, fetch from API
    updateStatusText("Loading messages...");
    isLoadingHistory = true;
    loadingHistoryFor = historyKey;

    fetchDirectMessageHistory(
        user.id,
        0,
        50,
        function (messages) {
            isLoadingHistory = false;

            // Check if we're still viewing this DM
            if (loadingHistoryFor !== historyKey) {
                console.debug("DM changed during load, ignoring results");
                return;
            }
            loadingHistoryFor = null;

            if (!messages || messages.length === 0) {
                updateStatusText("No messages yet. Start a conversation!");
                hasMoreMessages[historyKey] = false;
                return;
            }

            // Store messages in cache
            var historicalMessages = [];
            for (var i = messages.length - 1; i >= 0; i--) {
                var msg = messages[i];
                var nickname = currentUsername;
                if (msg.sender_id === user.id) {
                    nickname = user.name;
                }

                historicalMessages.push({
                    sender: { nick: nickname, id: msg.sender_id },
                    text: msg.message,
                    highlight: false,
                    time: new Date(msg.time)
                });
                displayHistoricalDirectMessage(msg, user);
            }

            storeHistoryMessages(historyKey, historicalMessages, false);
            hasMoreMessages[historyKey] = messages.length >= 50;

            updateStatusText("Type a message...");
            scrollChatToBottom();
        },
        function (xhr) {
            isLoadingHistory = false;
            loadingHistoryFor = null;
            console.error("Failed to load DM history:", xhr);
            updateStatusText("Failed to load message history.");
        }
    );
}

function createMessageElement(sender, text, highlight, time) {
    var messageElement = document.createElement("div");
    messageElement.className = "chat-message-entry";

    if (highlight) {
        messageElement.classList.add("highlighted");
    }

    var timestamp = time ? formatMessageTime(time) : formatMessageTime(new Date());
    var senderName = sender.nick || sender.name || "Unknown";
    var userId = sender.id;

    // Check if this is a /me command
    var isAction = false;
    if (text.startsWith("\u0001ACTION ") && text.endsWith("\u0001")) {
        // Remove \u0001ACTION and trailing \u0001
        text = text.substring(8, text.length - 1);
        isAction = true;
    }

    var timestampSpan = document.createElement("span");
    timestampSpan.className = "message-time";
    timestampSpan.textContent = timestamp;

    var senderLink = document.createElement("a");
    senderLink.className = "message-sender";
    senderLink.style.color = getUserColor(senderName);
    senderLink.style.fontWeight = "bold";
    senderLink.href = "#";

    if (userId) {
        senderLink.style.cursor = "pointer";
        senderLink.dataset.userId = userId;
        senderLink.href = "/u/" + userId;
    }

    var textSpan = document.createElement("span");
    textSpan.className = "message-text";

    // Parse and render message with links
    var parsedParts = parseMessageLinks(text);

    for (var i = 0; i < parsedParts.length; i++) {
        var part = parsedParts[i];
        if (part.type === "link") {
            var linkElement = document.createElement("a");
            linkElement.href = part.url;
            linkElement.textContent = part.text;
            linkElement.target = "_blank";
            linkElement.rel = "noopener noreferrer";
            textSpan.appendChild(linkElement);
        } else {
            textSpan.appendChild(document.createTextNode(part.content));
        }
    }

    messageElement.appendChild(timestampSpan);
    messageElement.appendChild(document.createTextNode(" "));

    if (isAction) {
        // For ACTION messages, display as: <username> text
        messageElement.appendChild(document.createTextNode("*"));
        senderLink.textContent = senderName;
        messageElement.appendChild(senderLink);
        messageElement.appendChild(document.createTextNode(" "));
    } else {
        // For normal messages, display as: username: text
        senderLink.textContent = senderName;
        messageElement.appendChild(senderLink);
        messageElement.appendChild(document.createTextNode(": "));
    }

    messageElement.appendChild(textSpan);
    return messageElement;
}

function displayMessage(sender, text, highlight, time) {
    var chatLog = document.querySelector(".chat-log");
    if (!chatLog) {
        console.error("Chat log element not found");
        return;
    }

    // Split message by newlines and display
    // each line as a separate message
    var lines = text.split("\n");
    for (var i = 0; i < lines.length; i++) {
        var messageElement = createMessageElement(sender, lines[i], highlight, time);
        chatLog.appendChild(messageElement);
    }

    scrollChatToBottom();
}

function getUserColor(username) {
    return `hsl(${hashCode(username) % 360}, 60%, 60%)`;
}

function hashCode(str) {
    var hash = 0;
    for (var i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    return hash;
}

function parseMessageLinks(text) {
    // First parse osu! links [<url> <text>]
    var parts = parseLinksWithRegex(text, osuLinkRegex, function (match) {
        return {
            type: "link",
            url: match[1],
            text: match[2]
        };
    });

    // Then parse regular URLs in remaining text parts
    var finalParts = [];

    for (var i = 0; i < parts.length; i++) {
        if (parts[i].type === "text") {
            var urlParts = parseLinksWithRegex(parts[i].content, urlRegex, function (match) {
                return {
                    type: "link",
                    url: match[0],
                    text: match[0]
                };
            });
            finalParts = finalParts.concat(urlParts);
        } else {
            finalParts.push(parts[i]);
        }
    }

    return finalParts;
}

function parseLinksWithRegex(text, regex, linkFactory) {
    var parts = [];
    var lastIndex = 0;
    var match;

    regex.lastIndex = 0;

    while ((match = regex.exec(text)) !== null) {
        // Add text before the match
        if (match.index > lastIndex) {
            parts.push({
                type: "text",
                content: text.substring(lastIndex, match.index)
            });
        }

        // Add the link using the factory function
        parts.push(linkFactory(match));
        lastIndex = regex.lastIndex;
    }

    // Add remaining text
    if (lastIndex < text.length) {
        parts.push({
            type: "text",
            content: text.substring(lastIndex)
        });
    }

    // If no matches, return original text
    if (parts.length === 0) {
        parts.push({
            type: "text",
            content: text
        });
    }

    return parts;
}

function displayHistoricalMessage(msg) {
    displayMessage({ nick: msg.sender.name, id: msg.sender.id }, msg.message, false, new Date(msg.time));
}

function displayHistoricalDirectMessage(msg, user) {
    var nickname = currentUsername;
    if (msg.sender_id === user.id) {
        nickname = user.name;
    }

    displayMessage({ nick: nickname, id: msg.sender_id }, msg.message, false, new Date(msg.time));
}

function formatMessageTime(time) {
    var date = time instanceof Date ? time : new Date(time);
    var hours = date.getHours().toString().padStart(2, "0");
    var minutes = date.getMinutes().toString().padStart(2, "0");
    return hours + ":" + minutes;
}

function clearChatLog() {
    var chatLog = document.querySelector(".chat-log");
    if (!chatLog) {
        return;
    }

    while (chatLog.firstChild) {
        chatLog.removeChild(chatLog.firstChild);
    }
}

function scrollChatToBottom() {
    var chatLog = document.querySelector(".chat-log");
    if (chatLog) {
        chatLog.scrollTop = chatLog.scrollHeight;
    }
}

function loadMoreChannelMessages(channel) {
    if (isLoadingHistory) {
        console.debug("Already loading history, ignoring request");
        return;
    }

    var historyKey = getChannelHistoryKey(channel);

    if (hasMoreMessages[historyKey] === false) {
        console.debug("No more messages to load for", channel.name);
        return;
    }

    var currentMessages = getHistoryMessages(historyKey);
    var offset = currentMessages.length;

    isLoadingHistory = true;
    updateStatusText("Loading more messages...");

    fetchChannelMessageHistory(
        channel.name,
        offset,
        50,
        function (messages) {
            isLoadingHistory = false;

            // Verify we're still on the same channel
            if (!activeChannel || activeChannel.id !== channel.id) {
                console.debug("Channel changed during loadMore, ignoring results");
                return;
            }

            if (!messages || messages.length === 0) {
                hasMoreMessages[historyKey] = false;
                updateStatusText("Type a message...");
                return;
            }

            var chatLog = document.querySelector(".chat-log");
            var scrollHeightBefore = chatLog.scrollHeight;
            var scrollTopBefore = chatLog.scrollTop;

            // Store and display older messages
            var olderMessages = [];

            for (var i = messages.length - 1; i >= 0; i--) {
                var msg = messages[i];
                olderMessages.push({
                    sender: { nick: msg.sender.name, id: msg.sender.id },
                    text: msg.message,
                    highlight: false,
                    time: new Date(msg.time)
                });

                // Insert at the beginning of the chat log
                var messageElement = createMessageElement(
                    { nick: msg.sender.name, id: msg.sender.id },
                    msg.message,
                    false,
                    new Date(msg.time)
                );
                if (chatLog.firstChild) {
                    chatLog.insertBefore(messageElement, chatLog.firstChild);
                } else {
                    chatLog.appendChild(messageElement);
                }
            }

            storeHistoryMessages(historyKey, olderMessages, true);
            hasMoreMessages[historyKey] = messages.length >= 50;

            // Maintain scroll position
            chatLog.scrollTop = scrollTopBefore + (chatLog.scrollHeight - scrollHeightBefore);

            updateStatusText("Type a message...");
        },
        function (xhr) {
            isLoadingHistory = false;
            console.error("Failed to load more channel messages:", xhr);
            updateStatusText("Failed to load more messages.");
        }
    );
}

function loadMoreDMMessages(userId) {
    if (isLoadingHistory) {
        console.debug("Already loading history, ignoring request");
        return;
    }

    var historyKey = getDMHistoryKey(userId);
    if (hasMoreMessages[historyKey] === false) {
        console.debug("No more messages to load for DM");
        return;
    }

    var currentMessages = getHistoryMessages(historyKey);
    var offset = currentMessages.length;

    isLoadingHistory = true;
    updateStatusText("Loading more messages...");

    fetchUserById(
        userId,
        function (user) {
            fetchDirectMessageHistory(
                userId,
                offset,
                50,
                function (messages) {
                    isLoadingHistory = false;

                    // Verify we're still on the same DM
                    if (activeDM !== userId) {
                        console.debug("DM changed during loadMore, ignoring results");
                        return;
                    }

                    if (!messages || messages.length === 0) {
                        hasMoreMessages[historyKey] = false;
                        updateStatusText("Type a message...");
                        return;
                    }

                    var chatLog = document.querySelector(".chat-log");
                    var scrollHeightBefore = chatLog.scrollHeight;
                    var scrollTopBefore = chatLog.scrollTop;

                    // Store and display older messages
                    var olderMessages = [];
                    for (var i = messages.length - 1; i >= 0; i--) {
                        var msg = messages[i];
                        var nickname = currentUsername;
                        if (msg.sender_id === user.id) {
                            nickname = user.name;
                        }

                        olderMessages.push({
                            sender: { nick: nickname, id: msg.sender_id },
                            text: msg.message,
                            highlight: false,
                            time: new Date(msg.time)
                        });

                        // Insert at the beginning of the chat log
                        var messageElement = createMessageElement(
                            { nick: nickname, id: msg.sender_id },
                            msg.message,
                            false,
                            new Date(msg.time)
                        );
                        if (chatLog.firstChild) {
                            chatLog.insertBefore(messageElement, chatLog.firstChild);
                        } else {
                            chatLog.appendChild(messageElement);
                        }
                    }

                    storeHistoryMessages(historyKey, olderMessages, true);
                    hasMoreMessages[historyKey] = messages.length >= 50;

                    // Maintain scroll position
                    chatLog.scrollTop = scrollTopBefore + (chatLog.scrollHeight - scrollHeightBefore);

                    updateStatusText("Type a message...");
                },
                function (xhr) {
                    isLoadingHistory = false;
                    console.error("Failed to load more DM messages:", xhr);
                    updateStatusText("Failed to load more messages.");
                }
            );
        },
        function (xhr) {
            isLoadingHistory = false;
            console.error("Failed to fetch user for DM:", xhr);
        }
    );
}

function scrollChatToBottom() {
    var chatLog = document.querySelector(".chat-log");
    if (chatLog) {
        chatLog.scrollTop = chatLog.scrollHeight;
    }
}

function updateChatTitle(title) {
    var statusText = document.querySelector(".chat-input .status-text");
    if (statusText) {
        statusText.textContent = title;
    }
}

function updateStatusText(text) {
    var statusText = document.querySelector(".chat-input .chat-message");
    if (statusText) {
        statusText.placeholder = text;
    }
}

function enableChatInput() {
    if (!connected) return;

    var inputField = document.querySelector(".chat-input .chat-message");
    var sendButton = document.querySelector(".chat-input .chat-send");

    if (inputField) {
        inputField.disabled = false;
        inputField.placeholder = "Type a message...";
    }

    if (sendButton) {
        sendButton.disabled = false;
    }
}

function disableChatInput() {
    var inputField = document.querySelector(".chat-input .chat-message");
    var sendButton = document.querySelector(".chat-input .chat-send");

    if (inputField) {
        inputField.disabled = true;
    }

    if (sendButton) {
        sendButton.disabled = true;
    }
}

function sendCurrentMessage() {
    var inputField = document.querySelector(".chat-input .chat-message");
    if (!inputField) {
        return;
    }

    var message = inputField.value.trim();
    if (!message) {
        return;
    }

    if (activeChannel) {
        // Don't add to cache here - let the server echo
        // handle it to avoid duplicate messages in cache
        sendChannelMessage(activeChannel.id, message);
        inputField.value = "";
        return;
    }

    if (activeDM) {
        var time = new Date();

        var historyKey = getDMHistoryKey(activeDM);
        if (!messageHistory[historyKey]) {
            messageHistory[historyKey] = [];
        }
        messageHistory[historyKey].push({
            sender: { nick: currentUsername, id: currentUser },
            text: message,
            highlight: false,
            time: time
        });

        displayMessage({ nick: currentUsername, id: currentUser }, message, false, time);
        inputField.value = "";

        var userObject = getUserById(activeDM);
        if (!userObject) {
            // User is most likely offline
            postDirectMessage(
                activeDM,
                message,
                function (response) {
                    console.log("DM sent successfully:", response);
                },
                function (xhr) {
                    console.error("Failed to send DM:", xhr);
                    updateStatusText("Failed to send message");
                }
            );
            return;
        }

        var userChannel = getChannelByName(userObject.nick);
        if (userChannel) {
            sendChannelMessage(userChannel.id, message);
            return;
        }
        return;
    }

    console.error("No active channel or DM selected");
}

function initializeChatHandlers() {
    var chatInput = document.querySelector(".chat-input .chat-message");
    var sendButton = document.querySelector(".chat-input .chat-send");
    var joinChannelBtn = document.getElementById("channel-join-btn");
    var joinDMApiBtn = document.getElementById("dm-join-btn");
    var chatLog = document.querySelector(".chat-log");

    if (sendButton) {
        sendButton.addEventListener("click", sendCurrentMessage);
    }

    if (chatInput) {
        chatInput.addEventListener("keypress", function (event) {
            if (event.key === "Enter") {
                sendCurrentMessage();
            }
        });
    }

    if (joinChannelBtn) {
        joinChannelBtn.addEventListener("click", handleJoinChannel);
    }

    if (joinDMApiBtn) {
        joinDMApiBtn.addEventListener("click", handleStartDMFromInput);
    }

    // Add listener for enter key on channel join input
    var channelJoinInput = document.getElementById("channel-join-input");
    if (channelJoinInput) {
        channelJoinInput.addEventListener("keyup", function (event) {
            if (event.key === "Enter") {
                handleJoinChannel();
            }
        });
    }

    // Add listener for enter key on dm join input
    var dmJoinInput = document.getElementById("dm-join-input");
    if (dmJoinInput) {
        dmJoinInput.addEventListener("keyup", function (event) {
            if (event.key === "Enter") {
                handleStartDMFromInput();
            }
        });
    }

    // Add scroll listener for infinite scroll
    if (chatLog) {
        chatLog.addEventListener("scroll", function () {
            // Check if scrolled to near the top (within 50px)
            if (chatLog.scrollTop < 50 && !isLoadingHistory) {
                if (activeChannel) {
                    loadMoreChannelMessages(activeChannel);
                } else if (activeDM) {
                    loadMoreDMMessages(activeDM);
                }
            }
        });
    }
}

function handleJoinChannel() {
    var channelInput = document.getElementById("channel-join-input");
    if (!channelInput) {
        console.error("Channel input field not found");
        return;
    }

    var channelName = channelInput.value.trim();
    if (!channelName) {
        return;
    }

    if (!channelName.startsWith("#")) {
        channelName = "#" + channelName;
    }

    // Check if we are already in this channel
    for (var id in channels) {
        if (channels[id].name === channelName) {
            switchToChannel(id);
            channelInput.value = "";
            return;
        }
    }

    joinChannel(channelName);
    channelInput.value = "";
}

function handleStartDMByName(username) {
    if (!username) {
        console.error("Username is required to start a DM");
        return;
    }

    fetchUserByName(
        username,
        function (user) {
            // Add user to DM list if not already present
            var dmContainer = document.getElementById("dm-container");
            var existing = dmContainer.querySelector(`.dm-entry[data-user-id="${user.id}"]`);
            if (!existing) {
                var dmEntry = document.createElement("div");
                dmEntry.className = "dm-entry";
                dmEntry.setAttribute("data-user-id", user.id);
                dmEntry.textContent = user.name;
                dmEntry.onclick = function () {
                    switchToDM(user.id);
                };
                dmContainer.appendChild(dmEntry);
            }

            // User found, switch to DM
            switchToDM(user.id);
        },
        function (xhr) {
            console.error("User '" + username + "' was not found.");
            alert("User '" + username + "' was not found.");
        }
    );
}

function handleStartDMById(userId) {
    if (!userId) {
        console.error("User ID is required to start a DM");
        return;
    }

    fetchUserById(
        userId,
        function (user) {
            // Add user to DM list if not already present
            var dmContainer = document.getElementById("dm-container");
            var existing = dmContainer.querySelector(`.dm-entry[data-user-id="${user.id}"]`);
            if (!existing) {
                var dmEntry = document.createElement("div");
                dmEntry.className = "dm-entry";
                dmEntry.setAttribute("data-user-id", user.id);
                dmEntry.textContent = user.name;
                dmEntry.onclick = function () {
                    switchToDM(user.id);
                };
                dmContainer.appendChild(dmEntry);
            }

            // User found, switch to DM
            switchToDM(user.id);
        },
        function (xhr) {
            console.error("User '" + userId + "' was not found.");
            alert("User '" + userId + "' was not found.");
        }
    );
}

function handleStartDMFromInput() {
    var usernameInput = document.getElementById("dm-join-input");
    if (!usernameInput) return;

    var username = usernameInput.value.trim();
    if (!username) return;

    usernameInput.value = "";
    handleStartDMByName(username);
}

function onIrcTokenResponse(xhr) {
    var response = JSON.parse(xhr.responseText);
    if (!response || !response.token) {
        console.error("Invalid response format or missing token:", response);
        alert("Failed to retrieve your IRC password. Please try again later!");
        return;
    }
    initializeSocket(currentUsername, response.token);
}

function onIrcTokenFailure(xhr) {
    console.error("Failed to retrieve IRC token:", xhr);
    alert("Failed to retrieve your IRC password. Please try again later!");
}

$(document).ready(function () {
    if (!isLoggedIn()) return;

    // Replace " " with "_" for IRC compatibility
    currentUsername = currentUsername.replace(/ /g, "_");

    // Check for target query parameter
    var target = getQueryParameter("target");
    if (target) {
        pendingTarget = target;
    }

    initializeChatHandlers();
    fetchIrcToken(onIrcTokenResponse, onIrcTokenFailure);
});
